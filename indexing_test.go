package main

import (
	"testing"
)

func TestBin(t *testing.T) {
	idx := RecvIndex("localhost")
	t.Log("Got index from localhost.")
	t.Log(RepIndex(idx))
}

/*func TestRep(t *testing.T) {
	print(RepIndex(Idx))
}*/

/*
func TestSave(t *testing.T) {
	index := NewIndex()
	index.save("/tmp/distru-index")
}*/
