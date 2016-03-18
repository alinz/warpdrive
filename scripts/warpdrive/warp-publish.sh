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

function jsonValue() {
  KEY=$1
  num=1
  awk -F"[,:}]" '{for(i=1;i<=NF;i++){if($i~/'$KEY'\042/){print $(i+1)}}}'      \
  | tr -d '"' | sed -n ${num}p
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

echo "Uploading started"

ALLFILES=$(find .release -type f | sed "s/^\.release\///" | \
           awk '{print "-F \"filename[]="$0"\""" -F \"file=@.release/"$0"\""}')
ALLFILES=$(echo "$ALLFILES" | tr "\n" ' ')

COMMAND="curl -sS -i -X POST -H 'Content-Type: multipart/form-data' "$ALLFILES" '$DOMAIN/apps/$APP_ID/cycles/$CYCLE_ID/releases?platform=$PLATFORM&version=$VERSION&note=$NOTE&jwt=$TOKEN'"

echo "Uploading ended"

RESULT=$(eval $COMMAND)
RELEASE_ID=$(echo $RESULT | jsonValue id)

echo "Locking release"

curl -X PATCH "$DOMAIN/apps/$APP_ID/cycles/$CYCLE_ID/releases/$RELEASE_ID/lock?jwt=$TOKEN"

echo "Done"
