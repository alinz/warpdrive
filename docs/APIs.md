# REST APIs

list of all available rest apis which can be used for external tools.
currently, some of the apis is being used by CLI tool, `warp`.

## Private APIs

Private APIs are protected and require jwt. jwt can be obtain by calling `session` apis in public section

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

- [x] GET       /apps/:appId/cycles/:cycleId/releases/:releaseId
- [x] POST      /apps/:appId/cycles/:cycleId/releases
- [x] PUT       /apps/:appId/cycles/:cycleId/releases/:releaseId
- [x] DELETE    /apps/:appId/cycles/:cycleId/releases/:releaseId
- [x] POST      /apps/:appId/cycles/:cycleId/releases/:releaseId/bundles
- [x] GET       /apps/:appId/cycles/:cycleId/releases/:releaseId/bundles?name=
- [x] POST      /apps/:appId/cycles/:cycleId/releases/:releaseId/lock
- [x] DELETE    /apps/:appId/cycles/:cycleId/releases/:releaseId/lock

## Public APIs

- [x] GET       /session
- [x] POST      /session/start
- [x] GET       /session/end

below APIs are used in react-native side

- [x] GET       /apps/:appId/cycles/:cycleId/releases?platforn=&version=&note=
- [x] GET       /apps/:appId/cycles/:cycleId/version/:version/platform/:platform/latest
- [x] POST      /apps/:appId/cycles/:cycleId/version/:version/platform/:platform/download