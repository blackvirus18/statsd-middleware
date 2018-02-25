package statsd_middleware

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitStatsDConfig(t *testing.T) {
	config := InitStatsDConfig("STATSD_APP_NAME", "localhost", "8125", true)
	assert.Equal(t, true, config.enabled)
}

func TestInitStatsDMetrics(t *testing.T) {
	config := InitStatsDConfig("STATSD_APP_NAME", "localhost", "8125", true)
	err := InitiateStatsDMetrics(config)
	assert.Nil(t, err)
}

func TestInitStatsDIncrement(t *testing.T) {
	config := InitStatsDConfig("STATSD_APP_NAME", "localhost", "8125", true)
	err := InitiateStatsDMetrics(config)
	assert.Nil(t, err)
	retBool := IncrementInStatsD("check")
	assert.True(t, retBool)
}

func TestInitStatsDIncrementWhenStasdNotInitialized(t *testing.T) {
	statsD = nil
	retBool := IncrementInStatsD("check")
	assert.False(t, retBool)
}
