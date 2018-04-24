package rkn

import (
	"testing"
)

func testSource(t *testing.T, source Source) {
	set, err := source()
	switch {
	case err != nil:
		t.Error(err)
	case len(set.IPs) == 0:
		t.Errorf("ips is empty")
	case len(set.Nets) == 0:
		t.Errorf("nets is empty")
	}
}
