test:
	APP_ENV=test go test --race ./src
	APP_ENV=test go test --race ./src/api
	APP_ENV=test go test --race ./src/auth
	APP_ENV=test go test --race ./src/ctl
	APP_ENV=test go test --race ./src/model
	APP_ENV=test go test --race ./src/logger
	APP_ENV=test go test --race ./src/schema
	APP_ENV=test go test --race ./src/service
	APP_ENV=test go test --race ./src/util

cover:
	rm -f *.coverprofile
	APP_ENV=test go test -coverprofile=src.coverprofile ./src
	APP_ENV=test go test -coverprofile=api.coverprofile ./src/api
	APP_ENV=test go test -coverprofile=auth.coverprofile ./src/auth
	APP_ENV=test go test -coverprofile=api.coverprofile ./src/ctl
	APP_ENV=test go test -coverprofile=dao.coverprofile ./src/model
	APP_ENV=test go test -coverprofile=logger.coverprofile ./src/logger
	APP_ENV=test go test -coverprofile=schema.coverprofile ./src/schema
	APP_ENV=test go test -coverprofile=service.coverprofile ./src/service
	APP_ENV=test go test -coverprofile=util.coverprofile ./src/util
	gover
	go tool cover -html=gover.coverprofile
	rm -f *.coverprofile

GO=$(shell which go)

assets:
	go-bindata -ignore=\\.DS_Store -o ./src/bindata.go -pkg src -prefix web/dist/ web/dist/...
clean:
	go-bindata -ignore=\\.* -o ./src/bindata.go -pkg src web/dist/...
web:
	cd web && rm -rf dist && npm run deploy:prod && cd -
dev: web assets
	go run cmd/kpass/kpass.go -dev

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build -a -installsuffix cgo -o dist/kpass_linux ./cmd/kpass
build-darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GO) build -a -installsuffix cgo -o dist/kpass ./cmd/kpass
build-win:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GO) build -a -installsuffix cgo -o dist/kpass.exe ./cmd/kpass
build: web assets build-darwin build-linux build-win clean

swagger:
	swaggo -s ./src/swagger.go
	staticgo

GIT_REVISION_SHORTCODE := $(shell git rev-parse --short HEAD)
GIT_REVISION := $(shell git describe --abbrev=0 --tags --exact-match 2> /dev/null || git rev-parse --short HEAD)
GIT_REVISION_DATE := $(shell git show -s --format=%ci $(GIT_REVISION_SHORTCODE))
REVISION_DATE := $(shell date -u -j -f "%F %T %z" "$(GIT_REVISION_DATE)" +"%Y%m%d.%H%M%S" 2>/dev/null || date -u -d "$(GIT_REVISION_DATE)" +"%Y%m%d.%H%M%S")
BUILD_DATE := $(shell date -u +%Y%m%d.%H%M%S)

require-version:
	@if [[ -z "$$VERSION" ]]; then echo "VERSION environment value is required."; exit 1; fi

package-darwin: require-version build-darwin
	@echo "Generating distribution package for darwin/amd64..." && \
	if [[ "$$(uname -s)" == "Darwin" ]]; then \
		INSTALLER_RESOURCES="installer_resources/darwin" && \
		mkdir -p installer && \
		rm -rf installer/Kpass.app && \
		cp -r $$INSTALLER_RESOURCES/Kpass.app_template installer/KPass.app && \
		mkdir installer/KPass.app/Contents/MacOS && \
		cp dist/kpass installer/KPass.app/Contents/MacOS/kpass && \
		cat installer/KPass.app/Contents/MacOS/kpass | bzip2 > dist/update_darwin_amd64.bz2 && \
		ls -l dist/kpass dist/update_darwin_amd64.bz2 && \
		rm -rf installer/KPass.dmg && \
		sed "s/__VERSION__/$$VERSION/g" $$INSTALLER_RESOURCES/kpass.dmg.json > $$INSTALLER_RESOURCES/kpass_versioned.dmg.json && \
		appdmg --quiet $$INSTALLER_RESOURCES/kpass_versioned.dmg.json installer/KPass.dmg && \
		mv installer/KPass.dmg installer/KPass.dmg.zlib && \
		hdiutil convert -quiet -format UDBZ -o installer/KPass.dmg installer/KPass.dmg.zlib && \
		rm installer/KPass.dmg.zlib && \
		rm $$INSTALLER_RESOURCES/kpass_versioned.dmg.json; \
	else \
		echo "-> Skipped: Can not generate a package on a non-OSX host."; \
	fi;

.PHONY: web assets test build cover clean
