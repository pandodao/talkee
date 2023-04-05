package core

import (
	"context"
	"strconv"
	"time"
)

type (
	Property struct {
		Key       string     `json:"key"`
		Value     string     `gorm:"type:text" json:"value"`
		UpdatedAt *time.Time `json:"updated_at"`
	}

	PropertyStore interface {
		// SELECT
		//  *
		// FROM @@table WHERE "key"=@key;
		Get(ctx context.Context, key string) (*Property, error)

		// UPDATE @@table
		//  {{set}}
		//    "value"=@value,
		//    "updated_at"=NOW()
		//  {{end}}
		// WHERE "key"=@key;
		Set(ctx context.Context, key string, value string) error
	}
)

func (p *Property) Time() time.Time {
	t, _ := time.Parse(time.RFC3339Nano, string(p.Value))
	return t
}

func (p *Property) String() string {
	return string(p.Value)
}

func (p *Property) Bytes() []byte {
	return []byte(p.Value)
}

func (p *Property) Int64() (int64, error) {
	return strconv.ParseInt(string(p.Value), 10, 64)
}
