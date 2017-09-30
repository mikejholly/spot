
build:
	go build -o likespot github.com/mikejholly/likespot/cmd

install: build
	 sudo mv likespot /usr/local/bin
