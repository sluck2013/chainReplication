Compile:
Compile with gccgo compiler. To compile the files, just run make under root directory.
The root directory should be set as $GOPATH/src
make clean will remove all executable files and log files

Run:
start single server - ./server.exe <configFileName> serverID
start multiple servers - ./launch.sh <configFileName> <minServerId> <maxServerId>
client - ./client.exe <configFileName>
* configFileName includes only the file name, should NOT inclue file path
* minServerId/maxServerId is the min/max serverID of those servers which are intentded to start at a time
