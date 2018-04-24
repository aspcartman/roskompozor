package rkn

import (
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/aspcartman/roskompozor/route"
)

type Source func() (route.IPSet, error)

func addip(str string, ips *[]net.IP, nets *[]net.IPNet) error {
	switch {
	case len(str) == 0:
		// sometimes it happens, just skip
	case strings.Contains(str, "/"):
		_, ipnet, err := net.ParseCIDR(str)
		if err != nil {
			return err
		}
		*nets = append(*nets, *ipnet)
	default:
		ip := net.ParseIP(str)
		if ip == nil {
			return errors.New(fmt.Sprint("invalid ipv ", str))
		}
		*ips = append(*ips, ip)
	}
	return nil
}
