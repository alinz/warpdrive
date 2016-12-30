import { NativeModules } from 'react-native'

const { Warpify } = NativeModules

// cycles returns an array of available cycles described in 
// WarpFile bundles inside the app.
// e.g. [{ id: 1, name: "dev" }, ...]
export const cycles = () => {
  return new Promise((resolve, reject) => {

  })
}

// remote asks the warpdrive server for all release versions available
// the order of array is based on latest first.
// e.g. [{ version: "1.0.10", note: "fixed some bugs" }, ...]
export const remote = (cycleId) => {
  return new Promise((relove, reject) => {

  })
}

// local checks the downloaded versions and returns the available versions
// e.g. [{ version: "1.0.9" }, ...]
// if you want more information, you have to ask `info` method to get more information
// about that particular version.
export const local = (cycleId) => {
  return new Promise((resolve, reject) => {

  })
}

// info asks the server for basic information about each version.
// e.g. { version: "1.0.10", note: "fixed some bugs" }
export const info = (cycleId, version) => {
  return new Promise((resolve, reject) => {

  })
}

// latest gets the cycleId and returns soft and hard versions if available
// e.g. { "soft": {}, "hard": {} }
export const latest = (cycleId) => {
  return new Promise((resolve, reject) => {

  })
}

// downlaod requests for download. if the bundle already downlaoded, it won't download the content
export const download = (cycleId, version) => {
  return new Promise((resolve, reject) => {

  })
}

// reload reloads the app to use the specific version, reload won't wont take affect if
// forceUpdate is true in native side. becuase as soon as new update pushes to defaultCycle, 
// all of the changes will be reverted.
export const reload = (cycleId, version) => {
  return new Promise((resolve, reject) => {

  })
}