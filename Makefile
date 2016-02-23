#
# Variables

VERSION=1.0.0
LDFLAGS+=-X github.com/pressly/warpdrive/warpdrive.Version="$(VERSION)"


#
# Cleaning

clean:
	@rm -rf ./bin
	@mkdir -p ./bin/data
	@mkdir -p ./bin/temp
	@mkdir -p ./bin/bundles


#
# Vendoring

update-deps:
	@echo "Updating Glockfile with package versions from GOPATH..."
	@rm -rf ./vendor
	@glock save github.com/pressly/warpdrive
	@$(MAKE) vendor

vendor:
	@echo "Syncing dependencies into vendor directory..."
	@rm -rf ./vendor
	@gv < Glockfile


#
# Killing

kill-fresh:
	@ps -ef | grep 'f[r]esh' | awk '{print $$2}' | xargs kill

kill-by-port:
	@lsof -t -i:8221 | xargs kill

kill: kill-fresh kill-by-port


#
# Development

db-reset:
	@(cd ./script && bash ./db-reset.sh)

dev: kill clean
	@(export CONFIG=$$PWD/etc/warpdrive.conf && \
		cd ./cmd/warpdrive && \
		fresh -c ../../etc/fresh-runner.conf -w=../..)


#
# Building

build:
	GOGC=off go build -i -ldflags "$(LDFLAGS)" -o ./bin/warpdrive ./cmd/warpdrive
