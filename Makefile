install:
	go install

test:
	go test -v

dependencies:
	# https://cloud.google.com/functions/docs/writing/specifying-dependencies-go#using_go_modules
	export GO111MODULE=on
	go mod init
	go mod tidy
