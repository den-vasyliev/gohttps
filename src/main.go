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
	w.Write([]byte("gohttps v1.0.0")) //should be in var
}

func healthz(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}

func main() {

	http.HandleFunc("/version", version)
	http.HandleFunc("/healthz", healthz)

	// handle `/` route
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// should be moved to var from conf
		//respM := RespM{"WARN_UNAVAILABLE_MAINTENANCE", "Maintenance", "The service is currently under regular maintenance, please, try again later"}
		respM := RespD{"WARN_UNAVAILABLE_DEPLOYMENT", "Deployment", "The service is currently being redeployed and will become available shortly"}

		js, err := json.Marshal(respM)
		if err != nil {
			log.Println("json:", err)
		}

		wss := r.Header.Get("Upgrade")
		if wss != "" {
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
				http.Error(w, "{'json-code':'json-payload'}", http.StatusServiceUnavailable)
				err = c.WriteMessage(mt, js)
				c.Close()
				if err != nil {
					log.Println("write:", err)
					break
				}
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		io.WriteString(w, string(js))

	})

	// run server on port "9000"
	// should be moved to var
	log.Fatal(http.ListenAndServeTLS(":9000", "/tmp/tls/tls.crt", "/tmp/tls/tls.key", nil))

}
