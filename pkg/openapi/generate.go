package openapi

//go:generate ogen --target apprun_dedicated --package apprun_dedicated --clean ../../openapis/apprun-dedicated.json
//go:generate oapi-codegen -generate types,client -package monitoring_suite -o monitoring_suite/client.gen.go ../../openapis/monitoring-suite.json
