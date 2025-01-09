run:
	go run .

load-data:
	go run ./cmd/loader

get-cassandra:
	docker pull cassandra

set-cassandra:
	docker network create cassandra-network

run-cassandra:
	docker run --name cassandra1 --network cassandra-network -d cassandra

run-cassandra2:
	docker run --name cassandra2 --network cassandra-network -e CASSANDRA_SEEDS=cassandra1 -d cassandra

run-cqlsh:
	docker run -it --network cassandra-network --rm cassandra cqlsh cassandra1

inspect-ip:
	docker inspect --format='{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' cassandra1

cleanup:
	docker kill cassandra1; docker rm cassandra1 &
	docker kill cassandra2; docker rm cassandra2 &
	docker network rm cassandra-network
