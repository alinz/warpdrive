# Warp Drive

## Design

- launched app creates an agent user agent:agent
- agent user can create admin|member user
- agent can creates apps
- agent is responsible to release the app
- admin is responsible to create member and release beta releases
- member can only read either production or beta.
- member tends to be used as the actual mobile devices.


Mobile app should know 2 thing

  1. Cycle id
  2. public key


App
  Cycles
    Releases
      Bundles

App
  id name

Cycles
  id app_id name public_key private_key

Releases
  id cycle_id platform version note locked

Bundles
  id release_id type name hash

Users
  id name email password

Permissions
  id app_id user_id permission


permission value can be
  - member: they can only upload and add more members
  - agent: they can lock the release and add more agents/members
  - admin: root user

App
  - Releases
    - Production
      - member with production type
    - Beta
      - member with beta type
    - Alpha
      - member with alpha type


member can have different type
type of member can be created by api


POST    /session/start
GET     /session/end


#


# mobile side
GET       /check/:cycleId/releases/:currentVersion
POST      /download/:cycleId/releases/:version
  -> images
  -> js/main.bundle


POST      /users
PATCH     /users/:userId/permissions
GET       /users
PATCH     /users/:userId
DELETE    /users/:userId

POST      /apps
GET       /apps
PATCH     /apps/:appId

POST      /cycles
GET       /cycles
PATCH     /cycles/:cycleId
GET       /cycles/:cycleId/config
POST      /cycles/:cycleId/releases
PATCH     /cycles/:cycleId/releases/:releaseId
POST      /cycles/:cycleId/releases/:releaseId
PATCH     /cycles/:cycleId/releases/:releaseId/lock
