package sniff

import (
	"errors"
	"fmt"
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func Sniff(name string, ipclb func(ip net.IP), dnsclb func(host string)) error {
	iface, err := interfaceNamed(name)
	if err != nil {
		return err
	}

	handle, err := pcap.OpenLive(iface.Name, 256, false, pcap.BlockForever)
	if err != nil {
		return err
	}
	defer handle.Close()

	var addr net.IP
	for _, a := range iface.Addresses {
		if len(a.IP) == net.IPv4len {
			addr = a.IP
		}
	}
	if addr == nil {
		return errors.New("no ipv4 addr")
	}

	flt := fmt.Sprintf("src host %v", addr)
	if err := handle.SetBPFFilter(flt); err != nil {
		return err
	}

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		if p, ok := packet.Layer(layers.LayerTypeDNS).(*layers.DNS); ok {
			for _, q := range p.Questions {
				dnsclb(string(q.Name))
			}
		} else if p, ok := packet.Layer(layers.LayerTypeIPv4).(*layers.IPv4); ok {
			ipclb(p.DstIP)
		}
	}

	return nil
}

func interfaceNamed(name string) (pcap.Interface, error) {
	devs, err := pcap.FindAllDevs()
	if err != nil {
		return pcap.Interface{}, err
	}

	for _, d := range devs {
		if d.Name == name {
			return d, nil
		}
	}

	return pcap.Interface{}, fmt.Errorf("interface not found: %s", name)
}
