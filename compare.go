package misc

import (
	"regexp"
	"strconv"
	"strings"
)

var regexpSigns = regexp.MustCompile(`[_\-+]`)
var regexpDotBeforeDigit = regexp.MustCompile(`([^.\d]+)`)
var regexpMultipleDots = regexp.MustCompile(`\.{2,}`)

var specialForms = map[string]int{
	"dev":   -6,
	"alpha": -5,
	"a":     -5,
	"beta":  -4,
	"b":     -4,
	"RC":    -3,
	"rc":    -3,
	"#":     -2,
	"p":     1,
	"pl":    1,
}

func CompareVersion(version1, version2, operator string) bool {
	compare := CompareVersionSimple(version1, version2)

	switch {
	case operator == ">" || operator == "gt":
		return compare > 0
	case operator == ">=" || operator == "ge":
		return compare >= 0
	case operator == "<=" || operator == "le":
		return compare <= 0
	case operator == "==" || operator == "=" || operator == "eq":
		return compare == 0
	case operator == "<>" || operator == "!=" || operator == "ne":
		return compare != 0
	case operator == "" || operator == "<" || operator == "lt":
		return compare < 0
	}

	return false
}

func CompareVersionSimple(version1, version2 string) int {
	var x, r, l int = 0, 0, 0

	v1, v2 := prepVersion(version1), prepVersion(version2)
	len1, len2 := len(v1), len(v2)

	if len1 > len2 {
		x = len1
	} else {
		x = len2
	}

	for i := 0; i < x; i++ {
		if i < len1 && i < len2 {
			if v1[i] == v2[i] {
				continue
			}
		}

		r = 0
		if i < len1 {
			r = numVersion(v1[i])
		}

		l = 0
		if i < len2 {
			l = numVersion(v2[i])
		}

		if r < l {
			return -1
		} else if r > l {
			return 1
		}
	}

	return 0
}

func prepVersion(version string) []string {
	if len(version) == 0 {
		return []string{""}
	}

	version = regexpSigns.ReplaceAllString(version, ".")
	version = regexpDotBeforeDigit.ReplaceAllString(version, ".$1.")
	version = regexpMultipleDots.ReplaceAllString(version, ".")

	return strings.Split(version, ".")
}

func numVersion(value string) int {
	if value == "" {
		return 0
	}

	if number, err := strconv.Atoi(value); err == nil {
		return number
	}

	if special, ok := specialForms[value]; ok {
		return special
	}

	return -7
}
