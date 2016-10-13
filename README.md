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