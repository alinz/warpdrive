#!/bin/bash

#
# Created by Ali Najafizadeh on 2016-02-29.
# Copyright © 2016 Pressly. All rights reserved.
#

rm -rf ./.warpdrive
mkdir ./.warpdrive

read -p "Email: " EMAIL;
read -s -p "Password: " PASSWORD;
echo;
read -p "Domain: " DOMAIN
read -p "App ID: " APPID
read -p "CYCLE ID: " CYCLEID

function jsonValue() {
  KEY=$1
  num=1
  awk -F"[,:}]" '{for(i=1;i<=NF;i++){if($i~/'$KEY'\042/){print $(i+1)}}}'      \
  | tr -d '"' | sed -n ${num}p
}

JWT=$(curl -sS -H "Content-Type: application/json" -X POST                     \
     -d "{\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\"}"                    \
     "$DOMAIN/session/start" | jsonValue jwt)

if [ -z "$JWT" ]; then
  echo "login failed"
else
  ACCESS=$(curl --write-out "%{http_code}\n" --silent --output /dev/null \
           "$DOMAIN/apps/$APPID/cycles/$CYCLEID?jwt=$JWT")

  if [ "$ACCESS" == "401" ] || [ "$ACCESS" == "400" ]; then
    echo "you don't access to this app"
    exit
  fi

  echo "$DOMAIN" > ./.warpdrive/.domain
  echo "$JWT" > ./.warpdrive/.token
  echo "$APPID" > ./.warpdrive/.appid
  echo "$CYCLEID" > ./.warpdrive/.cycleid
fi
