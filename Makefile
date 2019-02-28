all: deploy

deployDir = releases

define build
	# build scast
	GOOS=${1} GOARCH=${2} go build -o ${deployDir}/${1}/${3}/scast${4} -ldflags '-s -w'
	# build screencast
	GOOS=${1} GOARCH=${2} revel package github.com/MonkeyBuisness/screencast prod
	tar -xf screencast.tar.gz -C ./${deployDir}/${1}/${3}/
	rm -rf screencast.tar.gz
	# edit run files
	truncate -s -1 ./${deployDir}/${1}/${3}/run.sh
	echo " -q "'$$1'" -br "'$$2' >> ./${deployDir}/${1}/${3}/run.sh
	truncate -s -1 ./${deployDir}/${1}/${3}/run.bat
	echo " -q %1 -br %2" >> ./${deployDir}/${1}/${3}/run.bat
endef

deploy:
	# remove previous release if exist
	rm -rf ${deployDir}/*

	# build windows version 
	$(call build,windows,amd64,x64,.exe)
	$(call build,windows,386,x86,.exe)

	# build linux version
	$(call build,linux,amd64,x64,)
	$(call build,linux,386,x86,)

run:
	revel run github.com/MonkeyBuisness/screencast

install:
	sudo sh ./install.sh
