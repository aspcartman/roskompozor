package sets

import "sort"

func removeDuplicates(s sort.Interface) int {
	p, l := 0, s.Len()
	if l <= 1 {
		return l
	}

	for i := 1; i < l; i++ {
		if !s.Less(p, i) {
			continue
		}
		p++
		if p < i {
			s.Swap(p, i)
		}
	}
	return p + 1
}
