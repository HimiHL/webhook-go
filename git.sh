#!/bin/bash
WEB_PATH=$1
WEB_USER='www'
WEB_USERGROUP='www'

echo "Pulling Source Code"
cd $WEB_PATH
sudo git reset --hard $2
sudo git clean -f
sudo git pull
sudo git checkout $3
echo "Reset Permission"
sudo chown -R $WEB_USER:$WEB_USERGROUP $WEB_PATH
echo "Finished"