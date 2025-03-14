package main

import (
	"regexp"
)

var singleRE = regexp.MustCompile(`(?:.* / )*(\w+)`)
var matrixRE = regexp.MustCompile(`(?:.* / )*(\w+) \(([^)]+)\)`)

func JobToArtifactName(job string) string {
	if matrixRE.MatchString(job) {
		return matrixRE.ReplaceAllString(job, "$1-$2")
	}
	if singleRE.MatchString(job) {
		return singleRE.ReplaceAllString(job, "$1")
	}
	return job
}
