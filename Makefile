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

##
## Dependency mgmt
##
vendor-list:
	@govendor list

vendor-update:
	@govendor update +external

vendor-sync:
	@govendor sync +external

install: vendor-list vendor-update vendor-sync

##
## Building
##
clean:
	@rm -rf ./bin
	@mkdir -p ./bin

build: clean
	GOGC=off go build -i -gcflags="-e" -ldflags "$(LDFLAGS)" -o ./bin/warpdrive ./cmd/warpdrive

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
	fresh -c ./etc/fresh-runner.conf -p ./cmd/warpdrive -r '-config=./etc/warpdrive.conf' -o ./bin/warpdrive

kill:
	@lsof -t -i:8221 | xargs kill
