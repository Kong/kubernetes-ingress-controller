#!/bin/bash

# install redis as a daemon on a linux environment
echo "install gcc\n"
sudo yum install --assumeyes gcc make

echo "cd /usr/local/src\n"
sudo cd /usr/local/src

echo "wget redis\n"
sudo wget http://download.redis.io/redis-stable.tar.gz

echo "untar redis\n"
sudo tar xvzf redis-stable.tar.gz

echo "remove downloaded redis\n"
sudo rm -f redis-stable.tar.gzcd redis-stable

echo "groupinstall development tools\n"
sudo yum groupinstall --assumeyes "Development Tools"

echo "make distclean\n"
sudo make distclean

echo "make"
sudo make

echo "sudo install tcl \n"
sudo yum install -y tcl

echo "start redis as daemon on port 6379"
redis-server --daemonize yes

echo "ensure redis server is up"
redis-cli ping

if [ $? -eq 0 ]
then
  echo "Redis server has been up successfully."
else
  echo "Failed to checking redis status." >&2
fi