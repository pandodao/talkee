package site

import (
	"database/sql"
	"talkee/core"
	"talkee/store/db"

	"github.com/lib/pq"
)

func scanRow(scanner db.Scanner) (*core.Site, error) {
	defer scanner.Close()
	output := &core.Site{}
	if scanner.Next() {
		err := scanner.Scan(
			&output.ID, &output.UserID,
			pq.Array(&output.Origins),
			&output.Name,
			&output.UseArweave, &output.RewardStrategy,
			&output.CreatedAt, &output.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
	}
	if output.ID == 0 {
		return nil, sql.ErrNoRows
	}

	return output, nil
}

func scanRows(scanner db.Scanner) ([]*core.Site, error) {
	defer scanner.Close()
	outputs := make([]*core.Site, 0)
	for scanner.Next() {
		output := &core.Site{}
		err := scanner.Scan(
			&output.ID, &output.UserID,
			pq.Array(&output.Origins),
			&output.Name,
			&output.UseArweave, &output.RewardStrategy,
			&output.CreatedAt, &output.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		outputs = append(outputs, output)
	}

	return outputs, nil
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
