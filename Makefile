# Makefile for zwischentoncloud

VERSION=development
DATE=$(shell date +"%d-%m-%Y-%H-%M")
REF=refs/tags/development
export

build:
	go build -v -ldflags="-X 'github.com:TR-Projekt/ZwischentonCloud/server/status.ServerVersion=$(VERSION)' -X 'github.com:TR-Projekt/ZwischentonCloud/server/status.BuildTime=$(DATE)' -X 'github.com:TR-Projekt/ZwischentonCloud/server/status.GitRef=$(REF)'" -o zwischentoncloud main.go

install:
	cp zwischentoncloud /usr/local/bin/zwischentoncloud
	cp config_template.toml /etc/zwischentoncloud.conf
	cp operation/service_template.service /etc/systemd/system/zwischentoncloud.service

update:
	systemctl stop zwischentoncloud
	cp zwischentoncloud /usr/local/bin/zwischentoncloud
	systemctl start zwischentoncloud

uninstall:
	systemctl stop zwischentoncloud
	rm /usr/local/bin/zwischentoncloud
	rm /etc/zwischentoncloud.conf
	rm /etc/systemd/system/zwischentoncloud.service

run:
	./zwischentoncloud --debug

stop:
	killall zwischentoncloud

clean:
	rm -r zwischentoncloud