package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseTags_Empty(t *testing.T) {
	result := parseTags("")
	assert.Nil(t, result)
}

func TestParseTags_Single(t *testing.T) {
	result := parseTags("cloud")
	assert.Equal(t, []string{"cloud"}, result)
}

func TestParseTags_Multiple(t *testing.T) {
	result := parseTags("aws,cloud,production")
	assert.Equal(t, []string{"aws", "cloud", "production"}, result)
}

func TestParseTags_TrimsSpaces(t *testing.T) {
	result := parseTags(" aws , cloud ")
	assert.Equal(t, []string{"aws", "cloud"}, result)
}

func TestParseTags_SkipsEmpty(t *testing.T) {
	result := parseTags("aws,,cloud")
	assert.Equal(t, []string{"aws", "cloud"}, result)
}
