package helpers

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

const randomStringLetterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	randomStringletterIdxBits = 6
	randomStringletterIdxMask = 1<<randomStringletterIdxBits - 1
	randomStringletterIdxMax  = 63 / randomStringletterIdxBits
)

func RandomString(n int) string {
	// from : https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
	b := make([]byte, n)
	for i, cache, remain := n-1, rand.Int63(), randomStringletterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), randomStringletterIdxMax
		}
		if idx := int(cache & randomStringletterIdxMask); idx < len(randomStringLetterBytes) {
			b[i] = randomStringLetterBytes[idx]
			i--
		}
		cache >>= randomStringletterIdxBits
		remain--
	}

	return string(b)
}

const letterNumberBytes = "0123456789"

func RandomStringIntOnly(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letterNumberBytes[rand.Int63()%int64(len(letterNumberBytes))]
	}
	return string(b)
}

type Date time.Time

var _ json.Unmarshaler = &Date{}
var _ json.Marshaler = &Date{}

func (d *Date) Time() time.Time {
	return time.Time(*d)
}

func (d *Date) String() string {
	return d.Time().Format("2006-01-02")
}

func (d *Date) Scan(value interface{}) error {
	t, ok := value.(time.Time)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	tmpDate := Date(t)
	*d = tmpDate
	return nil
}

func (d Date) Value() (driver.Value, error) {
	return driver.Value(d.Time().Format("2006-01-02")), nil
}

func (Date) GormDataType() string {
	return "date"
}

func (d *Date) UnmarshalJSON(bs []byte) error {
	var s string
	err := json.Unmarshal(bs, &s)
	if err != nil {
		return err
	}
	t, err := time.ParseInLocation("2006-01-02", s, time.UTC)
	if err != nil {
		return err
	}
	*d = Date(t)
	return nil
}

func (d *Date) MarshalJSON() ([]byte, error) {
	return []byte(`"` + d.Time().UTC().Format("2006-01-02") + `"`), nil
}

type JSONB struct{}

func (m *JSONB) Value() (driver.Value, error) {
	val, err := json.Marshal(m)

	return driver.Value(string(val)), err
}

func (JSONB) GormDataType() string {
	return "JSONB"
}

func Scan(val interface{}, res interface{}) error {
	t, ok := val.([]byte)
	if !ok {
		return errors.New(fmt.Sprintf("Failed to unmarshal JSONB value: %v", val))
	}

	return json.Unmarshal(t, &res)
}
