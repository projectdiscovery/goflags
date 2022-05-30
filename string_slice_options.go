package goflags

var StringSliceOptions = Options{}

// No value normalization is happening.
var CommaSeparatedStringSliceOptions = Options{
	IsEmpty: isEmpty,
}

// No value normalization is happening.
var FileCommaSeparatedStringSliceOptions = Options{
	IsEmpty:    isEmpty,
	IsFromFile: func(s string) bool { return true },
}

var NormalizedOriginalStringSliceOptions = Options{
	IsEmpty:   isEmpty,
	Normalize: normalize,
}

// NormalizedStringSliceVarP adds a path slice flag with a shortname and longname.
// It supports comma separated values, that are normalized
// (lower-cased, stripped of any leading and trailing whitespaces and quotes)
var FileNormalizedStringSliceOptions = Options{
	IsEmpty:    isEmpty,
	Normalize:  normalizeLowercase,
	IsFromFile: func(s string) bool { return true },
}

var FileStringSliceOptions = Options{
	IsEmpty:    isEmpty,
	Normalize:  normalizeTrailingParts,
	IsFromFile: func(s string) bool { return true },
}

var NormalizedStringSliceOptions = Options{
	IsEmpty:   isEmpty,
	Normalize: normalizeLowercase,
}
