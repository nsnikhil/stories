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
	actualResult := NewCharacterTrie()
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

				ct := NewCharacterTrie()

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

				ct := NewCharacterTrie()

				sentence := "This is a sentence with small and BIG letters"
				res = append(res, ct.insert(sentence, "36982b87-be33-4683-aaaa-e69282a03c83")...)

				sentence = "another thing containing different things"
				res = append(res, ct.insert(sentence, "06bbae40-7048-47d2-a197-39cbf9e704e6")...)

				return res
			},
			expectedErrors: []error{},
		},
		{
			name: "test skip number errors",
			actualResult: func() []error {
				res := make([]error, 0)

				ct := NewCharacterTrie()

				sentence := "this a sentence which contains number 1 2 45"
				res = append(res, ct.insert(sentence, "36982b87-be33-4683-aaaa-e69282a03c83")...)

				return res
			},
			expectedErrors: []error{},
		},
		{
			name: "test skip symbol errors",
			actualResult: func() []error {
				res := make([]error, 0)

				ct := NewCharacterTrie()

				sentence := "this a sentence which, contains symbols %$ ^& @"
				res = append(res, ct.insert(sentence, "36982b87-be33-4683-aaaa-e69282a03c83")...)

				return res
			},
			expectedErrors: []error{},
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
				ct := NewCharacterTrie()
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
				ct := NewCharacterTrie()
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
				ct := NewCharacterTrie()

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
			name: "test get ids for big paragraph",
			actualResult: func() (map[string]bool, []error) {
				ct := NewCharacterTrie()

				sentence := "postman was asked to deliver atleast 10 letters to their respective location by the evening; He has a bike with a limited amount of fuel. For each delivery, he would earn some fixed amount.\nThe postman also gets a huge bonus if he:\ncan deliver 12 or more letters.\nhas fuel left in the bike.\n\nInstead of picking up 10 random letters and starting the trip, the postman looks at 50 different letters and based on the location, traffic prediction, familiarity, fuel and delivery time chooses 12 letters which maximises his revenue for the day.\nSimilarly, when you are presented with a problem statement, instead of straight away jumping into coding it's always better to first analyze the problem statement, the input data, look for any specific pattern, etc then design the solution.\nOnce you have a solution in hand, how do to verify if your solution is optimal, or how do you figure out is any space for improvement?\nThis is exactly what this and the next few articles in this series are going to cover, we will learn:\nThe tools you can use to analyze your algorithm.\nIn-depth analysis of these tools.\nHow can you use these tools to analyze an algorithm?\nAnalysis of two sorting algorithms from the above learning.\n\nAll the examples in the series are written in GO\n\nWe will start by looking at the tools we will use to analyze the efficiency of an algorithm in terms of time and space.\nIf you visit websites like Rotten Tomatoes or IMDb you will see ratings for movies, IMDb uses a scale of 1–10 and Rotten Tomatoes uses percentage but at the end, both websites give scores to movies based on various parameters like the story, music, cinematography, etc.\nSimilarly, there exists a scale to rate algorithms, this scale is independent of the hardware, hence it does not have any specific number to rate an algorithm like it would take X-sec or would consume Y-bytes, rather this scale has behaviours; which are used to interpret how an algorithm behaves with respect to a given input, how it affects the time and space taken with respect to a given input.\nWhen we rate an algorithm we bind it to a behaviour in the scale, this bound is an Asymptotic bound.\nWe will now try to understand what the above statement means, but before that let's understand the meaning of word Asymptote.\nWhat does Asmptote mean?\nA straight line that continually approaches a given curve but does not meet it at any finite distance.\nWhy do we care about them?\nWhen we see how a function behaves on a tail end or on an asymptotic end, we have a true understanding of its performance.\nFor example, The code below computes Fibonacci till nth term recursively:"
				ct.insert(sentence, "36982b87-be33-4683-aaaa-e69282a03c83")

				return ct.getIDs("analyze series GO IMDb")
			},
			expectedResult: func() map[string]bool {
				return map[string]bool{
					"36982b87-be33-4683-aaaa-e69282a03c83": true,
				}
			},
			expectedError: []error{},
		},
		{
			name: "test skip symbol error",
			actualResult: func() (map[string]bool, []error) {
				ct := NewCharacterTrie()
				sentence := "This is a sentence with small and BIG letters"
				ct.insert(sentence, "36982b87-be33-4683-aaaa-e69282a03c83")

				return ct.getIDs("BIG @")
			},
			expectedResult: func() map[string]bool {
				return map[string]bool{"36982b87-be33-4683-aaaa-e69282a03c83": true}
			},
			expectedError: []error{},
		},
		{
			name: "test get errors when words are not found",
			actualResult: func() (map[string]bool, []error) {
				ct := NewCharacterTrie()
				sentence := "This is a sentence with small and BIG letters"
				ct.insert(sentence, "36982b87-be33-4683-aaaa-e69282a03c83")

				return ct.getIDs("these i")
			},
			expectedResult: func() map[string]bool {
				return map[string]bool{}
			},
			expectedError: []error{errors.New("these not found"), errors.New("i not found")},
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

func TestSplitWords(t *testing.T) {
	testCases := []struct {
		name           string
		actualResult   func() []string
		expectedResult []string
	}{
		{
			name: "test split sentence with only small letters",
			actualResult: func() []string {
				sentence := "this is a sentence with only small letters"
				return splitQuery(sentence)
			},
			expectedResult: []string{"this", "is", "a", "sentence", "with", "only", "small", "letters"},
		},
		{
			name: "test split sentence with only small and big letters",
			actualResult: func() []string {
				sentence := "THIS is A sentence WiTh only SMALL and biG leTTers"
				return splitQuery(sentence)
			},
			expectedResult: []string{"THIS", "is", "A", "sentence", "WiTh", "only", "SMALL", "and", "biG", "leTTers"},
		},
		{
			name: "test split sentence with only small, big letters and numbers",
			actualResult: func() []string {
				sentence := "THIS is A sentence WiTh only SMALL and biG leTTers and nUmBeRs 1 23 478 378034"
				return splitQuery(sentence)
			},
			expectedResult: []string{"THIS", "is", "A", "sentence", "WiTh", "only", "SMALL", "and", "biG", "leTTers", "and", "nUmBeRs"},
		},
		{
			name: "test split sentence with only small, big letters, numbers and punctuations",
			actualResult: func() []string {
				sentence := "What does Asmptote mean?\nA straight line, that continually approaches a given curve but does not meet it at any finite distance.\nHow can you use this to analyze an algorithm?\nFor example, The code below computes Fibonacci till nth term recursively:\nAll the examples in the series are written in GO\nWhen we see how a function behaves on a tail end or on an asymptotic end, we have a true understanding of its performance."
				return splitQuery(sentence)
			},
			expectedResult: []string{
				"What", "does", "Asmptote", "mean", "A", "straight", "line", "that", "continually",
				"approaches", "a", "given", "curve", "but", "does", "not", "meet", "it", "at",
				"any", "finite", "distance", "How", "can", "you", "use", "this", "to", "analyze",
				"an", "algorithm", "For", "example", "The", "code", "below", "computes", "Fibonacci",
				"till", "nth", "term", "recursively", "All", "the", "examples", "in", "the", "series",
				"are", "written", "in", "GO", "When", "we", "see", "how", "a", "function", "behaves",
				"on", "a", "tail", "end", "or", "on", "an", "asymptotic", "end", "we", "have", "a",
				"true", "understanding", "of", "its", "performance",
			},
		},
		{
			name: "test split sentence with only numbers",
			actualResult: func() []string {
				sentence := "1 23 456 78910 11122314"
				return splitQuery(sentence)
			},
			expectedResult: []string(nil),
		},
		{
			name: "test split sentence with symbols",
			actualResult: func() []string {
				sentence := "ain't of b%^q *@1 fjn 778 A4 n*( (this) no-very closed"
				return splitQuery(sentence)
			},
			expectedResult: []string{"of", "fjn", "closed"},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, testCase.expectedResult, testCase.actualResult())
		})
	}
}
