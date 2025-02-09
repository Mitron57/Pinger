# Pinger
Microservice for pinging Docker containers

# Features
 - Polls every container on the host once per tick (specify period in seconds in [config](config/config.yaml))
 - Sends data about polled IPs to backend (specify link to PUT endpoint in [config](config/config.yaml))
 - Requires access to /var/run/docker.sock and NET_RAW capability to operate

# Standalone launch
1. Specify period of polling and API link to backend in config.yaml file in [config](config) or provide by yourself via -c flag
2. Build image: ```docker build -t monito-pinger-img```
3. Run container and provide your host network: ```docker run --name pinger --network host --cap-add NET_RAW -v /var/run/docker.sock:/var/run/docker.sock -d monito-pinger-img```
