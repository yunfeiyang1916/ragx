# 类型转换


### 转换失败，返回默认值0,nil,false等

```go
cast.ToString("mayonegg")         // "mayonegg"
cast.ToString(8)                  // "8"
cast.ToString(8.31)               // "8.31"
cast.ToString([]byte("one time")) // "one time"
cast.ToString(nil)                // ""

var foo interface{} = "one more time"
cast.ToString(foo)                // "one more time"

cast.ToInt(8)                  // 8
cast.ToInt(8.31)               // 8
cast.ToInt("8")                // 8
cast.ToInt(true)               // 1
cast.ToInt(false)              // 0

var eight interface{} = 8
cast.ToInt(eight)              // 8
cast.ToInt(nil)                // 0

cast.ToXXX()
...
```

### 转换失败，返回错误
```go
v,err:=cast.ToIntE(false)              // 0,nil

```