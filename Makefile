all: src/client/client src/server/server

src/client/client:
	cd src/client && go build

src/server/server:
	cd src/server && go build

clean:
	rm -f src/client/client
	rm -f src/server/server
