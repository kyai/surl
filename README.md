# 短网址解决方案

解决的问题是将`数字`*（通常是数据库自增ID）*转换成固定长度且无序的`字符串`

## Go示例

```go
import "github.com/kyai/surl"
```

```go
c := surl.NewCreator(4)

c.SetKey(c.NewKey())

arr := []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 100, 1000, 10000}

for _, v := range arr {
    s, _ := c.Encode(v)
    n, _ := c.Decode(s)
    fmt.Println(v, "\t", s, "\t", n)
}
```

示例结果

```
1        1T8d    1
2        K1CH    2
3        V8kO    3
4        0xm6    4
5        kbhJ    5
6        wuyj    6
7        eKNQ    7
8        p3VT    8
9        FFIk    9
10       JMr9    10
100      0xm5    100
1000     K1Cr    1000
10000    u9Co    10000
```

## 算法
