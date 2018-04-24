package rkn

import (
	"bufio"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/aspcartman/roskompozor/sets"
)

func AntiZapret() (sets.IPs, error) {
	res, err := http.Get("https://api.antizapret.info/group.php")
	if err != nil {
		return sets.IPs{}, err
	}

	defer res.Body.Close()

	var ips []net.IP
	var nets []net.IPNet

	r := bufio.NewReader(res.Body)
	for {
		str, err := r.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return sets.IPs{}, err
		}
		str = strings.TrimSpace(str)

		for _, ipstr := range strings.Split(str, ",") {
			if err := addip(ipstr, &ips, &nets); err != nil {
				return sets.IPs{}, err
			}
		}
	}

	return sets.NewIPSet(ips, nets), nil
}
