package main

import (
	"regexp"
)

var basicRE = regexp.MustCompile(`basic \(([^)]+)\)`)
var integrationRE = regexp.MustCompile(`integration \(([^)]+)\)`)
var matrixRE = regexp.MustCompile(`(\w+) \(([^)]+)\)`)

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
