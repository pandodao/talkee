package user

import (
	"database/sql"
	"talkee/core"
	"talkee/store/db"
)

func scanRow(scanner db.Scanner) (*core.User, error) {
	defer scanner.Close()
	output := &core.User{}
	if scanner.Next() {
		if err := scanner.StructScan(
			output,
		); err != nil {
			return nil, err
		}
	}

	if output.ID == 0 {
		return nil, sql.ErrNoRows
	}

	return output, nil
}

func scanReturnID(scanner db.Scanner) (uint64, error) {
	defer scanner.Close()
	var id uint64

	if scanner.Next() {
		if err := scanner.Scan(
			&id,
		); err != nil {
			return 0, err
		}
	}

	return id, nil
}
