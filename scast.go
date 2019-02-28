package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
)

const (
	osCastAddress = "SCREENCAST_ADDRESS"
	osCastPort    = "SCREENCAST_PORT"
)

func main() {
	// define flags
	port := flag.Int("port", 8080, "defines cast port")
	address := flag.String("addr", "localhost", "defines cast net address")
	quality := flag.Int("q", 40, "defines image quality")
	bitrate := flag.Int("br", 20, "defines cast bitrate")

	// parse flags
	flag.Parse()

	log.Println("Starting screencast...")
	log.Printf("Host:     %s:%d\n", *address, *port)
	log.Println("Quality: ", *quality)
	log.Println("Bitrate: ", *bitrate)

	// check file to run script
	var file string
	switch runtime.GOOS {
	case "windows":
		file = "./run.bat"
	case "linux":
		fallthrough
	default:
		file = "./run.sh"
	}

	// get absolute path to run script
	absPath, err := filepath.Abs(file)
	if err != nil {
		log.Fatal(err)
	}

	// set system variables
	os.Setenv(osCastAddress, *address)
	os.Setenv(osCastPort, strconv.Itoa(*port))

	// create command to start script
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command(absPath, strconv.Itoa(*quality), strconv.Itoa(*bitrate))
	case "linux":
		fallthrough
	default:
		cmd = exec.Command("/bin/bash", absPath, strconv.Itoa(*quality), strconv.Itoa(*bitrate))
	}

	// create stdout reader pipe
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	// create stderr reader pipe
	cmdErrReader, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}

	// start pipes listener
	go listen(cmdReader)
	go listen(cmdErrReader)

	// start command pipe
	if err = cmd.Start(); err != nil {
		log.Fatal(err)
	}

	// wait command to finished
	if err = cmd.Wait(); err != nil {
		log.Fatal(err)
	}
}

func listen(reader io.ReadCloser) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		log.Println(scanner.Text())
	}
}
