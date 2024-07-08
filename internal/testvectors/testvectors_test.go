package testvectors_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/abaxxtech/abaxx-id-go/internal/testvectors"
	"github.com/stretchr/testify/assert"
)

func TestLoadTestVectors(t *testing.T) {
	// Create a temporary file with test vectors
	tempFile, err := os.CreateTemp("", "testvectors-*.json")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())

	testData := testvectors.TestVectors[string, any]{
		Description: "Test Vectors",
		Vectors: []testvectors.TestVector[string, any]{
			{
				Description: "Test Vector 1",
				Input:       "input1",
				Output:      "output1",
				Errors:      false,
			},
			{
				Description: "Test Vector 2",
				Input:       "input2",
				Output:      "output2",
				Errors:      true,
			},
		},
	}

	data, err := json.Marshal(testData)
	assert.NoError(t, err)

	_, err = tempFile.Write(data)
	assert.NoError(t, err)

	// Load the test vectors from the temporary file
	loadedVectors, err := testvectors.LoadTestVectors[string, any](tempFile.Name())
	assert.NoError(t, err)

	// Verify the loaded test vectors
	assert.Equal(t, testData.Description, loadedVectors.Description)
	assert.Equal(t, len(testData.Vectors), len(loadedVectors.Vectors))

	for i, vector := range testData.Vectors {
		assert.Equal(t, vector.Description, loadedVectors.Vectors[i].Description)
		assert.Equal(t, vector.Input, loadedVectors.Vectors[i].Input)
		assert.Equal(t, vector.Output, loadedVectors.Vectors[i].Output)
		assert.Equal(t, vector.Errors, loadedVectors.Vectors[i].Errors)
	}
}
