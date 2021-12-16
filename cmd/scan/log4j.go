package scan

import (
	_ "embed"
	"encoding/json"
	"github.com/whitesource/log4j-detect/records"
)

type Vulnerability struct {
	CveId         string
	Summary       string
	FixResolution string
}

var fixes = map[string]map[string]string{
	"org.apache.logging.log4j": {
		"log4j-core": "Upgrade to version org.apache.logging.log4j:log4j-core:2.16.0",
	},
	"org.ops4j.pax.logging": {
		"pax-logging-log4j2": "Upgrade to version org.ops4j.pax.logging:pax-logging-log4j2:2.0.11",
	},
}

//go:embed assets/log4j_libs.json
var log4jLibsStr []byte
var log4jLibs []records.VulnerableLib

func init() {
	err := json.Unmarshal(log4jLibsStr, &log4jLibs)
	if err != nil {
		panic(err)
	}
}
