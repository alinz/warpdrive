LDFLAGS+=-X github.com/pressly/warpdrive.VERSION=$$(git describe --tags --abbrev=0 --always)
LDFLAGS+=-X github.com/pressly/warpdrive.LONGVERSION=$$(git describe --tags --long --always --dirty)

##
## Tools
##

tools:
	go get -u github.com/pressly/fresh
	go get -u github.com/pressly/goose
	go get -u github.com/pressly/sup
	go get -u github.com/kardianos/govendor
	go get -u github.com/jstemmer/gotags
	go get -u golang.org/x/tools/cmd/goimports
	go get -u golang.org/x/mobile/cmd/gomobile

##
## Dependency mgmt
## only use this if you want to resync new vendor package
## if you don't need to build for appengine, go to vendor/vendor.json
## and add `appengine` into `ignore` field. for example: `"ignore": "test appengine",` 
##


vendor-rebuild:
	@rm -rf vendor
	@govendor init

#@govendor add +e
#@govendor update +external 

# install:

##
## Building Server
##

clean-server:
	@rm -rf ./bin/warpdrive-server
	@mkdir -p ./bin

build-server: clean-server
	GOGC=off go build -i -gcflags="-e" -ldflags "$(LDFLAGS)" -o ./bin/warpdrive-server ./cmd/server

##
## Building cli
##
clean-cli:
	@rm -rf ./bin/warp
	@mkdir -p ./bin

build-cli: clean-cli
	GOGC=off go build -i -gcflags="-e" -ldflags "$(LDFLAGS)" -o ./bin/warp ./cmd/warp

##
## Building WarpDrive's Clients, ios and android
##

clean-ios:
	@rm -rf ./client/ios/Warpdrive.framework

build-ios: clean-ios
	@cd ./cmd/warpdrive && gomobile bind -target=ios . && mv -f Warpdrive.framework ../../client/ios

clean-android:
	@rm -rf client/android/warpdrive.aar

build-android: clean-android
	@cd ./cmd/warpdrive && gomobile bind -target=android . && mv -f warpdrive.aar ../../client/android

build-clients: build-ios build-android

##
## Building everything
##

build: build-server build-cli build-clients

##
## Database
##

db-create:
	@goose up

db-destroy:
	@goose down

db-reset: db-destroy db-create

##
## Development
##

build-dev-folder:
	@rm -rf ./bin
	@mkdir -p ./bin/tmp/warpdrive

run: build-dev-folder
	fresh -c ./etc/fresh-runner.conf -p ./cmd/server -r '-config=./etc/warpdrive.conf' -o ./bin/warpdrive-server

kill:
	@lsof -t -i:8221 | xargs kill

clean-ios-example:
	@rm -rf ./client/examples/Sample1/node_modules/react-native-warpdrive/ios/Warpdrive.framework

build-ios-example: clean-ios-example
	@cd ./cmd/client && gomobile bind -target=ios . && mv -f Warpdrive.framework ../../client/examples/Sample1/node_modules/react-native-warpdrive/ios

clean-cli-dev:
	@rm -rf ./client/examples/Sample1/warp

build-cli-dev: clean-cli-dev
	GOGC=off go build -i -gcflags="-e" -ldflags "$(LDFLAGS)" -o ./client/examples/Sample1/warp ./cmd/warp
