package comment

import (
	"database/sql"
	"talkee/core"
	"talkee/store/db"
)

func scanRow(scanner db.Scanner) (*core.Comment, error) {
	defer scanner.Close()
	output := &core.Comment{}
	if scanner.Next() {
		user := &core.User{}
		err := scanner.Scan(
			&output.ID, &output.UserID,
			&output.SiteID, &output.Slug,
			&output.FavorCount, &output.ReplyCount,
			&output.ArweaveTxHash,
			&output.Content,
			&user.MixinUserID, &user.MixinIdentityNumber, &user.FullName, &user.AvatarURL, &user.MvmPublicKey,
			&output.CreatedAt, &output.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		user.ID = output.UserID
		output.Creator = user
	}

	if output.ID == 0 {
		return nil, sql.ErrNoRows
	}

	return output, nil
}

func scanRowFav(scanner db.Scanner) (*core.CommentFavourite, error) {
	defer scanner.Close()
	output := &core.CommentFavourite{}
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

func scanRows(scanner db.Scanner) ([]*core.Comment, error) {
	defer scanner.Close()
	outputs := make([]*core.Comment, 0)
	for scanner.Next() {
		output := &core.Comment{}
		user := &core.User{}
		err := scanner.Scan(
			&output.ID, &output.UserID,
			&output.SiteID, &output.Slug,
			&output.FavorCount, &output.ReplyCount,
			&output.ArweaveTxHash,
			&output.Content,
			&user.MixinUserID, &user.MixinIdentityNumber, &user.FullName, &user.AvatarURL, &user.MvmPublicKey,
			&output.CreatedAt, &output.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		user.ID = output.UserID
		output.Creator = user
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
