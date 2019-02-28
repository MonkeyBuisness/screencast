package service

import (
	"bytes"
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
	jpegOptions := &jpeg.Options{Quality: config.Quality}

	for {
		if img, err := screenshot.CaptureRect(bounds); err != nil {
			log.Println(err)
		} else {
			// encode to jpeg image
			buf := new(bytes.Buffer)
			if err := jpeg.Encode(buf, img, jpegOptions); err != nil {
				log.Println(err)
			} else {
				if listener != nil {
					listener.NewScreenshot(buf.Bytes())
				}
			}
		}

		time.Sleep(rate)
	}
}
