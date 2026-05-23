package filter_test

import (
	"testing"

	"github.com/logslice/logslice/internal/filter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFilter_NoRules(t *testing.T) {
	f, err := filter.New(filter.Options{})
	require.NoError(t, err)
	assert.True(t, f.Match("any line passes"))
}

func TestFilter_IncludePlain(t *testing.T) {
	f, err := filter.New(filter.Options{
		Include: []string{"ERROR"},
	})
	require.NoError(t, err)
	assert.True(t, f.Match("2024-01-01 ERROR something broke"))
	assert.False(t, f.Match("2024-01-01 INFO all good"))
}

func TestFilter_ExcludePlain(t *testing.T) {
	f, err := filter.New(filter.Options{
		Exclude: []string{"DEBUG"},
	})
	require.NoError(t, err)
	assert.False(t, f.Match("DEBUG verbose output"))
	assert.True(t, f.Match("ERROR something bad"))
}

func TestFilter_IncludeAndExclude(t *testing.T) {
	f, err := filter.New(filter.Options{
		Include: []string{"ERROR"},
		Exclude: []string{"timeout"},
	})
	require.NoError(t, err)
	assert.True(t, f.Match("ERROR disk full"))
	assert.False(t, f.Match("ERROR timeout connecting"))
	assert.False(t, f.Match("INFO startup"))
}

func TestFilter_CaseInsensitiveDefault(t *testing.T) {
	f, err := filter.New(filter.Options{
		Include: []string{"error"},
	})
	require.NoError(t, err)
	assert.True(t, f.Match("ERROR: something failed"))
	assert.True(t, f.Match("Error: bad input"))
}

func TestFilter_CaseSensitive(t *testing.T) {
	f, err := filter.New(filter.Options{
		Include:       []string{"error"},
		CaseSensitive: true,
	})
	require.NoError(t, err)
	assert.False(t, f.Match("ERROR: something failed"))
	assert.True(t, f.Match("error: something failed"))
}

func TestFilter_RegexInclude(t *testing.T) {
	f, err := filter.New(filter.Options{
		Include:  []string{`level=(error|warn)`},
		UseRegex: true,
	})
	require.NoError(t, err)
	assert.True(t, f.Match(`ts=2024-01-01 level=error msg="disk full"`))
	assert.True(t, f.Match(`ts=2024-01-01 level=warn msg="low memory"`))
	assert.False(t, f.Match(`ts=2024-01-01 level=info msg="started"`))
}

func TestFilter_InvalidRegex(t *testing.T) {
	_, err := filter.New(filter.Options{
		Include:  []string{`[invalid`},
		UseRegex: true,
	})
	require.Error(t, err)
}
