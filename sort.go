package version

import (
	"sort"
)

func Sort(versionStrings []string) {
	versions := versionSlice(versionStrings)
	sort.Sort(versions)
}

type versionSlice []string

func (s versionSlice) Len() int {
	return len(s)
}

func (s versionSlice) Less(i, j int) bool {
	cmp := CompareSimple(s[i], s[j])
	if cmp == 0 {
		return s[i] < s[j]
	}
	return cmp < 0
}

func (s versionSlice) Swap(i, j int) {
	tmp := s[j]
	s[j] = s[i]
	s[i] = tmp
}
