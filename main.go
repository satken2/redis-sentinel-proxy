package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"
)

var (
	masterAddr *net.TCPAddr
	raddr      *net.TCPAddr
	saddr      *net.TCPAddr

	localAddr    = flag.String("listen", ":9999", "local address")
	sentinelAddr = flag.String("sentinel", ":26379", "remote address")
	masterName   = flag.String("master", "", "name of the master redis node")
)

func main() {
	flag.Parse()

	laddr, err := net.ResolveTCPAddr("tcp", *localAddr)
	if err != nil {
		log.Fatal("ERROR: Failed to resolve local address: %s", err)
	}
	saddr, err = net.ResolveTCPAddr("tcp", *sentinelAddr)
	if err != nil {
		log.Fatal("ERROR: Failed to resolve sentinel address: %s", err)
	}

	go master()

	listener, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Println(err)
			continue
		}

		go proxy(conn, masterAddr)
	}
}

func master() {
	var err error
	for {
		masterAddr, err = getMasterAddr(saddr, *masterName)
		if err != nil {
			log.Println(err)
		}
		time.Sleep(5 * time.Second)
	}
}

func pipe(r io.Reader, w io.WriteCloser) {
	io.Copy(w, r)
	w.Close()
}

func proxy(local io.ReadWriteCloser, remoteAddr *net.TCPAddr) {
  dialer := net.Dialer{Timeout: 10}
  remote, err := dialer.Dial("tcp", remoteAddr)
	if err != nil {
		log.Println(err)
		local.Close()
    log.Println("WARNING: Master connection failed. Closing connection...")
		return
	}
	go pipe(local, remote)
	go pipe(remote, local)
}

func getMasterAddr(sentinelAddress *net.TCPAddr, masterName string) (*net.TCPAddr, error) {
	conn, err := net.DialTCP("tcp", nil, sentinelAddress)
	if err != nil {
    log.Println("WARNING: Failed to connect to sentinel. Updating sentinel tcp address.")
    updateSentinelAddr()
		return nil, err
	}

	defer conn.Close()

	conn.Write([]byte(fmt.Sprintf("sentinel get-master-addr-by-name %s\n", masterName)))

	b := make([]byte, 256)
	_, err = conn.Read(b)
	if err != nil {
		log.Fatal(err)
	}

	parts := strings.Split(string(b), "\r\n")

	if len(parts) < 5 {
		err = errors.New("ERROR: Couldn't get master address from sentinel")
		return nil, err
	}

	//getting the string address for the master node
	stringaddr := fmt.Sprintf("%s:%s", parts[2], parts[4])
	addr, err := net.ResolveTCPAddr("tcp", stringaddr)

	if err != nil {
		return nil, err
	}

	//check that there's actually someone listening on that address
	conn2, err := net.DialTCP("tcp", nil, addr)
	if err == nil {
		defer conn2.Close()
	}

	return addr, err
}

func updateSentinelAddr()  {
  log.Println("INFO: Resolving sentinel address: %s", sentinelAddr)
  for {
    addr, err := net.ResolveTCPAddr("tcp", *sentinelAddr)
    if err != nil {
      log.Println("WARNING: Failed to resolve sentinel address. Retrying in 10 seconds...")
      time.Sleep(10 * time.Second)
    } else {
      saddr = addr
      log.Println("INFO: Successfully updated sentinel address.")
      break
    }
  }
}