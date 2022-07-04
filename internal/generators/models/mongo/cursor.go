package mongo

const cursorTmpl = `
package cursor

import (
	"bytes"
	"encoding/gob"
	"strings"
	"time"
)

func init() {
	gob.Register(time.Time{})

}

type Cursor struct {
	Version     int
	ID          string
	SortBy      string
	OffsetValue interface{}
}

func NewCursor(id, sortBy string, OffsetValue interface{}) *Cursor {
	return &Cursor{
		Version:     1,
		ID:          id,
		SortBy:      sortBy,
		OffsetValue: OffsetValue,
	}
}

func (c *Cursor) String() string {
	buffer := new(bytes.Buffer)
	if err := gob.NewEncoder(buffer).Encode(c); err != nil {
		panic(err)
	}
	return buffer.String()
}

func FromString(s string) *Cursor {
	c := Cursor{}
	if err := gob.NewDecoder(strings.NewReader(s)).Decode(&c); err != nil {
		panic(err)
	}
	return &c
}
`
