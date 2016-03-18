#!/bin/bash

#
# Created by Ali Najafizadeh on 2016-02-29.
# Copyright © 2016 Pressly. All rights reserved.
#

usage ()
{
  echo "warp build v1.0.0"
  echo ""
  echo "Usage:"
  echo ""
  echo " Options:"
  echo "    -h                        show usage/help"
  echo "    -p <string> required      platform 'ios' or 'android'"
  echo ""
}

PLATFORM=

while getopts "hp:" OPTION
do
  case $OPTION in
    h)
      usage
      exit 1
      ;;
    p)
      PLATFORM=$OPTARG
      ;;
    ?)
      usage
      exit
      ;;
  esac
done

if [ -z "$PLATFORM" ]; then
  usage
  exit
fi

rm -rf ./.release
mkdir ./.release

node --max_old_space_size=8192                                                 \
  node_modules/react-native/local-cli/cli.js bundle                            \
  --platform "$PLATFORM"                                                       \
  --entry-file "index.$PLATFORM.js"                                            \
  --bundle-output ./.release/main.jsbundle                                     \
  --assets-dest ./.release                                                     \
  --dev false
