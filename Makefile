#
# Variables

OS=darwin linux windows
ARCH=386 amd64

LDFLAGS+=-X github.com/pressly/warpdrive/warpdrive.Version=$$(scripts/version.sh --long)

#
# Cleaning

clean:
	@rm -rf ./bin


#
# Vendoring

tools:
	go get -u github.com/kardianos/govendor

vendor-list:
	@govendor list

vendor-update:
	@govendor update +external


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
	@(cd ./scripts && bash ./db-reset.sh)

dev: kill
	@(export CONFIG=$$PWD/etc/warpdrive.conf && \
		cd ./cmd/warpdrive && \
		fresh -c ../../etc/fresh-runner.conf -w=../..)


#
# Building

build-all: clean
	for GOOS in $(OS); do 																											\
		for GOARCH in $(ARCH); do 																								\
			echo "building $$GOOS $$GOARCH ..."; 																		\
			export GOGC=off;																												\
			export GOOS=$$GOOS; 																										\
			export GOARCH=$$GOARCH; 																								\
			go build -ldflags "$(LDFLAGS)" 																					\
			-o ./bin/warpdrive-$$GOOS-$$GOARCH ./cmd/warpdrive;       		          \
		done 																																			\
	done)

build: clean
	export GOGC=off;                                       	 		        				\
	go build -ldflags "$(LDFLAGS)" -o ./bin/warpdrive ./cmd/warpdrive;
