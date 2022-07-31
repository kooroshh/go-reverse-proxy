package main

import (
	"encoding/json"
	"net/http"
)

type Stats struct {
	TotalConnections    uint64 `json:"total_connections"`
	AcceptedConnections uint64 `json:"accepted_connections"`
	RejectedConnections uint64 `json:"rejected_connections"`
}

func StartStatsServer(address string) {
	tag := "[STATS]"
	http.HandleFunc("/stats", func(writer http.ResponseWriter, request *http.Request) {
		Log(LOG_DEBUG, tag, "New Stats Request Received")
		writer.Header().Set("Content-Type", "application/json")
		response := map[string]interface{}{}
		for _, p := range ProxyList {
			response[p.ListenAddr] = p.Stats
		}
		res, err := json.Marshal(response)
		if err != nil {
			writer.Write([]byte("{\"status\" : \"error\"}"))
		} else {
			writer.Write(res)
		}
	})
	Log(LOG_INFO, tag, "starting stats server")
	err := http.ListenAndServe(address, nil)
	if err != nil {
		panic(err)
	}
}
