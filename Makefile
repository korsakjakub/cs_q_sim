binary_dir=cmd/
binary_name=qsim.out
binary_path=${binary_dir}${binary_name}
internal_path=./internal
package=cli
dockerfile_path=build/Dockerfile
docker_name=cs_q_sim:1.0

export CONFIG_PATH=config/

.DEFAULT_GOAL := all
.PHONY: all build

all: build

build: ${binary_dir}${package}/main.go
	go mod download
	go build -o ${binary_path} ${binary_dir}${package}/*.go
	# run ${binary_path}

test:
	go test -v ${internal_path}/*** -cover

clean:
	go clean
	rm ${binary_path}

docker-build: ${dockerfile_path}
	docker build -t ${docker_name} -f ${dockerfile_path} .

docker-run:
	docker run --rm -it ${docker_name}
	
docker-clean:
	docker rmi ${docker_name}
