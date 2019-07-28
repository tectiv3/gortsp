package gortsp

import (
	"fmt"
	"log"

	"github.com/tectiv3/edrtsp/rtsp"
	"github.com/tectiv3/edrtsp/stats"
	"github.com/tectiv3/edrtsp/utils"
)

var rtspServer *rtsp.Server

func startRTSP() (err error) {
	if rtspServer == nil {
		err = fmt.Errorf("RTSP Server Not Found")
		return
	}

	link := fmt.Sprintf("rtsp://%s:%d", utils.LocalIP(), rtspServer.TCPPort)
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
	ip, err := utils.ExternalIP()
	if err != nil {
		ip = fmt.Sprint(err)
	}
	log.Println(ip)
	rtspServer = rtsp.GetServer()
	rtspServer.TCPPort = 8554
	go startRTSP()

	local := utils.LocalIP()
	return fmt.Sprintf("External: %s, Local: %s", ip, local)
}

//GetStats returns json encoded server stats
func GetStats() string {
	return string(stats.GetStats())
}

//GetPushers returns json encoded pushers stats
func GetPushers() string {
	return string(stats.GetPushersJSON())
}

//GetPlayers returns json encoded players stats
func GetPlayers() string {
	return string(stats.GetPlayersJSON())
}
