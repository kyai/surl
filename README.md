# 短网址解决方案

解决的问题是将`数字` *(通常是数据库自增ID)* 转换成固定长度且无序的`字符串`

可用于短网址、分享链接、订单号等业务场景

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



`数字ID`: 业务上使用的ID，如数据库自增ID

`字符ID`: 短网址后缀部分，如~~xxx.com/~~`abcd`

根据`数字ID`生成`字符ID`，需要满足以下条件：

* 长度可控
* 同一域名(业务)下唯一
* 能反向推导出该`数字ID`
* 当`数字ID`有规律变动时，对应的`字符ID`应无规律



### 步骤

#### 0. 字符ID长度

字符ID当然越短越好，同时要能满足业务所需的上限。通常使用`[0-9a-zA-Z]`来组成字符ID，即`62进制`，在明确这点后，可以得出各个长度的字符ID所能容纳的数字ID上限

```
字符ID长度    数字ID上限
4            62^4 ≈ 1400万
5            62^5 ≈ 9亿
6            62^6 ≈ 568亿
```

<font color=gray>由于数据库查询性能等原因，更大长度显得没有意义。</font>

以下步骤以长度4为例

#### 1. 将数字ID转换成62进制

`数字ID`为`100`，转换成62进制则是`1C`，根据长度补充位数则是`001C`

#### 2. 字典/KEY

上个步骤中得到的字符ID是有序的

```
100      001C
101      001D
102      001E
103      001F
...
```

打乱这个排列顺序，需要引入一个`字典`的机制

```
bjZprSIMLPQF1fDGT5c6iKezOxBNmhn32dsk4lwE9WY8oyCUHuAqVXJ0R7avtg
tkIzPvDYO40TLUiwo69prFZ5RcgNWnV2E1x8yBuSsKXfMQGdj3lhbemCqHJa7A
zKhYqQyNrVDdX1xIB9W3PSU0jLiAM4wETbcn7e2ClptFJG6sgRoHu58OZkvfma
4XyWJgNxeRBHESUtuT0PqCvDj6cLspdhVkan1AzFo8Ii3bmfY7GlKZQwO2M95r
```

生成这段字典的规则是

每一行包含组成字符ID的所有字符，且不重复；
列数和字符ID长度一致。

<font color=gray>对于短网址业务而言，域名和字典应该是一对一关系。</font>

#### 3. 打乱

有了上面的字典，那么可以根据不同的算法对字符ID进行处理，需要注意的是，该算法必须可逆。

几个关键词

`基准字符` 这里定义为`0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ`

`基准下标` 字符在基准字符里的下标；如`a`的基准下标为`10`

`字典下标` 字符的位置为行，字典该行中对应的下标；如`abcd`首位`a`的字典下标为`58`

`右移` 字符在基准字符里向右移动，超出范围则从首位计算；如`a`右移3位为`d`，`X`右移5位为`2`

这里使用的方法如下：

以100 -> 001C为例

* 除了末位，将每一位右移<u>下一位的字典下标值</u>

```
0(55)   0(10)   1(13)   C(21)   字典中对应的下标
+10     +13     +21
7       Y       m       C
```

* 再将末位右移<u>首位的字典下标值</u>

```
7(57)   Y       m       C
                        +57
7       Y       m       x
```

* 最后将每一位替换为<u>基准下标在字典中对应的字符</u>

```
7(7)    Y(60)   m(22)   x(33)   该字符的基准下标
│       │       │       │
M       7       U       k       字典中对应的字符
```

得到数字ID`100`的字符ID为`M7Uk`，对比顺序的字符ID：

```
100      001C      M7Uk
101      001D      Lljq
102      001E      VT1G
103      001F      S2lM
...
```

#### 4. 还原

将字符ID还原成数字ID，将以上算法逆向推导即可。



## 协议

[MIT](https://github.com/kyai/surl/blob/master/LICENSE)
