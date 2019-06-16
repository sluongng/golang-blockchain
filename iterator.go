package main

import "github.com/boltdb/bolt"
import "log"

type ChainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

func (bc *Blockchain) Iterator() *ChainIterator {
	ci := &ChainIterator{bc.tip, bc.db}

	return ci
}

func (i *ChainIterator) Next() *Block {
	var block *Block

	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(i.currentHash)
		block = DeserializeBlock(encodedBlock)

		return nil
	})
	if err != nil {
		log.Fatal("Failed iterating through chain", err)
	}

	i.currentHash = block.PrevBlockHash

	return block
}
