package app

import (
	"context"
	pb "draftgrpc/proto"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc"
)

// All in milliseconds
const HEARTBEAT_INTERVAL = 2000 * time.Millisecond
const ELECTION_TIMEOUT_MIN = 7000 * time.Millisecond
const ELECTION_TIMEOUT_MAX = 12000 * time.Millisecond
const RPC_TIMEOUT = 7000 * time.Millisecond

type RaftNodeRole int

const (
	Follower RaftNodeRole = iota
	Candidate
	Leader
)

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
)

type RaftNode struct {
	pb.CommandServiceServer
	pb.LogServiceServer
	pb.VoteServiceServer
	pb.MembershipServiceServer

	mtx sync.Mutex

	Address         *Address
	Contact_address *Address
	Store           *KVStore

	// election
	Role         RaftNodeRole
	Current_term int
	Voted_for    *Address
	VoterCount   int

	// Log related
	Log          []*LogEntry
	Commit_idx   int
	Last_applied int
	Next_idx     map[*Address]int // entry idx to send to follower
	Match_idx    map[*Address]int // last entry idx that follower has replicated
	commitCond   *sync.Cond

	// channels
	Peers           []*Address
	Timeout         *time.Timer
	TimeoutDuration time.Duration
	TimeoutStart    time.Time
}

func (r *RaftNode) Print() {
	log.Print("_______________________________________________" + Blue)
	log.Println("RaftNode: ", r.Address.Port)
	log.Println("Role: ", r.Role)
	log.Println("Current_term: ", r.Current_term)
	log.Println("Log: ", r.Log)
	log.Println("Commit_idx: ", r.Commit_idx)
	log.Println("Leader: ", r.Contact_address)
	log.Println("Store: ", r.Store)
	log.Println("Peers: ", r.Peers)
	log.Println("VoterCount: ", r.VoterCount)
	log.Println("Vote for: ", r.Voted_for, Reset)
	// log.Println("Voted_for: ", r.Voted_for)
	// if r.Timeout != nil {
	// 	remaining := r.TimeoutDuration - time.Since(r.TimeoutStart)
	// 	log.Println("Timeout: ", remaining.Seconds(), "s")
	// } else {
	// 	log.Println("Timeout: not set")
	// }
	log.Print("________________________________________________")
}

func NewRaftNode(addr *Address, contactAddr *Address) *RaftNode {
	// Implement the RaftNode initialization here
	return &RaftNode{
		Address:         addr,
		Contact_address: contactAddr,
		Store:           NewKVStore(),
		Role:            Follower,
		Current_term:    0,
		Voted_for:       addr,
		Log:             make([]*LogEntry, 0),
		Commit_idx:      -1,
		Last_applied:    -1,
		Next_idx:        make(map[*Address]int),
		Match_idx:       make(map[*Address]int),
		VoterCount:      0,
		commitCond:      sync.NewCond(&sync.Mutex{}),
		Peers:           make([]*Address, 0),
	}
}

func (r *RaftNode) Start() {
	r.StartElectionTimeout()
	if r.Contact_address == nil {
		// I am the first node, initialize as leader
		r.Current_term = 1
		r.Peers = append(r.Peers, r.Address)
		log.Printf("Node %v is the first node\n", r.Address)
		r.StartHeartbeat()
		r.BecomeLeader()
	} else {
		// I am not the first node, initialize as follower
		log.Printf("Node %v is not the first node\n", r.Address)
		r.SendMembershipRequest(r.Contact_address)
	}

	ticker := time.NewTicker(HEARTBEAT_INTERVAL)
	go func() {
		for range ticker.C {
			r.Print()
		}
	}()
}

func (r *RaftNode) BecomeLeader() {
	r.Role = Leader
	r.Contact_address = r.Address
	if r.Timeout != nil {
		r.Timeout.Stop()
	}
	r.Voted_for = nil
	r.VoterCount = 0

	r.StartHeartbeat()
}

func (r *RaftNode) BecomeFollower(leader *Address, term int) {
	r.Role = Follower
	if r.Timeout == nil {
		r.StartElectionTimeout()
	}
	r.Voted_for = nil
	r.VoterCount = 0
	r.Contact_address = leader
	r.Current_term = term
	r.RestartElectionTimeout()
}

func (r *RaftNode) BecomeCandidate() {
	if r.Timeout == nil {
		r.StartElectionTimeout()
	}
	r.Role = Candidate
	r.Current_term++
	r.Voted_for = r.Address
	r.VoterCount = 1
	r.RestartElectionTimeout()
}

func (r *RaftNode) StartHeartbeat() {
	ticker := time.NewTicker(HEARTBEAT_INTERVAL)
	go func() {
		for range ticker.C {
			if r.Role == Leader {
				log.Println(Green, "Heartbeat from: ", r.Address, " and term: ", r.Current_term, Reset)
				r.SendLogsToFollowers()
			}
		}
	}()
}

func (r *RaftNode) StartElectionTimeout() {

	if r.Timeout != nil {
		r.Timeout.Stop()
	}

	r.Timeout = time.AfterFunc(r.RandomElectionTimeout(), func() {
		log.Println(Yellow, "Election timeout, starting new election, becoming candidate...", Reset)
		r.BecomeCandidate()
		r.BroadcastVoteRequest()
	})
}

func (r *RaftNode) RestartElectionTimeout() {
	if r.Timeout != nil {
		r.Timeout.Reset(r.RandomElectionTimeout())
	}
}

func (r *RaftNode) RandomElectionTimeout() time.Duration {
	return time.Duration(rand.Int63n(int64(ELECTION_TIMEOUT_MAX-ELECTION_TIMEOUT_MIN))) + ELECTION_TIMEOUT_MIN
}

func (r *RaftNode) Ping(ctx context.Context, _ *pb.PingRequest) (*pb.PingResponse, error) {
	log.Println("Received ping from client")

	// Nyobain doang, tar harus dihapus
	r.Peers = []*Address{NewAddress("localhost", 5001), NewAddress("localhost", 5002)}
	r.SendLogsToFollowers()

	return &pb.PingResponse{Message: "Pong!"}, nil
}

// Get handles the Get RPC request
func (r *RaftNode) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	log.Println("Received get request from client")

	r.mtx.Lock()
	defer r.mtx.Unlock()

	value, err := r.Store.Get(req.Key)
	if err != nil {
		return &pb.GetResponse{Value: err.Error()}, err
	}

	return &pb.GetResponse{Value: value}, nil
}

// RequestLog handles the RequestLog RPC request
func (r *RaftNode) RequestLog(ctx context.Context, req *pb.LogRequest) (*pb.LogResponse, error) {
	log.Println("Received log request from client")
	r.mtx.Lock()
	defer r.mtx.Unlock()

	entries := make([]*pb.LogEntryResponse, 0)
	for _, entry := range r.Log {
		entries = append(entries, &pb.LogEntryResponse{
			Term:    int32(entry.GetTerm()),
			Command: int32(entry.GetCommand()),
			Key:     entry.GetKey(),
			Value:   entry.GetValue(),
		})
	}

	return &pb.LogResponse{Entries: entries}, nil
}

// Set handles the Set RPC request
func (r *RaftNode) Set(ctx context.Context, req *pb.SetRequest) (*pb.SetResponse, error) {

	if r.Role != Leader {
		return &pb.SetResponse{Status: ""}, fmt.Errorf("connect to leader [%s] to set key", r.Contact_address)
	}

	log.Println("Received set request from client")
	r.mtx.Lock()

	r.Log = append(r.Log, &LogEntry{r.Current_term, Set, req.Key, req.Value})
	index := len(r.Log) - 1

	r.mtx.Unlock()

	r.commitCond.L.Lock()
	for r.Commit_idx < index {
		r.commitCond.Wait()
	}
	r.commitCond.L.Unlock()

	res, err := r.Store.Set(req.Key, req.Value)
	r.Last_applied = index
	if err != nil {
		return &pb.SetResponse{Status: ""}, err
	}

	return &pb.SetResponse{Status: res}, nil
}

// Strln handles the Strln RPC request
func (r *RaftNode) Strln(ctx context.Context, req *pb.StrlnRequest) (*pb.StrlnResponse, error) {

	log.Println("Received strln request from client")
	r.mtx.Lock()
	defer r.mtx.Unlock()

	length, err := r.Store.Strln(req.Key)
	if err != nil {
		return &pb.StrlnResponse{Length: -1}, err
	}

	return &pb.StrlnResponse{Length: int32(length)}, nil
}

// Del handles the Del RPC request
func (r *RaftNode) Del(ctx context.Context, req *pb.DelRequest) (*pb.DelResponse, error) {
	if r.Role != Leader {
		return &pb.DelResponse{Value: ""}, fmt.Errorf("connect to leader [%s] to delete key", r.Contact_address)
	}

	log.Println("Received del request from client")
	r.mtx.Lock()

	if !r.Store.KeyAvailable(req.Key) {
		r.mtx.Unlock()
		return &pb.DelResponse{Value: ""}, fmt.Errorf("key not found")
	}

	r.Log = append(r.Log, &LogEntry{r.Current_term, Del, req.Key, ""})
	index := len(r.Log) - 1

	r.mtx.Unlock()

	r.commitCond.L.Lock()
	for r.Commit_idx < index {
		r.commitCond.Wait()
	}
	r.commitCond.L.Unlock()

	value, err := r.Store.Delete(req.Key)
	if err != nil {
		return &pb.DelResponse{Value: err.Error()}, err
	}

	return &pb.DelResponse{Value: value}, nil
}

// Append handles the Append RPC request
func (r *RaftNode) Append(ctx context.Context, req *pb.AppendRequest) (*pb.AppendResponse, error) {

	if r.Role != Leader {
		return &pb.AppendResponse{Status: ""}, fmt.Errorf("connect to leader [%s] to append value", r.Contact_address)
	}

	log.Println("Received append request from client")
	r.mtx.Lock()

	r.Log = append(r.Log, &LogEntry{r.Current_term, Append, req.Key, req.Value})
	index := len(r.Log) - 1

	r.mtx.Unlock()

	r.commitCond.L.Lock()
	for r.Commit_idx < index {
		r.commitCond.Wait()
	}
	r.commitCond.L.Unlock()

	status, err := r.Store.Append(req.Key, req.Value)
	if err != nil {
		return &pb.AppendResponse{Status: ""}, err
	}
	return &pb.AppendResponse{Status: status}, nil
}

func (r *RaftNode) AppendEntries(ctx context.Context, req *pb.AppendEntriesRequest) (*pb.AppendEntriesResponse, error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.RestartElectionTimeout()

	host := req.SenderIp
	port, _ := strconv.ParseInt(req.SenderPort, 10, 0)

	r.BecomeFollower(NewAddress(host, int(port)), int(req.Term))

	if req.Term < int32(r.Current_term) {
		// stale/outdated leader, ignore
		return &pb.AppendEntriesResponse{Term: int32(r.Current_term), Success: false, Message: "Leader outdated"}, nil
	}

	if req.Term > int32(r.Current_term) {
		// new leader has been elected, salute the new leader
		r.Current_term = int(req.Term)
		r.Role = Follower
	}

	// leader's thought this node is more up-to-date than it actually is
	// leader have corrupted log for this node
	if req.PrevLogIndex >= 0 && (len(r.Log) <= int(req.PrevLogIndex) ||
		r.Log[req.PrevLogIndex].GetTerm() != int(req.PrevLogTerm)) {
		return &pb.AppendEntriesResponse{Term: int32(r.Current_term), Success: false, Message: "Leader corrupted"}, nil
	}

	// if entry conflicts, replace the conflicting entry and all that follow it
	for i := req.PrevLogIndex + 1; i < int32(len(r.Log)); i++ {
		if i-req.PrevLogIndex-1 >= int32(len(req.Entries)) || int32(r.Log[i].GetTerm()) != req.Entries[i-req.PrevLogIndex-1].Term {
			r.Log = r.Log[:i]
			break
		}
	}

	// append any new entries not already in the log
	for i := req.PrevLogIndex + 1; i-req.PrevLogIndex-1 < int32(len(req.Entries)); i++ {
		if i >= int32(len(r.Log)) {
			data := req.Entries[i-req.PrevLogIndex-1]
			entry := LogEntry{
				int(data.Term), CommandType(data.Command), data.Key, data.Value,
			}
			r.Log = append(r.Log, &entry)
		}
	}

	// update commit index
	if req.LeaderCommit > int32(r.Commit_idx) {
		log.Println("Commiting it...")
		r.Commit_idx = int(req.LeaderCommit)
		r.applyUncommitedLogs()
	}

	// log.Printf("Update log: %v\n", r.Log)
	log.Println(Yellow, "Received Entries: ", req.Entries, Reset)
	log.Println(Yellow, "Leader Commit: ", req.LeaderCommit, Reset)

	return &pb.AppendEntriesResponse{Term: int32(r.Current_term), Success: true, Message: "Success"}, nil
}

func (r *RaftNode) applyUncommitedLogs() {
	for i := r.Last_applied + 1; i <= r.Commit_idx; i++ {
		entry := r.Log[i]
		switch entry.GetCommand() {
		case Set:
			r.Store.Set(entry.GetKey(), entry.GetValue())
		case Del:
			r.Store.Delete(entry.GetKey())
		case Append:
			r.Store.Append(entry.GetKey(), entry.GetValue())
		}
	}
	r.Last_applied = r.Commit_idx
}

func (r *RaftNode) SendLogsToFollowers() {
	for _, peer := range r.Peers {
		if peer.Equal(r.Address) {
			continue
		}
		go func(peer *Address) {

			r.mtx.Lock()
			prevLogIdx := r.Next_idx[peer] - 1
			prevLogTerm := -1
			if prevLogIdx >= 0 {
				prevLogTerm = r.Log[prevLogIdx].GetTerm()
			}
			entries := make([]*pb.LogEntry, 0)
			for i := r.Next_idx[peer]; i < len(r.Log); i++ {
				entry := r.Log[i]
				entries = append(entries, &pb.LogEntry{
					Term:    int32(entry.GetTerm()),
					Command: int32(entry.GetCommand()),
					Key:     entry.GetKey(),
					Value:   entry.GetValue(),
				})
			}

			r.mtx.Unlock()

			req := &pb.AppendEntriesRequest{
				Term:         int32(r.Current_term),
				PrevLogIndex: int32(prevLogIdx),
				PrevLogTerm:  int32(prevLogTerm),
				Entries:      entries,
				LeaderCommit: int32(r.Commit_idx),
				SenderIp:     r.Address.IP,
				SenderPort:   strconv.Itoa(r.Address.Port),
			}

			resp, err := r.AppendEntriesRPC(peer, req)
			if err != nil {
				log.Println(Red, "Error sending AppendEntries to %v: %v\n", peer, err, Reset)
				return
			}

			r.mtx.Lock()
			defer r.mtx.Unlock()

			if resp.Success {
				log.Println("Succesfully replicate logs to", peer)
				r.Next_idx[peer] = prevLogIdx + len(entries) + 1
				r.Match_idx[peer] = r.Next_idx[peer] - 1

				for i := len(r.Log) - 1; i > r.Commit_idx; i-- {
					count := 1
					for _, matchIdx := range r.Match_idx {
						if matchIdx >= i {
							count++
						}
					}

					if count > (len(r.Peers)-1)/2 {
						r.Commit_idx = i
						r.commitCond.Broadcast()
						break
					}
				}
			} else {
				log.Println("Failed to replicate log because", resp.Message)
				r.Next_idx[peer] = max(1, r.Next_idx[peer]-1)
			}
		}(peer)
	}
}

func (r *RaftNode) AppendEntriesRPC(peer *Address, req *pb.AppendEntriesRequest) (*pb.AppendEntriesResponse, error) {
	addr := peer.IP + ":" + strconv.Itoa(peer.Port)
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}

	client_server := pb.NewLogServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), RPC_TIMEOUT)
	defer cancel()

	resp, err := client_server.AppendEntries(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (r *RaftNode) BroadcastVoteRequest() {

	r.StartElectionTimeout()
	log.Println(Yellow, "Broadcasting VoteRequest, current term: %v\n", r.Current_term, Reset)
	for _, peer := range r.Peers {
		if peer.Equal(r.Address) {
			continue
		}
		go func(peer *Address) {

			var lastLogIndex int32
			var lastLogTerm int32

			if len(r.Log) > 0 {
				lastLogIndex = int32(len(r.Log) - 1)
				lastLogTerm = int32(r.Log[len(r.Log)-1].GetTerm())
			} else {
				lastLogIndex = -1 //no entries
				lastLogTerm = -1  //no entries
			}
			voteReq := &pb.VoteRequest{
				Term:         int32(r.Current_term),
				LastLogIndex: lastLogIndex,
				LastLogTerm:  lastLogTerm,
				SenderIp:     r.Address.IP,
				SenderPort:   strconv.Itoa(r.Address.Port),
			}

			resp, err := r.sendVoteRequest(peer, voteReq)
			if resp == nil && err != nil {
				log.Printf("Error sending VoteRequest to %v: %v\n", peer, err)
				return
			}

			r.mtx.Lock()
			defer r.mtx.Unlock()

			if resp.VoteGranted {
				r.VoterCount++
			}

			if r.VoterCount > len(r.Peers)/2 {
				log.Println(Green, "Election won with vote ", r.VoterCount, "/", len(r.Peers), Reset)
				r.BecomeLeader()
			}

		}(peer)
	}

}

func (r *RaftNode) sendVoteRequest(peer *Address, req *pb.VoteRequest) (*pb.VoteResponse, error) {
	addr := peer.IP + ":" + strconv.Itoa(peer.Port)
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}

	client_server := pb.NewVoteServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), RPC_TIMEOUT)
	defer cancel()

	resp, err := client_server.InitiateVote(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// there is one node that initiates the vote and call this function
// in every other node. (all must comply and vote)
func (r *RaftNode) InitiateVote(ctx context.Context, req *pb.VoteRequest) (*pb.VoteResponse, error) {

	// condition for vote: if candidate term higher, if canditate lastLogIndex is higher or the same
	vote_granted := true

	// candidate rejected
	if req.LastLogIndex != -1 && req.LastLogTerm != -1 {
		if req.LastLogIndex < int32(len(r.Log)-1) {
			log.Println(Red, "Vote rejected cause last log INDEX is smaller than me, req: ", req, Reset)
			vote_granted = false
		}

		if req.LastLogTerm < int32(r.Log[len(r.Log)-1].term) {
			vote_granted = false
			log.Println(Red, "Vote rejected cause last log TERM is smaller than me, req: ", req, Reset)
		}

		if r.Role == Candidate && r.Current_term >= int(req.Term) {
			vote_granted = false
			log.Println(Red, "Vote rejected cause I already candidate, req: ", req, Reset)
		}
	}

	if vote_granted {
		port, _ := strconv.ParseInt(req.SenderPort, 10, 0)
		r.Voted_for = NewAddress(req.SenderIp, int(port))
		r.BecomeFollower(r.Voted_for, int(req.Term))
	}

	return &pb.VoteResponse{
		Term:        int32(r.Current_term),
		VoteGranted: vote_granted,
	}, nil
}

func (r *RaftNode) SendMembershipRequest(address *Address) {
	resp, err := r.applyMembershipRPC(address, &pb.MembershipRequest{RequesterAddr: r.Address.IP + ":" + strconv.Itoa(r.Address.Port)})
	if err != nil {
		log.Printf("Error applying membership: %v\n", err)
		return
	}
	r.mtx.Lock()
	defer r.mtx.Unlock()
	log.Printf("Membership response: %v\n", resp)
	if resp.Status == "success" {
		// convert resp.LeaderAddr ip:port string to Address
		parts := strings.Split(resp.LeaderAddr, ":")
		if len(parts) != 2 {
			log.Printf("Invalid leader address format: %s", resp.LeaderAddr)
			return
		}
		leader_addr := parts[0]
		leader_port, err := strconv.Atoi(parts[1])
		// log.Printf("[MEMBERSHIP] Leader address: %v:%v\n", leader_addr, leader_port)
		if err != nil {
			log.Printf("Failed to convert port number for leader address %s: %v", resp.LeaderAddr, err)
			return
		}
		r.BecomeFollower(NewAddress(leader_addr, leader_port), int(resp.Term))
		log.Printf("Successfully applied membership to %v\n", address)
	} else {
		log.Printf("Failed to apply membership: %v\n", resp.Status)
		log.Printf("Try to apply to the leader address: %v\n", resp.LeaderAddr)
	}
}

func (r *RaftNode) ApplyMembership(ctx context.Context, req *pb.MembershipRequest) (*pb.MembershipResponse, error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	if r.Role == Leader {
		// split req.RequesterAddr into IP and port
		parts := strings.Split(req.RequesterAddr, ":")
		if len(parts) != 2 {
			log.Printf("Invalid requester address format: %s", req.RequesterAddr)
			return &pb.MembershipResponse{Status: "failed", LeaderAddr: r.Contact_address.IP + ":" + strconv.Itoa(r.Contact_address.Port)}, nil
		}
		ip := parts[0]
		port, err := strconv.Atoi(parts[1])
		if err != nil {
			log.Printf("Failed to convert port number for requester address %s: %v", req.RequesterAddr, err)
			return &pb.MembershipResponse{Status: "failed", LeaderAddr: r.Contact_address.IP + ":" + strconv.Itoa(r.Contact_address.Port)}, nil
		}

		// Check if the peer already exists
		is_new_peer := true
		newPeer := NewAddress(ip, port)
		for _, peer := range r.Peers {
			if peer.IP == newPeer.IP && peer.Port == newPeer.Port {
				log.Printf("Peer already exists: %v\n", newPeer)
				is_new_peer = false
			}
		}

		// Add new peer if it doesn't already exist
		if is_new_peer {
			r.Peers = append(r.Peers, newPeer)
			log.Printf("New peer added: %v\n", newPeer)
		}
		r.copyPeersToPeers()
		return &pb.MembershipResponse{Status: "success", LeaderAddr: r.Contact_address.IP + ":" + strconv.Itoa(r.Contact_address.Port), Term: int32(r.Current_term)}, nil
	}
	return &pb.MembershipResponse{Status: "failed", LeaderAddr: r.Contact_address.IP + ":" + strconv.Itoa(r.Contact_address.Port), Term: int32(r.Current_term)}, nil
}

func (r *RaftNode) applyMembershipRPC(address *Address, req *pb.MembershipRequest) (*pb.MembershipResponse, error) {
	addr := address.IP + ":" + strconv.Itoa(address.Port)
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}

	client_server := pb.NewMembershipServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), RPC_TIMEOUT)
	defer cancel()

	resp, err := client_server.ApplyMembership(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *RaftNode) CopyPeers(ctx context.Context, req *pb.CopyPeersRequest) (*pb.CopyPeersResponse, error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	// convert req.Peers string to Adress
	r.Peers = make([]*Address, 0)
	for _, addrStr := range req.Peers {
		// Split addrStr into IP and port
		log.Printf("Peer address: %v\n", addrStr)
		parts := strings.Split(addrStr, ":")
		if len(parts) != 2 {
			log.Printf("Invalid peer address format: %s", addrStr)
			continue
		}
		ip := parts[0]
		port, err := strconv.Atoi(parts[1])
		if err != nil {
			log.Printf("Failed to convert port number for peer address %s: %v", addrStr, err)
			continue
		}
		r.Peers = append(r.Peers, &Address{
			IP:   ip,
			Port: port,
		})
	}
	return &pb.CopyPeersResponse{Status: "success"}, nil
}

func (r *RaftNode) sendCopyPeersRPC(peer *Address, req *pb.CopyPeersRequest) (*pb.CopyPeersResponse, error) {
	addr := peer.IP + ":" + strconv.Itoa(peer.Port)
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}

	client_server := pb.NewMembershipServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), RPC_TIMEOUT)
	defer cancel()

	resp, err := client_server.CopyPeers(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (r *RaftNode) copyPeersToPeers() {
	for _, peer := range r.Peers {
		go func(peer *Address) {
			// convert address to string
			var addressStrings []string
			for _, addr := range r.Peers {
				addressStrings = append(addressStrings, fmt.Sprintf("%s:%d", addr.IP, addr.Port))
			}
			req := &pb.CopyPeersRequest{
				Peers: addressStrings,
			}
			resp, err := r.sendCopyPeersRPC(peer, req)
			if err != nil {
				log.Printf("Error sending CopyPeers to %v: %v\n", peer, err)
				return
			}
			r.mtx.Lock()
			defer r.mtx.Unlock()
			if resp.Status == "success" {
				log.Printf("Successfully copied peers to %v\n", peer)
			} else {
				log.Printf("Failed to copy peers to %v\n", peer)
			}
		}(peer)
	}
}

func (r *RaftNode) CancelMembership(ctx context.Context, req *pb.MembershipRequest) (*pb.MembershipResponse, error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	if r.Role == Leader {
		// split req.RequesterAddr into IP and port
		parts := strings.Split(req.RequesterAddr, ":")
		if len(parts) != 2 {
			log.Printf("Invalid requester address format: %s", req.RequesterAddr)
			return &pb.MembershipResponse{Status: "failed", LeaderAddr: r.Contact_address.IP + ":" + strconv.Itoa(r.Contact_address.Port)}, nil
		}
		ip := parts[0]
		port, err := strconv.Atoi(parts[1])
		if err != nil {
			log.Printf("Failed to convert port number for requester address %s: %v", req.RequesterAddr, err)
			return &pb.MembershipResponse{Status: "failed", LeaderAddr: r.Contact_address.IP + ":" + strconv.Itoa(r.Contact_address.Port)}, nil
		}
		for i, peer := range r.Peers {
			if peer.IP == ip && peer.Port == port {
				r.Peers = append(r.Peers[:i], r.Peers[i+1:]...)
				break
			}
		}
		log.Printf("Membership cancelled for %v\n", req.RequesterAddr)
		log.Printf("Current peers: %v\n", r.Peers)
		r.copyPeersToPeers()
		return &pb.MembershipResponse{Status: "success", LeaderAddr: r.Contact_address.IP + ":" + strconv.Itoa(r.Contact_address.Port)}, nil
	}
	return &pb.MembershipResponse{Status: "failed", LeaderAddr: r.Contact_address.IP + ":" + strconv.Itoa(r.Contact_address.Port)}, nil
}

func (r *RaftNode) ApplyCancelMembershipRPC(address *Address, req *pb.MembershipRequest) (*pb.MembershipResponse, error) {
	addr := address.IP + ":" + strconv.Itoa(address.Port)
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}

	client_server := pb.NewMembershipServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), RPC_TIMEOUT)
	defer cancel()

	resp, err := client_server.CancelMembership(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
