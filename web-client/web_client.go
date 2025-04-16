package main

import (
	"context"
	pb "draftgrpc/proto"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"google.golang.org/grpc"
)

var (
	ip         string
	port       string
	clientConn pb.CommandServiceClient
)

const RPC_TIMEOUT = 5 * time.Second

func main() {
	// Set up a connection to the server.
	if len(os.Args) < 3 {
		log.Fatalf(`

Usage: go run client.go <ip> <port> [<contact_ip> <contact_port>]

example:
go run server.go localhost 5000
		`)
	}

	log.Println("Args: ", os.Args)
	local_ip := os.Args[1]
	local_port, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatalf("Invalid port number: %v", err)
	}

	portStr := strconv.Itoa(local_port)
	address := local_ip + ":" + portStr
	log.Println("Connecting to server at", address)
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Printf("Failed to connect to %v: %v\n", address, err)
	}

	clientConn = pb.NewCommandServiceClient(conn)
	log.Println("Connected to server at", address)

	ip = local_ip
	port = portStr

	http.HandleFunc("/connect", connectHandler)
	http.HandleFunc("/ping", pingHandler)
	http.HandleFunc("/get", getHandler)
	http.HandleFunc("/set", setHandler)
	http.HandleFunc("/strln", strlnHandler)
	http.HandleFunc("/del", delHandler)
	http.HandleFunc("/append", appendHandler)
	http.HandleFunc("/switch", switchHandler)
	http.HandleFunc("/request_logs", requestLogsHandler)

	log.Println("Starting web client at localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func setupCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func connectHandler(w http.ResponseWriter, r *http.Request) {
	setupCORS(&w)
	if r.Method == "OPTIONS" {
		return
	}

	local_ip := r.URL.Query().Get("ip")
	if local_ip == "" {
		local_ip = "localhost"
	}

	local_port := r.URL.Query().Get("port")
	if local_port == "" {
		http.Error(w, "Port is required", http.StatusBadRequest)
		return
	}

	if _, err := strconv.Atoi(local_port); err != nil {
		http.Error(w, "Invalid port number", http.StatusBadRequest)
		return
	}

	if local_ip != ip || local_port != port {
		http.Error(w, "Not connected to "+local_ip+":"+local_port, http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "%s:%s", local_ip, local_port)
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	setupCORS(&w)
	if r.Method == "OPTIONS" {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), RPC_TIMEOUT)
	defer cancel()

	rsp, err := clientConn.Ping(ctx, &pb.PingRequest{})
	if err != nil {
		http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, rsp.Message)
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	setupCORS(&w)
	if r.Method == "OPTIONS" {
		return
	}

	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "Key is required", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), RPC_TIMEOUT)
	defer cancel()

	rsp, err := clientConn.Get(ctx, &pb.GetRequest{Key: key})
	if err != nil {
		http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, rsp.Value)
}

func setHandler(w http.ResponseWriter, r *http.Request) {
	setupCORS(&w)
	if r.Method == "OPTIONS" {
		return
	}

	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")
	if key == "" {
		http.Error(w, "Key is required", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), RPC_TIMEOUT)
	defer cancel()

	rsp, err := clientConn.Set(ctx, &pb.SetRequest{Key: key, Value: value})
	if err != nil {
		http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Set response from %s:%s: %s", ip, port, rsp.Status)
}

func strlnHandler(w http.ResponseWriter, r *http.Request) {
	setupCORS(&w)
	if r.Method == "OPTIONS" {
		return
	}

	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "Key is required", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), RPC_TIMEOUT)
	defer cancel()

	rsp, err := clientConn.Strln(ctx, &pb.StrlnRequest{Key: key})
	if err != nil {
		http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Strln response from %s:%s: %d", ip, port, rsp.Length)
}

func delHandler(w http.ResponseWriter, r *http.Request) {
	setupCORS(&w)
	if r.Method == "OPTIONS" {
		return
	}

	key := r.URL.Query().Get("key")

	ctx, cancel := context.WithTimeout(context.Background(), RPC_TIMEOUT)
	defer cancel()

	rsp, err := clientConn.Del(ctx, &pb.DelRequest{Key: key})
	if err != nil {
		http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Del response from %s:%s: %s", ip, port, rsp.Value)
}

func appendHandler(w http.ResponseWriter, r *http.Request) {
	setupCORS(&w)
	if r.Method == "OPTIONS" {
		return
	}

	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")
	if key == "" {
		http.Error(w, "Key is required", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), RPC_TIMEOUT)
	defer cancel()

	rsp, err := clientConn.Append(ctx, &pb.AppendRequest{Key: key, Value: value})
	if err != nil {
		http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Append response from %s:%s: %s", ip, port, rsp.Status)
}

func switchHandler(w http.ResponseWriter, r *http.Request) {
	setupCORS(&w)
	if r.Method == "OPTIONS" {
		return
	}

	ip := r.URL.Query().Get("ip")
	port := r.URL.Query().Get("port")
	if ip == "" || port == "" {
		http.Error(w, "ip and port are required", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), RPC_TIMEOUT)
	defer cancel()

	rsp, err := clientConn.Switch(ctx, &pb.SwitchRequest{Ip: ip, Port: port})
	if err != nil {
		http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Switch response from %s:%s: %s", ip, port, rsp.Status)
}

func requestLogsHandler(w http.ResponseWriter, r *http.Request) {
	setupCORS(&w)
	if r.Method == "OPTIONS" {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), RPC_TIMEOUT)
	defer cancel()

	rsp, err := clientConn.RequestLog(ctx, &pb.LogRequest{})
	if err != nil {
		http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if len(rsp.Entries) == 0 {
		fmt.Fprintf(w, "No log entries found\n")
	} else {
		commandMap := map[int32]string{
			0: "Set",
			1: "Del",
			2: "Append",
		}
		
		for _, entry := range rsp.Entries {
			command := commandMap[entry.Command]
			fmt.Fprintf(w, "Term: %d, Command: %s, Key: %s, Value: %s\n", entry.Term, command, entry.Key, entry.Value)
		}
	}
}
