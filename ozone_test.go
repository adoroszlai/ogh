package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJobToArtifactName(t *testing.T) {
	job := "CI / acceptance (ozonesecure-ha/test-s3g-virtual-host.sh) / acceptance (ozonesecure-ha/test-s3g-virtual-host.sh)"
	artifact := JobToArtifactName(job)
	assert.Equal(t, "acceptance-ozonesecure-ha-test-s3g-virtual-host", artifact)

	job = "CI / junit (ozone-mini-cluster) / junit (ozone-mini-cluster)"
	artifact = JobToArtifactName(job)
	assert.Equal(t, "junit-ozone-mini-cluster", artifact)

	job = "CI / repro / repro"
	artifact = JobToArtifactName(job)
	assert.Equal(t, "repro", artifact)
}
