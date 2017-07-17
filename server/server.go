package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/oscarpfernandez/idgen/ssid"
	"github.com/sirupsen/logrus"
)

// SSIDResponse servers reponse
type SSIDResponse struct {
	Config ssid.Config `json:"config"`
	SSIDs  []uint64    `json:"ssids"`
}

// ErrorResponse servers error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// CreateServer creates a server router and initializes the endpoints
func CreateServer() *httprouter.Router {
	router := httprouter.New()
	router.GET("/ssid", GenerateID)
	router.GET("/ssid/:count/:generator", GenerateIDs)

	return router
}

// GenerateID generates a single ID endpoint
func GenerateID(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	logrus.Infof("Generating single SSID")
	startTime := time.Now()
	config := ssid.Config{}
	newInstance, err := ssid.NewSSID(&config)
	if err != nil {
		errResponse := ErrorResponse{
			Error: err.Error(),
		}
		json.NewEncoder(w).Encode(errResponse)
		return
	}

	resultID, err := newInstance.GenerateIDs(1)
	if err != nil {
		errResponse := ErrorResponse{
			Error: err.Error(),
		}
		json.NewEncoder(w).Encode(errResponse)
		return
	}

	response := SSIDResponse{
		Config: config,
		SSIDs:  resultID,
	}

	json.NewEncoder(w).Encode(response)
	elapsed := time.Since(startTime)
	logrus.Infof("Finished in %s", elapsed)
}

// GenerateIDs generates a set of SSIDs provided a `count` and `generator ID`
func GenerateIDs(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	countStr := ps.ByName("count")
	count, err := strconv.ParseUint(countStr, 10, 64)
	if err != nil {
		errResponse := ErrorResponse{
			Error: err.Error(),
		}
		json.NewEncoder(w).Encode(errResponse)
		return
	}

	generatorIDStr := ps.ByName("generator")
	generatorID, err := strconv.ParseUint(generatorIDStr, 10, 64)
	if err != nil {
		errResponse := ErrorResponse{
			Error: err.Error(),
		}
		json.NewEncoder(w).Encode(errResponse)
		return
	}

	logrus.Infof("Generating %d IDs using generatorID %d", count, generatorID)
	startTime := time.Now()

	config := ssid.Config{GeneratorID: uint16(generatorID)}
	ssidInstance, err := ssid.NewSSID(&config)
	if err != nil {
		errResponse := ErrorResponse{
			Error: err.Error(),
		}
		json.NewEncoder(w).Encode(errResponse)
		return
	}

	generatedIDs, err := ssidInstance.GenerateIDs(uint16(count))
	if err != nil {
		errResponse := ErrorResponse{
			Error: err.Error(),
		}
		json.NewEncoder(w).Encode(errResponse)
		return
	}

	response := SSIDResponse{
		Config: config,
		SSIDs:  generatedIDs,
	}

	json.NewEncoder(w).Encode(response)

	elapsed := time.Since(startTime)
	logrus.Infof("Finished in %s", elapsed)
}
