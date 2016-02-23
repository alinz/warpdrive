# Cycles API


## Create a cycle for an app

POST  /apps/:appId/cycles
GET   /apps/:appId/cycles
PATCH /apps/:appId/cycles/:cycleId

## Downloads the config file
GET   /apps/:appId/cycles/:cycleId/config

## Creates a new release
POST  /apps/:appId/cycles/:cycleId/releases?platform=&version=&notes=

## Changes of release fields
PATCH /apps/:appId/cycles/:cycleId/releases/:releaseId

## Get all releases
GET  /apps/:appId/cycles/:cycleId/releases?currentVersion=1.2.0

## Locked the version so client can download
PATCH /apps/:appId/cycles/:cycleId/releases/:releaseId/lock

## Check if any download is available
GET   /apps/:appId/cycles/:cycleId/releases/version/:version

{
  "latest_version": "1.2.3"
}

## Download bundles

POST /apps/:appId/cycles/:cycleId/releases/version/:version/download

##

## Get all cycles for specific app

## Update cycles for specific app

##

POST      /cycles
GET       /cycles
PATCH     /cycles/:cycleId
GET       /cycles/:cycleId/config
POST      /cycles/:cycleId/releases
PATCH     /cycles/:cycleId/releases/:releaseId
POST      /cycles/:cycleId/releases/:releaseId
PATCH     /cycles/:cycleId/releases/:releaseId/lock
