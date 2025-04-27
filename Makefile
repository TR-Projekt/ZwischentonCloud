# Makefile for zwischentoncloud

VERSION=development
DATE=$(shell date +"%d-%m-%Y-%H-%M")
REF=refs/tags/development
DEV_PATH_MAC=$(shell echo ~/Library/Containers/org.zwischenton.project)
export

build:
	go build -v -ldflags="-X 'github.com:TR-Projekt/ZwischentonCloud/server/status.ServerVersion=$(VERSION)' -X 'github.com:TR-Projekt/ZwischentonCloud/server/status.BuildTime=$(DATE)' -X 'github.com:TR-Projekt/ZwischentonCloud/server/status.GitRef=$(REF)'" -o zwischentoncloud main.go

install:
	mkdir -p $(DEV_PATH_MAC)/usr/local/bin
	mkdir -p $(DEV_PATH_MAC)/etc
	mkdir -p $(DEV_PATH_MAC)/var/log
	mkdir -p $(DEV_PATH_MAC)/usr/local/zwischentoncloud
	
	cp operation/local/ca.crt  $(DEV_PATH_MAC)/usr/local/zwischentoncloud/ca.crt
	cp operation/local/server.crt  $(DEV_PATH_MAC)/usr/local/zwischentoncloud/server.crt
	cp operation/local/server.key  $(DEV_PATH_MAC)/usr/local/zwischentoncloud/server.key
	cp operation/local/authentication.publickey.pem  $(DEV_PATH_MAC)/usr/local/zwischentoncloud/authentication.publickey.pem
	cp operation/local/authentication.privatekey.pem  $(DEV_PATH_MAC)/usr/local/zwischentoncloud/authentication.privatekey.pem
	cp zwischentoncloud $(DEV_PATH_MAC)/usr/local/bin/zwischentoncloud
	chmod +x $(DEV_PATH_MAC)/usr/local/bin/zwischentoncloud
	cp operation/local/config_template_dev.toml $(DEV_PATH_MAC)/etc/zwischentoncloud.conf

run:
	./zwischentoncloud --container="$(DEV_PATH_MAC)"

stop:
	killall zwischentoncloud

clean:
	rm -r zwischentoncloud