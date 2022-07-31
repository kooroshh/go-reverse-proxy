package main

type Config struct {
	Servers []struct {
		LocalAddress  string `json:"local_address"`
		RemoteAddress string `json:"remote_address"`
	} `json:"servers"`
	Users struct {
		Mode     string   `json:"mode"`
		Urls     []string `json:"urls"`
		Periodic bool     `json:"periodic"`
		Interval int      `json:"interval"`
	} `json:"users"`
	Ampq struct {
		Enable           bool   `json:"enable"`
		Exchange         string `json:"exchange"`
		ConnectionString string `json:"connection_string"`
	} `json:"ampq"`
	Node     string `json:"node"`
	Secret   string `json:"secret"`
	LogLevel string `json:"log_level"`
	Stats    struct {
		ListenAddress string `json:"listen_address"`
		Enable        bool   `json:"enable"`
	} `json:"stats"`
}
