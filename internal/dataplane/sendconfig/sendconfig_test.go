package sendconfig

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateReportingUtilities(t *testing.T) {
	assert.False(t, hasSHAUpdateAlreadyBeenReported([]byte("fake-sha")))
	assert.True(t, hasSHAUpdateAlreadyBeenReported([]byte("fake-sha")))
	assert.False(t, hasSHAUpdateAlreadyBeenReported([]byte("another-fake-sha")))
	assert.True(t, hasSHAUpdateAlreadyBeenReported([]byte("another-fake-sha")))
	assert.False(t, hasSHAUpdateAlreadyBeenReported([]byte("yet-another-fake-sha")))
	assert.True(t, hasSHAUpdateAlreadyBeenReported([]byte("yet-another-fake-sha")))
	assert.True(t, hasSHAUpdateAlreadyBeenReported([]byte("yet-another-fake-sha")))
	assert.True(t, hasSHAUpdateAlreadyBeenReported([]byte("yet-another-fake-sha")))
}
