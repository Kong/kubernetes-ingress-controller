echo "update\n"
sudo apt update

echo "sudo apt install redis-server\n"
sudo apt install redis-server -y

echo "sudo systemctl restart redis.service\n"
sudo systemctl restart redis.service

echo "sudo systemctl status redis"
sudo systemctl status redis

echo "ensure redis server is up"
redis-cli ping

if [ $? -eq 0 ]
then
  echo "Redis server has been up successfully."
else
  echo "Failed to checking redis status." >&2
fi

