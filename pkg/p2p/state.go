package p2p

type NodeState string

const (
	StateLeader   NodeState = "Leader"
	StateFollower NodeState = "Follower"
	StateSyncing  NodeState = "Syncing"
)
