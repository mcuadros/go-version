package version

import (
	"regexp"
	"strconv"
	"strings"
)

type Constraint struct {
	Sign    string
	Version string
}

func (self *Constraint) String() string {
	return strings.Trim(self.Sign+" "+self.Version, " ")
}

func ParseConstraints(constraint string) []*Constraint {
	result := RegFind(`(?i)^([^,\s]*?)@(stable|RC|beta|alpha|dev)$`, constraint)
	if result != nil {
		constraint = result[1]
		if constraint == "" {
			constraint = "*"
		}
	}

	result = RegFind(`(?i)^(dev-[^,\s@]+?|[^,\s@]+?\.x-dev)#.+$`, constraint)
	if result != nil {
		if result[1] != "" {
			constraint = result[1]
		}
	}

	constraints := RegSplit(`\s*,\s*`, strings.Trim(constraint, " "))

	if len(constraints) > 1 {
		output := make([]*Constraint, 0)
		for _, part := range constraints {
			output = append(output, parseConstraint(part)...)
		}

		return output
	}

	return parseConstraint(constraints[0])
}

func parseConstraint(constraint string) []*Constraint {

	stabilityModifier := ""

	result := RegFind(`(?i)^([^,\s]+?)@(stable|RC|beta|alpha|dev)$`, constraint)
	if result != nil {
		constraint = result[1]
		if result[2] != "stable" {
			stabilityModifier = result[2]
		}
	}

	result = RegFind(`^[x*](\.[x*])*$`, constraint)
	if result != nil {
		return make([]*Constraint, 0)
	}

	highVersion := ""
	lowVersion := ""

	result = RegFind(`(?i)^~(\d+)(?:\.(\d+))?(?:\.(\d+))?(?:\.(\d+))?`+modifierRegex+`?$`, constraint)
	if result != nil {
		if len(result) > 4 && result[4] != "" {
			last, _ := strconv.Atoi(result[3])
			highVersion = result[1] + "." + result[2] + "." + strconv.Itoa(last+1) + ".0-dev"
			lowVersion = result[1] + "." + result[2] + "." + result[3] + "." + result[4]
		} else if len(result) > 3 && result[3] != "" {
			last, _ := strconv.Atoi(result[2])
			highVersion = result[1] + "." + strconv.Itoa(last+1) + ".0.0-dev"
			lowVersion = result[1] + "." + result[2] + "." + result[3] + ".0"
		} else {
			last, _ := strconv.Atoi(result[1])
			highVersion = strconv.Itoa(last+1) + ".0.0.0-dev"
			if len(result) > 2 && result[2] != "" {
				lowVersion = result[1] + "." + result[2] + ".0.0"
			} else {
				lowVersion = result[1] + ".0.0.0"
			}
		}

		if len(result) > 5 && result[5] != "" {
			lowVersion = lowVersion + "-" + expandStability(result[5])

		}

		if len(result) > 6 && result[6] != "" {
			lowVersion = lowVersion + result[6]
		}

		if len(result) > 7 && result[7] != "" {
			lowVersion = lowVersion + "-dev"
		}

		return []*Constraint{
			{">=", lowVersion},
			{"<", highVersion},
		}
	}

	result = RegFind(`^(\d+)(?:\.(\d+))?(?:\.(\d+))?\.[x*]$`, constraint)
	if result != nil {
		if len(result) > 3 && result[3] != "" {
			highVersion = result[1] + "." + result[2] + "." + result[3] + ".9999999"
			if result[3] == "0" {
				last, _ := strconv.Atoi(result[2])
				lowVersion = result[1] + "." + strconv.Itoa(last-1) + ".9999999.9999999"
			} else {
				last, _ := strconv.Atoi(result[3])
				lowVersion = result[1] + "." + result[2] + "." + strconv.Itoa(last-1) + ".9999999"
			}

		} else if len(result) > 2 && result[2] != "" {
			highVersion = result[1] + "." + result[2] + ".9999999.9999999"
			if result[2] == "0" {
				last, _ := strconv.Atoi(result[1])
				lowVersion = strconv.Itoa(last-1) + ".9999999.9999999.9999999"
			} else {
				last, _ := strconv.Atoi(result[2])
				lowVersion = result[1] + "." + strconv.Itoa(last-1) + ".9999999.9999999"
			}

		} else {
			highVersion = result[1] + ".9999999.9999999.9999999"
			if result[1] == "0" {
				return []*Constraint{{"<", highVersion}}
			} else {
				last, _ := strconv.Atoi(result[1])
				lowVersion = strconv.Itoa(last-1) + ".9999999.9999999.9999999"
			}
		}

		return []*Constraint{
			{">", lowVersion},
			{"<", highVersion},
		}
	}

	// match operators constraints
	result = RegFind(`^(<>|!=|>=?|<=?|==?)?\s*(.*)`, constraint)
	if result != nil {
		version := Normalize(result[2])

		if stabilityModifier != "" && parseStability(version) == "stable" {
			version = version + "-" + stabilityModifier
		} else if result[1] == "<" {
			match := RegFind(`(?i)-stable$`, result[2])
			if match == nil {
				version = version + "-dev"
			}
		}

		if len(result) > 1 && result[1] != "" {
			return []*Constraint{{result[1], version}}
		} else {
			return []*Constraint{{"=", version}}

		}
	}

	return []*Constraint{{constraint, stabilityModifier}}
}

func parseStability(version string) string {
	version = regexp.MustCompile(`(?i)#.+$`).ReplaceAllString(version, " ")
	version = strings.ToLower(version)

	if strings.HasPrefix(version, "dev-") || strings.HasSuffix(version, "-dev") {
		return "dev"
	}

	result := RegFind(`(?i)^v?(\d{1,3})(\.\d+)?(\.\d+)?(\.\d+)?`+modifierRegex+`$`, version)
	if result != nil {
		if len(result) > 3 {
			return "dev"
		}
	}

	if result[1] != "" {
		if "beta" == result[1] || "b" == result[1] {
			return "beta"
		}
		if "alpha" == result[1] || "a" == result[1] {
			return "alpha"
		}
		if "rc" == result[1] {
			return "RC"
		}
	}

	return "stable"
}

func RegFind(pattern, subject string) []string {
	reg := regexp.MustCompile(pattern)
	matched := reg.FindAllStringSubmatch(subject, -1)

	if matched != nil {
		return matched[0]
	}

	return nil
}

func RegSplit(pattern, subject string) []string {
	reg := regexp.MustCompile(pattern)
	indexes := reg.FindAllStringIndex(subject, -1)

	laststart := 0
	result := make([]string, len(indexes)+1)

	for i, element := range indexes {
		result[i] = subject[laststart:element[0]]
		laststart = element[1]
	}

	result[len(indexes)] = subject[laststart:len(subject)]
	return result
}
