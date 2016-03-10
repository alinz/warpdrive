#!/bin/bash

#
# Created by Ali Najafizadeh on 2016-02-29.
# Copyright © 2016 Pressly. All rights reserved.
#

usage ()
{
  echo "warp publish v1.0.0"
  echo ""
  echo "Usage:"
  echo ""
  echo " Options:"
  echo "    -h                        show usage/help"
  echo "    -p <string> required      platform 'ios' or 'android'"
  echo "    -v <x.y.z>  required      bundle's version"
  echo "    -r <string> optional      release's note"
  echo ""
}

PLATFORM=
VERSION=
NOTE=

while getopts "hp:v:a:c:t:d:r" OPTION
do
  case $OPTION in
    h)
      usage
      exit 1
      ;;
    p)
      PLATFORM=$OPTARG
      ;;
    v)
      VERSION=$OPTARG
      ;;
    r)
      NOTE=$OPTARG
      ;;
    ?)
      usage
      exit
      ;;
  esac
done

if [ -z "$PLATFORM" ] || [ -z "$VERSION" ]; then
  usage
  exit
fi

if [ ! -f ./.warpdrive/.token ]; then
  echo "you need to login first"
  exit
fi


APP_ID=$(cat ./.warpdrive/.appid)
CYCLE_ID=$(cat ./.warpdrive/.cycleid)
TOKEN=$(cat ./.warpdrive/.token)
DOMAIN=$(cat ./.warpdrive/.domain)

#
# upload single file bundle

upload () {
  local FILENAME = $1
  curl -i                                                                      \
       -X POST                                                                 \
       -H "Content-Type: multipart/form-data"                                  \
       -F "filename=$FILENAME"                                                 \
       -F "file=@.release/$FILENAME"                                           \
       "$DOMAIN/apps/$APP_ID/cycles/$CYCLE_ID/releases?platform=$PLATFORM"     \
       "&version=$VERSION&note=$NOTE&jwt=$TOKEN"
}

for filename in ./release/*
do
  if [ -f "$filename" ];then
    upload "$filename"
  fi
done
