Compile:
Compile with gccgo compiler. To compile the files, just run make under root directory.
The root directory should be put set as $GOPATH/src
make clean will remove all executable files and log files

Run:
server - ./server <configFileName> serverID
client - ./client <configFileName>
* configFileName includes only the file name, should NOT inclue file path
