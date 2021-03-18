#!/bin/bash
WEB_PATH=$1
WEB_USER='www'
WEB_USERGROUP='www'

echo "Pulling Source Code"
cd $WEB_PATH
if [ ! -d ".git" ]; then
    sudo git clone $3 $1
else
    sudo git reset --hard origin/$2
    sudo git clean -f
    sudo git pull
    sudo git checkout $2
fi
echo "Reset Permission"
sudo chown -R $WEB_USER:$WEB_USERGROUP $WEB_PATH
if [ -f "deploy.json" ]; then
    echo "Auto Deploying.."
    sudo /home/deploy/deploy -path=$WEB_PATH -config=deploy.json
    echo "Deployed"
fi