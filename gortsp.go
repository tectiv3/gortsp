package gortsp

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"html"
	"image"
	"image/jpeg"
	"log"
	"net"
	"net/http"

	"github.com/saljam/mjpeg"
	// "github.com/tectiv3/dorsvr/rtspserver"
)

var encoded string
var rgba image.Image
var stream *mjpeg.Stream

// var enc *x264.Encoder

func startWebServer() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	http.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hi")
	})
	http.HandleFunc("/image", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/jpeg")
		// w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
		jpeg.Encode(w, rgba, nil) // Write to the ResponseWriter
	})

	stream = mjpeg.NewStream()
	http.Handle("/camera", stream)

	log.Fatal(http.ListenAndServe(":8081", nil))

	// server := rtspserver.New(nil)

	// portNum := 8554
	// err := server.Listen(portNum)
	// if err != nil {
	// 	log.Printf("Failed to bind port: %d\n", portNum)
	// 	return
	// }

	// if !server.SetupTunnelingOverHTTP(80) ||
	// 	!server.SetupTunnelingOverHTTP(8000) ||
	// 	!server.SetupTunnelingOverHTTP(8080) {
	// 	log.Printf("We use port %d for optional RTSP-over-HTTP tunneling, "+
	// 		"or for HTTP live streaming (for indexed Transport Stream files only).\n", server.HTTPServerPortNum())
	// } else {
	// 	log.Println("(RTSP-over-HTTP tunneling is not available.)")
	// }

	// urlPrefix := server.RtspURLPrefix()
	// log.Println("This server's URL: " + urlPrefix + "<filename>.")

	// server.Start()

	// select {}
}

// func initX264() {
// 	opts := &x264.Options{
// 		Width:     320,
// 		Height:    240,
// 		FrameRate: 10,
// 		Tune:      "zerolatency",
// 		Preset:    "veryfast",
// 		Profile:   "baseline",
// 		LogLevel:  x264.LogDebug,
// 	}
// 	var err error

// 	buf := bytes.NewBuffer(make([]byte, 0))
// 	enc, err = x264.NewEncoder(buf, opts)
// 	if err != nil {
// 		log.Printf("%s\n", err.Error())
// 		return
// 	}

// 	defer enc.Close()
// }

//StartServer starts webserver
func StartServer(name string) string {
	ip, err := externalIP()
	if err != nil {
		ip = fmt.Sprint(err)
	}
	log.Println(ip)
	go startWebServer()

	return fmt.Sprintf("IP: %s for %s.", ip, name)
}

//DumpByteArray just dumps array length and encodes it into string
func DumpByteArray(img []byte) {
	log.Printf("Len: %d\n", len(img))
	encoded = base64.StdEncoding.EncodeToString(img)
}

//PushImage pushes YUV image down to encoder
func PushImage(y, u, v []byte, width, height int) {
	log.Printf("Len: %d,%d,%d. %dx%d", len(y), len(u), len(v), width, height)
	// encoded = base64.StdEncoding.EncodeToString(img)
	rect := image.Rectangle{image.Point{0, 0}, image.Point{width, height}}
	res := image.NewYCbCr(rect, image.YCbCrSubsampleRatio420)
	res.Y = y
	res.Cb = u
	res.Cr = v
	// b := res.Bounds()
	// m := image.New (image.Rect(0, 0, b.Dx(), b.Dy()))
	// draw.Draw(m, m.Bounds(), res, b.Min, draw.Src)
	rgba = res

	buf := &bytes.Buffer{}
	jpeg.Encode(buf, res, nil)

	stream.UpdateJPEG(buf.Bytes())
	// if err := enc.Encode(res); err != nil {
	// 	log.Printf("%s\n", err.Error())
	// }
}

func toH264(img []byte, width, height int) {
	// w := 400
	// h := 400
	// var nal [][]byte

	// c, _ := codec.NewH264Encoder(w, h, image.YCbCrSubsampleRatio420)
	// nal = append(nal, c.Header)

	// for i := 0; i < 60; i++ {
	// 	img := image.NewYCbCr(image.Rect(0, 0, w, h), image.YCbCrSubsampleRatio420)
	// 	p, _ := c.Encode(img)
	// 	if len(p.Data) > 0 {
	// 		nal = append(nal, p.Data)
	// 	}
	// }
	// for {
	// 	// flush encoder
	// 	p, err := c.Encode(nil)
	// 	if err != nil {
	// 		break
	// 	}
	// 	nal = append(nal, p.Data)
	// }
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
