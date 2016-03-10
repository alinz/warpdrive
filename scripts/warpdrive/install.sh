#!/bin/sh

rm -f /usr/local/bin/warp                                                      \
      /usr/local/bin/warp-build.sh                                             \
      /usr/local/bin/warp-login.sh                                             \
      /usr/local/bin/warp-publish.sh

cp ./warp.sh /usr/local/bin/warp
cp ./warp-build.sh /usr/local/bin/warp-build.sh
cp ./warp-login.sh /usr/local/bin/warp-login.sh
cp ./warp-publish.sh /usr/local/bin/warp-publish.sh

chmod +x /usr/local/bin/warp
