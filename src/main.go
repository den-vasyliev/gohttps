package main

import (
    "encoding/json"
    "io"
    "log"
    "net/http"

    "github.com/gorilla/websocket"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
    versionPath         = "/version"
    healthzPath         = "/healthz"
    defaultPath         = "/"
    defaultPort         = ":9000"
    certFile            = "/tmp/tls/tls.crt"
    privateKey          = "/tmp/tls/tls.key"
    maintenanceMessage  = "The service is currently under regular maintenance, please try again later"
    maintenanceCode     = "WARN_UNAVAILABLE_MAINTENANCE"
    namespace           = "gohttps"
    subsystem           = "http"
    maintenanceCodeDesc = "Indicates if the service is unavailable due to maintenance"
)

var (
    upgrader = websocket.Upgrader{} // use default options
    maintenanceResponse = jsonResponse(RespM{
        Code:    maintenanceCode,
        Message: "Maintenance",
        Details: maintenanceMessage,
    })

    maintenanceCodeMetric = prometheus.NewGauge(prometheus.GaugeOpts{
        Namespace: namespace,
        Subsystem: subsystem,
        Name:      "maintenance_code",
        Help:      maintenanceCodeDesc,
    })

	// Define a Prometheus counter to track the number of requests made to the service
    requestCounter = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of requests received",
        },
        []string{"method", "path", "status"},
    )

    // Define a Prometheus gauge to track the number of active connections to the websocket endpoint
    activeConnectionsGauge = prometheus.NewGauge(
        prometheus.GaugeOpts{
            Name: "active_connections",
            Help: "Number of active websocket connections",
        },
    )
)

//RespM message
type RespM struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    Details string `json:"details"`
}

func init() {
    prometheus.MustRegister(maintenanceCodeMetric)
}

func jsonResponse(resp interface{}) []byte {
    js, err := json.Marshal(resp)
    if err != nil {
        log.Println("json:", err)
        return []byte{}
    }
    return js
}
// Handler for /version endpoint

func version(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("gohttps v1.0.0"))
}
// Handler for /healthz endpoint

func healthz(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("ok"))
}
// Handler for default endpoint

func handleDefault(w http.ResponseWriter, r *http.Request) {
    wss := r.Header.Get("Upgrade")
    if wss != "" {
        upgrader.CheckOrigin = func(r *http.Request) bool { return true }
        c, err := upgrader.Upgrade(w, r, nil)
        if err != nil {
            log.Print("upgrade:", err)
            return
        }
        defer c.Close()

		        // Increase the active connections gauge and register it with Prometheus
				activeConnectionsGauge.Inc()
				prometheus.Register(activeConnectionsGauge)
				
        for {
            mt, message, err := c.ReadMessage()
            if err != nil {
                log.Println("read:", err)
                break
            }
            log.Printf("recv: %s %s", message)
            err = c.WriteMessage(mt, maintenanceResponse)
            if err != nil {
                log.Println("write:", err)
                break
            }
        }
        maintenanceCodeMetric.Set(1)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusServiceUnavailable)
    io.WriteString(w, string(maintenanceResponse))
    maintenanceCodeMetric.Set(1)
}

func main() {
    http.HandleFunc(versionPath, version)
    http.HandleFunc(healthzPath, healthz)
    http.HandleFunc(defaultPath, handleDefault)
    http.Handle("/metrics", promhttp.Handler())

    log.Fatal(http.ListenAndServeTLS(defaultPort, certFile, privateKey, nil))
}
