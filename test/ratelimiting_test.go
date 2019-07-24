// +build integration

package test

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSample(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(10, 20, "WTF is that")
}
