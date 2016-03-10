#!/bin/sh

# 
# Created by Ali Najafizadeh on 2016-02-29.
# Copyright © 2016 Pressly. All rights reserved.
#

args=("$@")
COMMAND=${args[0]}
REST="${@:2}"

usage ()
{
  echo "warp v1.0.0"
  echo ""
  echo "Usage:"
  echo ""
  echo "warp [COMMAND] <OPTIONS>"
  echo ""
  echo "  COMMAND"
  echo "    login                           login to specific warpdrive server"
  echo "    build   <OPTIONS>               build and prepare before publishing"
  echo "    publish <OPTIONS>               publish to logged in server"
  echo ""
}

if [ "$#" -eq 0 ]; then
  usage
  exit
fi

case $COMMAND in
  login)
    bash warp-login.sh
    ;;
  build)
    bash warp-build.sh "$REST"
    ;;
  publish)
    bash warp-publish.sh "$REST"
    ;;
  ?)
    usage
    exit
    ;;
esac
