package main

type Config struct {
	Servers []struct {
		LocalAddress  string `json:"local_address"`
		RemoteAddress string `json:"remote_address"`
	} `json:"servers"`
	LogLevel string `json:"log_level"`
}
