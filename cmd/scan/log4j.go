package scan

import (
	"embed"
	_ "embed"
	"encoding/json"
	"github.com/whitesource/log4j-detect/records"
	"io/fs"
	"path/filepath"
	"strings"
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

//go:embed cve/*.json
var cveFiles embed.FS

var cve2Lib []records.VulnerableLib

func init() {
	files, err := fs.ReadDir(cveFiles, "cve")
	if err != nil {
		panic(err)
	}

	for _, f := range files {
		var libs []records.VulnerableLib
		data, err := cveFiles.ReadFile(filepath.Join("cve", f.Name()))
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(data, &libs)
		if err != nil {
			panic(err)
		}
		cve := strings.ReplaceAll(f.Name(), ".json", "")
		for i := range libs {
			libs[i].CVE = cve
		}
		cve2Lib = append(cve2Lib, libs...)
	}
}
