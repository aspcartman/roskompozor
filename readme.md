# roskompozor

The utility fetches banned IPs and routes them through VPN, allowing user to use VPN exclusively for banned services.


### Requirements
OSX only, Linux & Windows wip.
Requires /dev/bpf* devices presented. Instructions on those will be added soon. 
As a quickfix, install Wireshark.

### Usage

Start your VPN with "Send all traffic over VPN connection" unchecked and

```bash
$ sudo go run main.go en0 ppp0
``` 

Replace en0 and ppp0 with actual interfaces. 


