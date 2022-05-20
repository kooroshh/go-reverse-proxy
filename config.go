package main

type Config struct {
	Servers []struct {
		LocalAddress  string `json:"local_address"`
		RemoteAddress string `json:"remote_address"`
	} `json:"servers"`
	Ampq struct {
		Enable           bool   `json:"enable"`
		Exchange         string `json:"exchange"`
		ConnectionString string `json:"connection_string"`
	} `json:"ampq"`
	NodeName string `json:"node"`
	LogLevel string `json:"log_level"`
}
