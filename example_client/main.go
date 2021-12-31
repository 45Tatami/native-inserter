package main

import (
	"bufio"
	"encoding/binary"
	"log"
	"math"
	"net"
	"os"
)

const remote_addr = "localhost:30501"

func main() {
	c, err := net.Dial("tcp", remote_addr)
	if err != nil {
		log.Panic("Error connecting:", err)
	}
	defer c.Close()

	log.Printf("Connected to '%s'. Please input lines\n", remote_addr)

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		msg := []byte(s.Text())

		msg_len := len(msg)
		if msg_len > math.MaxUint32 {
			log.Println("Message too long")
			continue
		}

		len_buf := make([]byte, 4)
		binary.LittleEndian.PutUint32(len_buf, uint32(msg_len))

		if _, err := c.Write(len_buf); err != nil {
			log.Panic("Error sending len:", err)
		}
		if _, err := c.Write(msg); err != nil {
			log.Panic("Error sending payload:", err)
		}
	}
	if err := s.Err(); err != nil {
		log.Println("Error reading:", err)
	}
}
