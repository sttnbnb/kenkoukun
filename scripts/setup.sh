#!/bin/bash

# install dependencies
sudo apt install -y make

# install docker
echo "Installing Docker..."
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER

# setting tokens
echo -n "DISCORD_BOT_TOKEN: "
read BOT_TOKEN
sed -i -e "/BOT_TOKEN/c BOT_TOKEN=$BOT_TOKEN" .env

echo -n "DISCORD_DEFAULT_GUILD_ID: "
read GUILD_ID
sed -i -e "/DEFAULT_GUILD_ID/c DEFAULT_GUILD_ID=$GUILD_ID" .env

echo -n "DISCORD_DEFAULT_CHANNEL_ID: "
read CHANNEL_ID
sed -i -e "/DEFAULT_CHANNEL_ID/c DEFAULT_CHANNEL_ID=$CHANNEL_ID" .env
