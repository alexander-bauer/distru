package main

import "testing"

func TestIndex(t *testing.T) {
	index := NewIndex()
	print(RepIndex(index))
}
