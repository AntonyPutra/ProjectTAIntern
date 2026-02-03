package utils

import (
	"encoding/json"
	"log"
	"strconv"
)

// LogError logs error messages with specific format
func LogError(svcName, refNo, methodName string, err error, data ...string) {
	str := svcName
	if refNo != "" {
		str += " , referenceNo :: " + refNo
	}
	if methodName != "" {
		str += " , error " + methodName
	}
	if err != nil {
		str += " :: " + err.Error()
	}
	for i, s := range data {
		str += ",\n Data " + strconv.Itoa(i+1) + " ::" + s
	}
	log.Println(str)
}

// LogInfo logs information messages with specific format
func LogInfo(svcName, refNo, methodName string, data ...string) {
	str := svcName
	if refNo != "" {
		str += " , referenceNo :: " + refNo
	}
	if methodName != "" {
		str += " , " + methodName
	}
	for i, s := range data {
		str += ",\n Data " + strconv.Itoa(i+1) + " ::" + s
	}
	log.Println(str)
}

// LogActivity logs the start and end of an API interaction (Middleware style)
func LogActivity(endpoint string, headers interface{}, body interface{}, response interface{}) {
	log.Println("[Start]")
	log.Printf("EndPoint : %s\n", endpoint)
	log.Printf("Header : %v\n", headers)

	// Format Body
	bodyStr := ""
	if b, ok := body.(string); ok {
		bodyStr = b
	} else {
		bJSON, _ := json.MarshalIndent(body, "", "  ")
		bodyStr = string(bJSON)
	}
	log.Printf("Body : %s\n", bodyStr)

	// Format Response
	respStr := ""
	if r, ok := response.(string); ok {
		respStr = r
	} else {
		rJSON, _ := json.Marshal(response)
		respStr = string(rJSON)
	}
	log.Printf("Response : %s\n", respStr)

	log.Println("[End]")
}
