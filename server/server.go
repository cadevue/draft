package main

import (
	"log"
	"net"
	"os"
	"strconv"
	"syscall"
	"unsafe"

	app "draftgrpc/lib/app"
	pb "draftgrpc/proto"

	"google.golang.org/grpc"
)

const enableVirtualTerminalProcessing = 0x0004

func ANSIColorsInConsole() {
	stdout := syscall.Stdout
	var originalMode uint32

	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	getConsoleMode := kernel32.NewProc("GetConsoleMode")
	getConsoleMode.Call(uintptr(stdout), uintptr(unsafe.Pointer(&originalMode)))

	setConsoleMode := kernel32.NewProc("SetConsoleMode")
	setConsoleMode.Call(uintptr(stdout), uintptr(originalMode|enableVirtualTerminalProcessing))
}

func start_serving(node *app.RaftNode) {
	log.Printf("Starting Raft Server at %s:%d\n", node.Address.IP, node.Address.Port)

	listener, err := net.Listen("tcp", ":"+strconv.FormatInt(int64(node.Address.Port), 10))
	if err != nil {
		node.ApplyCancelMembershipRPC(
			node.Contact_address,
			&pb.MembershipRequest{
				RequesterAddr: node.Address.IP + ":" + strconv.FormatInt(int64(node.Address.Port), 10),
			},
		)
		log.Fatalf("Error starting server: %v", err)
		return
	}

	s := grpc.NewServer()
	pb.RegisterCommandServiceServer(s, node)
	pb.RegisterLogServiceServer(s, node)
	pb.RegisterVoteServiceServer(s, node)
	pb.RegisterMembershipServiceServer(s, node)

	log.Println("Server started, listening for requests...")
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
		node.ApplyCancelMembershipRPC(
			node.Contact_address,
			&pb.MembershipRequest{
				RequesterAddr: node.Address.IP + ":" + strconv.FormatInt(int64(node.Address.Port), 10),
			},
		)
		return
	}
}

func main() {
	// kalo pake mac fungsi ini dikomen aja
	ANSIColorsInConsole()

	if len(os.Args) < 3 {
		log.Fatalf(`

Usage: go run server.go <ip> <port> [<contact_ip> <contact_port>]

example:
go run server.go localhost 5000
		`)
	}

	ip := os.Args[1]
	port, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatalf("Invalid port number: %v", err)
	}

	addr := app.NewAddress(ip, port)
	node := app.NewRaftNode(addr, nil)

	// if port == 5000 {
	// 	node.Peers = append(node.Peers, app.NewAddress(ip, 5001))
	// 	node.Peers = append(node.Peers, app.NewAddress(ip, 5002))
	// } else if port == 5001 {
	// 	node.Peers = append(node.Peers, app.NewAddress(ip, 5000))
	// 	node.Peers = append(node.Peers, app.NewAddress(ip, 5002))
	// } else {
	// 	node.Peers = append(node.Peers, app.NewAddress(ip, 5000))
	// 	node.Peers = append(node.Peers, app.NewAddress(ip, 5001))
	// }

	// if port == 5000 {
	// 	node.Peers = append(node.Peers, app.NewAddress(ip, 5001))
	// } else if port == 5001 {
	// 	node.Peers = append(node.Peers, app.NewAddress(ip, 5000))
	// }

	// Demo doang, ntar dihapus
	// if port == 5000 {
	// 	node.Role = app.Leader
	// 	node.Current_term = 8
	// 	node.Log = []*app.LogEntry{
	// 		app.NewLogEntry(node.Current_term, app.Set, "x", "10"),
	// 		app.NewLogEntry(node.Current_term, app.Set, "y", "20"),
	// 		app.NewLogEntry(node.Current_term, app.Set, "z", "30"),
	// 	}
	// }

	if len(os.Args) == 5 {
		port, err := strconv.Atoi(os.Args[4])
		if err != nil {
			log.Fatalf("Invalid port number: %v", err)
		}

		node.Contact_address = app.NewAddress(os.Args[3], port)
		node.Start()
		start_serving(node)
	} else {
		node.Start()
		start_serving(node)
	}
}
