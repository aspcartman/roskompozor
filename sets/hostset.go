package sets

import (
	"sort"
	"strings"
)

type Hosts []string

func NewHostSet(hosts []string) Hosts {
	s := Hosts(hosts)
	s.Trim()
	s.Normalize()
	return s
}

func (s *Hosts) Normalize() {
	s.Sort()
	s.RemoveDuplicates()
}

func (s *Hosts) Add(hosts ...string) {
	*s = append(*s, NewHostSet(hosts)...)
	s.Normalize()
}

func (s Hosts) Contains(host string) bool {
	host = reverseString(host)
	i := sort.SearchStrings(s, host)
	return (i < len(s) && hostMatchesEntry(host, s[i])) || (i > 0 && hostMatchesEntry(host, s[i-1]))
}

func hostMatchesEntry(host, entry string) bool {
	return strings.HasPrefix(host, entry) && (len(host) == len(entry) || host[len(entry)] == '.')
}

func (s Hosts) Trim() {
	for i, h := range s {
		s[i] = reverseString(trimHost(h))
	}
}

func (s Hosts) Sort() {
	sort.Strings(s)
}

func (s *Hosts) RemoveDuplicates() {
	*s = (*s)[:removeDuplicates(s)]
}

func (s Hosts) Len() int {
	return len(s)
}

func (s Hosts) Less(i, j int) bool {
	return s[i] < s[j]
}
func (s Hosts) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func trimHost(host string) string {
	if strings.HasPrefix(host, "*.") {
		host = host[2:]
	}
	return host
}

func reverseString(str string) string {
	s := make([]rune, len(str))
	for i, r := range str {
		s[len(s)-i-1] = r
	}

	return string(s)
}
