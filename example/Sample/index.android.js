/**
 * Sample React Native App
 * https://github.com/facebook/react-native
 * @flow
 */

import React, { Component } from 'react'
import { AppRegistry, StyleSheet, Text, View, NativeModules, Alert } from 'react-native'

const { WarpdriveManager } = NativeModules

type Release = {
  version: string,
  notes: string,
  at: string
}

export const isAnyUpdate = (cb: (release: ?Release) => void) => {
  typeof cb === 'function' &&
    WarpdriveManager.isAnyUpdate((value: string) => {
      if (value === '') {
        cb(null)
      } else {
        cb(JSON.parse(value))
      }
    })
}

export const update = (cb: (err: ?Error) => void) => {
  typeof cb === 'function' &&
    WarpdriveManager.update((err: ?string) => {
      cb(err ? new Error(err) : null)
    })
}

export const reload = () => {
  setTimeout(() => {
    WarpdriveManager.reload()
  })
}

export default class Sample extends Component {
  onCheck = () => {
    isAnyUpdate(release => {
      if (release) {
        Alert.alert(
          'Info',
          `an update ${release.version} is available, do you want to update?`,
          [
            {
              text: 'Yes',
              onPress: () => {
                update(err => {
                  if (err) {
                    Alert.alert('Error', err)
                    return
                  }

                  Alert.alert(
                    'Info',
                    'update is completed',
                    [{ text: 'Reload', onPress: () => reload() }, { text: 'Not Yet' }],
                    { cancelable: false }
                  )
                })
              }
            },
            { text: 'No' }
          ],
          { cancelable: false }
        )
      }
    })
  }

  render() {
    return (
      <View style={styles.container}>
        <Text style={styles.instructions}>
          To get started, edit index.android.js
        </Text>
        <Text style={styles.instructions}>
          Double tap R on your keyboard to reload,{'\n'}
          Shake or press menu button for dev menu
        </Text>
      </View>
    )
  }
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    backgroundColor: '#F5FCFF'
  },
  welcome: {
    fontSize: 20,
    textAlign: 'center',
    margin: 10
  },
  instructions: {
    textAlign: 'center',
    color: '#333333',
    marginBottom: 5
  }
})

AppRegistry.registerComponent('Sample', () => Sample)
