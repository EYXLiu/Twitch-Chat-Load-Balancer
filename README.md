# Twitch Chat Load Balancer
**Tech Stack:** Go, Gin, Websockets, Multithreading, Twitch API, Redis, Make  
Scalable multi-threaded Twitch Chat ingestion and processing system to learn load balancing and scaling in a high-throughput streaming environment  
The system uses Twitch chat traffic to simulate real-world load patterns similar to production-scale streaming architectures  
It has been stress-tested with ~30,000 concurrent users along with continuous streams of gifted subscriptions and bits  

## Features
- Auto load balancing architecture that spins up newly threaded instances as required
- Producer Consumer event pipeline for asynchronous processing of messages
- Real-time Twitch Chat Integration using the official Websocket connection
- Redis-backed sliding window cache
- Concurrent worker processing using goroutines and channels
- Fault tolerant event distribution where workers are pinged to ensure constant uptime 

## How to run
`export TWITCH_CHANNELS={channel1},{channel2},{channel3},...`  
`export REDIS={redis endpoint}` (normally localhost:6379)  
`make`

## Redis (make clean)
`brew services start redis`  
`redis-cli lrange twitch:messages 0 -1`  
`redis-cli del twitch:messages`  
`redis-cli xgroup destroy twitch:events workers`  
`redis-cli del twitch:events`  
`redis-cli keys '*'`  
`brew services stop redis`  

## Links
[twitch websocket](https://dev.twitch.tv/docs/eventsub/handling-websocket-events)  
