package route_test

import (
	"net"
	"testing"

	"github.com/aspcartman/roskompozor/rkn"
	"github.com/aspcartman/roskompozor/route"
)

func TestContains(t *testing.T) {
	set, err := rkn.GetBanned()
	switch {
	case err != nil:
		t.Error(err)
	case !set.Contains(net.ParseIP("108.174.10.10")):
		t.Errorf("set doesn't contain an IP that it should contain in ipv6 form")
	case !set.Contains(route.IPv4(net.ParseIP("108.174.10.10"))):
		t.Errorf("set doesn't contain an IP that it should contain in ipv4 form")
	case set.Contains(net.ParseIP("108.174.10.11")):
		t.Errorf("set contains an IP that it shouldn't contain in ipv6 form")
	case set.Contains(route.IPv4(net.ParseIP("108.174.10.11"))):
		t.Errorf("set contains an IP that it shouldn't contain in ipv4 form")
	}
}

func BenchmarkContains(b *testing.B) {
	set, _ := rkn.GetBanned()
	b.ResetTimer()
	ip := net.ParseIP("108.174.10.10")
	for i := 0; i < b.N; i++ {
		set.Contains(ip)
	}
}

func TestAdd(t *testing.T) {
	set := route.IPSet{}
	set.AddIP(net.ParseIP("108.174.10.10"))
	set.AddIP(net.ParseIP("108.174.10.11"))
	switch {
	case set.Contains(net.ParseIP("108.174.10.9")):
		t.Errorf("set contains not added element")
	case !set.Contains(net.ParseIP("108.174.10.10")):
		fallthrough
	case !set.Contains(net.ParseIP("108.174.10.11")):
		t.Errorf("set doesn't contain added element")
	}
}
