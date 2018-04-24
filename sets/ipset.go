package sets

import (
	"bytes"
	"net"
	"sort"
)

type IPs struct {
	IPs  []net.IP
	Nets []net.IPNet
}

func NewIPSet(ips []net.IP, nets []net.IPNet) IPs {
	s := IPs{ips, nets}
	s.V4()
	s.Sort()
	s.RemoveDuplicates()
	return s
}

func (s *IPs) AddIP(ip net.IP) {
	if s.Contains(ip) {
		return
	}
	s.IPs = append(s.IPs, IPv4Form(ip))
	s.Sort()
}

func (s IPs) V4() {
	for i := range s.IPs {
		s.IPs[i] = IPv4Form(s.IPs[i])
	}
}

func (s IPs) Contains(ip net.IP) bool {
	ip = IPv4Form(ip)
	pos := sort.Search(len(s.IPs), func(i int) bool {
		return bytes.Compare(s.IPs[i], ip) != -1
	})

	switch {
	case pos == -1 || pos == len(s.IPs):
		//return false
	case s.IPs[pos].Equal(ip):
		return true
	}

	for _, n := range s.Nets {
		if n.Contains(ip) {
			return true
		}
	}

	return false
}

func (s IPs) Sort() {
	sort.Sort(s)
}

func (s *IPs) RemoveDuplicates() {
	s.IPs = s.IPs[:removeDuplicates(s)]
}

func (s IPs) Len() int {
	return len(s.IPs)
}

func (s IPs) Less(i, j int) bool {
	return bytes.Compare(s.IPs[i], s.IPs[j]) == -1
}

func (s IPs) Swap(i, j int) {
	s.IPs[i], s.IPs[j] = s.IPs[j], s.IPs[i]
}

func IPv4Form(ip net.IP) net.IP {
	if len(ip) == net.IPv4len {
		return ip
	}
	return net.IP{ip[12], ip[13], ip[14], ip[15]}
}
