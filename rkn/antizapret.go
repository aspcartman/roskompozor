package rkn

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/aspcartman/roskompozor/route"
)

func GetBanned() (route.IPSet, error) {
	res, err := http.Get("https://api.antizapret.info/group.php")
	if err != nil {
		return route.IPSet{}, err
	}

	defer res.Body.Close()

	var ips []net.IP
	var nets []net.IPNet

	s := bufio.NewScanner(res.Body)
	for s.Scan() {
		for _, ipstr := range strings.Split(s.Text(), ",") {
			switch {
			case len(ipstr) < 1:
				continue
			case strings.Contains(ipstr, "/"):
				_, ipnet, err := net.ParseCIDR(ipstr)
				if err != nil {
					return route.IPSet{}, err
				}
				nets = append(nets, *ipnet)
			default:
				ip := net.ParseIP(ipstr)
				if ip == nil {
					return route.IPSet{}, errors.New(fmt.Sprint("invalid ipv ", ipstr))
				}
				ips = append(ips, ip)
			}
		}
	}

	return route.NewIPSet(ips, nets), nil
}
