package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

func getHome(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	response.Write([]byte(`{ "message": "Welcome! To use this system please query /ip/$ip" }`))
}

func getIp(response http.ResponseWriter, request *http.Request) {
	if !validApiKey(request, false) {
		writeError(response, "Sorry, this API requires a key")
		return
	}

	ipString := request.PathValue("ip")
	jsonBytes, err := fetchIPJson(ipString)
	if err != nil {
		writeError(response, err.Error())
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.Write(jsonBytes)
}

func getRandomIp(response http.ResponseWriter, request *http.Request) {
	if !validApiKey(request, true) {
		writeError(response, "Sorry, this API requires a key")
		return
	}

	var ipString string
	ipVersion := request.PathValue("ipVersion")
	if ipVersion == "4" {
		ipString = randomIpv4()
	} else {
		ipString = randomIpv6();
	}

	jsonBytes, err := fetchIPJson(ipString)
	if err != nil {
		writeError(response, err.Error())
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.Write(jsonBytes)
}

func getBenchmark(response http.ResponseWriter, request *http.Request) {
	if !validApiKey(request, true) {
		writeError(response, "Sorry, this API requires a key")
		return
	}

	var ipString string

	ipVersion	:= request.PathValue("ipVersion")
	times		:= request.PathValue("times")

	timesInt, err := strconv.Atoi(times)
	if err != nil || timesInt < 1 {
		writeError(response, "URL must contain a positive numeric number of times to run")
		return
	}

	// Don't want this to be part of the benchmark
	var testIps []string
	for i := 0; i < timesInt; i++ {
		if ipVersion == "4" {
			ipString = randomIpv4()
		} else {
			ipString = randomIpv6();
		}

		testIps = append(testIps, ipString)
	}

	start := time.Now()
	for _, ipString := range testIps {
		_, err := fetchIP(ipString)
		if err != nil {
			writeError(response, "Error encountered during run (" + ipString + ")")
			return
		}
	}

	ms := int(time.Now().Sub(start).Milliseconds())
	us := int(time.Now().Sub(start).Microseconds())

	msPerOp := ms / timesInt;
	usPerOp := us / timesInt;

	response.Header().Set("Content-Type", "application/json")
	response.Write([]byte(`{
		"times": `		+ strconv.Itoa(timesInt) + `, 
		"ms": `			+ strconv.Itoa(ms) + `, 
		"μs": `			+ strconv.Itoa(us) + `, 
		"ms_per_op": `	+ strconv.Itoa(msPerOp) + `, 
		"μs_per_op": `	+ strconv.Itoa(usPerOp) + ` 
	}`))
}

func writeError(response http.ResponseWriter, message string) {
	jsonBytes, _ := json.Marshal(map[string]string{ "error": message })

	response.Header().Set("Content-Type", "application/json")
	response.Write(jsonBytes)
}