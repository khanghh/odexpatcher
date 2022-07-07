odexpatcher:
	go build -ldflags "-w -s \
		-X 'main.CommitHash=$(shell git rev-list -1 HEAD)' \
		-X 'main.Version=$(shell git describe --tags)' \
		-X 'main.BuiltTime=$(shell date)'" \
		-o bin/$@ cmd/odexpatcher/main.go \
		cmd/odexpatcher/version.go \
		cmd/odexpatcher/patch_cmd.go \
		cmd/odexpatcher/info_cmd.go 

all: odexpatcher
