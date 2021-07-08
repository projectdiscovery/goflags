package goflags

import (
	"github.com/pkg/errors"
	"strings"
)

type Severity int

const (
	Info Severity = iota
	Low
	Medium
	High
	Critical
	limit
)

var severityMappings = map[Severity]string{
	Info:     "info",
	Low:      "low",
	Medium:   "medium",
	High:     "high",
	Critical: "critical",
}

func toSeverity(valueToMap string) (Severity, error) {
	for key, currentValue := range severityMappings {
		if normalizeValue(valueToMap) == currentValue {
			return key, nil
		}
	}
	return -1, errors.New("Invalid severity: " + valueToMap)
}

func GetSupportedSeverities() []Severity {
	var result []Severity
	for index := Severity(0); index < limit; index++ {
		result = append(result, index)
	}
	return result
}

func normalizeValue(value string) string {
	return strings.TrimSpace(strings.ToLower(value))
}

func (severity Severity) normalize() string {
	return strings.TrimSpace(strings.ToLower(severity.String()))
}

func (severity Severity) String() string {
	return severityMappings[severity]
}
