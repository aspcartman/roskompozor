package rkn

import (
	"testing"
)

func TestGetBanned(t *testing.T) {
	set, err := GetBanned()
	switch {
	case err != nil:
		t.Error(err)
	case len(set.IPs) == 0:
		t.Errorf("ips is empty")
	case len(set.Nets) == 0:
		t.Errorf("nets is empty")
	}
}
