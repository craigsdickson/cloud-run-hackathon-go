package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"

	"leaderboard-service/shared"

	"github.com/gomodule/redigo/redis"
)

var redisPool *redis.Pool

func init() {

}

func main() {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)

	const maxConnections = 10
	redisPool = &redis.Pool{
		MaxIdle: maxConnections,
		Dial:    func() (redis.Conn, error) { return redis.Dial("tcp", redisAddr) },
	}

	http.HandleFunc("/", EventProcessor)
	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	// Start HTTP server.
	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func EventProcessor(w http.ResponseWriter, req *http.Request) {
	var arenaUpdate shared.ArenaUpdate
	var pubsubMessageEvent shared.PubSubMessageEvent
	defer req.Body.Close()
	d := json.NewDecoder(req.Body)
	// d.DisallowUnknownFields()
	if err := d.Decode(&pubsubMessageEvent); err != nil {
		log.Printf("WARN: failed to decode PubSubMessageEvent in response body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err := json.Unmarshal(pubsubMessageEvent.Message.Data, &arenaUpdate)
	if err != nil {
		log.Fatalf("fatal error parsing ArenaUpdate: %v", err)
	}
	var leaderboard []shared.PlayerState
	for k, v := range arenaUpdate.Arena.State {
		v.Id = k
		leaderboard = append(leaderboard, v)
	}
	// now sort the leaderboard
	log.Printf("unsorted leaderboard is: %v", leaderboard)
	sort.Slice(leaderboard, func(i, j int) bool {
		return leaderboard[i].Score > leaderboard[j].Score
	})
	log.Printf("sorted leaderboard is: %v", leaderboard)
	conn := redisPool.Get()
	defer conn.Close()
	leaderboardAsByteArray, err := json.Marshal(leaderboard)
	if err != nil {
		log.Fatalf("fatal error marshalling leaderboard: %v", err)
	}
	log.Printf("updating leaderboard in redis to: %v", string(leaderboardAsByteArray))
	conn.Do("SET", "leaderboard", string(leaderboardAsByteArray))
}
