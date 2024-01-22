package ormx

import (
	"context"
	"testing"
	"time"
)

func TestNewLog(t *testing.T) {
	lg := NewLog("info")
	t.Log(lg)

	lg.Trace(context.Background(), time.Now().Add(time.Millisecond*-100), func() (sql string, rowsAffected int64) {
		return "select * from users limit 10", 10
	}, nil)
}
