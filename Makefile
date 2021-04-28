discord:
	go build -o bin/discord && ./bin/discord

danmu:
	cd cmd/danmu && go build -o ../../bin/danmu && ../../bin/danmu
