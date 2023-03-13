package core

import (
	"fmt"
	"strings"
)

const (
	ActionDeposit = "deposit"
)

type (
	SnapshotAction struct {
		Action   string `json:"a"`
		Param1   string `json:"1"`
		Param2   string `json:"2"`
		Param3   string `json:"3"`
		Param4   string `json:"4"`
		Param5   string `json:"5"`
		FollowID string `json:"f"`
	}
)

func (act *SnapshotAction) Unmarshal(input []byte) error {
	s := string(input)
	parts := strings.Split(s, ",")
	if len(parts) < 2 {
		return ErrCorruptedSnapshotAction
	}

	// the 1st part should be a valid action name
	if !isSupportedActionName(parts[0]) {
		return ErrCorruptedSnapshotAction
	}

	// the 2nd part should be the follow id, which is an uuid
	if len(parts[1]) != 36 {
		return ErrCorruptedSnapshotAction
	}

	act.Action = strings.ToLower(parts[0])
	act.FollowID = strings.ToLower(parts[1])

	// the following are params
	for i := 2; i < len(parts); i++ {
		switch i {
		case 2:
			act.Param1 = parts[i]
		case 3:
			act.Param2 = parts[i]
		case 4:
			act.Param3 = parts[i]
		case 5:
			act.Param4 = parts[i]
		case 6:
			act.Param5 = parts[i]
		}
	}

	return nil
}

func (act *SnapshotAction) Marshal() ([]byte, error) {
	inner := ""
	if act.Param1 != "" {
		inner += act.Param1 + ","
	}
	if act.Param2 != "" {
		inner += act.Param2 + ","
	}
	if act.Param3 != "" {
		inner += act.Param3 + ","
	}
	if act.Param4 != "" {
		inner += act.Param4 + ","
	}
	if act.Param5 != "" {
		inner += act.Param5 + ","
	}
	// remove the "," at the end
	if len(inner) > 0 {
		inner = inner[:len(inner)-1]
	}
	ret := fmt.Sprintf("%s,%s,%s", act.Action, act.FollowID, inner)
	return []byte(ret), nil
}

func isSupportedActionName(name string) bool {
	names := []string{
		ActionDeposit,
	}
	for _, n := range names {
		if n == name {
			return true
		}
	}
	return false
}
