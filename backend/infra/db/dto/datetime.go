package dto

import (
	"database/sql"
	"database/sql/driver"
	"time"

	"github.com/aereal/nikki/backend/types"
)

// refs. https://www.sqlite.org/lang_datefunc.html#time_values
var parser = types.DateTimeParser("2006-01-02T15:04:05.999Z07:00")

type DateTime time.Time

var (
	_ driver.Valuer = DateTime{}
	_ sql.Scanner   = (*DateTime)(nil)
)

func (dt DateTime) Value() (driver.Value, error) {
	return parser.Format(time.Time(dt)), nil
}

func (dt *DateTime) Scan(src any) error {
	s, err := types.Cast[string](src)
	if err != nil {
		return err
	}
	t, err := parser.Parse(s)
	if err != nil {
		return err
	}
	*dt = DateTime(t)
	return nil
}
