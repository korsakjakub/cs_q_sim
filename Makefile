binary_dir=cmd/
binary_name=qsim.out
binary_path=${binary_dir}${binary_name}
package=cli

export CONFIG_PATH=config/

.DEFAULT_GOAL := all
.PHONY: all build

all: build

build: ${binary_dir}${package}/main.go
	go mod download
	go build -o ${binary_path} ${binary_dir}${package}/*.go
	# run ${binary_path}

test:
	go test -v ./pkg/*** -cover

clean:
	go clean
	rm ${binary_path}
