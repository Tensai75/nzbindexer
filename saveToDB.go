package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

func saveToDB(message *Message) error {

	var result sql.Result
	var posterID int64

	newFile := 0
	newHeader := 0

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		fmt.Printf("Database error when beginning transaction: %v\n", err)
		return err
	}

	// update poster table if poster not found
	// may give deadlock so try for tree times
	for i := 0; i < 3; i++ {
		result, err = tx.ExecContext(
			ctx,
			"INSERT IGNORE INTO `poster` (`poster`) VALUES(?)",
			message.from,
		)
		if err == nil || (errors.As(err, &mysqlError) && mysqlError.Number != 1213) {
			break
		}
		time.Sleep(1 * time.Second)
	}
	if err != nil {
		fmt.Printf("Database error at insert poster: %v\n", err)
		tx.Rollback()
		return err
	}
	// get posterID
	row := tx.QueryRowContext(
		ctx,
		"SELECT `id` FROM `poster` WHERE `poster` = ?",
		message.from,
	)
	if err = row.Scan(&posterID); err == sql.ErrNoRows {
		fmt.Printf("Database error at select poster: %v\n", err)
		tx.Rollback()
		return err
	}

	// update groups_to_files table
	result, err = tx.ExecContext(
		ctx,
		"INSERT IGNORE INTO `groups_to_files` (`group_id`, `file`) VALUES(?, ?)",
		message.groupId,
		message.fileHash,
	)
	if err != nil {
		fmt.Printf("Database error at update groups_to_files: %v\n", err)
		tx.Rollback()
		return err
	}

	// update segment table
	// may give deadlock so try for tree times
	for i := 0; i < 3; i++ {
		result, err = tx.ExecContext(
			ctx,
			"INSERT INTO `segments` (`file_hash`, `segment_id`, `segment_no`, `size`, `date`, `poster`) VALUES(?, ?, ?, ?, ?, ?)",
			message.fileHash,
			message.messageId,
			message.segmentNo,
			message.bytes,
			message.date,
			posterID,
		)
		if err == nil || (errors.As(err, &mysqlError) && mysqlError.Number != 1213) {
			break
		}
		time.Sleep(1 * time.Second)
	}
	if errors.As(err, &mysqlError) && mysqlError.Number != 1062 {
		fmt.Printf("Database error at update segments: %v\n", err)
		tx.Rollback()
		return err
	} else if errors.As(err, &mysqlError) && mysqlError.Number == 1062 {
		// simply commit and return if segment is already indexed
		err = tx.Commit()
		if err != nil {
			fmt.Printf("Database error when commiting changes: %v\n", err)
			return err
		}
		return nil
	}

	// update files table
	result, err = tx.ExecContext(
		ctx,
		"INSERT INTO `files` (`hash`, `header_hash`, `subject`, `file_no`, `segments`, `total_segments`, `size`, `date`, `poster`) VALUES(?, ?, ?, ?, 1, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE `segments` = `segments` + 1, `size` = `size` + ?, `date` = ?",
		message.fileHash,
		message.headerHash,
		message.subject,
		message.fileNo,
		message.totalSegments,
		message.bytes,
		message.date,
		posterID,
		message.bytes,
		message.date,
	)
	if err != nil {
		fmt.Printf("Database error at update files: %v\n", err)
		tx.Rollback()
		return err
	}
	if rows, _ := result.RowsAffected(); rows == 1 {
		newFile = 1
	}

	// update header table
	result, err = tx.ExecContext(
		ctx,
		"INSERT INTO `headers` (`hash`, `files`, `total_files`, `size`, `date`, `poster`) VALUES(?, 1, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE `files` = `files` + ?, `size` = `size` + ?, `date` = ?",
		message.headerHash,
		message.totalFiles,
		message.bytes,
		message.date,
		posterID,
		newFile,
		message.bytes,
		message.date,
	)
	if err != nil {
		fmt.Printf("Database error at update headers: %v\n", err)
		tx.Rollback()
		return err
	}
	if rows, _ := result.RowsAffected(); rows == 1 {
		newHeader = 1
	}

	// update group table
	_, err = tx.ExecContext(
		ctx,
		"UPDATE `groups` SET `current_message_id` = ?,`headers` = `headers` + ?,`files` = `files` + ?, `segments` = `segments` + 1,  `size` = `size` + ?, `date` = ?  WHERE `group_name` = ?",
		message.messageNo,
		newHeader,
		newFile,
		message.bytes,
		message.date,
		message.groupName,
	)
	if err != nil {
		fmt.Printf("Database error at update group: %v\n", err)
		tx.Rollback()
		return err
	}

	// commit changes
	err = tx.Commit()
	if err != nil {
		fmt.Printf("Database error when commiting changes: %v\n", err)
		return err
	}

	return nil

}
