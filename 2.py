#coding=utf-8 
import sys
import redis

if (len(sys.argv) < 3):
	print("Usage:pid(must)|count(must)|host|port")
if (len(sys.argv) < 5):
	host = '127.0.0.1'
	port = 6379
if (len(sys.argv) == 5):
	host = sys.argv[2]
	port = sys.argv[3]
pid = sys.argv[1]
count = sys.argv[2]
red = redis.Redis(host='127.0.0.1', port=6379, db=0)
i = red.flushdb()
print i
count_pid = "count_" + pid
res = red.set(count_pid,count)
print res
