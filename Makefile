version?=v0.0.1
name=kube-namespace-deleter
name_local=kube-namespace-deleter
pkg=github.com/plantoncloud/kube-namespace-deleter
build_dir=build
LDFLAGS=-ldflags "-X ${pkg}/internal/cli/version.Version=${version}"
build_cmd=go build -v ${LDFLAGS}

.PHONY: deps
deps:
	go mod download
	go mod tidy

.PHONY: build_darwin
build_darwin: vet
	GOOS=darwin ${build_cmd} -o ${build_dir}/${name}-darwin .

.PHONY: build
build: ${build_dir}/${name}

${build_dir}/${name}: deps vet
	GOOS=darwin ${build_cmd} -o ${build_dir}/${name}-darwin .
	GOOS=darwin GOARCH=amd64 ${build_cmd} -o ${build_dir}/${name}-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 ${build_cmd} -o ${build_dir}/${name}-darwin-arm64 .
	GOOS=linux GOARCH=amd64 ${build_cmd} -o ${build_dir}/${name}-linux .
	openssl dgst -sha256 ${build_dir}/${name}-darwin
	openssl dgst -sha256 ${build_dir}/${name}-linux

.PHONY: test
test:
	go test -race -v -count=1 ./...

.PHONY: run
run: build
	${build_dir}/${name}

.PHONY: vet
vet:
	go vet ./...

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: clean
clean:
	rm -rf ${build_dir}

checksum_darwin:
	@openssl dgst -sha256 ${build_dir}/${name}-darwin

checksum_linux:
	openssl dgst -sha256 ${build_dir}/${name}-linux

checksum: checksum_darwin checksum_linux

local: build_darwin
	rm -f ${HOME}/bin/${name_local}
	cp ./${build_dir}/${name}-darwin ${HOME}/bin/${name_local}
	chmod +x ${HOME}/bin/${name_local}
