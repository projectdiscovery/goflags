package goflags

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"
	"text/tabwriter"

	"github.com/cnf/structhash"
)

// FlagSet is a list of flags for an application
type FlagSet struct {
	description string
	configFile  string
	flagKeys    map[string]*flagData
}

type flagData struct {
	usage string
	short string
	long  string
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

// SetConfigFile sets a config file to read values from.
func (f *FlagSet) SetConfigFile(file string) {
	f.configFile = file
}

// Parse parses the flags provided to the library.
func (f *FlagSet) Parse() error {
	flag.CommandLine.Usage = f.usageFunc
	flag.Parse()

	if f.configFile != "" {
		return nil
	}
	return nil
}

// readConfigFile reads the config file and returns any flags
// that might have been set by the config file.
//
// Command line flags however always take prcedence over config file ones.
func (f *FlagSet) readConfigFile() error {
	return nil
}

// StringVarP adds a string flag with a shortname and longname
func (f *FlagSet) StringVarP(field *string, short, long, defaultValue, usage string) {
	flag.StringVar(field, short, defaultValue, usage)
	flag.StringVar(field, long, defaultValue, usage)

	flagData := &flagData{
		usage: usage,
		short: short,
		long:  long,
	}
	f.flagKeys[short] = flagData
	f.flagKeys[long] = flagData
}

// StringVar adds a string flag with a longname
func (f *FlagSet) StringVar(field *string, long, defaultValue, usage string) {
	flag.StringVar(field, long, defaultValue, usage)

	flagData := &flagData{
		usage: usage,
		long:  long,
	}
	f.flagKeys[long] = flagData
}

// BoolVarP adds a bool flag with a shortname and longname
func (f *FlagSet) BoolVarP(field *bool, long, short string, defaultValue bool, usage string) {
	flag.BoolVar(field, short, defaultValue, usage)
	flag.BoolVar(field, long, defaultValue, usage)

	flagData := &flagData{
		usage: usage,
		short: short,
		long:  long,
	}
	f.flagKeys[short] = flagData
	f.flagKeys[long] = flagData
}

// BoolVar adds a bool flag with a longname
func (f *FlagSet) BoolVar(field *bool, long string, defaultValue bool, usage string) {
	flag.BoolVar(field, long, defaultValue, usage)

	flagData := &flagData{
		usage: usage,
		long:  long,
	}
	f.flagKeys[long] = flagData
}

// IntVarP adds a int flag with a shortname and longname
func (f *FlagSet) IntVarP(field *int, long, short string, defaultValue int, usage string) {
	flag.IntVar(field, short, defaultValue, usage)
	flag.IntVar(field, long, defaultValue, usage)

	flagData := &flagData{
		usage: usage,
		short: short,
		long:  long,
	}
	f.flagKeys[short] = flagData
	f.flagKeys[long] = flagData
}

// IntVar adds a int flag with a longname
func (f *FlagSet) IntVar(field *int, long string, defaultValue int, usage string) {
	flag.IntVar(field, long, defaultValue, usage)

	flagData := &flagData{
		usage: usage,
		long:  long,
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

	flagData := &flagData{
		usage: usage,
		short: short,
		long:  long,
	}
	f.flagKeys[short] = flagData
	f.flagKeys[long] = flagData
}

// StringSliceVar adds a string slice flag with a longname
func (f *FlagSet) StringSliceVar(field *StringSlice, long string, defaultValue []string, usage string) {
	for _, item := range defaultValue {
		_ = field.Set(item)
	}
	flag.Var(field, long, usage)

	flagData := &flagData{
		usage: usage,
		long:  long,
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
