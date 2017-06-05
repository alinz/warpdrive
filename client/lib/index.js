// @flow

import { NativeModules } from 'react-native'

const { WarpdriveManager } = NativeModules

type Release = {
  version: string,
  notes: string,
  at: string
}

export const isAnyUpdate = (cb: (release: ?Release) => void) => {
  typeof cb === 'function' && WarpdriveManager.isAnyUpdate((value: string) => {
    if (value === '') {
      cb(null)
    } else {
      cb(JSON.parse(value))
    }
  })
}

export const update = (cb: (err: ?Error) => void) => {
  typeof cb === 'function' && WarpdriveManager.update((err? string) => {
    cb(err ? new Error(err) : null)
  })
}

export const reload = () => {
  WarpdriveManager.reload()
}
