# Warpdrive

this is a single repo contains all the required and missing pieces for setup a system for updating react-native apps just like `Code-Push` or `AppHub`. 

So Why a new system? 

`Code-Push` is missing server component, `Walmart lab` is open sourced the server side but you always have to make sure the new release of `Code-Push` is compatiable by server side. Also for now `Code-Push` is free and we hope it stays that way but who knows what happens next in the future.

`AppHub` has the same probelm as `Code-Push` and also it cost money.

`Warpdrive` solves all of the above plus it uses `golang` under the hood. `golang` is a mature language and has been used for server side for quite a while and as soon as `gomobile` releases we decided to power `warpdrive` for both `android` and `ios`.

it has couple of major benefits:

- it has a powerful streaming capability which enables us to do zip and unzip over network on the fly. (good luck trying to implement it on both `Java` and `Objective-c`)
- low memory foot print
- one code base for both android and ios
- share code between server and client

<p align="center">
  We think Gopher can have a react-native tattoo as well!
</p>
<p align="center">
  <img width="200" src="https://raw.githubusercontent.com/pressly/warpdrive/master/docs/assets/gopher-tattoo.jpg" />
</p>

Cheers,
Pressly - Team

