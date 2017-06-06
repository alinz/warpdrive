# Warpdrive

This is a single repo containg all the pieces required to build a Wardrive pipeline. It has been rebuilt from ground up for 2 main reasons

1 - Performance
  - using GRPC to transfer data safe and fast, with HTTP2 by default
  - using Golang to unlock stream power for both android and ios in one language and compile to multiple os such as darwin, windows and linux

2 - Simplicity
  - using Golang allow us to maintain one code base and write less code. The entire server side is around 500 loc and the client which does decryption, download, file management is around 300 loc.

# Usage

The warpdrive project consitsts of 3 main components,

  1 - Warpdrive Server
  2 - Warpdrive Client
  3 - Warpdrive Cli aka `warp`

### Warpdrive Server

This is the Server component of the pipeline. Behind the scene we are using `storm` which uses `boltdb` under the hood to power our database system. Because of the simplicity reason we decided to go with key/valu store system.
For security reasons, we are using certificate approach as a pose to tradtional username/password. This makes it easy, secure on none HTTPS server and compatiable with GRPC. Also it removes the headache of maintain users list.

### Warpdrive Client

### Warpdrive Cli
