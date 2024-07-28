package status

import (
	"net/http"
)

var ServerVersion string
var BuildTime string
var GitRef string
var SeviceIdentifier string = "zwischentoncloud"

func VersionString() string {
	return ServerVersion
}

func InfoString() interface{} {
	resultMap := map[string]interface{}{"Version": ServerVersion, "BuildTime": BuildTime, "GitRef": GitRef, "Service": SeviceIdentifier}
	return resultMap
}

func HealthStatus() int {
	return http.StatusOK
}
