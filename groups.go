package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"
)

var groups []string
var serverGroups []string

func scanGroups() error {

	if conf.Groups != "ALL" && conf.Groups != "BINARIES" {
		err := readGroups(conf.Groups)
		if err != nil {
			fmt.Printf("Error while reading group file '%s': %v\n", conf.Groups, err)
			return err
		}
	}
	fmt.Println("Connecting to usenet server")
	conn, err := ConnectNNTP()
	if err != nil {
		fmt.Printf("Error while connecting to usenet server: %v\n", err)
		return err
	}
	fmt.Println("Requesting the list of groups")
	filter := ""
	if conf.Groups == "BINARIES" {
		filter = "alt.binaries.*"
	}
	groupsList, err := conn.List("ACTIVE", filter)
	if err != nil {
		fmt.Printf("Error while requesting list of groups: %v\n", err)
		return err
	}
	fmt.Println("Processing the groups")
	for _, group := range groupsList {
		groupData := strings.Split(string(group), " ")
		first, _ := strconv.Atoi(groupData[2])
		last, _ := strconv.Atoi(groupData[1])
		if conf.Groups == "ALL" || conf.Groups == "BINARIES" || slices.Contains(groups, groupData[0]) {
			serverGroups = append(serverGroups, groupData[0])
			if last > first {
				_, err := db.Exec(
					"INSERT INTO `groups` (`group_name`, `first_message_id`, `last_message_id`, `current_message_id`) VALUES (?, ?, ?, ?) ON DUPLICATE KEY UPDATE `first_message_id` = ?, `last_message_id` = ?",
					groupData[0],
					first,
					last,
					first,
					first,
					last,
				)
				if err != nil {
					fmt.Printf("Database error while processing '%s': %v\n", groupData[0], err)
					return err
				}
			}
		}
	}

	return nil

}

func readGroups(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		groups = append(groups, scanner.Text())
	}
	return scanner.Err()
}
