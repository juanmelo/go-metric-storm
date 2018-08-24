package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jasonlvhit/gocron"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	var databaseUrl = os.Getenv("DATABASE_URL")
	interval, _ := strconv.ParseUint(os.Getenv("STORM_INTERVAL_CHECK"), 10, 64)

	log.Printf("Reading topolgies %v from %v ( database: %v ) \n", os.Getenv("TOPOLOGY_NAME"), os.Getenv("STORM_CONSOLE"), databaseUrl)

	gocron.Every(interval).Minutes().Do(task)

	<-gocron.Start()
}

func task() {
	var stormConsole = os.Getenv("STORM_CONSOLE")
	var topologyNames = strings.Split(os.Getenv("TOPOLOGY_NAME"), ",")

	for _, topologyName := range topologyNames {
		fmt.Println("Getting metrics for " + topologyName)

		topologyId, err := getTopologyId(stormConsole+"/api/v1/topology/summary", topologyName)
		if err != nil {
			fmt.Println("Error getting metrics ", err)
			return
		}
		topology, err := getTopology(stormConsole + "/api/v1/topology/" + topologyId)
		if err != nil {
			fmt.Println("Error getting metrics ", err)
			return
		}
		fmt.Println(topology)

		topologyLag, err := getTopologyLag(stormConsole + "/api/v1/topology/" + topologyId + "/lag")
		if err != nil {
			fmt.Println("Error getting metrics ", err)
			return
		}
		fmt.Println(topologyLag)

		stormLog := StormLog{
			TopologyId: topology.Id,
			Workers:    topology.Workers,
		}

		var databaseUrl = os.Getenv("DATABASE_URL")
		db := sqlx.MustConnect("mysql", databaseUrl)
		db.Ping()

		defer db.Close()
		query := "insert into go_storm_log (topology_id,workers,created_date) " +
			"values(:topology_id,:workers,sysdate())"
		res, err := db.NamedExec(query, &stormLog)

		if err != nil {
			fmt.Println("Error persisting ", err)
			return
		}

		seq, _ := res.LastInsertId()

		for _, bolt := range topology.Bolts {
			boltCapacity, _ := strconv.ParseFloat(bolt.Capacity, 64)
			boltExecuteLatency, _ := strconv.ParseFloat(bolt.ExecuteLatency, 64)
			stormDetail := StormLogDetail{
				StormLogId:     seq,
				BoltId:         bolt.BoldId,
				Executors:      bolt.Executors,
				Tasks:          bolt.Tasks,
				Capacity:       boltCapacity,
				ExecuteLatency: boltExecuteLatency,
			}
			query = "insert into go_storm_log_detail (go_storm_log_id, bolt_id, executors, tasks, capacity, execute_latency) " +
				"values (:go_storm_log_id, :bolt_id, :executors, :tasks, :capacity, :execute_latency)"
			_, err := db.NamedExec(query, &stormDetail)
			if err != nil {
				fmt.Println("Error persisting ", err)
			}
		}

		if spout, ok := topologyLag["ZookeeperKafkaSpout"]; ok {
			for topic, topicValue := range spout.SpoutLagResult {
				for partition, topicLag := range topicValue {
					stormLogLagDetail := StormLogLagDetail{
						StormLogId:      seq,
						SpoutId:         "ZookeeperKafkaSpout",
						Topic:           topic,
						Partition:       partition,
						LogHeadOffset:   topicLag.LogHeadOffset,
						CommittedOffset: topicLag.ConsumerCommittedOffset,
						Lag:             topicLag.Lag,
					}

					query = "insert into go_storm_log_lag (go_storm_log_id, spout_id, topic, topic_partition, log_head_offset, consumer_committed_offset, lag) " +
						"values (:go_storm_log_id, :spout_id, :topic, :topic_partition, :log_head_offset, :consumer_committed_offset, :lag)"
					_, err := db.NamedExec(query, &stormLogLagDetail)
					if err != nil {
						fmt.Println("Error persisting ", err)
					}
				}
			}
		}
	}
}

var myClient = &http.Client{Timeout: 10 * time.Second}

func getTopologyId(url string, topologyName string) (string, error) {
	summary := TopologySummary{}
	r, err := myClient.Get(url)
	if err != nil {
		return "", err
	}
	defer r.Body.Close()

	err = json.NewDecoder(r.Body).Decode(&summary)
	for _, topology := range summary.Topologies {
		if topology.Name == topologyName {
			return topology.Id, nil
		}
	}
	return "", err
}

func getTopology(url string) (Topology, error) {
	topology := Topology{}
	r, err := myClient.Get(url)
	if err != nil {
		return topology, err
	}
	defer r.Body.Close()

	err = json.NewDecoder(r.Body).Decode(&topology)
	return topology, err
}

func getTopologyLag(url string) (TopologyLag, error) {
	topologyLag := TopologyLag{}
	r, err := myClient.Get(url)
	if err != nil {
		return topologyLag, err
	}
	defer r.Body.Close()

	err = json.NewDecoder(r.Body).Decode(&topologyLag)
	return topologyLag, err
}
