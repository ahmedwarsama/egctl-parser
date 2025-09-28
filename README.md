# Description
A script that parses egctl output and only prints relevant info in human readable format.

# Depndencies
In order to run this script you would first need the egctl cli tool and access to a kubernetes cluster that has deployed envoy-gateway. 
Follow this guide to install egctl https://gateway.envoyproxy.io/docs/install/install-egctl/. 
The only egctl output that this script can parse atm is `egctl config envoy-proxy endpoint`.

# Build the binary
cd egctl-parser
go build -o egctl-parser main.go

# Run the script
When running this script you would need to pipe the output of egctl to the egctl-parser script like so:
```
egctl config envoy-proxy endpoint | egctl-parser
```

# TODO
- Add support to choose what httproute to print.
- Fix a better looking output.
