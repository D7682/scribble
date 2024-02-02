package integration

import (
	"os"
	"testing"

	"github.com/D7682/scribble"
	"github.com/stretchr/testify/assert"
)

func TestScribbleIntegration(t *testing.T) {
	// Set up a temporary directory for the test database
	dir := "./scribble_test_db"
	defer func() {
		// Clean up the temporary directory after the test
		err := os.RemoveAll(dir)
		assert.NoError(t, err)
	}()

	// Create a new scribble driver instance for the test database
	db, err := scribble.New(dir, nil)
	assert.NoError(t, err)

	// Perform integration tests with the scribble package
	// Example: Write to the database
	err = db.Write("example", "key", map[string]interface{}{"field": "value"})
	assert.NoError(t, err)

	// Example: Read from the database
	var result map[string]interface{}
	err = db.Read("example", "key", &result)
	assert.NoError(t, err)
	assert.Equal(t, "value", result["field"])
}
