package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func echoServer(c net.Conn) {
	for {
		buf := make([]byte, 10240)
		nr, err := c.Read(buf)
		if err != nil {
			return
		}

		data := buf[0:nr]
		println("Server got:", string(data))
		_, err = c.Write([]byte("SUCCESS\n"))
		if err != nil {
			log.Fatal("Writing client error: ", err)
		}
	}
}

func main() {

	socket_path := "/tmp/go.sock"
	log.Println("Starting unix server at", socket_path)
	ln, err := net.Listen("unix", socket_path)
	if err != nil {
		log.Fatal("Listen error: ", err)
	}

	if err := os.Chmod(socket_path, 755); err != nil {
        log.Fatal("Unable to set perms to 755", err)
    }

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, syscall.SIGTERM)
	go func(ln net.Listener, c chan os.Signal) {
		sig := <-c
		log.Printf("Caught signal %s: shutting down.", sig)
		ln.Close()
		os.Exit(0)
	}(ln, sigc)

	for {
		fd, err := ln.Accept()
		if err != nil {
			log.Fatal("Accept error: ", err)
		}

		go echoServer(fd)
	}
}