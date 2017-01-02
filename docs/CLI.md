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

#### version

Prints the version. It is usually used to make sure the cli and server side has the same version.

```bash
> warp version

Warp v1.0.1
```

#### bundle

it bundles react-native project for both ios and android and store them in `.bundles` folder inside your react-native project.
every time you execute this command, it will erase the `.bundles` folder and recreate it again. So you don't have to delete this
folder everytime. This command needs to be called before `publish` command is called for either of platforms.

it comes with one optional argument `-p` or `--platform`. if you don't provide the value, it will build and bundle for both platform.

```bash
> warp bundle -p android

android bundle started
android bundle finished
```

#### publish

> requires `setup` command being set before calling this command

`publish` pushes newly created bundles to warpdrive. it accepts 3 flags, version `-v`, platform `-p` and release note `-n`.
Only version flag is required. if platform flag not provided, it tries to push ios and android from `.bundles` folder which created from `bundle` command.

```bash
> warp publish -v 1.0.1-prod -p android -n "fixed couple of bugs"

published new version for android
```

#### setup

in order for warpdrive to work, a configuration file needs to be created and placed inside the both ios and android. This configuration tells the warpdrive
how and where it needs to get the updates. The name of the file is `WarpFile` and it automatically being placed inside `ios` and `android/app/src/main/assets`.
for ios, developer needs to import the WarpFile from ios folder into the xcode project. For android, there is no action required since anything inside `assets`
folder will be bundled.

setup command let's you view the current configuration, if a `WarpFile` already exists inside the project file

```bash
> warp setup -l

Server Address: localhost:8221

App's ID: 1
App's Name: app

Cycles:

ID: 1
Name: dev
Key: -----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA92Qh/p4xlotcviMNk0Kf
LTWFf1B0aCBSm4UbpxzyvrBeY6suG45ueXPEgoSRzY9e8sQiiFtgliEKU/19NyYf
7WgH8lmKcg9h7m0H0GJNglS5A/h9Wa3V4uKGPkQHfKxoqAagXLAnVZ3SrOkLHla6
rLJSaMwGXNAr9qa8mFAVBx2Hl6slDAkH+EW5AGDSl7ckdFreAivDMPKpyQhBqL78
VAsZLtAvK4JLVCxoB0LqvndiAlHSZomv8kYFeLOHtZtxQJbgaSnb9pfR+8btoMhg
kBx28V+aJA8r5r8A27YxLV5pcsy7cVEjMGo+dvyFOnWMIsV1lefx38SbSYDqYUkb
AwIDAQAB
-----END PUBLIC KEY-----
```

if you plan to add add a new cycle to your project, you can call `-c`
and if you want to completely change the configuration, you can call the command with `-r`
