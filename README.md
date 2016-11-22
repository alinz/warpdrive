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

# APIs

- [x] GET       /session
- [x] POST      /session/start
- [x] GET       /session/end

- [x] GET       /users?name=&email=
- [x] GET       /users/:userId
- [X] POST      /users
- [X] PUT       /users

- [x] GET       /apps?name=
- [x] GET       /apps/:appId
- [x] POST      /apps
- [X] PUT       /apps/:appId

- [x] GET       /apps/:appId/users?name=&email=
- [x] POST      /apps/:appId/users/:userId
- [x] DELETE    /apps/:appId/users/:userId

- [x] GET       /apps/:appId/cycles?name=
- [x] GET       /apps/:appId/cycles/:cycleId
- [x] POST      /apps/:appId/cycles
- [X] GET       /apps/:appId/cycles/:cycleId/key
- [X] PUT       /apps/:appId/cycles/:cycleId
- [X] DELETE    /apps/:appId/cycles/:cycleId

- [x] GET       /apps/:appId/cycles/:cycleId/releases?platforn=&version=&note=
- [x] GET       /apps/:appId/cycles/:cycleId/releases/:releaseId
- [x] POST      /apps/:appId/cycles/:cycleId/releases
- [x] PUT       /apps/:appId/cycles/:cycleId/releases/:releaseId
- [x] DELETE    /apps/:appId/cycles/:cycleId/releases/:releaseId
- [x] POST      /apps/:appId/cycles/:cycleId/releases/:releaseId/bundles
- [ ] GET       /apps/:appId/cycles/:cycleId/releases/:releaseId/bundles?name=
- [x] POST      /apps/:appId/cycles/:cycleId/releases/:releaseId/lock
- [x] DELETE    /apps/:appId/cycles/:cycleId/releases/:releaseId/lock

## the following apis call in clinet native, also they are public

- [ ] GET       /apps/:appId/cycles/:cycleId/releases/latest/:version
- [ ] POST      /apps/:appId/cycles/:cycleId/releases/:releaseId/download

## this following api calls for any auditing

- [ ] GET       /apps/:appId/logs

# Rollback
you have to unlock those version that you don't want