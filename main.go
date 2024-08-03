package main

import (
	"encoding/json"
	"io"
	"github.com/sirupsen/logrus"
	"net/http"
	"orderDetails/impl"
	"orderDetails/model"
	"os"
	"time"

)
var log = logrus.New()
func init() {
	var err error
	logFile, err := os.OpenFile("execution-logs.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Error("Error opening log file :", err)
	}
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)
	log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	log.SetLevel(logrus.DebugLevel)
	model.Log = log
	log.Debug("Initialization started")
	err = impl.DbConnect()
	if err != nil {
		log.Error("Mongodb connection Failed. Check after Sometime")
	}
	err = impl.ReadCSV()
	if err != nil {
		log.Error("Error in Reading the xlsx file")
	}
}

func Reload(w http.ResponseWriter, r *http.Request) {
	if err := impl.ReadCSV(); err != nil {
		log.Error("Error in Reading xlsx file")
		w.Write([]byte(err.Error()))
	}
	w.Write([]byte("Order Data is Reading it will take sometime to load in db. Please check the logs if you face any issue"))
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func Revenue(w http.ResponseWriter, r *http.Request) {
	var reqBody model.Request
	bytes, _ := io.ReadAll(r.Body)
	err := json.Unmarshal(bytes, &reqBody)
	if err != nil {
		log.Error("Error in unmarshalling ::: ", err)
	}
	response, _ := impl.GetRevenue(reqBody)
	json.NewEncoder(w).Encode(response)

}
func main() {
	log.Debug("Server start Listening ...")
	http.HandleFunc("/healthCheck", HealthCheck)
	http.HandleFunc("/reloadAll", Reload)
	http.HandleFunc("/revenue", Revenue)
	go func() {
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Println("Fail to start server :8080" + err.Error())
		}
	}()
	scheduler := time.NewTicker(1 * time.Hour)
	defer scheduler.Stop()
	for {
		select {
		case <-scheduler.C:
			impl.ReadCSV()
		}
	}
}
