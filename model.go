package main

type TopologySummary struct {
	Topologies []Topology `json:"topologies"`
}

type Topology struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Workers int    `json:"workersTotal"`
	Bolts   []Bolt `json:"bolts"`
}

type Bolt struct {
	BoldId         string `json:"boltId"`
	Executors      int    `json:"executors"`
	Tasks          int    `json:"tasks"`
	Capacity       string `json:"capacity"`
	ExecuteLatency string `json:"executeLatency"`
}

type TopologyLag map[string]SpoutLag

type SpoutLag struct {
	SpoutId        string                               `json:"spoutId"`
	SpoutType      string                               `json:"spoutType"`
	SpoutLagResult map[string]map[string]SpoutLagResult `json:"spoutLagResult"`
}

type SpoutLagResult struct {
	ConsumerCommittedOffset int64 `json:"consumerCommittedOffset"`
	LogHeadOffset           int64 `json:"logHeadOffset"`
	Lag                     int64 `json:"lag"`
}

type StormLog struct {
	Id         int    `db:"go_storm_log_id"`
	TopologyId string `db:"topology_id"`
	Workers    int    `db:"workers"`
}

type StormLogDetail struct {
	StormLogId     int64   `db:"go_storm_log_id"`
	BoltId         string  `db:"bolt_id"`
	Executors      int     `db:"executors"`
	Tasks          int     `db:"tasks"`
	Capacity       float64 `db:"capacity"`
	ExecuteLatency float64 `db:"execute_latency"`
}

type StormLogLagDetail struct {
	StormLogId      int64  `db:"go_storm_log_id"`
	SpoutId         string `db:"spout_id"`
	Topic           string `db:"topic"`
	Partition       string `db:"topic_partition"`
	CommittedOffset int64  `db:"consumer_committed_offset"`
	LogHeadOffset   int64  `db:"log_head_offset"`
	Lag             int64  `db:"lag"`
}
