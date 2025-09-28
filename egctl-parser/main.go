package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

type EnvoyOutput struct {
	EnvoyGatewaySystem EnvoyGatewaySystem `json:"envoy-gateway-system"`
}

type EnvoyGatewaySystem map[string]EnvoyPod

type EnvoyPod struct {
	DynamicEndpointConfigs []DynamicEndpointConfig `json:"dynamicEndpointConfigs"`
}

type DynamicEndpointConfig struct {
	EndpointConfig EndpointConfig `json:"endpointConfig"`
}

type EndpointConfig struct {
	ClusterName string      `json:"clusterName"`
	Endpoints   []Endpoints `json:"endpoints"`
}

type Endpoints struct {
	LbEndpoints []LbEndpoint `json:"lbEndpoints"`
}

type LbEndpoint struct {
	Endpoints Endpoint `json:"endpoint"`
}

type Endpoint struct {
	Address Address `json:"address"`
}

type Address struct {
	SocketAddress SocketAddress `json:"socketAddress"`
}

type SocketAddress struct {
	Address   string `json:"address"`
	PortValue int    `json:"portValue"`
}

type EndpointInfo struct {
	HttpRouteName string
	Address       string
}

func fileToString() string {
	filename := os.Args[1]
	rawData, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	jsonString := string(rawData)
	return jsonString
}

func stdinToString() string {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	stringData := string(data)
	return stringData
}

func endpointData(s string) []map[string]EndpointInfo {
	var endpointSlice []map[string]EndpointInfo
	var envoyOutput EnvoyOutput

	if err := json.Unmarshal([]byte(s), &envoyOutput); err != nil {
		panic(err)
	}

	for envoyPod, value := range envoyOutput.EnvoyGatewaySystem {
		for _, cfg := range value.DynamicEndpointConfigs {
			for _, ep := range cfg.EndpointConfig.Endpoints {
				for _, lb := range ep.LbEndpoints {
					endpointInfo := EndpointInfo{
						HttpRouteName: cfg.EndpointConfig.ClusterName,
						Address:       lb.Endpoints.Address.SocketAddress.Address,
					}
					endpointMap := map[string]EndpointInfo{
						envoyPod: endpointInfo,
					}
					endpointSlice = append(endpointSlice, endpointMap)
				}
			}
		}
	}
	return endpointSlice
}

func main() {
	stat, _ := os.Stdin.Stat()
	var rawJson string

	// Check if filename has been added as an argument or input comes from stdin (piped)
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		rawJson = stdinToString()
	} else if len(os.Args) == 2 {
		rawJson = fileToString()
	} else if len(os.Args) > 2 {
		log.Fatal("Too many arguments, only 1 supported.")
	} else {
		log.Fatal("No stdin or arguments provided.")
	}

	endpointDataSlice := endpointData(rawJson)
	for _, endpointMap := range endpointDataSlice {
		for envoyPod, endpointInfo := range endpointMap {
			fmt.Println("Envoy-Pod: ", envoyPod)
			fmt.Println("  Httproute-Name: ", endpointInfo.HttpRouteName)
			fmt.Println("  Endpoint-Address: ", endpointInfo.Address)
		}
	}
}
