package main

import (
	"context"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

func UpdateUserDB() {
	for {
		userDB.Flush()
		tag := "[UserDB]"
		httpClient := http.Client{}
		for _, url := range conf.Users.Urls {
			ctx, timeout := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
			defer timeout()
			req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
			if err != nil {
				Log(LOG_ERROR, tag, "Unable to make request", err)
				continue
			}
			req.Header.Add("Authorization", conf.Secret)
			res, err := httpClient.Do(req)
			if err != nil {
				Log(LOG_ERROR, tag, "Unable to get response")
				continue
			}
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				Log(LOG_ERROR, tag, "Unable to read response", err)
				continue
			}
			records := strings.Split(string(body), "\n")
			for _, row := range records {
				if net.ParseIP(row) != nil {
					userDB.Add(row, true, -1)
				}
			}
			Log(LOG_INFO, tag, "Source", url, "Total Count=", len(records))
		}
		if !conf.Users.Periodic {
			break
		}
		time.Sleep(time.Duration(conf.Users.Interval) * time.Minute)
	}
}
