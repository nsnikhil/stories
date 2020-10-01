package grpc

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsValidUUID(t *testing.T) {
	testCases := []struct {
		name           string
		actualResult   func() bool
		expectedResult bool
	}{
		{
			name: "test return true when uuid is valid",
			actualResult: func() bool {
				return isValidUUID("adbca278-7e5c-4831-bf90-15fadfda0dd1")
			},
			expectedResult: true,
		},
		{
			name: "test return false when uuid is invalid",
			actualResult: func() bool {
				return isValidUUID("this-is-an-invalid-uuid")
			},
			expectedResult: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, testCase.expectedResult, testCase.actualResult())
		})
	}
}
