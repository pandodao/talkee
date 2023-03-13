package reward

import (
	"talkee/store/db"
)

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
