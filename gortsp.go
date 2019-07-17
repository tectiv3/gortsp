package gortsp

import (
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/tectiv3/edrtsp/rtsp"
)

var rtspServer *rtsp.Server

func startRTSP() (err error) {
	if rtspServer == nil {
		err = fmt.Errorf("RTSP Server Not Found")
		return
	}

	link := fmt.Sprintf("rtsp://%s:%d", localIP(), rtspServer.TCPPort)
	log.Println("rtsp server started -->", link)
	go func() {
		if err := rtspServer.Start(); err != nil {
			log.Println("start rtsp server error", err)
		}
		log.Println("rtsp server stopped")
	}()

	select {}
	return
}

func stopRTSP() (err error) {
	if rtspServer == nil {
		err = fmt.Errorf("RTSP Server Not Found")
		return
	}
	rtspServer.Stop()
	return
}

//StartServer starts webserver
func StartServer() string {
	ip, err := externalIP()
	if err != nil {
		ip = fmt.Sprint(err)
	}
	log.Println(ip)
	rtspServer = rtsp.GetServer()
	rtspServer.TCPPort = 8554
	go startRTSP()

	local := localIP()
	return fmt.Sprintf("External: %s, Local: %s", ip, local)
}

func localIP() string {
	ip := ""
	if addrs, err := net.InterfaceAddrs(); err == nil {
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && !ipnet.IP.IsMulticast() && !ipnet.IP.IsLinkLocalUnicast() && !ipnet.IP.IsLinkLocalMulticast() && ipnet.IP.To4() != nil {
				ip = ipnet.IP.String()
			}
		}
	}
	return ip
}

func isPortInUse(port int) bool {
	if conn, err := net.DialTimeout("tcp", net.JoinHostPort("", fmt.Sprintf("%d", port)), 3*time.Second); err == nil {
		conn.Close()
		return true
	}
	return false
}

func externalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("are you connected to the network?")
}
