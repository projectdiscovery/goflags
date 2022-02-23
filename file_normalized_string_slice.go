package goflags

// FileNormalizedStringSlice is a slice of strings
type FileNormalizedStringSlice []string

// Set appends a value to the string slice.
func (fileNormalizedStringSlice *FileNormalizedStringSlice) Set(value string) error {
	slice, err := ToFileNormalizedStringSlice(value)
	if err != nil {
		return err
	}
	*fileNormalizedStringSlice = append(*fileNormalizedStringSlice, slice...)
	return nil
}

func (fileNormalizedStringSlice FileNormalizedStringSlice) String() string {
	return ToString(fileNormalizedStringSlice)
}

// FileOriginalNormalizedStringSlice is a slice of strings without normalization
type FileOriginalNormalizedStringSlice []string

// Set appends a value to the string slice.
func (fileNormalizedStringSlice *FileOriginalNormalizedStringSlice) Set(value string) error {
	slice, err := ToFileOriginalNormalizedStringSlice(value)
	if err != nil {
		return err
	}
	*fileNormalizedStringSlice = append(*fileNormalizedStringSlice, slice...)
	return nil
}

func (fileNormalizedStringSlice FileOriginalNormalizedStringSlice) String() string {
	return ToString(fileNormalizedStringSlice)
}

var DefaultFileNormalizedStringSliceOptions = Options{
	IsEmpty:    isEmpty,
	Normalize:  normalizeLowercase,
	IsFromFile: func(s string) bool { return true },
}

func ToFileNormalizedStringSlice(value string) ([]string, error) {
	return toStringSlice(value, DefaultFileNormalizedStringSliceOptions)
}

var DefaultFileOriginalNormalizedStringSliceOptions = Options{
	IsEmpty:    isEmpty,
	Normalize:  normalize,
	IsFromFile: func(s string) bool { return true },
}

func ToFileOriginalNormalizedStringSlice(value string) ([]string, error) {
	return toStringSlice(value, DefaultFileOriginalNormalizedStringSliceOptions)
}
