# Warpdrive

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