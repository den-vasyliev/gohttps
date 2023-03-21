package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

//RespM message
type RespM struct {
	Code    string
	Message string
	Details string
}

//RespD message
type RespD struct {
	Code    string
	Message string
	Details string
}

var upgrader = websocket.Upgrader{} // use default options

func version(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("gohttps v1.0.0")) // should be moved to var from conf
}

func healthz(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok")) // should be moved to var from conf
}

func main() {

	// Define handlers for the "/version" and "/healthz" routes
	http.HandleFunc("/version", version)
	http.HandleFunc("/healthz", healthz)

	// Handle the root route "/"
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Define a message to be returned in the response
		respM := RespM{"WARN_UNAVAILABLE_MAINTENANCE", "Maintenance", "The service is currently under regular maintenance, please, try again later"}
		//respD := RespD{"WARN_UNAVAILABLE_DEPLOYMENT", "Deployment", "The service is currently being redeployed and will become available shortly"}

		// Convert the message to JSON
		js, err := json.Marshal(respM)
		if err != nil {
			log.Println("json:", err)
		}

		// Check if the request header contains "Upgrade"
		wss := r.Header.Get("Upgrade")
		if wss != "" {
			// Upgrade the HTTP connection to a WebSocket connection
			upgrader.CheckOrigin = func(r *http.Request) bool { return true }
			c, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				log.Print("upgrade:", err)
				return
			}
			defer c.Close()
			for {
				mt, message, err := c.ReadMessage()
				if err != nil {
					log.Println("read:", err)
					break
				}
				log.Printf("recv: %s %s", message)
				//http.Error(w, "{'json-code':'json-payload'}", http.StatusServiceUnavailable)
				err = c.WriteMessage(mt, js)
				c.Close()
				if err != nil {
					log.Println("write:", err)
					break
				}
			}
		}

		// Set the response header to indicate that the response is in JSON format
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable) // Set the response status code to "Service Unavailable"
		io.WriteString(w, string(js)) // Write the JSON message to the response body

	})

	// Run the server on port 9000 using TLS
	log.Fatal(http.ListenAndServeTLS(":9000", "/tmp/tls/tls.crt", "/tmp/tls/tls.key", nil))

}
