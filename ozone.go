package main

import (
	"regexp"
)

var basicRE = regexp.MustCompile(`(?:CI / )?(?:basic|unit) \(([^)]+)\)`)
var integrationRE = regexp.MustCompile(`(?:CI / )?integration \(([^)]+)\)`)
var matrixRE = regexp.MustCompile(`(?:CI / )?(\w+) \(([^)]+)\)`)

func JobToArtifactName(job string) string {
	if basicRE.MatchString(job) {
		return basicRE.ReplaceAllString(job, "$1")
	}
	if integrationRE.MatchString(job) {
		return integrationRE.ReplaceAllString(job, "it-$1")
	}
	if matrixRE.MatchString(job) {
		return matrixRE.ReplaceAllString(job, "$1-$2")
	}
	return job
}
