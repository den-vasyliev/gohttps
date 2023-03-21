# gohttps

### This code defines a simple HTTP server using Go's net/http package, with two additional routes (/version and /healthz) and a handler for the root route (/). 

### The root handler returns a JSON response with a predefined message and a status code of "Service Unavailable". 

### Additionally, if the HTTP request header contains "Upgrade", the handler upgrades the connection to a WebSocket connection and sends the JSON response over the WebSocket connection. 

### The server runs on port 9000 using TLS.

## Change log

Here are the changes I made to the code:

I moved constants and variables to the top of the file and gave them descriptive names.

I added an error handler to the jsonResponse function to prevent the program from crashing if there is an error.

I created a handleDefault function to handle the default route. This function returns a maintenance response if the client is requesting an upgrade, and returns a 503 error otherwise.

I updated the main function to use the handleDefault function to handle the default route.

Overall, these changes should make the code more readable and easier to maintain.