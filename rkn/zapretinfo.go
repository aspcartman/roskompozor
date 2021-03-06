package rkn

import (
	"bufio"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/aspcartman/roskompozor/sets"
)

func ZapretInfo() (sets.IPs, error) {
	res, err := http.Get("https://raw.githubusercontent.com/zapret-info/z-i/master/dump.csv")
	if err != nil {
		return sets.IPs{}, err
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
			return sets.IPs{}, err
		}
		str = strings.TrimSpace(str)

		for _, ipstr := range strings.Split(strings.Split(str, ";")[0], " | ") {
			if err := addip(ipstr, &ips, &nets); err != nil {
				return sets.IPs{}, err
			}
		}
	}

	return sets.NewIPSet(ips, nets), nil
}
