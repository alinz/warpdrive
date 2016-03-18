#!/bin/bash

#
# Created by Ali Najafizadeh on 2016-02-29.
# Copyright © 2016 Pressly. All rights reserved.
#

usage ()
{
  echo "warp build v1.0.0"
  echo ""
  echo "you need to login first"
  echo "Usage:"
  echo ""
  echo " Options:"
  echo "    -h                        show usage/help"
  echo "    -p <string> optional      path to save config"
  echo "    -a <number> optional      override app_id"
  echo "    -c <number> optional      override cycle_id"
  echo ""
  echo "e.g."
  echo "    warp config -p etc/config"
  echo ""
}

if [ ! -f ./.warpdrive/.token ]; then
  echo "you need to login first"
  exit
fi

#default is current localtion
CONFIG_PATH=.

APP_ID=$(cat ./.warpdrive/.appid)
CYCLE_ID=$(cat ./.warpdrive/.cycleid)
TOKEN=$(cat ./.warpdrive/.token)
DOMAIN=$(cat ./.warpdrive/.domain)

while getopts "hp:a:c:" OPTION
do
  case $OPTION in
    h)
      usage
      exit 1
      ;;
    p)
      CONFIG_PATH=$OPTARG
      ;;
    a)
      APP_ID=$OPTARG
      ;;
    c)
      CYCLE_ID=$OPTARG
      ;;
    ?)
      usage
      exit
      ;;
  esac
done

CONFIG_FILE=warpdrive.config

# TODO: check if user can access app_id and cycle_id

curl "$DOMAIN/apps/$APP_ID/cycles/$CYCLE_ID/config?jwt=$TOKEN"                 \
     --silent -o "$CONFIG_PATH/$CONFIG_FILE"
