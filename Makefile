LDFLAGS+=-X github.com/pressly/warpdrive/warpdrive.VERSION=$$(git describe --tags --abbrev=0 --always)
LDFLAGS+=-X github.com/pressly/warpdrive/warpdrive.LONGVERSION=$$(git describe --tags --long --always --dirty)

# to create a root certificate. secret key must be kept secret and
# protected with strong password. Make sure to backup `ca-*.*` files.
#
# e.g. make create-root-certificate password=1234
create-root-certificate:
	@certstrap --depot-path "cert" init --passphrase "$(password)" --common-name "ca-query" &&					\
	certstrap --depot-path "cert" init --passphrase "$(password)" --common-name "ca-command";

# before you deploy the warpdrive, you have to call this to generate service certificates
create-service-certificates:
	@certstrap --depot-path "cert" request-cert --passphrase "" --common-name "query" &&								\
	certstrap --depot-path "cert" sign "query" --CA "ca-query" &&																				\
	certstrap --depot-path "cert" request-cert --passphrase "" --common-name "command" &&								\
	certstrap --depot-path "cert" sign "command" --CA "ca-command";

create-cli-certificate:
	@certstrap --depot-path "cert" request-cert --passphrase "" --common-name "cli" &&									\
	certstrap --depot-path "cert" sign "cli" --CA "ca-command";

# follwing command will generate certificate for mobile which let's them connect to
# query service
create-device-certificate:
	@certstrap --depot-path "cert" request-cert --passphrase "" --common-name "device" &&								\
	certstrap --depot-path "cert" sign "device" --CA "ca-query";

# to compile the proto file to go code and then run
# go-inject to inject some code to work properly with storm package
compile-protobuf:
	@protoc --go_out=plugins=grpc:. ./proto/*.proto;													  										 		\
	protoc-go-inject-tag -input=./proto/warpdrive.pb.go;

compile-server-linux: compile-protobuf
	@export GOGC=off;																																										\
	export GOOS=linux;																																									\
	export GOARCH=amd64;																																								\
	go build -ldflags "$(LDFLAGS)" 																																			\
	-o ./bin/server/warpdrive-server ./cmd/server;	

build-docker: compile-server-linux
	docker build -t warpdrive .

build-cli: compile-protobuf
	@export GOGC=off;																																										\
	export GOOS=darwin;																																									\
	export GOARCH=amd64;																																								\
	go build -ldflags "$(LDFLAGS)" 																																			\
	-o ./bin/cli/warp ./cmd/warp && mv -f ./bin/cli/warp ./example/Sample/warp; 

clean-ios:
	@rm -rf ./client/ios/Warpdrive.framework;																														\
	mkdir -p ./client/ios;

build-ios: clean-ios
	@cd ./cmd/warpdrive && gomobile bind -target=ios -ldflags="-s -w" . && 															\
	mv -f Warpdrive.framework ../../client/ios;

clean-android:
	@rm -rf client/android/lib/warpdrive.aar

build-android: clean-android
	@cd ./cmd/warpdrive && gomobile bind -target=android -ldflags="-s -w" . && 													\
	mkdir -p ../../client/android/lib &&																																\
	mv -f warpdrive.aar ../../client/android/lib;

build-clients: build-ios build-android

server-cleanup:
	@docker-compose down; docker rmi warpdrive; rm -rf ./tmp; mkdir -p ./tmp/bundles; mkdir -p ./tmp/db;
