package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
)

type Message struct {
	Body string `json:"body"`
}

const MSG_SIZE_LIMIT = 1000 * 1000 // Unclear if 1MiB or 1MB, using lower bound
const LISTEN_INTERFACE = ":30501"

func main() {
	log.Println("Listening on", LISTEN_INTERFACE)

	l, err := net.Listen("tcp", LISTEN_INTERFACE)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	msg_pipe := make(chan string)

	// Output routine
	go func() {
		for {
			msg, err := json.Marshal(Message{<-msg_pipe})
			if err != nil {
				log.Println("Error encoding message:", err)
				continue
			}

			msg_len := len(msg)
			if msg_len >= MSG_SIZE_LIMIT {
				log.Println("Message over webextension limit. Discarding")
				continue
			}

			len_buf := make([]byte, 4)
			binary.LittleEndian.PutUint32(len_buf, uint32(msg_len))

			log.Println(len_buf, string(msg))
			fmt.Print(string(len_buf))
			fmt.Print(string(msg))
		}
	}()

	// Accept connections
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("Error accepting conn:", err)
		}
		go conn_hndlr(conn, msg_pipe)
	}
}

func conn_hndlr(c net.Conn, ch chan<- string) {
	defer c.Close()
	log.Println("Connection opened:", c)
	for {
		msg, err := read_msg(c)
		if err != nil {
			break
		}
		ch <- msg
	}
	log.Println("Closing connection:", c)
}

func read_msg(c net.Conn) (string, error) {
	len_buf := make([]byte, 4)

	if _, err := io.ReadFull(c, len_buf); err != nil {
		log.Println("Len read error:", err)
		return "", err
	}

	msg_len := binary.LittleEndian.Uint32(len_buf)
	msg_buf := make([]byte, msg_len)
	if _, err := io.ReadFull(c, msg_buf); err != nil {
		log.Println("Msg read error:", err)
		return "", err
	}

	return string(msg_buf), nil
}
