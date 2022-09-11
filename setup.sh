#!/bin/bash

# install docker
echo "Installing Docker...\n"
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo apt-get install -y uidmap
dockerd-rootless-setuptool.sh install
echo 'export PATH=/usr/bin:$PATH' > ~/.bashrc
echo 'export DOCKER_HOST=unix:///run/user/1000/docker.sock' > ~/.bashrc
source ~/.bashrc
sudo setcap cap_net_bind_service=ep $HOME/bin/rootlesskit
sudo echo 'net.ipv4.ip_unprivileged_port_start=0' > sudo /etc/sysctl.conf
sudo sysctl --system

# setting tokens
echo -n "DISCORD_BOT_TOKEN: "
read BOT_TOKEN
sed -i -e "/BOT_TOKEN/c BOT_TOKEN = $BOT_TOKEN" .env

echo -n "DISCORD_GUILD_ID: "
read GUILD_ID
sed -i -e "/GUILD_ID/c GUILD_ID = $GUILD_ID" .env

echo -n "DISCORD_CHANNEL_ID: "
read CHANNEL_ID
sed -i -e "/CHANNEL_ID/c CHANNEL_ID = $CHANNEL_ID" .env

# setting for systemd
echo "Setting for systemd...\n"
sudo cp kenkoukun.service /etc/systemd/system/kenkoukun.service
sudo systemctl daemon-reload
sudo systemctl enable kenkoukun
sudo systemctl start kenkoukun
