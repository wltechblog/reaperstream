package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"time"
	"unicode"
)

// For each host, we attempt to connect on ports 5000-5010
func startStream(host string, cchan chan Car) {
	if c.debug {
		log.Printf("Starting streams on %s", host)
	}
	for x := 5000; x < 5010; x++ {
		go runStream(host+":"+strconv.Itoa(x), cchan)
	}
}

// This is where we connect to a stream and do all the parsing...
func runStream(host string, cchan chan Car) {
	if c.debug {
		log.Printf("Running stream on %s", host)
	}
	sock, err := net.Dial("tcp", host)
	if err != nil {
		if c.debug {
			log.Println("No con on port ", host)
		}
		return
	}
	os.MkdirAll("data", 0755)
	uuid := ""
	for {
		block, err := ReadBytes(sock, 4)
		if err != nil {
			log.Printf("Reading from sock for %s on %s", err, host)
			return
		}
		if block[0] == 0x08 && block[1] == 0x04 { // keep-alive
			if c.debug {
				log.Printf("Keepalive from %s", host)
			}
			continue
		}
		if block[0] == 0xbb && block[1] == 0x0b { // frame start
			if c.debug {
				log.Printf("Frame start from %s", host)
			}
			var c Car
			datalength, err := ReadBytes(sock, 4)
			if err != nil {
				log.Printf("datalength 3 read: ", err)
				continue
			}
			// parse 4 byte block into 32 bit int
			n := int(int(datalength[3])<<24 | int(datalength[2])<<16 | int(datalength[1])<<8 | int(datalength[0]))
			frame, err := ReadBytes(sock, n)
			if err != nil {
				log.Printf("frame read: ", err)
				continue
			}

			framebuffer := bufio.NewReader(bytes.NewReader(frame))

			command, err := ReadBytes(framebuffer, 4)
			if err != nil {
				log.Printf("command read: ", err)
				continue
			}

			if command[0] == 0x3d && command[1] == 00 { // license plate start
				// 96 bytes for license plate data
				platebuf, err := ReadBytes(framebuffer, 96)
				if err != nil {
					log.Printf("platebuf read: ", err)
					continue
				}

				plate := ""
				for ptr := 0; unicode.IsLetter(rune(platebuf[ptr])) || unicode.IsNumber(rune(platebuf[ptr])); ptr++ {
					plate = plate + string(platebuf[ptr])
				}
				c.LicensePlate = plate
				u, _ := ReadBytes(framebuffer, 96)
				uuid = ""
				for ptr := 0; unicode.IsLetter(rune(u[ptr])) || unicode.IsNumber(rune(u[ptr])) || u[ptr] == '-'; ptr++ {
					uuid = uuid + string(u[ptr])
				}
				c.UUID = uuid
				// churn 2324 more bytes with unknown purpose
				_, err = ReadBytes(framebuffer, 2324)
				if err != nil {
					log.Printf("churn read: ", err)
					continue
				}

			}
			command, err = ReadBytes(framebuffer, 4)
			if err != nil {
				log.Printf("command read: ", err)
				continue
			}

			if command[0] == 0x02 && command[1] == 0x00 { // jpg image
				datalength, err = ReadBytes(framebuffer, 4)
				if err != nil {
					log.Printf("datalength read: ", err)
					continue
				}

				n := int(int(datalength[3])<<24 | int(datalength[2])<<16 | int(datalength[1])<<8 | int(datalength[0]))
				jpg, err := ReadBytes(framebuffer, n)
				if err != nil {
					log.Printf("jpg read: ", err)
					continue
				}
				// make data directory with year/month/dat
				savepath := "data/" + time.Now().Format("2006/01/02")
				os.MkdirAll(savepath, 0755)
				fn := savepath + "/frame-" + uuid + ".jpg"
				c.Filename = fn
				os.WriteFile(fn, jpg, 0777)
				datalength, err = ReadBytes(framebuffer, 4) // Unknown payload with dynamic length
				if err != nil {
					log.Printf("datalength 2 read: ", err)
					continue
				}

				n = int(int(datalength[3])<<24 | int(datalength[2])<<16 | int(datalength[1])<<8 | int(datalength[0]))
				_, err = ReadBytes(framebuffer, n)
				if err != nil {
					log.Fatal("unknown variable data read: ", err)
				}
				_, err = ReadBytes(framebuffer, 12) // Unknown payload with static length
				if err != nil {
					log.Fatal("unknown 12 byte read: ", err)
				}

				datalength, err = ReadBytes(framebuffer, 4) // json
				if err != nil {

					log.Printf("datalength 1 read: ", err)
					continue
				}

				n = int(int(datalength[3])<<24 | int(datalength[2])<<16 | int(datalength[1])<<8 | int(datalength[0]))
				j, err := ReadBytes(framebuffer, n)
				if err != nil {
					log.Fatal("json read: ", err)
				}
				err = json.Unmarshal(j, &c)
				if err != nil {
					log.Fatal(err) // this should never happen, if you hit it let me know!
				}
				c.Server = host
				cchan <- c
				// Save metadata including the stuff we added
				jm, err := json.Marshal(c)
				if err != nil {
					log.Fatal(err)
				}
				fn = savepath + "/meta-" + uuid + ".json"
				//
				os.WriteFile(fn, jm, 0777)
			}
		}
	}
}

// read given number of bytes, return an error if short read.
func ReadBytes(s io.Reader, n int) ([]byte, error) {
	var b = make([]byte, n)
	r, err := io.ReadFull(s, b[:n])
	if r != n {
		err = errors.New("short read")
	}
	return b, err
}
