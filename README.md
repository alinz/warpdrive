# Warpdrive

 COMMAND_CA=cert/ca-command.crt COMMAND_CRT=cert/cli.crt COMMAND_KEY=cert/cli.key COMMAND_ADDR=command:10000 ./warp publish -a share -p ios -r dev -v 1.1.3 -n test

# Dev Setup

in order to compile the code for android and ios, you need to have gomobile install

first install the gomobile

```
go get -u golang.org/x/mobile/cmd/gomobile 
```

make sure you have install ndk as well if you are not sure, look into this tutorial

``` 
https://developer.android.com/ndk/guides/index.html
```

then you have to initialize gomobile. this is one time operation

```
gomobile init
gomobile init -ndk ~/Library/Android/sdk/ndk-bundle/
```

### Postgres Setup

before running warpdrive, make sure you have a right role and database. you can run the follwoing sql in Postgres terminal

```bash
CREATE USER warpdrive WITH PASSWORD 'warpdrive';
CREATE DATABASE warpdrivedb;
```

and make sure to set the correct username, password and database in warpdrive.conf.

# Warpfile

server:
  addr: 192.168.0.1:3000
cycles:
  production:
    build: react-native build $PLATFORM



# Rollback
you have to unlock those version that you don't want


# Android

if you getting an error in Android Studio

```
Unsupported method: AndroidProject.getPluginGeneration().
The version of Gradle you connect to does not support that method.
To resolve the problem you can change/upgrade the target version of Gradle you connect to.
Alternatively, you can ignore this exception and read other information from the model.
```

then you should do this, go to `File / Settings/ Build, Execution, Deployment / Instant Run.` and Uncheck `Enable Instant Run to hot swap code...`
