package asset

import (
	"talkee/core"
	"talkee/store/db"
)

func scanRow(scanner db.Scanner, output *core.Asset) error {
	defer scanner.Close()

	if scanner.Next() {
		if err := scanner.StructScan(
			output,
		); err != nil {
			return err
		}
	}

	return nil
}

func scanRows(scanner db.Scanner) ([]*core.Asset, error) {
	defer scanner.Close()
	outputs := make([]*core.Asset, 0)
	for scanner.Next() {
		output := &core.Asset{}
		err := scanner.StructScan(output)
		if err != nil {
			return nil, err
		}
		outputs = append(outputs, output)
	}

	return outputs, nil
}
