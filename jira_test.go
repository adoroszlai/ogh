package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTweakTitle(t *testing.T) {
	assert.Equal(t, "Bump actions/checkout to 6.0.3", tweakDependabotTitle("Bump actions/checkout from 6.0.2 to 6.0.3"))
	assert.Equal(t, "Bump asm to 9.10.1", tweakDependabotTitle("Bump asm.version from 9.10 to 9.10.1"))
	assert.Equal(t, "Bump awssdk to 2.45.1", tweakDependabotTitle("Bump software.amazon.awssdk:bom from 2.44.12 to 2.45.1"))
	assert.Equal(t, "Bump commons-configuration2 to 2.15.1", tweakDependabotTitle("Bump org.apache.commons:commons-configuration2 from 2.15.0 to 2.15.1"))
	assert.Equal(t, "Bump jacoco-maven-plugin to 0.8.15", tweakDependabotTitle("Bump org.jacoco:jacoco-maven-plugin from 0.8.14 to 0.8.15"))
	assert.Equal(t, "Bump netty to 4.2.15", tweakDependabotTitle("Bump io.netty:netty-bom from 4.2.14.Final to 4.2.15.Final"))
	assert.Equal(t, "Bump dropwizard to 4.2.39", tweakDependabotTitle("Bump shaded.dropwizard.version from 4.2.38 to 4.2.39"))
	assert.Equal(t, "Bump slf4j to 2.0.18", tweakDependabotTitle("Bump org.slf4j:slf4j-bom from 2.0.17 to 2.0.18"))
}
