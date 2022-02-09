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

func ToFileNormalizedStringSlice(value string) ([]string, error) {
	return toStringSlice(value, DefaultFileNormalizedStringSliceOptions)
}

var DefaultFileNormalizedStringSliceOptions = Options{
	IsEmpty:    isEmpty,
	Normalize:  normalize,
	IsFromFile: func(s string) bool { return true },
}
