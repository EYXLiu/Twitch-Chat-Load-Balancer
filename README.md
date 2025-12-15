# Twitch Chat Load Balancer

## How to run
`export TWITCH_CHANNEL={channel}`  
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
