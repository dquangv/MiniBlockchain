package p2p

import "time"

type NodeState string

const (
	StateLeader   NodeState = "Leader"
	StateFollower NodeState = "Follower"
	StateSyncing  NodeState = "Syncing"
)

var CurrentLeader string
var LastHeartbeat time.Time
