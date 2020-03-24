package main

import (
	"bufio"
	"log"
	"net"
	"strings"
)

func main() {
	listen, _ := net.Listen("tcp", ":25")

	for {
		conn, _ := listen.Accept()
		log.Println("accepted")

		conn.Write([]byte("220 mock ESMTP\n"))

		reader := bufio.NewReader(conn)
		buf, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalln(err)
		}
		msg := string(buf)
		log.Printf("mock msg: %s", msg)

		conn.Write([]byte("250\n"))

		buf, err = reader.ReadString('\n')
		if err != nil {
			log.Fatalln(err)
		}
		msg = string(buf)
		log.Printf("mock msg: %s", msg)

		conn.Write([]byte("250 localhost\n"))

		buf, err = reader.ReadString('\n')
		if err != nil {
			log.Fatalln(err)
		}
		msg = string(buf)
		log.Printf("mock msg: %s", msg)
		email := strings.Split(strings.Split(msg, "<")[1], ">")[0]
		log.Printf("mock email: %s", email)

		conn.Write([]byte("250 " + email + " sender accepted\n"))

		buf, err = reader.ReadString('\n')
		if err != nil {
			log.Fatalln(err)
		}
		msg = string(buf)
		log.Printf("mock msg: %s", msg)
		email = strings.Split(strings.Split(msg, "<")[1], ">")[0]
		log.Printf("mock email: %s", email)

		conn.Write([]byte("250 " + email + " ok\n"))

		buf, err = reader.ReadString('\n')
		if err != nil {
			log.Fatalln(err)
		}
		msg = string(buf)
		log.Printf("mock msg: %s", msg)
		conn.Write([]byte("354 Enter mail, end with \".\" on a line by itself\n"))

		for {
			buf, err = reader.ReadString('\n')
			if err != nil {
				log.Fatalln(err)
			}
			msg = strings.TrimRight(string(buf), "\r\n")
			log.Printf("mock: len: %d,  msg: '%s'", len(msg), msg)
			if msg == "." {
				log.Println("mock: exit circle")
				break
			}
		}

		conn.Write([]byte("250 769947 message accepted for delivery\n"))

		buf, err = reader.ReadString('\n')
		if err != nil {
			log.Fatalln(err)
		}
		msg = string(buf)
		log.Printf("mock msg: %s", msg)

		conn.Write([]byte("221 mock CommuniGate Pro SMTP closing connection\n"))

		conn.Close()
	}
}
