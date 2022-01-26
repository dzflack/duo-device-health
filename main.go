package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"regexp"
	"strconv"
	"time"
)

// Get env var or a default value
func getEnv(key, defaultVal string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = defaultVal
	}
	return value
}

type healthAppPayload struct {
	HealthCheckStartTimestamp string        `json:"health_check_start_timestamp"`
	IsPasswordSet             bool          `json:"is_password_set"`
	HealthCheckEndTimestamp   string        `json:"health_check_end_timestamp"`
	IsAutoLoginEnabled        bool          `json:"is_auto_login_enabled"`
	Os                        string        `json:"os"`
	IsEncryptionEnabled       bool          `json:"is_encryption_enabled"`
	DeviceID                  string        `json:"device_id"`
	Txid                      string        `json:"txid"`
	DeploymentTrack           string        `json:"deployment_track"`
	IsFirewallEnabled         bool          `json:"is_firewall_enabled"`
	DuoClientVersion          string        `json:"duo_client_version"`
	OsVersion                 string        `json:"os_version"`
	DeviceName                string        `json:"device_name"`
	HealthCheckLengthMillis   float64       `json:"health_check_length_millis"`
	OsBuild                   string        `json:"os_build"`
	CommunicationScheme       string        `json:"communication_scheme"`
	SecurityAgents            []interface{} `json:"security_agents"`
}

func setHeaders(w http.ResponseWriter, r *http.Request) {
	location, _ := time.LoadLocation("GMT")
	theTime := time.Now().In(location)

	w.Header().Add("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.Header().Add("Connection", "keep-alive")
	w.Header().Add("Transfer-Encoding", "chunked")
	w.Header().Add("Access-Control-Allow-Headers", "x-xsrftoken")
	w.Header().Add("Date", theTime.Format("Mon, 02 Jan 2006 15:04:05 MST"))
	w.Header().Add("Cache-Control", "no-store, must-revalidate")
	w.Header().Add("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.WriteHeader(204)
}

func checkAliveHandler(w http.ResponseWriter, r *http.Request) {
	setHeaders(w, r)
}

func generateReportHandler(w http.ResponseWriter, r *http.Request) {
	setHeaders(w, r)

	client := &http.Client{}
	currentTime := time.Now()
	finishedTime := currentTime.Add(time.Second)
	var agents []interface{}

	deviceId := getEnv("DEVICE_ID", "DE13DE13-23F28-57A7-B81E-B81E7AE23F28")

	body := healthAppPayload{
		HealthCheckStartTimestamp: currentTime.Format("2006-01-02 15:04:05"),
		HealthCheckEndTimestamp:   finishedTime.Format("2006-01-02 15:04:05"),
		IsPasswordSet:             true,
		IsAutoLoginEnabled:        false,
		Os:                        "macOS",
		IsEncryptionEnabled:       true,
		DeviceID:                  deviceId,
		Txid:                      r.URL.Query().Get("txid"),
		DeploymentTrack:           "release",
		IsFirewallEnabled:         true,
		DuoClientVersion:          "3.4.0.0",
		OsVersion:                 "11.15.4",
		DeviceName:                "work-mbp",
		HealthCheckLengthMillis:   540.57899999999995,
		OsBuild:                   "20E287",
		CommunicationScheme:       "https",
		SecurityAgents:            agents,
	}

	jsonData := new(bytes.Buffer)
	json.NewEncoder(jsonData).Encode(body)

	re := regexp.MustCompile(`https://2.endpointhealth.duosecurity.com/v1/healthapp/device/health\?_req_trace_group=(.*)\?`)
	results := re.FindSubmatch([]byte(r.URL.Query().Get("eh_service_url")))

	url := "https://2.endpointhealth.duosecurity.com/v1/healthapp/device/health?_req_trace_group=" +
		string(results[1]) + "?_=" +
		strconv.FormatInt(finishedTime.Unix(), 10)

	req, err := http.NewRequest("POST", url, jsonData)
	if err != nil {
		log.Println("Failed to create POST request: ", err)
	}

	req.Host = "2.endpointhealth.duosecurity.com"
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Connection", "close")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("User-Agent", "Duo%20Device%20Health/3.4.0.0 CFNetwork/1125.2 Darwin/20.4.0 (x86_64)")
	req.Header.Set("Accept-Language", "en-au")
	req.Header.Set("Accept-Encoding", "gzip, deflate")

	requestDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		log.Println("Failed to dump POST request before sending it: ", err)
	}

	fmt.Println(string(requestDump), "\n\n ")

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Failed to send POST to health endpoint: ", err)
	}

	defer resp.Body.Close()
}

func main() {
	http.HandleFunc("/report", generateReportHandler)
	http.HandleFunc("/alive", checkAliveHandler)
	tlsCertFile, tlsCertProvided := os.LookupEnv("DUO_LOCAL_TLS_CERT")
	tlsKeyFile, tlsKeyProvided := os.LookupEnv("DUO_LOCAL_TLS_KEY")
	port := getEnv("DUO_LOCAL_PORT", "53106")
	host := fmt.Sprintf("127.0.0.1:%s", port)
	log.Println("Listening on", host)
	if tlsCertProvided != tlsKeyProvided {
		log.Fatalln("Please provide neither or both of DUO_LOCAL_TLS_KEY, DUO_LOCAL_TLS_CERT")
	} else if tlsCertProvided {
		log.Println("Listening for http")
		ret := http.ListenAndServe(host, nil)
		log.Println(ret)
	} else {
		log.Println("Listening for https using", tlsCertFile, tlsKeyFile)
		ret := http.ListenAndServeTLS(host, tlsCertFile, tlsKeyFile, nil)
		log.Println(ret)
	}
}
