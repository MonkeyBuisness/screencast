package service

import (
	"bytes"
	//"encoding/base64"
	"github.com/kbinani/screenshot"
	"image/jpeg"
	"log"
	"time"
)

const (
	millisInSecond = 1000
)

var (
	screen   chan []byte
	listener Listener
)

type ScreenCastConfig struct {
	BitRate int
	Quality int
}

type Listener interface {
	NewScreenshot(data []byte)
}

func StartScreencast(config ScreenCastConfig) {
	screen = make(chan []byte, 5)

	go screenCapture(config)
}

func SetListener(l Listener) {
	listener = l
}

func StopScreencast() {
	close(screen)
}

func screenCapture(config ScreenCastConfig) {
	bounds := screenshot.GetDisplayBounds(0)
	rate := time.Millisecond * time.Duration(millisInSecond/config.BitRate)

	log.Printf("Service succesfully started at: %d\n", config.BitRate)

	for {
		if img, err := screenshot.CaptureRect(bounds); err != nil {
			log.Println(err)
		} else {
			// decode to base64 string
			buf := new(bytes.Buffer)
			if err := jpeg.Encode(buf, img, &jpeg.Options{Quality: config.Quality}); err != nil {
				log.Println(err)
			} else {
				if listener != nil {
					//s := base64.RawStdEncoding.EncodeToString(buf.Bytes())
					listener.NewScreenshot(buf.Bytes())
				}
			}
		}

		time.Sleep(rate)
	}
}
