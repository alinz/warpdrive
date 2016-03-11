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
  echo ""
  echo "e.g."
  echo "    warp config -p etc/config"
  echo ""
}

#default is current localtion
CONFIG_PATH=.

while getopts "hp:" OPTION
do
  case $OPTION in
    h)
      usage
      exit 1
      ;;
    p)
      CONFIG_PATH=$OPTARG
      ;;
    ?)
      usage
      exit
      ;;
  esac
done

if [ ! -f ./.warpdrive/.token ]; then
  echo "you need to login first"
  exit
fi

APP_ID=$(cat ./.warpdrive/.appid)
CYCLE_ID=$(cat ./.warpdrive/.cycleid)
TOKEN=$(cat ./.warpdrive/.token)
DOMAIN=$(cat ./.warpdrive/.domain)

CONFIG_FILE=warpdrive.config

curl "$DOMAIN/apps/$APP_ID/cycles/$CYCLE_ID/config?jwt=$TOKEN" -o "$CONFIG_PATH/$CONFIG_FILE"
