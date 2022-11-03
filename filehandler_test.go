package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_readIntFromFile(t *testing.T) {
	i, err := readIntFromFile("testdata/last_num_files_test")

	assert.NoError(t, err)
	assert.Equal(t, 1, i)
}
