package main

import (
	"github.com/denisbrodbeck/machineid"
	"log"
	"net"
	"os/user"
	"runtime"
	"strings"
)

type FingerPrint struct {
	Uid       string `json:"Uid"`
	Gid       string `json:"Gid"`
	Username  string `json:"Username"`
	Name      string `json:"Name"`
	HomeDir   string `json:"HomeDir"`
	OS        string `json:"OS"`
	MachineId string `json:"MachineId"`
	Mac       string `json:"Mac"`
	Ip        string `json:"Ip"`
	Interface string `json:"Interface"`
}

func GetFingerPrint() FingerPrint {
	fingerprint := FingerPrint{}
	u, _ := user.Current()
	fingerprint.Uid = u.Uid
	fingerprint.Gid = u.Gid
	fingerprint.Username = u.Username
	fingerprint.Name = u.Name
	fingerprint.HomeDir = u.HomeDir
	fingerprint.OS = runtime.GOOS
	fingerprint.Ip = fingerprint.getOutboundIP().String()
	fingerprint.Mac, fingerprint.Interface = fingerprint.getMacAddrByIp(fingerprint.Ip)
	fingerprint.MachineId, _ = machineid.ID()

	return fingerprint
}

func (f *FingerPrint) getOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func (f *FingerPrint) getMacAddrByIp(ip string) (string, string) {
	ifas, err := net.Interfaces()
	if err != nil {
		log.Fatal(err)
		return "nil", ""
	}
	resevedMac := ""
	for _, ifa := range ifas {
		a := ifa.HardwareAddr.String()
		if a != "" {
			addrs, _ := ifa.Addrs()
			for _, addr := range addrs {
				if strings.Contains(addr.String(), "/24") {
					resevedMac = a
				}
				i := strings.Split(addr.String(), "/")
				if i[0] == ip {
					return a, ifa.Name
				}
			}
		}
	}

	//log.Fatal("no found", resevedMac)
	return resevedMac, "default"
}
