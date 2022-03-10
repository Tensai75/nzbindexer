package main

import (
	"fmt"
	"strconv"

	"github.com/Tensai75/nntp"
)

func ConnectNNTP() (*nntp.Conn, error) {

	conn, err := nntp.Dial("tcp", conf.Server.Host+":"+strconv.Itoa(conf.Server.Port))
	if err != nil {
		fmt.Printf("Connection to usenet server failed: %v\n", err)
		return conn, err
	}
	if err := conn.Authenticate(conf.Server.User, conf.Server.Password); err != nil {
		fmt.Printf("Authentication with usenet server failed: %v\n", err)
		return conn, err
	}
	fmt.Println("Connection to usenet server established")

	return conn, nil

}
