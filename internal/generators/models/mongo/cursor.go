package mongo

const cursorTmpl = `
package cursor

import (
	"encoding/base64"
	"fmt"
	"time"
)

const (
	cursorFormat = "cursor:v%d:%x:%s"
)

type Cursor struct {
	Version   int
	TimeStamp time.Time
	ID        string
}

func NewCursor(id string, t time.Time) *Cursor {
	return &Cursor{
		Version:   1,
		TimeStamp: t,
		ID:        id,
	}
}

func (c *Cursor) String() string {
	str := fmt.Sprintf(cursorFormat, c.Version, c.TimeStamp.UnixMicro(), c.ID)
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func FromString(s string) *Cursor {
	bts, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil
	}
	c := Cursor{}
	var epoch int64
	fmt.Sscanf(string(bts), cursorFormat, &c.Version, &epoch, &c.ID)
	c.TimeStamp = time.UnixMicro(epoch)
	return &c
}
`
