package gortsp

import (
	"encoding/base64"
	"errors"
	"fmt"
	"html"
	"image"
	"image/jpeg"
	"log"
	"net"
	"net/http"
	// codec "github.com/sikang99/codec"
)

var encoded string
var rgba image.Image

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

	log.Fatal(http.ListenAndServe(":8081", nil))
}

func StartServer(name string) string {
	ip, err := externalIP()
	if err != nil {
		ip = fmt.Sprint(err)
	}
	log.Println(ip)
	go startWebServer()
	return fmt.Sprintf("IP: %s for %s.", ip, name)
}

func DumpByteArray(img []byte) {
	log.Printf("Len: %d\n", len(img))
	encoded = base64.StdEncoding.EncodeToString(img)
}

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
