# warp aka `warpdrive cli`

warp is a cli tool which connects to warpdrive over GRPC and does 2 things

1 - build the react-native and make it ready for upload to warpdrive server

for example, the following command will bundle your react-native app for ios and put the created bundle to `.bundles/ios` folder of react-native project.

```bash
warp bundle -p ios
```

`android` and `all` can be used instead of `ios`. `all` will build both ios and android bundles and put them into their target folder.

> NOTE: make sure to put `.bundles` into your `.gitignore` so it won't accidently added into your git history. Also you do not need to delete `.bundles` folder. the command will remove the content of the `.bundles` folder upon execution.

2 - package and publish the new bundle to warpdrive server



