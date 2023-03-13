package attachment

import (
	"talkee/core"
	"talkee/store/db"
)

func scanRow(scanner db.Scanner) (*core.Attachment, error) {
	defer scanner.Close()
	output := &core.Attachment{}
	if scanner.Next() {
		if err := scanner.StructScan(
			output,
		); err != nil {
			return nil, err
		}
	}

	return output, nil
}

func scanRows(scanner db.Scanner) ([]*core.Attachment, error) {
	defer scanner.Close()
	outputs := make([]*core.Attachment, 0)
	for scanner.Next() {
		output := &core.Attachment{}
		err := scanner.StructScan(output)
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
