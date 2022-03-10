package main

import (
	"context"
	"fmt"
	"strings"
)

type Message struct {
	messageNo     int
	subject       string
	messageId     string
	from          string
	bytes         int
	date          int64
	header        string
	filename      string
	basefilename  string
	fileNo        int
	totalFiles    int
	segmentNo     int
	totalSegments int
	headerHash    string
	fileHash      string
	groupId       int
	groupName     string
}

func indexer(group string, ctx context.Context) error {

	defer wg.Done()

	var groupID, currentMessageID int

	fmt.Printf("Start indexing group '%s'\n", group)
	row := db.QueryRow(
		"SELECT `id`, `current_message_id` FROM `groups` WHERE `group_name` = ?",
		group,
	)
	if err := row.Scan(&groupID, &currentMessageID); err != nil {
		fmt.Printf("Database error when getting currentMessageID while processing '%s': %v\n", group, err)
		return err
	}

	conn, err := ConnectNNTP()
	defer DisconnectNNTP(conn)
	if err != nil {
		fmt.Printf("Error connecting to the usenet server while indexing '%s': %v\n", group, err)
		return err
	}
	_, firstMessageID, lastMessageID, err := conn.Group(group)
	if err != nil {
		fmt.Printf("Error retrieving group information form the usenet server while indexing '%s': %v\n", group, err)
		return err
	}
	DisconnectNNTP(conn)

	// return if nothing to index...
	if currentMessageID >= lastMessageID {
		fmt.Printf("No new messages in group '%s'\n", group)
		return nil
	}

	startMessageID := currentMessageID + 1
	if firstMessageID > currentMessageID {
		startMessageID = firstMessageID
	}

	firstMessageID = startMessageID
	currentMessageID = firstMessageID

	// for testing only
	// goal is to let one group fully run through to the lastMessageID
	lastMessageID = startMessageID + conf.Step

	endMessageID := lastMessageID
	if startMessageID+conf.Step < lastMessageID {
		endMessageID = startMessageID + conf.Step
	}

	for startMessageID < lastMessageID {

		conn, err := ConnectNNTP()
		defer DisconnectNNTP(conn)
		if err != nil {
			fmt.Printf("Error connecting to the usenet server while indexing '%s': %v\n", group, err)
			return err
		}
		if _, _, _, err := conn.Group(group); err != nil {
			fmt.Printf("Error selecting group while indexing '%s': %v\n", group, err)
			return err
		}
		results, err := conn.Overview(startMessageID, endMessageID)
		if err != nil {
			fmt.Printf("Error retrieving message overview from the usenet server while indexing '%s': %v\n", group, err)
			return err
		}
		DisconnectNNTP(conn)

		for id, overview := range results {

			select {
			case <-ctx.Done():
				fmt.Printf("Stopped indexing group '%s'\n", group)
				fmt.Printf("Messages %d to %d were indexed and saved to the DB\n", firstMessageID, currentMessageID)
				return nil
			default:
				var message Message
				message.messageNo = overview.MessageNumber
				message.subject = strings.ToValidUTF8(overview.Subject, "")
				message.messageId = strings.Trim(overview.MessageId, "<>")
				message.from = strings.ToValidUTF8(overview.From, "")
				message.bytes = overview.Bytes
				if date := overview.Date.Unix(); date < 0 {
					message.date = 0
				} else {
					message.date = date
				}
				message.fileNo = 1
				message.totalFiles = 1
				message.segmentNo = 1
				message.totalSegments = 1
				message.groupId = groupID
				message.groupName = group

				if err := parseSubject(&message); err != nil {
					// message probably did not contain a yEnc encoded file?

					// update group table
					_, err = db.Exec(
						"UPDATE `groups` SET `current_message_id` = ? WHERE `group_name` = ?",
						overview.MessageNumber,
						group,
					)
					if err != nil {
						fmt.Printf("Database error at update group when skipping segment while indexing '%s':: %v\n", group, err)
						return err
					}
				} else {
					if err := saveToDB(&message); err != nil {
						// message could not be saved to the DB
						fmt.Printf("Error saving message to the DB while indexing '%s': %v\n", group, err)
						fmt.Println("ID:", id, "Message No:", overview.MessageNumber, "/ Subject:", message.subject)
						return err
					}
				}

				currentMessageID = overview.MessageNumber
				counter = counter + 1
			}
		}

		// update start and end message id for next request
		startMessageID = endMessageID + 1
		if startMessageID+conf.Step < lastMessageID {
			endMessageID = startMessageID + conf.Step
		} else {
			endMessageID = lastMessageID
		}

	}

	fmt.Printf("Finished indexing group '%s'\n", group)
	fmt.Printf("Messages %d to %d were indexed and saved to the DB\n", firstMessageID, currentMessageID)

	return nil
}
