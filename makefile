.PHONY: server1 server2 all client

IP = localhost
PORT1 = 5000
PORT2 = 5001

DIR_SERVER = server/
DIR_CLIENT = client/

server1:
    start cmd /k go run $(DIR_SERVER)server.go $(IP) $(PORT1)

server2:
    start cmd /k go run $(DIR_SERVER)server.go $(IP) 5001 $(IP) $(PORT1)

server3:
    start cmd /k go run $(DIR_SERVER)server.go $(IP) 5002 $(IP) $(PORT1)

server4:
    start cmd /k go run $(DIR_SERVER)server.go $(IP) 5003 $(IP) $(PORT1)

server5:
    start cmd /k go run $(DIR_SERVER)server.go $(IP) 5004 $(IP) $(PORT1)

server6:
    start cmd /k go run $(DIR_SERVER)server.go $(IP) 5005 $(IP) $(PORT1)

client: FORCE
	start cmd /k go run $(DIR_CLIENT)client.go $(IP) 5000

cluster1: server1 server2 server3 client

all: cluster1 client

FORCE: