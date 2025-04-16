package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	pb "draftgrpc/proto"

	"google.golang.org/grpc"
)

func main() {
	/*
		// Set up a connection to the server.
		ports := []string{"5000", "5001"} // List of ports
		client_server_instances := make(map[string]pb.CommandServiceClient)
		for _, port := range ports {
			address := "localhost:" + port
			log.Println("Connecting to server at", address)
			conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
			if err == nil {
				log.Println("Connected to server at", address)
				client_server_instances[port] = pb.NewCommandServiceClient(conn)
			} else {
				log.Printf("Failed to connect to %v: %v\n", address, err)
			}
		}
	*/
	if len(os.Args) < 3 {
		log.Fatalf(`

Usage: go run client.go <ip> <port> [<contact_ip> <contact_port>]

example:
go run server.go localhost 5000
		`)
	}

	log.Println("Args: ", os.Args)
	ip := os.Args[1]
	portStr := os.Args[2]
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Invalid port number: %v", err)
	}

	address := ip + ":" + portStr
	log.Println("Connecting to server at", address)
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Printf("Failed to connect to %v: %v\n", address, err)
	}

	client_server := pb.NewCommandServiceClient(conn)
	log.Println("Connected to server at", address)

	for {
		// Prompt the user for input
		fmt.Print(">> ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		parts := strings.Split(strings.TrimSpace(input), " ")

		if len(parts) < 1 {
			log.Println("Bad Usage! Usage: <command>")
			continue
		}

		cmd := strings.ToLower(parts[0])

		// client_server, ok := client_server_instances[port]
		// if !ok {
		// 	fmt.Println("Not connected to server at port", port)
		// 	continue
		// }

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		switch cmd {
		case "ping":
			r, err := client_server.Ping(ctx, &pb.PingRequest{})
			if err != nil {
				log.Fatalf("Could not greet: %v", err)
				continue
			}
			log.Printf("Response from %s: %s", strconv.Itoa(port), r.GetMessage())

		case "get":
			if len(parts) < 2 {
				log.Println("Usage: get <key>")
				continue
			}
			key := parts[1]
			r, err := client_server.Get(ctx, &pb.GetRequest{Key: key})
			if err != nil {
				log.Printf("Could not get value: %v", err)
				continue
			}
			log.Printf("Get response from %s: %s", portStr, r.Value)

		case "set":
			if len(parts) < 3 {
				log.Println("Usage: set <key> <value>")
				continue
			}
			key := parts[1]
			value := parts[2]
			r, err := client_server.Set(ctx, &pb.SetRequest{Key: key, Value: value})
			if err != nil {
				log.Printf("Could not set value: %v", err)
				continue
			}
			log.Printf("Set response from %s: %s", portStr, r.Status)

		case "strln":
			if len(parts) < 2 {
				log.Println("Usage: strln <key>")
				continue
			}
			key := parts[1]
			r, err := client_server.Strln(ctx, &pb.StrlnRequest{Key: key})
			if err != nil {
				log.Printf("Could not get length: %v", err)
				continue
			}
			log.Printf("Strln response from %s: %d", portStr, r.Length)

		case "del":
			if len(parts) < 2 {
				log.Println("Usage: del <key>")
				continue
			}
			key := parts[1]
			r, err := client_server.Del(ctx, &pb.DelRequest{Key: key})
			if err != nil {
				log.Printf("Could not delete value: %v", err)
				continue
			}
			log.Printf("Del response from %s: %s", portStr, r.Value)

		case "append":
			if len(parts) < 3 {
				log.Println("Usage: append <key> <value>")
				continue
			}
			key := parts[1]
			value := parts[2]
			r, err := client_server.Append(ctx, &pb.AppendRequest{Key: key, Value: value})
			if err != nil {
				log.Printf("Could not append value: %v", err)
				continue
			}
			log.Printf("Append response from %s: %s", portStr, r.Status)

		case "request_log":
			r, err := client_server.RequestLog(ctx, &pb.LogRequest{})
			if err != nil {
				log.Fatalf("Could not greet: %v", err)
				continue
			}

			log.Println("Log entries:")

			for _, entry := range r.Entries {
				log.Printf("Term: %d, Command: %d, Key: %s, Value: %s\n", entry.Term, entry.Command, entry.Key, entry.Value)
			}

		case "switch":
			if len(parts) < 3 {
				log.Println("Usage: switch <ip> <port>")
				continue
			}
			ip = parts[1]
			portStr = parts[2]
			port, err = strconv.Atoi(portStr)
			if err != nil {
				log.Fatalf("Invalid port number: %v", err)
			}

			// Switch to the new server
			address := ip + ":" + portStr
			log.Println("Connecting to server at", address)
			conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
			if err != nil {
				log.Printf("Failed to connect to %v: %v\n", address, err)
			}

			client_server = pb.NewCommandServiceClient(conn)
			log.Println("Connected to server at", address)

			if err != nil {
				log.Printf("Could not switch server: %v", err)
				continue
			}
			log.Printf("Switch response from %s: %s", portStr, "OK")

		default:
			log.Println("Unknown command:", cmd)
		}
	}
}
