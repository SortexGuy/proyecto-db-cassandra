package config

import (
	"log"
	"os"

	"github.com/gocql/gocql"
)

var SESSION *gocql.Session

func GetClusterConfig(root bool) *gocql.ClusterConfig {
	cass_ip := os.Getenv("CASSANDRA_IPADDRESS")
	log.Println("Trying to connect to container at ", cass_ip)
	cluster := gocql.NewCluster(cass_ip)
	if !root {
		cluster.Keyspace = "app"
	}
	cluster.Consistency = gocql.Quorum
	return cluster
}
