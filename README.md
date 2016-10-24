采用bloomfilter实现大单量的去重. 对于电商场景, 一个单据的事件可能被发多次.

安装
```
go get git.ishopex.cn/stat/uniq
make
```

运行:
```
  % ./uniq
  Usage of ./uniq:
    -datadir string
      	数据文件夹, 文件名为{datadir}/{day}.blf, day=1970/01/01到今天的天数. (default "data")
    -days int
      	去重最多存储天数
    -number uint
      	布隆过滤器数量, 通常是总条数的20倍. (default 2147483648)
    -port string
      	服务端口号 (default "6532")
```

使用memcache协议访问.   get指令检查某id是否存在,  第一次为空, 第二次的检查则会返回"true".