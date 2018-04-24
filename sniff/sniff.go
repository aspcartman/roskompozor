package sniff

import (
	"errors"
	"fmt"
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)


func Sniff(name string, clb func(ip net.IP)) error {
	iface, err := interfaceNamed(name)
	if err != nil {
		return err
	}

	handle, err := pcap.OpenLive(iface.Name, 1024, false, pcap.BlockForever)
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
		if p, ok := packet.Layer(layers.LayerTypeIPv4).(*layers.IPv4); ok {
			clb(p.DstIP)
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
