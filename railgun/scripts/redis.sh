#!/bin/bash
# install redis as a daemon on a linux environment
echo "install gcc"
sudo yum -y install gcc make
echo "cd /usr/local/src"
sudo cd /usr/local/src
echo "wget redis"
sudo wget http://download.redis.io/redis-stable.tar.gz
echo "untar redis"
sudo tar xvzf redis-stable.tar.gz
echo "remove downloaded redis"
sudo rm -f redis-stable.tar.gzcd redis-stable
echo "groupinstall development tools"
sudo yum groupinstall "Development Tools"
echo "make distclean"
sudo make distclean
echo "make"
sudo make
echo "install tcl"
sudo yum install -y tcl
echo "start redis as daemon on port 6379"
redis-server --daemonize yes
echo "ensure redis server is up"
$redis = new Redis();
$redis->connect('127.0.0.1', 6379);

echo $redis->ping();
$redis->disconnect()