CC=go

all: client.exe server.exe
client.exe: client/client.go common/config/config.go common/reply/reply.go common/logger/logger.go
	$(CC) build -o client.exe client/client.go
server.exe: server/server.go common/config/config.go common/bank/bank.go common/account/account.go common/reply/reply.go common/logger/logger.go
	$(CC) build -o server.exe server/server.go
clean:
	echo "Removing executable files..."
	rm -f client.exe server.exe
	echo "Removing log files..."
	rm -f logs/*
