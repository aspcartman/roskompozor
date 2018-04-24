package rkn

import (
	"bufio"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/aspcartman/roskompozor/route"
)

func ZapretInfo() (route.IPSet, error) {
	res, err := http.Get("https://raw.githubusercontent.com/zapret-info/z-i/master/dump.csv")
	if err != nil {
		return route.IPSet{}, err
	}

	defer res.Body.Close()

	var ips []net.IP
	var nets []net.IPNet

	r := bufio.NewReader(res.Body)
	r.ReadString('\n') // skip first line
	for {
		str, err := r.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return route.IPSet{}, err
		}
		str = strings.TrimSpace(str)

		for _, ipstr := range strings.Split(strings.Split(str, ";")[0], " | ") {
			if err := addip(ipstr, &ips, &nets); err != nil {
				return route.IPSet{}, err
			}
		}
	}

	return route.NewIPSet(ips, nets), nil
}
