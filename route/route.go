package route

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"os/exec"
	"strings"
)

func Add(dst net.IP, iface string) error {
	out, err := exec.Command("route", "add", "-host", fmt.Sprint(dst), "-interface", iface).CombinedOutput()
	if err != nil && len(out) > 0 {
		err = fmt.Errorf("%s: %s", err.Error(), string(out))
	}
	return err
}

func Routed() (IPSet, error) {
	out, err := exec.Command("netstat", "-rnf", "inet").CombinedOutput()
	if err != nil && len(out) > 0 {
		err = fmt.Errorf("%s: %s", err.Error(), string(out))
	}

	var ips []net.IP
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), " ")[0]
		if ip := net.ParseIP(s); ip != nil {
			ips = append(ips, ip)
		}
	}

	return NewIPSet(ips, nil), err
}
