package models

import (
	"fmt"
	"strings"
	"time"
)

type DateOnly time.Time

func ParseDateOnly(value string) (DateOnly, error) {
	t, err := time.Parse("2006-01-02", value)
	if err != nil {
		return DateOnly{}, err
	}
	return DateOnly(t), nil
}

func (d DateOnly) MarshalJSON() ([]byte, error) {
	t := time.Time(d)
	return []byte(fmt.Sprintf("\"%s\"", t.Format("2006-01-02"))), nil
}

func (d *DateOnly) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	*d = DateOnly(t)
	return nil
}

func (d DateOnly) ToTime() time.Time {
	return time.Time(d)
}
func (d DateOnly) String() string {
	return d.ToTime().Format("2006-01-02")
}
func (d DateOnly) IsZero() bool {
	return d.ToTime().IsZero()
}
