# WARPDRIVE

is a server module of the whole process. other services such as CodePush or AppHub are host your code on their end. WarpDrive is a full package 
which lets you in charge of your server. With the help of Docker and cheap servers you can easily manage your internal server for any apps and scale it as you
go.

The content of bundle transmited between WarpDrive and client sides are being g-ziped and encrypted with AES-256 by default. So you can run your server as non secure HTTP.

> we don't suggest that, and we are encourging you to use HTTPS over HTTP as much as you can. You can always get the free certificate from `letsencrypt.org`

## Design

Warpdrive is built with `Golang`. The entire rest apis are built with our popular open source `router` called [chi](https://github.com/pressly/chi) and the ORM behind is another huge open source sponser by `Pressly`, called [Upper DB](upper.io).
The entire server side can be compiled and executed as a single binary file on any platform, Mac, Windows or Linux.

The Warpdrive contains 3 main concepts

- Apps
- Cycles
- Releases

So any react-native project is considered as an App. Each App can have multiple cycles. For example, `Dev`, `Alpha`, `Beta`, `RC` and `Prod` can be considered as cycles. Each app can be configured to use all or one of the cycles.
so now you can publish releases to each cycles as you wish and the proper configuration on mobile devices will grab the update and restart the app for you. 

## Usage

The design of warpdrive went through a lot of reversions to simplify the execution. There are 2 things that you need,

first, you need to grab the warpdrive executable suites to your target platform,
second, you need a configuration file which tells the warpdrive how to connect to database and some security stuff. you can find a sample configuration in `etc/warpdrive.sample.conf`

## Configuration

let's take a look at the sample configuration

```
[server]
addr              = "0.0.0.0:8221"
data_dir          = "./bin/tmp/warpdrive"

[jwt]
secret_key        = "secret"
max_age           = 3600
path              = "/"
domain            = ""
secure            = false

[db]
database          = "warpdrivedb"
hosts             = "127.0.0.1"
username          = "warpdrive"
password          = "warpdrive"

[security]
key_size          = 2048

[file_upload]
file_max_size     = 8388608
```

Usually you need to change couple of them, for example, you need to change the `data_dir` to reflect the directory patch which it needs to save bundles, and `secert_key` for your encrypt and decrypt JWT tokens and `db` section to connects to a postgres database.
if you need to increase the security, you can always increase the `key_size`. But beawre of performance penalty on mobile side if you increase the key too hight.