package goflags

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/cnf/structhash"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// FlagSet is a list of flags for an application
type FlagSet struct {
	description string
	flagKeys    map[string]*flagData
}

type flagData struct {
	usage        string
	short        string
	long         string
	defaultValue string
}

// Hash returns the unique hash for a flagData structure
// NOTE: Hash panics when the structure cannot be hashed.
func (f *flagData) Hash() string {
	hash, _ := structhash.Hash(f, 1)
	return hash
}

// New creates a new flagset structure for the application
func New() *FlagSet {
	return &FlagSet{flagKeys: make(map[string]*flagData)}
}

// SetDescription sets the description field for a flagset to a value.
func (f *FlagSet) SetDescription(description string) {
	f.description = description
}

// MergeConfigFile reads a config file to merge values from.
func (f *FlagSet) MergeConfigFile(file string) error {
	return f.readConfigFile(file)
}

// Parse parses the flags provided to the library.
func (f *FlagSet) Parse() {
	flag.CommandLine.Usage = f.usageFunc
	flag.Parse()

	appName := os.Args[0]
	homepath, err := os.UserHomeDir()
	if err != nil {
		return
	}

	config := path.Join(homepath, ".config", appName, "default-config.yaml")
	if _, err := os.Stat(config); os.IsNotExist(err) {
		configData := f.generateDefaultConfig()
		_ = ioutil.WriteFile(config, configData, os.ModePerm)
	}
}

// generateDefaultConfig generates a default YAML config file for a flagset.
func (f *FlagSet) generateDefaultConfig() []byte {
	hashes := make(map[string]struct{})
	configBuffer := &bytes.Buffer{}

	for _, v := range f.flagKeys {
		dataHash := v.Hash()
		if _, ok := hashes[dataHash]; ok {
			continue
		}
		hashes[dataHash] = struct{}{}

		configBuffer.WriteString("# ")
		configBuffer.WriteString(v.usage)
		configBuffer.WriteString("\n")
		configBuffer.WriteString("#")
		configBuffer.WriteString(v.long)
		configBuffer.WriteString(": ")
		configBuffer.WriteString(v.defaultValue)
		configBuffer.WriteString("\n")
		configBuffer.WriteString("\n")
	}
	return bytes.TrimSuffix(configBuffer.Bytes(), []byte("\n\n"))
}

// readConfigFile reads the config file and returns any flags
// that might have been set by the config file.
//
// Command line flags however always take prcedence over config file ones.
func (f *FlagSet) readConfigFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return errors.Wrap(err, "could not open config file")
	}
	defer file.Close()

	data := make(map[string]string)
	err = yaml.NewDecoder(file).Decode(&data)
	if err != nil {
		return errors.Wrap(err, "could not unmarshal config file")
	}
	flag.CommandLine.VisitAll(func(fl *flag.Flag) {
		item, ok := data[fl.Name]
		value := fl.Value.String()
		if strings.EqualFold(fl.DefValue, value) && ok {
			_ = fl.Value.Set(item)
		}
	})
	return nil
}

// StringVarP adds a string flag with a shortname and longname
func (f *FlagSet) StringVarP(field *string, short, long, defaultValue, usage string) {
	flag.StringVar(field, short, defaultValue, usage)
	flag.StringVar(field, long, defaultValue, usage)

	flagData := &flagData{
		usage:        usage,
		short:        short,
		long:         long,
		defaultValue: defaultValue,
	}
	f.flagKeys[short] = flagData
	f.flagKeys[long] = flagData
}

// StringVar adds a string flag with a longname
func (f *FlagSet) StringVar(field *string, long, defaultValue, usage string) {
	flag.StringVar(field, long, defaultValue, usage)

	flagData := &flagData{
		usage:        usage,
		long:         long,
		defaultValue: defaultValue,
	}
	f.flagKeys[long] = flagData
}

// BoolVarP adds a bool flag with a shortname and longname
func (f *FlagSet) BoolVarP(field *bool, long, short string, defaultValue bool, usage string) {
	flag.BoolVar(field, short, defaultValue, usage)
	flag.BoolVar(field, long, defaultValue, usage)

	flagData := &flagData{
		usage:        usage,
		short:        short,
		long:         long,
		defaultValue: strconv.FormatBool(defaultValue),
	}
	f.flagKeys[short] = flagData
	f.flagKeys[long] = flagData
}

// BoolVar adds a bool flag with a longname
func (f *FlagSet) BoolVar(field *bool, long string, defaultValue bool, usage string) {
	flag.BoolVar(field, long, defaultValue, usage)

	flagData := &flagData{
		usage:        usage,
		long:         long,
		defaultValue: strconv.FormatBool(defaultValue),
	}
	f.flagKeys[long] = flagData
}

// IntVarP adds a int flag with a shortname and longname
func (f *FlagSet) IntVarP(field *int, long, short string, defaultValue int, usage string) {
	flag.IntVar(field, short, defaultValue, usage)
	flag.IntVar(field, long, defaultValue, usage)

	flagData := &flagData{
		usage:        usage,
		short:        short,
		long:         long,
		defaultValue: strconv.Itoa(defaultValue),
	}
	f.flagKeys[short] = flagData
	f.flagKeys[long] = flagData
}

// IntVar adds a int flag with a longname
func (f *FlagSet) IntVar(field *int, long string, defaultValue int, usage string) {
	flag.IntVar(field, long, defaultValue, usage)

	flagData := &flagData{
		usage:        usage,
		long:         long,
		defaultValue: strconv.Itoa(defaultValue),
	}
	f.flagKeys[long] = flagData
}

// StringSliceVarP adds a string slice flag with a shortname and longname
func (f *FlagSet) StringSliceVarP(field *StringSlice, long, short string, defaultValue []string, usage string) {
	for _, item := range defaultValue {
		_ = field.Set(item)
	}

	flag.Var(field, short, usage)
	flag.Var(field, long, usage)

	defaultBuilder := &strings.Builder{}
	defaultBuilder.WriteString("[")
	for i, k := range *field {
		defaultBuilder.WriteString("\"")
		defaultBuilder.WriteString(k)
		defaultBuilder.WriteString("\"")
		if i != len(*field)-1 {
			defaultBuilder.WriteString(", ")
		}
	}
	defaultBuilder.WriteString("]")

	flagData := &flagData{
		usage:        usage,
		short:        short,
		long:         long,
		defaultValue: defaultBuilder.String(),
	}
	f.flagKeys[short] = flagData
	f.flagKeys[long] = flagData
}

// StringSliceVar adds a string slice flag with a longname
func (f *FlagSet) StringSliceVar(field *StringSlice, long string, defaultValue []string, usage string) {
	for _, item := range defaultValue {
		_ = field.Set(item)
	}

	defaultBuilder := &strings.Builder{}
	defaultBuilder.WriteString("[")
	for i, k := range *field {
		defaultBuilder.WriteString("\"")
		defaultBuilder.WriteString(k)
		defaultBuilder.WriteString("\"")
		if i != len(*field)-1 {
			defaultBuilder.WriteString(", ")
		}
	}
	defaultBuilder.WriteString("]")

	flag.Var(field, long, usage)

	flagData := &flagData{
		usage:        usage,
		long:         long,
		defaultValue: defaultBuilder.String(),
	}
	f.flagKeys[long] = flagData
}

func (f *FlagSet) usageFunc() {
	hashes := make(map[string]struct{})

	fmt.Fprintf(flag.CommandLine.Output(), "%s\n\n", f.description)
	fmt.Fprintf(flag.CommandLine.Output(), "Usage:\n  %s [flags]\n\n", os.Args[0])
	fmt.Fprintf(flag.CommandLine.Output(), "Flags:\n")

	w := tabwriter.NewWriter(flag.CommandLine.Output(), 0, 0, 1, ' ', 0)
	flag.CommandLine.VisitAll(func(fl *flag.Flag) {
		data, _ := f.flagKeys[fl.Name]

		dataHash := data.Hash()
		if _, ok := hashes[dataHash]; ok {
			return // Don't print the value if printed previously
		}
		hashes[dataHash] = struct{}{}

		var s string
		if data.short != "" {
			s = fmt.Sprintf("  \t-%s, -%s", data.short, data.long) // Two spaces before -; see next two comments.
		} else {
			s = fmt.Sprintf("  \t-%s", data.long)
		}
		name, usage := flag.UnquoteUsage(fl)
		if len(name) > 0 {
			s += " " + name
		}
		s += "\t\t"
		s += strings.ReplaceAll(usage, "\n", "\n    \t")

		if !isZeroValue(fl, fl.DefValue) {
			if _, ok := fl.Value.(*stringValue); ok {
				// put quotes on the value
				s += fmt.Sprintf(" (default %q)", fl.DefValue)
			} else {
				s += fmt.Sprintf(" (default %v)", fl.DefValue)
			}
		}
		fmt.Fprint(w, s, "\n")
	})
	w.Flush()
}

// isZeroValue determines whether the string represents the zero
// value for a flag.
func isZeroValue(f *flag.Flag, value string) bool {
	// Build a zero value of the flag's Value type, and see if the
	// result of calling its String method equals the value passed in.
	// This works unless the Value type is itself an interface type.
	typ := reflect.TypeOf(f.Value)
	var z reflect.Value
	if typ.Kind() == reflect.Ptr {
		z = reflect.New(typ.Elem())
	} else {
		z = reflect.Zero(typ)
	}
	return value == z.Interface().(flag.Value).String()
}
