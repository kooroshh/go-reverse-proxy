[![release](https://github.com/kooroshh/go-reverse-proxy/actions/workflows/release.yml/badge.svg)](https://github.com/kooroshh/go-reverse-proxy/actions/workflows/release.yml)
# GO Reverse Proxy 

Go Reverse Proxy is a simple layer 4 tcp reverse proxy that can listen on multiple ports and forward traffic to multiple upstream servers.  

Logging can be toggled using log_level flag in the json config file.  

The Following features will be implemented in the future:  

* Load Balancing  
* Rate Limiting

# Configuration
### Servers 
* Local Address : Listening Address
* Remote Address : Proxy Address

### Users
* Mode : 
    * blacklist
    * whitelist 
    * disabled
* Urls : array of remote \n seperated ip list
* Periodic : True | False (should check remote list periodic ?)
* Interval : uint (Minuets)
### AMPQ
* enable : True | False (should enable ampq ?)
* exchange : Exchange name to bind the queue 
* connection string : AMPQ Connection string

### Node
in case of using AMPQ node name will be used to create UNIQUE queue name to bind on FANOUT exchange.

### Secret
authorization header secret for fetching user db via external api

### Log Level 
log level supports debug and none for now 


# How it works ?
![How it works](/.assets/Whiteboard.png?raw=true "White Board")
