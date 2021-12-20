package scan

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/whitesource/log4j-detect/records"
)

type Vulnerability struct {
	CveId         string
	Summary       string
	FixResolution string
}

var fixes = map[string]map[string]string{
	"org.apache.logging.log4j": {
		"log4j-core": "Upgrade to version org.apache.logging.log4j:log4j-core:2.17.0",
	},
	"org.ops4j.pax.logging": {
		"pax-logging-log4j2": "Upgrade to version org.ops4j.pax.logging:pax-logging-log4j2:2.0.12",
	},
}

//go:embed cve/libs.json
var cveFiles []byte

var cve2Lib []records.VulnerableLib

func init() {
	err := json.Unmarshal(cveFiles, &cve2Lib)
	if err != nil {
		panic(fmt.Sprintf("failed to unmarshal libraries: %v", err))
	}
}
