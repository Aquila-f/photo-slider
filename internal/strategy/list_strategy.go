package strategy

import (
	"context"
	"math/rand"
)

type RandomListStrategy struct{}

func NewRandomListStrategy() *RandomListStrategy {
	return &RandomListStrategy{}
}

func (s *RandomListStrategy) Arrange(_ context.Context, tokens []string) ([]string, error) {
	// Fisher-Yates shuffle
	rand.Shuffle(len(tokens), func(i, j int) {
		tokens[i], tokens[j] = tokens[j], tokens[i]
	})
	return tokens, nil
}
