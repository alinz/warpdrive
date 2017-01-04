/**
 * Sample React Native App
 * https://github.com/facebook/react-native
 * @flow
 */

import React, { Component } from 'react';
import {
  AppRegistry,
  StyleSheet,
  Text,
  View
} from 'react-native';

import Warpify from 'react-native-warpdrive'

export default class Sample1 extends Component {

  async componentDidMount() {
    try {
      const cycles = await Warpify.cycles()
      console.log(cycles)
      const localVersions = await Warpify.local(4)
      console.log(localVersions)
      const remoteVersions = await Warpify.remote(4)
      console.log(remoteVersions)

      const latestVersion = await Warpify.latest(4)
      console.log(latestVersion)
      
      await Warpify.download(4, latestVersion.soft.version)

      await Warpify.reload(4, latestVersion.soft.version)

    } catch (e) {
      console.log(e)
    }
  }

  render() {
    return (
      <View style={styles.container}>
        <Text style={styles.welcome}>
          Welcome to React Native!
        </Text>
        <Text style={styles.instructions}>
          To get started, edit index.ios.js
        </Text>
        <Text style={styles.instructions}>
          Press Cmd+R to reload,{'\n'}
          Cmd+D or shake for dev menu
        </Text>
      </View>
    );
  }
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    backgroundColor: '#F5FCFF',
  },
  welcome: {
    fontSize: 20,
    textAlign: 'center',
    margin: 10,
  },
  instructions: {
    textAlign: 'center',
    color: '#333333',
    marginBottom: 5,
  },
});

AppRegistry.registerComponent('Sample1', () => Sample1);
