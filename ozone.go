package main

import (
	"regexp"
	"strings"
)

var singleRE = regexp.MustCompile(`(?:.* / )*(\w+)$`)
var matrixRE = regexp.MustCompile(`(?:.* / )*(\w+) \(([^)]+)\)$`)

func JobToArtifactName(job string) string {
	if matrixRE.MatchString(job) {
		var name = matrixRE.ReplaceAllString(job, "$1-$2")
		return strings.ReplaceAll(strings.ReplaceAll(name, "/", "-"), ".sh", "")
	}
	if singleRE.MatchString(job) {
		return singleRE.ReplaceAllString(job, "$1")
	}
	return job
}
