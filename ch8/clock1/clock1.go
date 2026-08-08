package clock1

import (
	"io"
	"log"
	"net"
	"time"

	"golangbook.io/ch8/reverb1"
)

func Exec() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}

		// go handleConn(conn)
		go reverb1.HandleConn(conn)
	}
}

func handleConn(c net.Conn) {
	defer c.Close()

	for {
		_, err := io.WriteString(c, time.Now().Format("15:04:05\n"))
		if err != nil {
			return
		}

		time.Sleep(1 * time.Second)
	}
}
