#
# Variables

VERSION=1.0.0
LDFLAGS+=-X github.com/pressly/warpdrive/warpdrive.Version="$(VERSION)"

OS=darwin linux windows
ARCH=386 amd64
BUILD_PATH=../../bin/build

#
# Cleaning

clean:
	@rm -rf ./bin
	@mkdir -p ./bin/data
	@mkdir -p ./bin/temp
	@mkdir -p ./bin/bundles
	@mkdir -p ./bin/build


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

dev: kill
	@(export CONFIG=$$PWD/etc/warpdrive.conf && \
		cd ./cmd/warpdrive && \
		fresh -c ../../etc/fresh-runner.conf -w=../..)


#
# Building

build-all: clean
	@(cd cmd/warpdrive; 																												\
	for GOOS in $(OS); do 																											\
		for GOARCH in $(ARCH); do 																								\
			echo "building $$GOOS $$GOARCH ..."; 																		\
			export GOGC=off;																												\
			export GOOS=$$GOOS; 																										\
			export GOARCH=$$GOARCH; 																								\
			go build -ldflags "$(LDFLAGS)" 																					\
			-o $(BUILD_PATH)/warpdrive-$$GOOS-$$GOARCH; 														\
		done 																																			\
	done)

build-docker: clean
	@(cd cmd/warpdrive;																													\
	export export GOGC=off; export GOOS=linux; export GOARCH=amd64;	 						\
	go build -ldflags "$(LDFLAGS)" -o $(BUILD_PATH)/warpdrive-linux-amd64;			\
	cd ../..;																																		\
	docker build -t warpdrive .)
