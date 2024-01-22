# ormx

## 链路追踪（v0.0.1）

```go
package test

import (
	"context"
	"git.cht-group.net/cht-group-go/engine/tracex"
	"git.cht-group.net/cht-group-go/ormx"
	"testing"
	"time"
)

func TestOrmTrace(t *testing.T) {
	lg := ormx.NewLog("info")
	//设置回调函数
	lg.SetTrace(tracex.OrmTrace)
	//打开开关
	//ormx.JaegerEnable = true
	lg.Trace(context.Background(), time.Now().Add(time.Millisecond*-100), func() (sql string, rowsAffected int64) {
		return "select * from users limit 10", 10
	}, nil)
}
```
