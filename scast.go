package main

import (
	"bufio"
	"flag"
	//"fmt"
	//"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	//"strings"
)

func main() {
	port := flag.Int("port", 8080, "defines cast port")
	address := flag.String("addr", "localhost", "defines cast net address")
	quality := flag.Int("q", 40, "defines image quality")
	bitrate := flag.Int("br", 20, "defines cast bitrate")

	flag.Parse()

	log.Println("Starting screencast...")
	log.Printf("Host:     %s:%d\n", *address, *port)
	log.Println("Quality: ", *quality)
	log.Println("Bitrate: ", *bitrate)

	var file string
	switch runtime.GOOS {
	case "windows":
		file = "./run.bat"
	case "linux":
		fallthrough
	default:
		file = "./run.sh"
	}

	absPath, err := filepath.Abs(file)
	if err != nil {
		log.Fatal(err)
	}

	///
	os.Setenv("hello", "world!!!!")
	///

	cmd := exec.Command("/bin/bash", absPath, "77")
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			log.Println(scanner.Text())
		}
	}()
	if err = cmd.Start(); err != nil {
		log.Fatal(err)
	}
	if err = cmd.Wait(); err != nil {
		log.Fatal(err)
	}

	/*var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("/bin/bash", absPath)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	if stdout.Len() > 0 {
		log.Println(string(stdout.Bytes()))
	}

	if stderr.Len() > 0 {
		log.Fatal(string(stderr.Bytes()))
	}*/
}
