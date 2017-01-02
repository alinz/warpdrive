# WARPIFY

`warpify` is a client side for both `iOS` and `android`. The core logic of warpify is written in `Golang`, and uses samll `Java` and `Objective-C` code 
to wrap the functionality which exposes them to js side of react-native.

One of the main reason for choosing `Golang` was easy to write very complex operations and it has a comprehensive stream functionality. At the very begining, 
Warpify was designed to uses as little memory as possible. All the download, decreypt and unpacking the bundles on device was written to be stream compatiable.
Also by using Golang, we can support multiple target devices without spending time on two or more different languages.

> Note: if you plan to use the provided examples, make sure to use `Yarn`. Because Yarn local installation path is different than npm one.

## Installation

```
npm install react-native-warpdrive
```

### iOS setup

There are 2 things we need to import into `xcode`, `Warpify.framework` and `Warpify.xcodeproj`.

- create a group call Frameworks

<p align="center">
    <img src ="https://github.com/pressly/warpdrive/raw/master/docs/images/ios-step1.png" />
</p>

> this is optional but it is recommended

- right click on Frameworks group and select `Add Files to ...`

<p align="center">
    <img src ="https://github.com/pressly/warpdrive/raw/master/docs/images/ios-step2.png" />
</p>

- navigate to `node_modules/react-native-warpdrive/ios` and select `Warpify.frmework` file

<p align="center">
    <img src ="https://github.com/pressly/warpdrive/raw/master/docs/images/ios-step3.png" />
</p>

- right click on Libraries group and select `Add Files to ...`

<p align="center">
    <img src ="https://github.com/pressly/warpdrive/raw/master/docs/images/ios-step4.png" />
</p>

- navigate to `node_modules/react-native-warpdrive/ios` and select `Warpify.xcodeproj`

<p align="center">
    <img src ="https://github.com/pressly/warpdrive/raw/master/docs/images/ios-step5.png" />
</p>

- select the project name and go to `General` tab, in `Linked Frameworks and Libraries` at the bottom click on `+` button and select `libWarpify.a` and click on `Add` button

<p align="center">
    <img src ="https://github.com/pressly/warpdrive/raw/master/docs/images/ios-step6.png" />
</p>

<p align="center">
    <img src ="https://github.com/pressly/warpdrive/raw/master/docs/images/ios-step7.png" />
</p>

- we also need to add the `framework` as well, so go ahead and click on `+` button again and click on `Add Other...` button and navigate to `node_modules/react-native-warpdrive/ios` and select `Warpify.framework`

- if you build the project, you should get the linker error. In order to resolve this we need to add a path to framework.

- go to your project configuration and select the `Build Settings` and search for `Framework Search Paths` and add this path `$(SRCROOT)/../node_modules/react-native-warpdrive/ios`.

<p align="center">
    <img src ="https://github.com/pressly/warpdrive/raw/master/docs/images/ios-step8.png" />
</p>

- by now you should be able to compile and build the your project with warpify successfully.

- we also need to add the Warpify Header path to your project. So select the `Build Settings` and search for `Header Search Paths` and add this path `$(SRCROOT)/../node_modules/react-native-warpdrive/ios`.

<p align="center">
    <img src ="https://github.com/pressly/warpdrive/raw/master/docs/images/ios-step9.png" />
</p>

- very React-Native project at runtime requires a path to source bundle. `Warpify` is going to take over that responsibility but first we have to tell React-Native to use warpdrive.

`WarpifyManager` exposes only one class method. This method requires 3 arguments.

- defaultCycleName: is the name of the cycle name which you have provided by cli tool. This name will be used along side of forceUpdate
- groupName: if you plan to use react-native for share extension and you want to update the share extension as well, then you need to configure the groupName, otherwise pass `nil`.
- forceUpdate: forceUpdate is a built-in functionality which takes over updating app without using javascript.

- go to `AppDelegate.m` file. add a new header `#import "WarpifyManager.h"` and replace the 

```obj-c
jsCodeLocation = [[RCTBundleURLProvider sharedSettings] jsBundleURLForBundleRoot:@"index.ios" fallbackResource:nil];
```

to

```obj-c
jsCodeLocation = [WarpifyManager sourceBundleWithDefaultCycle:@"prod" groupName:nil forceUpdate:NO];
```

- the last part is to include the `WarpFile`. `WarpFile` is created by cli tool provided with `warpdrive`. Please refer to cli doc.

### android Setup