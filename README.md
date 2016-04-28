# 秒杀系统

## 简要信息

设计并实现一个简单的秒杀系统，主要包括如下功能：

+ 商品秒杀，支持多种类商品

http://serverip:8001/killing/seckilling?userid=123456&productid=111

+ 用户秒杀结果查询

http://serverip:8001/killing/queryUserSeckillingInfo?userid=123456&productid=111

+ 商品秒杀结果查询

http://serverip:8001/killing/queryProductSeckillingInfo?productid=111

+ 商品初始信息添加

http://serverip:9001/killing/addproduct?productid=111&productnum=100&starttime=20160429102103

+ 商品初始信息查询

http://serverip:9001/killing/getproductlist

+ 秒杀数据清空

http://serverip:9001/killing/cleandb?productid=111

+ 简单的WEB UI

## 使用方法

命令行下：

```shell
make
./load.sh
```





