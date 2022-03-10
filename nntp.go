package main

import (
	"fmt"
	"strconv"

	"github.com/Tensai75/nntp"
)

var connectionGuard chan struct{}

func ConnectNNTP() (*nntp.Conn, error) {

	if connectionGuard == nil {
		connectionGuard = make(chan struct{}, conf.Server.Connections)
	}
	connectionGuard <- struct{}{} // will block if guard channel is already filled
	conn, err := nntp.Dial("tcp", conf.Server.Host+":"+strconv.Itoa(conf.Server.Port))
	if err != nil {
		fmt.Printf("Connection to usenet server failed: %v\n", err)
		return conn, err
	}
	if err := conn.Authenticate(conf.Server.User, conf.Server.Password); err != nil {
		fmt.Printf("Authentication with usenet server failed: %v\n", err)
		return conn, err
	}

	return conn, nil

}

func DisconnectNNTP(conn *nntp.Conn) {
	if conn != nil {
		conn.Quit()
		select {
		case <-connectionGuard:
			// go on
		default:
			// go on
		}
	}
	conn = nil
}
