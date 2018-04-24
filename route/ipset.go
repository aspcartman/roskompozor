package route

import (
	"bytes"
	"net"
	"sort"
)

type IPSet struct {
	IPs  []net.IP
	Nets []net.IPNet
}

func NewIPSet(ips []net.IP, nets []net.IPNet) IPSet {
	s := IPSet{ips, nets}
	s.V4()
	s.Sort()
	s.RemoveDuplicates()
	return s
}

func (b IPSet) V4() {
	for i := range b.IPs {
		b.IPs[i] = IPv4(b.IPs[i])
	}
}

func (b IPSet) Contains(ip net.IP) bool {
	ip = IPv4(ip)
	pos := sort.Search(len(b.IPs), func(i int) bool {
		return bytes.Compare(b.IPs[i], ip) != -1
	})

	switch {
	case pos == -1 || pos == len(b.IPs):
		//return false
	case b.IPs[pos].Equal(ip):
		return true
	}

	for _, n := range b.Nets {
		if n.Contains(ip) {
			return true
		}
	}

	return false
}

func (b *IPSet) AddIP(ip net.IP) {
	if b.Contains(ip) {
		return
	}
	b.IPs = append(b.IPs, IPv4(ip))
	b.Sort()
}

func (b *IPSet) Sort() {
	sort.Sort(ipsetsortips{b})
}

func (b *IPSet) RemoveDuplicates() {
	w := ipsetsortips{b}
	p, l := 0, w.Len()
	if l <= 1 {
		return
	}

	for i := 1; i < l; i++ {
		if !w.Less(p, i) {
			continue
		}
		p++
		if p < i {
			w.Swap(p, i)
		}
	}
	w.IPs = w.IPs[:p+1]
}

type ipsetsortips struct {
	*IPSet
}

func (b ipsetsortips) Len() int {
	return len(b.IPs)
}

func (b ipsetsortips) Less(i, j int) bool {
	return bytes.Compare(b.IPs[i], b.IPs[j]) == -1
}

func (b ipsetsortips) Swap(i, j int) {
	b.IPs[i], b.IPs[j] = b.IPs[j], b.IPs[i]
}

func IPv4(ip net.IP) net.IP {
	if len(ip) == net.IPv4len {
		return ip
	}
	return net.IP{ip[12], ip[13], ip[14], ip[15]}
}
