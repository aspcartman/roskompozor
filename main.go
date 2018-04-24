package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/aspcartman/roskompozor/rkn"
	"github.com/aspcartman/roskompozor/route"
	"github.com/aspcartman/roskompozor/sniff"
	"github.com/sirupsen/logrus"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: roskompozor <main_iface> <vpn_iface>")
		os.Exit(-1)
	}

	router{Main: os.Args[1], VPN: os.Args[2]}.routeLoop()
}

type router struct {
	Main string
	VPN  string

	bannedRefresh time.Time
	banned        route.IPSet

	routedRefresh time.Time
	routed        route.IPSet
}

func (r router) routeLoop() {
	r.refreshBanned()
	r.refreshRouted()

	logrus.WithFields(logrus.Fields{"main": r.Main, "vpn": r.VPN}).Info("listening and routing")

	if err := sniff.Sniff(r.Main, func(ip net.IP) {
		switch {
		case r.banned.Contains(ip) && !r.routed.Contains(ip):
			r.route(ip)
		case time.Since(r.routedRefresh) > 10*time.Second:
			r.refreshRouted()
		case time.Since(r.bannedRefresh) > 10*time.Minute:
			r.refreshBanned()
		}
	}); err != nil {
		logrus.WithError(err).Error("failed sniffing network traffic")
		os.Exit(-1)
	}
}

func (r *router) route(ip net.IP) {
	if err := route.Add(ip, r.VPN); err != nil {
		logrus.WithError(err).Error("failed adding route")
		return
	}
	r.routed.AddIP(ip)

	logrus.WithField("ip", ip).WithField("dest", r.VPN).Info("routed")
}

func (r *router) refreshBanned() {
	logrus.Info("refreshing banned")
retry:
	set, err := rkn.GetBanned()
	if err != nil {
		logrus.WithError(err).Error("failed fetching banned ips")
		time.Sleep(2 * time.Second)
		goto retry
	}

	r.banned, r.bannedRefresh = set, time.Now()
	logrus.WithFields(logrus.Fields{"ips": len(set.IPs), "networks": len(set.Nets)}).Info("fetched banned ips")
}

func (r *router) refreshRouted() {
	logrus.Info("refreshing routed")
retry:
	set, err := route.Routed()
	if err != nil {
		logrus.WithError(err).Error("failed fetching routed ips")
		time.Sleep(2 * time.Second)
		goto retry
	}

	r.routed, r.routedRefresh = set, time.Now()
	logrus.WithFields(logrus.Fields{"ips": len(set.IPs), "networks": len(set.Nets)}).Info("fetched routed ips")
}
