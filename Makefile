APP_SERVER_VERSION := 1
DEB_PACKAGE_NAME := aproxy
BINARY_PATH := cmd/aproxy/aproxy
DEB_BINARY_PATH := deb/usr/local/sbin/
gobuild = go build -ldflags "-X main.buildVersion=$1 -X 'main.buildDate=$(shell date +'%Y/%m/%d %H:%M:%S')' -X main.buildCommit=$(shell git rev-parse HEAD)" -v -o $2 $3


.PHONY: build
build:
	$(call gobuild,${APP_SERVER_VERSION}, ${BINARY_PATH}, "cmd/aproxy/main.go")

.PHONY: build_linux
build_linux:
	GOOS=linux GOARCH=amd64 $(call gobuild,${APP_SERVER_VERSION}, ${BINARY_PATH}, "cmd/aproxy/main.go")

.PHONY: build_windows
build_windows:
	GOOS=windows GOARCH=amd64 $(call gobuild,${APP_SERVER_VERSION}, "cmd/aproxy/aproxy.exe", "cmd/aproxy/main.go")

.PHONY: build_macos
build_macos:
	GOOS=darwin GOARCH=amd64 $(call gobuild,${APP_SERVER_VERSION}, "cmd/aproxy/aproxy.app", "cmd/aproxy/main.go")

.PHONY: compose
compose:build_linux
	docker-compose up --build --force-recreate --no-deps -d

.PHONY: decompose
decompose:
	docker-compose down

.PHONY: deb
deb:
	mkdir -p ${DEB_BINARY_PATH}
	mkdir -p deb/etc
	cp conf/aproxy.yml deb/etc
	mv ${BINARY_PATH} ${DEB_BINARY_PATH}
	cp -r DEBIAN deb/
	sed -i.orig "s/DEB_PACKAGE_NAME/${DEB_PACKAGE_NAME}/g" deb/DEBIAN/control
	sed -i.orig "s|DEB_PACKAGE_VERSION|${DEB_PACKAGE_VERSION}|g" deb/DEBIAN/control
	dpkg-deb -b deb ${DEB_PACKAGE_NAME}_${DEB_PACKAGE_VERSION}_amd64.deb

.PHONY: debclean
debclean:
	rm -rf deb/
	rm -f *.deb

.PHONY: pushdeb
pushdeb:
	ls|grep "${DEB_PACKAGE_NAME}_${DEB_PACKAGE_VERSION}_amd64.deb" | while read l; do curl -X POST -F file=@"$$l" --user ${APTUSER}:${APTPASS} ${APTURL}/files/${APT_DIR};done
	curl -X POST --user ${APTUSER}:${APTPASS} ${APTURL}/repos/${APT_LOCAL_REPO}/file/${APT_DIR}
	curl -X PUT -H "Content-Type:application/json" \
	--data '{"Sources":[{"Component":"non-free"}], "Architectures":["amd64"], "Distribution":"${DISTR}"}' \
	--user ${APTUSER}:${APTPASS} ${APTURL}/publish/${APT_REPO}/${DISTR}

.DEFAULT_GOAL := build
