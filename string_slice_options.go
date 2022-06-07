package goflags

// StringSliceOptions represents the default string slice (list of items)
// Tokenization: None
// Normalization: None
// Type: []string
// Example: -flag value1 -flag value2 => {"value1", "value2"}
var StringSliceOptions = Options{}

// CommaSeparatedStringSliceOptions represents a list of comma separated items
// Tokenization: Comma
// Normalization: None
// Type: []string
// Example: -flag value1,value2 => {"value1", "value2"}
var CommaSeparatedStringSliceOptions = Options{
	IsEmpty: isEmpty,
}

// FileCommaSeparatedStringSliceOptions represents a list of comma separated files containing items
// Tokenization: Comma
// Normalization: None
// Type: []string
// Example: -flag test.txt # with test.txt containing on multiple lines value1 and nvalue2 => {"value1", "value2"}
var FileCommaSeparatedStringSliceOptions = Options{
	IsEmpty:    isEmpty,
	IsFromFile: func(s string) bool { return true },
}

// NormalizedOriginalStringSliceOptions represents a list of items
// Tokenization: None
// Normalization: Standard
// Type: []string
// Example: -flag /value/1 -flag value2 => {"/value/1", "value2"}
var NormalizedOriginalStringSliceOptions = Options{
	IsEmpty:   isEmpty,
	Normalize: normalize,
}

// FileNormalizedStringSliceOptions represents a list of path items
// Tokenization: Comma
// Normalization: Standard
// Type: []string
// Example: -flag /value/1 -flag value2 => {"/value/1", "value2"}
var FileNormalizedStringSliceOptions = Options{
	IsEmpty:    isEmpty,
	Normalize:  normalizeLowercase,
	IsFromFile: func(s string) bool { return true },
}

// FileStringSliceOptions represents a list of items stored in a file
// Tokenization: Standard
// Normalization: Standard
var FileStringSliceOptions = Options{
	IsEmpty:    isEmpty,
	Normalize:  normalizeTrailingParts,
	IsFromFile: func(s string) bool { return true },
}

// NormalizedStringSliceOptions represents a list of items
// Tokenization: Comma
// Normalization: Standard
var NormalizedStringSliceOptions = Options{
	IsEmpty:   isEmpty,
	Normalize: normalizeLowercase,
}
