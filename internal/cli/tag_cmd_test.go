package cli_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTagCmd_AddAndList(t *testing.T) {
	cmdr, _ := tempCommander(t)

	require.NoError(t, cmdr.Set("prod", "KEY", "val"))
	require.NoError(t, cmdr.AddTag("prod", "production"))
	require.NoError(t, cmdr.AddTag("prod", "aws"))

	tags := cmdr.GetTags("prod")
	assert.ElementsMatch(t, []string{"aws", "production"}, tags)
}

func TestTagCmd_ListByTag(t *testing.T) {
	cmdr, _ := tempCommander(t)

	require.NoError(t, cmdr.Set("prod", "K", "v"))
	require.NoError(t, cmdr.Set("staging", "K", "v"))
	require.NoError(t, cmdr.Set("dev", "K", "v"))

	require.NoError(t, cmdr.AddTag("prod", "cloud"))
	require.NoError(t, cmdr.AddTag("staging", "cloud"))
	require.NoError(t, cmdr.AddTag("dev", "local"))

	cloud, err := cmdr.ListByTag("cloud")
	require.NoError(t, err)
	assert.Equal(t, []string{"prod", "staging"}, cloud)

	local, err := cmdr.ListByTag("local")
	require.NoError(t, err)
	assert.Equal(t, []string{"dev"}, local)
}

func TestTagCmd_RemoveTag(t *testing.T) {
	cmdr, _ := tempCommander(t)

	require.NoError(t, cmdr.Set("prod", "K", "v"))
	require.NoError(t, cmdr.AddTag("prod", "cloud"))
	require.NoError(t, cmdr.AddTag("prod", "aws"))
	require.NoError(t, cmdr.RemoveTag("prod", "cloud"))

	tags := cmdr.GetTags("prod")
	assert.Equal(t, []string{"aws"}, tags)
}

func TestTagCmd_RemoveLastTag(t *testing.T) {
	cmdr, _ := tempCommander(t)

	require.NoError(t, cmdr.Set("prod", "K", "v"))
	require.NoError(t, cmdr.AddTag("prod", "cloud"))
	require.NoError(t, cmdr.RemoveTag("prod", "cloud"))

	tags := cmdr.GetTags("prod")
	assert.Empty(t, tags)
}

func TestTagCmd_DuplicateTag(t *testing.T) {
	cmdr, _ := tempCommander(t)

	require.NoError(t, cmdr.Set("prod", "K", "v"))
	require.NoError(t, cmdr.AddTag("prod", "cloud"))
	err := cmdr.AddTag("prod", "cloud")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
}

func TestTagCmd_MissingProfile(t *testing.T) {
	cmdr, _ := tempCommander(t)

	err := cmdr.AddTag("ghost", "cloud")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestTagCmd_ListByTag_Empty(t *testing.T) {
	cmdr, _ := tempCommander(t)

	result, err := cmdr.ListByTag("nonexistent")
	require.NoError(t, err)
	assert.Empty(t, result)
}
