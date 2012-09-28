package main

import (
	"testing"
	"log"
)

func TestBin(t *testing.T) {
	idx := RecvIndex("localhost")
	log.Println("Recieved binary index.")
	log.Println(RepIndex(idx))
}

/*func TestRep(t *testing.T) {
	print(RepIndex(Idx))
}*/

/*
func TestSave(t *testing.T) {
	index := NewIndex()
	index.save("/tmp/distru-index")
}*/
