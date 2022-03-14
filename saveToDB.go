package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

func saveToDB(message *Message) error {

	var headerID, newHeader, fileID, newFile, posterID int64

	// start transaction
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		fmt.Printf("Database error when beginning transaction: %v\n", err)
		return err
	}

	// get headerID
	headerID, err = getID(tx, ctx, "header_hashes", "hash", message.headerHash)
	if err != nil {
		fmt.Printf("Database error at select headerID: %v\n", err)
		tx.Rollback()
		return err
	}
	if headerID == 0 {
		// update header_hashes table if hash is not yet in it
		headerID, err = updateTable(tx, ctx, "header_hashes", "hash", message.headerHash)
		if err != nil || headerID == 0 {
			if err == nil {
				err = errors.New("return value for id is zero")
			}
			fmt.Printf("Database error at insert headerHash: %v\n", err)
			tx.Rollback()
			return err
		}
		newHeader = 1
	}

	// get fileID
	fileID, err = getID(tx, ctx, "file_hashes", "hash", message.fileHash)
	if err != nil {
		fmt.Printf("Database error at select fileID: %v\n", err)
		tx.Rollback()
		return err
	}
	if fileID == 0 {
		// update header_hashes table if hash is not yet in it
		fileID, err = updateTable(tx, ctx, "file_hashes", "hash", message.fileHash)
		if err != nil || fileID == 0 {
			if err == nil {
				err = errors.New("return value for id is zero")
			}
			fmt.Printf("Database error at insert fileHash: %v\n", err)
			tx.Rollback()
			return err
		}
		newFile = 1
	}

	// get posterID
	posterID, err = getID(tx, ctx, "poster", "poster", message.from)
	if err != nil {
		fmt.Printf("Database error at select posterID: %v\n", err)
		tx.Rollback()
		return err
	}
	if posterID == 0 {
		// update header_hashes table if hash is not yet in it
		posterID, err = updateTable(tx, ctx, "poster", "poster", message.from)
		if err != nil || posterID == 0 {
			if err == nil {
				err = errors.New("return value for id is zero")
			}
			fmt.Printf("Database error at insert poster: %v\n", err)
			tx.Rollback()
			return err
		}
	}

	// update groups_to_files table
	_, err = tx.ExecContext(
		ctx,
		"INSERT IGNORE INTO `groups_to_files` (`group_id`, `file_id`) VALUES(?, ?)",
		message.groupId,
		fileID,
	)
	if err != nil {
		fmt.Printf("Database error at update groups_to_files: %v\n", err)
		tx.Rollback()
		return err
	}

	// update segment table
	_, err = tx.ExecContext(
		ctx,
		"INSERT INTO `segments` (`file_id`, `segment_id`, `segment_no`, `size`, `date`, `poster`) VALUES(?, ?, ?, ?, ?, ?)",
		fileID,
		message.messageId,
		message.segmentNo,
		message.bytes,
		message.date,
		posterID,
	)
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
	_, err = tx.ExecContext(
		ctx,
		"INSERT INTO `files` (`id`, `header_id`, `subject`, `file_no`, `segments`, `total_segments`, `size`, `date`, `poster`) VALUES(?, ?, ?, ?, 1, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE `segments` = `segments` + 1, `size` = `size` + ?, `date` = ?",
		fileID,
		headerID,
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

	// update headers table
	_, err = tx.ExecContext(
		ctx,
		"INSERT INTO `headers` (`id`, `files`, `total_files`, `size`, `date`, `poster`) VALUES(?, 1, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE `files` = `files` + ?, `size` = `size` + ?, `date` = ?",
		headerID,
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

	// update groups table
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

func getID(tx *sql.Tx, ctx context.Context, table string, name string, value string) (id int64, err error) {
	row := tx.QueryRowContext(
		ctx,
		"SELECT `id` FROM `"+table+"` WHERE `"+name+"` = ?",
		value,
	)
	if err = row.Scan(&id); err != nil && err != sql.ErrNoRows {
		return 0, err
	} else if err == sql.ErrNoRows {
		return 0, nil
	}
	return id, nil
}

func updateTable(tx *sql.Tx, ctx context.Context, table string, name string, value string) (id int64, err error) {
	_, err = tx.ExecContext(
		ctx,
		"INSERT IGNORE INTO `"+table+"` (`"+name+"`) VALUES(?)",
		value,
	)
	if err != nil {
		return 0, err
	}
	return getID(tx, ctx, table, name, value)
}
