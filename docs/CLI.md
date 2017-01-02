# CLI

cli is a stand-alone `Golang` project which wraps http-client to talk to `warpdrive`'s server.

## Requirement

cli tool is call `warp` and must be called in `react-native`'s project folder. 

## Installation

`warp` file can be downloaded from release section and the file is an stand-alone executable file.
It is better to placed in a folder and refer to it in your PATH file so you can access the warp anywhere 
in your system.

if you prefer to compile and build the cli executable, you can clone the project and make sure you have latest
`Golang` installed, and simply call `make build-cli`. if everything goes well, you would be able to see a warp executable file 
sitting in bin folder of your clone project.

## Usage

`warp` must be executed inside `react-native`'s project. by providing `-h` or `--help` you can access to
details descriptions of all commands.

`warp` cli contains 5 main commands

##### version

Prints the version. It is usually used to make sure the cli and server side has the same version.

```bash
> warp version

Warp v1.0.1
```

##### bundle

it bundles react-native project for both ios and android and store them in `.bundles` folder inside your react-native project.
every time you execute this command, it will erase the `.bundles` folder and recreate it again. So you don't have to delete this
folder everytime. This command needs to be called before `publish` command is called for either of platforms.

it comes with one optional argument `-p` or `--platform`. if you don't provide the value, it will build and bundle for both platform.

```bash
> warp bundle -p android

android bundle started
android bundle finished
```

##### publish


