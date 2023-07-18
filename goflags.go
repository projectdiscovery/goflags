package goflags

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/cnf/structhash"
	fileutil "github.com/projectdiscovery/utils/file"
	folderutil "github.com/projectdiscovery/utils/folder"
	"golang.org/x/exp/maps"
	"gopkg.in/yaml.v3"
)

// FlagSet is a list of flags for an application
type FlagSet struct {
	CaseSensitive  bool
	Marshal        bool
	description    string
	flagKeys       InsertionOrderedMap
	groups         []groupData
	CommandLine    *flag.FlagSet
	configFilePath string

	// OtherOptionsGroupName is the name for all flags not in a group
	OtherOptionsGroupName string
	configOnlyKeys        InsertionOrderedMap
}

type groupData struct {
	name        string
	description string
}

type FlagData struct {
	usage        string
	short        string
	long         string
	group        string // unused unless set later
	defaultValue interface{}
	skipMarshal  bool
	field        flag.Value
}

// Group sets the group for a flag data
func (flagData *FlagData) Group(name string) {
	flagData.group = name
}

// NewFlagSet creates a new flagSet structure for the application
func NewFlagSet() *FlagSet {
	flag.CommandLine.ErrorHandling()
	return &FlagSet{
		flagKeys:              newInsertionOrderedMap(),
		OtherOptionsGroupName: "other options",
		CommandLine:           flag.NewFlagSet(os.Args[0], flag.ExitOnError),
		configOnlyKeys:        newInsertionOrderedMap(),
	}
}

func newInsertionOrderedMap() InsertionOrderedMap {
	return InsertionOrderedMap{values: make(map[string]*FlagData)}
}

// Hash returns the unique hash for a flagData structure
// NOTE: Hash panics when the structure cannot be hashed.
func (flagData *FlagData) Hash() string {
	hash, _ := structhash.Hash(flagData, 1)
	return hash
}

// SetDescription sets the description field for a flagSet to a value.
func (flagSet *FlagSet) SetDescription(description string) {
	flagSet.description = description
}

// SetGroup sets a group with name and description for the command line options
//
// The order in which groups are passed is also kept as is, similar to flags.
func (flagSet *FlagSet) SetGroup(name, description string) {
	flagSet.groups = append(flagSet.groups, groupData{name: name, description: description})
}

// MergeConfigFile reads a config file to merge values from.
func (flagSet *FlagSet) MergeConfigFile(file string) error {
	return flagSet.readConfigFile(file)
}

// Parse parses the flags provided to the library.
func (flagSet *FlagSet) Parse() error {
	flagSet.CommandLine.SetOutput(os.Stdout)
	flagSet.CommandLine.Usage = flagSet.usageFunc
	_ = flagSet.CommandLine.Parse(os.Args[1:])
	configFilePath, _ := flagSet.GetConfigFilePath()

	// migrate data from old config dir to new one
	// Ref: https://github.com/projectdiscovery/nuclei/issues/3576
	flagSet.migrateConfigDir()

	// if config file doesn't exist, create one
	if !fileutil.FileExists(configFilePath) {
		configData := flagSet.generateDefaultConfig()
		configFileDir := flagSet.GetToolConfigDir()
		if !fileutil.FolderExists(configFileDir) {
			_ = fileutil.CreateFolder(configFileDir)
		}
		return os.WriteFile(configFilePath, configData, os.ModePerm)
	}

	_ = flagSet.MergeConfigFile(configFilePath) // try to read default config after parsing flags
	return nil
}

func (flagSet *FlagSet) migrateConfigDir() {
	// migration condition
	// 1. old config dir exists
	// 2. new config dir doesn't exist
	// 3. old config dir is not same as new config dir

	toolConfigDir := flagSet.GetToolConfigDir()
	if toolConfigDir != oldAppConfigDir && fileutil.FolderExists(oldAppConfigDir) && !fileutil.FolderExists(toolConfigDir) {
		_ = fileutil.CreateFolder(toolConfigDir)
		// move old config dir to new one
		_ = folderutil.SyncDirectory(oldAppConfigDir, toolConfigDir)
	}
}

// generateDefaultConfig generates a default YAML config file for a flagset.
func (flagSet *FlagSet) generateDefaultConfig() []byte {
	hashes := make(map[string]struct{})
	configBuffer := &bytes.Buffer{}
	configBuffer.WriteString("# ")
	configBuffer.WriteString(path.Base(os.Args[0]))
	configBuffer.WriteString(" config file\n# generated by https://github.com/projectdiscovery/goflags\n\n")

	// Attempts to marshal natively if proper flag is set, in case of errors fallback to normal mechanism
	if flagSet.Marshal {
		flagsToMarshall := make(map[string]interface{})

		flagSet.flagKeys.forEach(func(key string, data *FlagData) {
			if !data.skipMarshal {
				flagsToMarshall[key] = data.defaultValue
			}
		})

		flagSetBytes, err := yaml.Marshal(flagsToMarshall)
		if err == nil {
			configBuffer.Write(flagSetBytes)
			return configBuffer.Bytes()
		}
	}

	flagSet.flagKeys.forEach(func(key string, data *FlagData) {
		dataHash := data.Hash()
		if _, ok := hashes[dataHash]; ok {
			return
		}
		hashes[dataHash] = struct{}{}

		configBuffer.WriteString("# ")
		configBuffer.WriteString(strings.ToLower(data.usage))
		configBuffer.WriteString("\n")
		configBuffer.WriteString("#")
		configBuffer.WriteString(data.long)
		configBuffer.WriteString(": ")
		switch dv := data.defaultValue.(type) {
		case string:
			configBuffer.WriteString(dv)
		case flag.Value:
			configBuffer.WriteString(dv.String())
		case StringSlice:
			configBuffer.WriteString(dv.String())
		}

		configBuffer.WriteString("\n\n")
	})

	return bytes.TrimSuffix(configBuffer.Bytes(), []byte("\n\n"))
}

// CreateGroup within the flagset
func (flagSet *FlagSet) CreateGroup(groupName, description string, flags ...*FlagData) {
	flagSet.SetGroup(groupName, description)
	for _, currentFlag := range flags {
		currentFlag.Group(groupName)
	}
}

// readConfigFile reads the config file and returns any flags
// that might have been set by the config file.
//
// Command line flags however always take precedence over config file ones.
func (flagSet *FlagSet) readConfigFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	data := make(map[string]interface{})
	err = yaml.NewDecoder(file).Decode(&data)
	if err != nil {
		return err
	}
	flagSet.CommandLine.VisitAll(func(fl *flag.Flag) {
		item, ok := data[fl.Name]
		value := fl.Value.String()

		if strings.EqualFold(fl.DefValue, value) && ok {
			switch itemValue := item.(type) {
			case string:
				_ = fl.Value.Set(itemValue)
			case bool:
				_ = fl.Value.Set(strconv.FormatBool(itemValue))
			case int:
				_ = fl.Value.Set(strconv.Itoa(itemValue))
			case time.Duration:
				_ = fl.Value.Set(itemValue.String())
			case []interface{}:
				for _, v := range itemValue {
					vStr, ok := v.(string)
					if ok {
						_ = fl.Value.Set(vStr)
					}
				}
			}
		}
	})

	flagSet.configOnlyKeys.forEach(func(key string, flagData *FlagData) {
		item, ok := data[key]
		if ok {
			fl := flag.Lookup(key)
			if fl == nil {
				flag.Var(flagData.field, key, flagData.usage)
				fl = flag.Lookup(key)
			}

			switch data := item.(type) {
			case string:
				_ = fl.Value.Set(data)
			case bool:
				_ = fl.Value.Set(strconv.FormatBool(data))
			case int:
				_ = fl.Value.Set(strconv.Itoa(data))
			case []interface{}:
				for _, v := range data {
					vStr, ok := v.(string)
					if ok {
						_ = fl.Value.Set(vStr)
					}
				}
			}
		}
	})
	return nil
}

// VarP adds a Var flag with a shortname and longname
func (flagSet *FlagSet) VarP(field flag.Value, long, short, usage string) *FlagData {
	flagData := &FlagData{
		usage:        usage,
		long:         long,
		defaultValue: field,
	}
	if short != "" {
		flagData.short = short
		flagSet.CommandLine.Var(field, short, usage)
		flagSet.flagKeys.Set(short, flagData)
	}
	flagSet.CommandLine.Var(field, long, usage)
	flagSet.flagKeys.Set(long, flagData)
	return flagData
}

// Var adds a Var flag with a longname
func (flagSet *FlagSet) Var(field flag.Value, long, usage string) *FlagData {
	return flagSet.VarP(field, long, "", usage)
}

// StringVarEnv adds a string flag with a shortname and longname with a default value read from env variable
// with a default value fallback
func (flagSet *FlagSet) StringVarEnv(field *string, long, short, defaultValue, envName, usage string) *FlagData {
	if envValue, exists := os.LookupEnv(envName); exists {
		defaultValue = envValue
	}
	return flagSet.StringVarP(field, long, short, defaultValue, usage)
}

// StringVarP adds a string flag with a shortname and longname
func (flagSet *FlagSet) StringVarP(field *string, long, short, defaultValue, usage string) *FlagData {
	flagData := &FlagData{
		usage:        usage,
		long:         long,
		defaultValue: defaultValue,
	}
	if short != "" {
		flagData.short = short
		flagSet.CommandLine.StringVar(field, short, defaultValue, usage)
		flagSet.flagKeys.Set(short, flagData)
	}
	flagSet.CommandLine.StringVar(field, long, defaultValue, usage)
	flagSet.flagKeys.Set(long, flagData)
	return flagData
}

// StringVar adds a string flag with a longname
func (flagSet *FlagSet) StringVar(field *string, long, defaultValue, usage string) *FlagData {
	return flagSet.StringVarP(field, long, "", defaultValue, usage)
}

// BoolVarP adds a bool flag with a shortname and longname
func (flagSet *FlagSet) BoolVarP(field *bool, long, short string, defaultValue bool, usage string) *FlagData {
	flagData := &FlagData{
		usage:        usage,
		long:         long,
		defaultValue: strconv.FormatBool(defaultValue),
	}
	if short != "" {
		flagData.short = short
		flagSet.CommandLine.BoolVar(field, short, defaultValue, usage)
		flagSet.flagKeys.Set(short, flagData)
	}
	flagSet.CommandLine.BoolVar(field, long, defaultValue, usage)
	flagSet.flagKeys.Set(long, flagData)
	return flagData
}

// BoolVar adds a bool flag with a longname
func (flagSet *FlagSet) BoolVar(field *bool, long string, defaultValue bool, usage string) *FlagData {
	return flagSet.BoolVarP(field, long, "", defaultValue, usage)
}

// IntVarP adds a int flag with a shortname and longname
func (flagSet *FlagSet) IntVarP(field *int, long, short string, defaultValue int, usage string) *FlagData {
	flagData := &FlagData{
		usage:        usage,
		short:        short,
		long:         long,
		defaultValue: strconv.Itoa(defaultValue),
	}
	if short != "" {
		flagData.short = short
		flagSet.CommandLine.IntVar(field, short, defaultValue, usage)
		flagSet.flagKeys.Set(short, flagData)
	}
	flagSet.CommandLine.IntVar(field, long, defaultValue, usage)
	flagSet.flagKeys.Set(long, flagData)
	return flagData
}

// IntVar adds a int flag with a longname
func (flagSet *FlagSet) IntVar(field *int, long string, defaultValue int, usage string) *FlagData {
	return flagSet.IntVarP(field, long, "", defaultValue, usage)
}

// StringSliceVarP adds a string slice flag with a shortname and longname
// Use options to customize the behavior
func (flagSet *FlagSet) StringSliceVarP(field *StringSlice, long, short string, defaultValue StringSlice, usage string, options Options) *FlagData {
	optionMap[field] = options
	for _, defaultItem := range defaultValue {
		values, _ := ToStringSlice(defaultItem, options)
		for _, value := range values {
			_ = field.Set(value)
		}
	}
	optionDefaultValues[field] = *field
	flagData := &FlagData{
		usage:        usage,
		long:         long,
		defaultValue: defaultValue,
	}
	if short != "" {
		flagData.short = short
		flagSet.CommandLine.Var(field, short, usage)
		flagSet.flagKeys.Set(short, flagData)
	}
	flagSet.CommandLine.Var(field, long, usage)
	flagSet.flagKeys.Set(long, flagData)
	return flagData
}

// StringSliceVar adds a string slice flag with a longname
// Supports ONE value at a time. Adding multiple values require repeating the argument (-flag value1 -flag value2)
// No value normalization is happening.
func (flagSet *FlagSet) StringSliceVar(field *StringSlice, long string, defaultValue []string, usage string, options Options) *FlagData {
	return flagSet.StringSliceVarP(field, long, "", defaultValue, usage, options)
}

// StringSliceVarConfigOnly adds a string slice config value (without flag) with a longname
func (flagSet *FlagSet) StringSliceVarConfigOnly(field *StringSlice, long string, defaultValue []string, usage string) *FlagData {
	for _, item := range defaultValue {
		_ = field.Set(item)
	}
	flagData := &FlagData{
		usage:        usage,
		long:         long,
		defaultValue: defaultValue,
		field:        field,
	}
	flagSet.configOnlyKeys.Set(long, flagData)
	flagSet.flagKeys.Set(long, flagData)
	return flagData
}

// RuntimeMapVarP adds a runtime only map flag with a longname
func (flagSet *FlagSet) RuntimeMapVar(field *RuntimeMap, long string, defaultValue []string, usage string) *FlagData {
	return flagSet.RuntimeMapVarP(field, long, "", defaultValue, usage)
}

// RuntimeMapVarP adds a runtime only map flag with a shortname and longname
func (flagSet *FlagSet) RuntimeMapVarP(field *RuntimeMap, long, short string, defaultValue []string, usage string) *FlagData {
	for _, item := range defaultValue {
		_ = field.Set(item)
	}

	flagData := &FlagData{
		usage:        usage,
		long:         long,
		defaultValue: defaultValue,
		skipMarshal:  true,
	}

	if short != "" {
		flagData.short = short
		flagSet.CommandLine.Var(field, short, usage)
		flagSet.flagKeys.Set(short, flagData)
	}
	flagSet.CommandLine.Var(field, long, usage)
	flagSet.flagKeys.Set(long, flagData)
	return flagData
}

// PortVar adds a port flag with a longname
func (flagSet *FlagSet) PortVar(field *Port, long string, defaultValue []string, usage string) *FlagData {
	return flagSet.PortVarP(field, long, "", defaultValue, usage)
}

// PortVarP adds a port flag with a shortname and longname
func (flagSet *FlagSet) PortVarP(field *Port, long, short string, defaultValue []string, usage string) *FlagData {
	for _, item := range defaultValue {
		_ = field.Set(item)
	}
	portOptionDefaultValues[field] = maps.Clone(field.kv)

	flagData := &FlagData{
		usage:        usage,
		long:         long,
		defaultValue: defaultValue,
		skipMarshal:  true,
	}

	if short != "" {
		flagData.short = short
		flagSet.CommandLine.Var(field, short, usage)
		flagSet.flagKeys.Set(short, flagData)
	}
	flagSet.CommandLine.Var(field, long, usage)
	flagSet.flagKeys.Set(long, flagData)
	return flagData
}

// EnumVar adds a enum flag with a longname
func (flagSet *FlagSet) EnumVar(field *string, long string, defaultValue EnumVariable, usage string, allowedTypes AllowdTypes) *FlagData {
	return flagSet.EnumVarP(field, long, "", defaultValue, usage, allowedTypes)
}

// EnumVarP adds a enum flag with a shortname and longname
func (flagSet *FlagSet) EnumVarP(field *string, long, short string, defaultValue EnumVariable, usage string, allowedTypes AllowdTypes) *FlagData {
	var hasDefaultValue bool
	for k, v := range allowedTypes {
		if v == defaultValue {
			hasDefaultValue = true
			*field = k
		}
	}
	if !hasDefaultValue {
		panic("undefined default value")
	}
	flagData := &FlagData{
		usage:        usage,
		long:         long,
		defaultValue: *field,
	}
	if short != "" {
		flagData.short = short
		flagSet.CommandLine.Var(&EnumVar{allowedTypes, field}, short, usage)
		flagSet.flagKeys.Set(short, flagData)
	}
	flagSet.CommandLine.Var(&EnumVar{allowedTypes, field}, long, usage)
	flagSet.flagKeys.Set(long, flagData)
	return flagData
}

func (flagSet *FlagSet) usageFunc() {
	var helpAsked bool

	// Only show help usage if asked by user
	for _, arg := range os.Args {
		argStripped := strings.Trim(arg, "-")
		if argStripped == "h" || argStripped == "help" {
			helpAsked = true
		}
	}
	if !helpAsked {
		return
	}

	cliOutput := flagSet.CommandLine.Output()
	fmt.Fprintf(cliOutput, "%s\n\n", flagSet.description)
	fmt.Fprintf(cliOutput, "Usage:\n  %s [flags]\n\n", os.Args[0])
	fmt.Fprintf(cliOutput, "Flags:\n")

	writer := tabwriter.NewWriter(cliOutput, 0, 0, 1, ' ', 0)

	// If user has specified group with help and we have groups, return
	// with it's usage function
	if len(flagSet.groups) > 0 && len(os.Args) == 3 {
		group := flagSet.getGroupbyName(strings.ToLower(os.Args[2]))
		if group.name != "" {
			flagSet.displayGroupUsageFunc(newUniqueDeduper(), group, cliOutput, writer)
			return
		}
		flag := flagSet.getFlagByName(os.Args[2])
		if flag != nil {
			flagSet.displaySingleFlagUsageFunc(os.Args[2], flag, cliOutput, writer)
			return
		}
	}

	if len(flagSet.groups) > 0 {
		flagSet.usageFuncForGroups(cliOutput, writer)
	} else {
		flagSet.usageFuncInternal(writer)
	}
}

func (flagSet *FlagSet) getGroupbyName(name string) groupData {
	for _, group := range flagSet.groups {
		if strings.EqualFold(group.name, name) || strings.EqualFold(group.description, name) {
			return group
		}
	}
	return groupData{}
}

func (flagSet *FlagSet) getFlagByName(name string) *FlagData {
	var flagData *FlagData
	flagSet.flagKeys.forEach(func(key string, data *FlagData) {
		// check if the items are equal
		// - Case sensitive
		equal := flagSet.CaseSensitive && (data.long == name || data.short == name)
		// - Case insensitive
		equalFold := !flagSet.CaseSensitive && (strings.EqualFold(data.long, name) || strings.EqualFold(data.short, name))
		if equal || equalFold {
			flagData = data
			return
		}
	})
	return flagData
}

// usageFuncInternal prints usage for command line flags
func (flagSet *FlagSet) usageFuncInternal(writer *tabwriter.Writer) {
	uniqueDeduper := newUniqueDeduper()

	flagSet.flagKeys.forEach(func(key string, data *FlagData) {
		if currentFlag := flagSet.CommandLine.Lookup(key); currentFlag != nil {
			if !uniqueDeduper.isUnique(data) {
				return
			}
			result := createUsageString(data, currentFlag)
			fmt.Fprint(writer, result, "\n")
		}
	})
	writer.Flush()
}

// usageFuncForGroups prints usage for command line flags with grouping enabled
func (flagSet *FlagSet) usageFuncForGroups(cliOutput io.Writer, writer *tabwriter.Writer) {
	uniqueDeduper := newUniqueDeduper()

	var otherOptions []string
	for _, group := range flagSet.groups {
		otherOptions = append(otherOptions, flagSet.displayGroupUsageFunc(uniqueDeduper, group, cliOutput, writer)...)
	}

	// Print Any additional flag that may have been left
	if len(otherOptions) > 0 {
		fmt.Fprintf(cliOutput, "%s:\n", normalizeGroupDescription(flagSet.OtherOptionsGroupName))

		for _, option := range otherOptions {
			fmt.Fprint(writer, option, "\n")
		}
		writer.Flush()
	}
}

// displayGroupUsageFunc displays usage for a group
func (flagSet *FlagSet) displayGroupUsageFunc(uniqueDeduper *uniqueDeduper, group groupData, cliOutput io.Writer, writer *tabwriter.Writer) []string {
	fmt.Fprintf(cliOutput, "%s:\n", normalizeGroupDescription(group.description))

	var otherOptions []string
	flagSet.flagKeys.forEach(func(key string, data *FlagData) {
		if currentFlag := flagSet.CommandLine.Lookup(key); currentFlag != nil {
			if data.group == "" {
				if !uniqueDeduper.isUnique(data) {
					return
				}
				otherOptions = append(otherOptions, createUsageString(data, currentFlag))
				return
			}
			// Ignore the flag if it's not in our intended group
			if !strings.EqualFold(data.group, group.name) {
				return
			}
			if !uniqueDeduper.isUnique(data) {
				return
			}
			result := createUsageString(data, currentFlag)
			fmt.Fprint(writer, result, "\n")
		}
	})
	writer.Flush()
	fmt.Printf("\n")
	return otherOptions
}

// displaySingleFlagUsageFunc displays usage for a single flag
func (flagSet *FlagSet) displaySingleFlagUsageFunc(name string, data *FlagData, cliOutput io.Writer, writer *tabwriter.Writer) {
	if currentFlag := flagSet.CommandLine.Lookup(name); currentFlag != nil {
		result := createUsageString(data, currentFlag)
		fmt.Fprint(writer, result, "\n")
		writer.Flush()
	}
}

type uniqueDeduper struct {
	hashes map[string]interface{}
}

func newUniqueDeduper() *uniqueDeduper {
	return &uniqueDeduper{hashes: make(map[string]interface{})}
}

// isUnique returns true if the flag is unique during iteration
func (u *uniqueDeduper) isUnique(data *FlagData) bool {
	dataHash := data.Hash()
	if _, ok := u.hashes[dataHash]; ok {
		return false // Don't print the value if printed previously
	}
	u.hashes[dataHash] = struct{}{}
	return true
}

func isNotBlank(value string) bool {
	return len(strings.TrimSpace(value)) != 0
}

func createUsageString(data *FlagData, currentFlag *flag.Flag) string {
	valueType := reflect.TypeOf(currentFlag.Value)

	result := createUsageFlagNames(data)
	result += createUsageTypeAndDescription(currentFlag, valueType)
	result += createUsageDefaultValue(data, currentFlag, valueType)

	return result
}

func createUsageDefaultValue(data *FlagData, currentFlag *flag.Flag, valueType reflect.Type) string {
	if !isZeroValue(currentFlag, currentFlag.DefValue) {
		defaultValueTemplate := " (default "
		switch valueType.String() { // ugly hack because "flag.stringValue" is not exported from the parent library
		case "*flag.stringValue":
			defaultValueTemplate += "%q"
		default:
			defaultValueTemplate += "%v"
		}
		defaultValueTemplate += ")"
		return fmt.Sprintf(defaultValueTemplate, data.defaultValue)
	}
	return ""
}

func createUsageTypeAndDescription(currentFlag *flag.Flag, valueType reflect.Type) string {
	var result string

	flagDisplayType, usage := flag.UnquoteUsage(currentFlag)
	if len(flagDisplayType) > 0 {
		if flagDisplayType == "value" { // hardcoded in the goflags library
			switch valueType.Kind() {
			case reflect.Ptr:
				pointerTypeElement := valueType.Elem()
				switch pointerTypeElement.Kind() {
				case reflect.Slice, reflect.Array:
					switch pointerTypeElement.Elem().Kind() {
					case reflect.String:
						flagDisplayType = "string[]"
					default:
						flagDisplayType = "value[]"
					}
				}
			}
		}
		result += " " + flagDisplayType
	}

	result += "\t\t"
	result += strings.ReplaceAll(usage, "\n", "\n"+strings.Repeat(" ", 4)+"\t")
	return result
}

func createUsageFlagNames(data *FlagData) string {
	flagNames := strings.Repeat(" ", 2) + "\t"

	var validFlags []string
	addValidParam := func(value string) {
		if isNotBlank(value) {
			validFlags = append(validFlags, fmt.Sprintf("-%s", value))
		}
	}

	addValidParam(data.short)
	addValidParam(data.long)

	if len(validFlags) == 0 {
		panic("CLI arguments cannot be empty.")
	}

	flagNames += strings.Join(validFlags, ", ")
	return flagNames
}

// isZeroValue determines whether the string represents the zero
// value for a flag.
func isZeroValue(f *flag.Flag, value string) bool {
	// Build a zero value of the flag's Value type, and see if the
	// result of calling its String method equals the value passed in.
	// This works unless the Value type is itself an interface type.
	valueType := reflect.TypeOf(f.Value)
	var zeroValue reflect.Value
	if valueType.Kind() == reflect.Ptr {
		zeroValue = reflect.New(valueType.Elem())
	} else {
		zeroValue = reflect.Zero(valueType)
	}
	return value == zeroValue.Interface().(flag.Value).String()
}

// normalizeGroupDescription returns normalized description field for group
func normalizeGroupDescription(description string) string {
	return strings.ToUpper(description)
}
