package main

import (
	"encoding/json"
	"fmt"
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

var endpointSlice []map[string]EndpointInfo

func main() {
	filename := os.Args[1]
	rawData, err := os.ReadFile(filename)
	rawJson := string(rawData)
	if err != nil {
		log.Fatal(err)
	}
	var envoyOutput EnvoyOutput
	if err := json.Unmarshal([]byte(rawJson), &envoyOutput); err != nil {
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
	for _, endpointMap := range endpointSlice {
		for envoyPod, endpointInfo := range endpointMap {
			fmt.Println("Envoy-Pod: ", envoyPod)
			fmt.Println("  Httproute-Name: ", endpointInfo.HttpRouteName)
			fmt.Println("  Endpoint-Address: ", endpointInfo.Address)
		}
	}
}
