package scan

import (
	"fmt"
	"github.com/go-logr/logr"
	"github.com/spf13/cobra"
	"github.com/whitesource/log4j-detect/cmd/clioptions"
	"github.com/whitesource/log4j-detect/cmd/clioptions/settings"
	cmdutil "github.com/whitesource/log4j-detect/cmd/util"
	"github.com/whitesource/log4j-detect/operations"
	"github.com/whitesource/log4j-detect/records"
	"github.com/whitesource/log4j-detect/screening"
	"github.com/whitesource/log4j-detect/supplements"
	"github.com/whitesource/log4j-detect/utils"
	"github.com/whitesource/log4j-detect/utils/exec"
	"path/filepath"
	"strings"
)

type Options struct {
	Directory string
	commander exec.Commander
	*settings.Settings
	settings.Flags
	logr.Logger
	clioptions.IOStreams
}

func NewCmdScan(logger logr.Logger, streams clioptions.IOStreams) *cobra.Command {
	o := &Options{
		commander: exec.New(),
		Logger:    logger.WithName("scan"),
		IOStreams: streams,
	}

	cmd := &cobra.Command{
		Use:   "scan",
		Short: "Discover projects that contain the log4j exploit",
		Run: func(_ *cobra.Command, _ []string) {
			cmdutil.CheckErr(o.Complete())
			cmdutil.CheckErr(o.Run())
		},
	}

	cmd.Flags().StringVarP(&o.Directory, "directory", "d", ".", "the directory to scan")

	return cmd
}

// Complete verifies command line arguments and loads data from the command environment
func (o *Options) Complete() (err error) {
	if o.Settings, err = o.ToSettings(o.Logger); err != nil {
		return err
	}

	if abs, err := filepath.Abs(o.Directory); err == nil {
		o.Directory = abs
	}

	return nil
}

func (o *Options) Run() error {
	_, _ = fmt.Fprintf(o.Out, "Scanning %s for vulnerabilities...\n", o.Directory)

	detected, err := screening.ScreenDirectory(o.Logger, o.Directory, o.Settings.Resolvers.ManifestQueries(), o.Settings.GlobalExcludes())
	if err != nil {
		return err
	}

	if len(detected) == 0 {
		return fmt.Errorf("no maven/gradle projects or jar files were detected in %s\n", o.Directory)
	}

	operationResults := operations.Perform(o.Logger, detected, o.Settings.Resolvers.Surgeons(o.Logger, o.commander))
	enhancedResults := supplements.Supplement(o.Logger, operationResults)
	cves := o.addVulnerabilities(enhancedResults, cve2Sha2Lib)

	_, _ = fmt.Fprintln(o.Out)

	if len(cves) > 0 {
		o.displayVulnerabilities(enhancedResults)
		_, _ = fmt.Fprintf(o.Out, `
One or more of your projects contain the %s exploit.
%s
Learn more about the vulnerability and it's remediation:
%s

©WhiteSource
`,
			utils.MakeOrangeText(fmt.Sprintf("log4j %s", strings.Join(cves, "/"))),
			o.generateRemediationSteps(enhancedResults),
			utils.MakeBlueText("https://www.whitesourcesoftware.com/resources/blog/log4j-vulnerability-cveFiles-2021-44228/"))
	} else {
		_, _ = fmt.Fprintln(o.Out, utils.MakeGreenText("No vulnerabilities were detected"))
	}

	return nil
}

func (o *Options) addVulnerabilities(results []records.EnhancedResult, cve2Sha12Lib CveToSha1ToLib) []string {
	cveMap := map[string]bool{}
	for i := range results {
		r := &results[i]
		r.DepId2VulnerableLib = map[records.Id]records.VulnerableLib{}
		for id, dep := range *r.Deps {
			for cve, sha12Lib := range cve2Sha12Lib {
				if lib, ok := sha12Lib[dep.Sha1]; ok {
					r.DepId2VulnerableLib[id] = lib
					cveMap[cve] = true
				}
			}
		}
	}
	var result []string
	for cve := range cveMap {
		result = append(result, cve)
	}
	return result
}

func (o *Options) displayVulnerabilities(results []records.EnhancedResult) {
	for _, r := range results {
		if len(r.DepId2VulnerableLib) == 0 || r.LType == records.LTFs {
			continue
		}

		_, _ = fmt.Fprintf(o.Out, "%s %s\n", utils.MakeBlueText("Manifest:"), r.ManifestFile)
		_, _ = fmt.Fprintln(o.Out, utils.MakeBlueText("Vulnerable Paths:"))
		r.Print(o.Out)
		_, _ = fmt.Fprintf(o.Out, "\n")
	}

	for _, r := range results {
		if r.LType == records.LTFs {
			if len(r.DepId2VulnerableLib) == 0 {
				break
			}

			_, _ = fmt.Fprintln(o.Out, utils.MakeBlueText("Vulnerable Files: "))
			for id := range r.DepId2VulnerableLib {
				path := (*r.Libraries)[id].SystemPath
				if abs, err := filepath.Abs(path); err == nil {
					path = abs
				}
				_, _ = fmt.Fprintf(o.Out, "\t%s %s\n", path, utils.MakeRedText("Vulnerable"))
			}
			break
		}
	}
}

func (o *Options) generateRemediationSteps(results []records.EnhancedResult) string {
	set := map[string]bool{}
	for _, r := range results {
		for _, v := range r.DepId2VulnerableLib {
			if artifact2Fix, found := fixes[v.GroupId]; found {
				if fix, fixFound := artifact2Fix[v.ArtifactId]; fixFound {
					set[fix] = true
				}
			}
		}
	}

	if len(set) == 0 {
		return ""
	}

	var steps strings.Builder
	steps.WriteString("\nRemediation Steps:\n")
	for fix := range set {
		steps.WriteString(fmt.Sprintf("\t* %s\n", fix))
	}

	return steps.String()
}
