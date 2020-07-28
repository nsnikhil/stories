package store

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateNewTrieNode(t *testing.T) {
	actualResult := newTrieNode()
	expectedResult := &trieNode{}
	assert.Equal(t, expectedResult, actualResult)
}

func TestCreateNewTrie(t *testing.T) {
	actualResult := newCharacterTrie()
	expectedResult := &characterTrie{root: newTrieNode()}
	assert.Equal(t, expectedResult, actualResult)
}

func TestTrieInsert(t *testing.T) {
	testCases := []struct {
		name           string
		actualResult   func() []error
		expectedErrors []error
	}{
		{
			name: "test get ids when only small letters are present",
			actualResult: func() []error {
				res := make([]error, 0)

				ct := newCharacterTrie()

				sentence := "this is a sentence with only small letters"
				res = append(res, ct.insert(sentence, "36982b87-be33-4683-aaaa-e69282a03c83")...)

				sentence = "another thing containing different things"
				res = append(res, ct.insert(sentence, "06bbae40-7048-47d2-a197-39cbf9e704e6")...)

				return res
			},
			expectedErrors: []error{},
		},
		{
			name: "test get ids when both small and large letters are present",
			actualResult: func() []error {
				res := make([]error, 0)

				ct := newCharacterTrie()

				sentence := "This is a sentence with small and BIG letters"
				res = append(res, ct.insert(sentence, "36982b87-be33-4683-aaaa-e69282a03c83")...)

				sentence = "another thing containing different things"
				res = append(res, ct.insert(sentence, "06bbae40-7048-47d2-a197-39cbf9e704e6")...)

				return res
			},
			expectedErrors: []error{},
		},
		{
			name: "test get errors when sentence contains alphanumeric characters",
			actualResult: func() []error {
				res := make([]error, 0)

				ct := newCharacterTrie()

				sentence := "this a sentence which contains number 1 2 45"
				res = append(res, ct.insert(sentence, "36982b87-be33-4683-aaaa-e69282a03c83")...)

				return res
			},
			expectedErrors: []error{
				errors.New("1 is not a valid character"),
				errors.New("2 is not a valid character"),
				errors.New("4 is not a valid character"),
			},
		},
		{
			name: "test get errors when sentence contains symbols",
			actualResult: func() []error {
				res := make([]error, 0)

				ct := newCharacterTrie()

				sentence := "this a sentence which, contains symbols %$ ^& @"
				res = append(res, ct.insert(sentence, "36982b87-be33-4683-aaaa-e69282a03c83")...)

				return res
			},
			expectedErrors: []error{
				errors.New("% is not a valid character"),
				errors.New("^ is not a valid character"),
				errors.New("@ is not a valid character"),
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, testCase.expectedErrors, testCase.actualResult())
		})
	}
}

func TestTrieGetIds(t *testing.T) {
	testCases := []struct {
		name           string
		actualResult   func() (map[string]bool, []error)
		expectedResult func() map[string]bool
		expectedError  []error
	}{
		{
			name: "test get ids when only small letters are present",
			actualResult: func() (map[string]bool, []error) {
				ct := newCharacterTrie()
				sentence := "this is a sentence with only small letters"
				ct.insert(sentence, "36982b87-be33-4683-aaaa-e69282a03c83")

				sentence = "another thing containing different things"
				ct.insert(sentence, "06bbae40-7048-47d2-a197-39cbf9e704e6")

				return ct.getIDs("things with")

			},
			expectedResult: func() map[string]bool {
				return map[string]bool{"06bbae40-7048-47d2-a197-39cbf9e704e6": true, "36982b87-be33-4683-aaaa-e69282a03c83": true}
			},
			expectedError: []error{},
		},
		{
			name: "test get ids when both small and large letters are present",
			actualResult: func() (map[string]bool, []error) {
				ct := newCharacterTrie()
				sentence := "This is a sentence with small and BIG letters"
				ct.insert(sentence, "36982b87-be33-4683-aaaa-e69282a03c83")

				sentence = "another thing containing different things"
				ct.insert(sentence, "06bbae40-7048-47d2-a197-39cbf9e704e6")

				return ct.getIDs("BIG different")

			},
			expectedResult: func() map[string]bool {
				return map[string]bool{"06bbae40-7048-47d2-a197-39cbf9e704e6": true, "36982b87-be33-4683-aaaa-e69282a03c83": true}
			},
			expectedError: []error{},
		},
		{
			name: "test get ids punctuations are present",
			actualResult: func() (map[string]bool, []error) {
				ct := newCharacterTrie()

				sentence := "this! one"
				ct.insert(sentence, "36982b87-be33-4683-aaaa-e69282a03c83")

				sentence = "another, two"
				ct.insert(sentence, "06bbae40-7048-47d2-a197-39cbf9e704e6")

				sentence = "LETTER."
				ct.insert(sentence, "439502c3-8830-4df8-a2d1-2d4acbbe2645")

				sentence = "hence: say"
				ct.insert(sentence, "484c68b8-0870-4a46-90be-10b2faf958f7")

				sentence = "ends; HERE"
				ct.insert(sentence, "17752987-24fa-4a7d-b0fb-a86e6b5cc08f")

				sentence = "okay?"
				ct.insert(sentence, "2eaa0697-2572-47f9-bcff-0bdf0c7c6432")

				return ct.getIDs("this another letter hence ends okay")

			},
			expectedResult: func() map[string]bool {
				return map[string]bool{
					"36982b87-be33-4683-aaaa-e69282a03c83": true,
					"06bbae40-7048-47d2-a197-39cbf9e704e6": true,
					"439502c3-8830-4df8-a2d1-2d4acbbe2645": true,
					"484c68b8-0870-4a46-90be-10b2faf958f7": true,
					"17752987-24fa-4a7d-b0fb-a86e6b5cc08f": true,
					"2eaa0697-2572-47f9-bcff-0bdf0c7c6432": true,
				}
			},
			expectedError: []error{},
		},
		{
			name: "test get errors when sentence contains symbols",
			actualResult: func() (map[string]bool, []error) {
				ct := newCharacterTrie()
				sentence := "This is a sentence with small and BIG letters"
				ct.insert(sentence, "36982b87-be33-4683-aaaa-e69282a03c83")

				return ct.getIDs("BIG @")

			},
			expectedResult: func() map[string]bool {
				return map[string]bool{"36982b87-be33-4683-aaaa-e69282a03c83": true}
			},
			expectedError: []error{errors.New("@ is not a valid character")},
		},
		{
			name: "test get errors when words are not found",
			actualResult: func() (map[string]bool, []error) {
				ct := newCharacterTrie()
				sentence := "This is a sentence with small and BIG letters"
				ct.insert(sentence, "36982b87-be33-4683-aaaa-e69282a03c83")

				return ct.getIDs("these i")

			},
			expectedResult: func() map[string]bool {
				return map[string]bool{}
			},
			expectedError: []error{errors.New("these not found"),errors.New("i not found")},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			res, err := testCase.actualResult()

			assert.Equal(t, testCase.expectedError, err)
			assert.Equal(t, testCase.expectedResult(), res)
		})
	}
}
