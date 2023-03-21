# gohttps

### This code defines a simple HTTP server using Go's net/http package, with two additional routes (/version and /healthz) and a handler for the root route (/). 

### The root handler returns a JSON response with a predefined message and a status code of "Service Unavailable". 

### Additionally, if the HTTP request header contains "Upgrade", the handler upgrades the connection to a WebSocket connection and sends the JSON response over the WebSocket connection. 

### The server runs on port 9000 using TLS.