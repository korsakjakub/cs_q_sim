binary_dir=cmd/
binary_name=qsim.out
binary_path=${binary_dir}${binary_name}
internal_path=internal/quantum_simulator
package=quantum_simulator

build: ${binary_dir}${package}/main.go
	go build -o ${binary_path} ${binary_dir}${package}/main.go

test:
	go test -v ${internal_path}

run: ${binary_path}
	${binary_path}

clean:
	go clean
	rm ${binary_path}