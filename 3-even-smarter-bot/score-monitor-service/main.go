package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"contrib.go.opencensus.io/exporter/stackdriver"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"

	"score-monitor-service/shared"
)

var (
	score    = stats.Int64("score", "The current score of the bot", stats.UnitDimensionless)
	botId, _ = tag.NewKey("botId")
)

func init() {
	v := &view.View{
		Name:        "score_history",
		Measure:     score,
		Description: "Bot scores ploted as a time series",
		TagKeys:     []tag.Key{botId},
		Aggregation: view.LastValue(),
	}
	if err := view.Register(v); err != nil {
		log.Fatalf("Failed to register the view: %v", err)
	}
}

func main() {
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
	// log.Printf("event data: byte[] %v", pubsubMessageEvent.Message.Data)
	// log.Printf("event data: string %v", string(pubsubMessageEvent.Message.Data))

	err := json.Unmarshal(pubsubMessageEvent.Message.Data, &arenaUpdate)
	if err != nil {
		log.Fatalf("fatal error parsing ArenaUpdate: %v", err)
	}
	exporter, err := stackdriver.NewExporter(stackdriver.Options{})
	if err != nil {
		log.Fatal(err)
	}
	// Flush must be called before main() exits to ensure metrics are recorded.
	defer exporter.Flush()

	if err := exporter.StartMetricsExporter(); err != nil {
		log.Fatalf("Error starting metric exporter: %v", err)
	}
	defer exporter.StopMetricsExporter()
	// now populate squares and leaderboard with players
	for k, v := range arenaUpdate.Arena.State {
		ctx, err := tag.New(context.Background(), tag.Insert(botId, k))
		if err != nil {
			log.Fatal(err)
		}
		stats.Record(ctx, score.M(int64(v.Score)))
	}
}
