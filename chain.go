package main

import (
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

const dbFile = "blockchain.db"
const blocksBucket = "blocks"

type Blockchain struct {
	tip []byte
	db  *bolt.DB
}

func NewBlockChain() *Blockchain {
	var tip []byte

	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic("Could not init DB", err)
	}

	err = db.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(blocksBucket))

		// If not found Blockchain
		if b == nil {
			fmt.Println("No previous blockchain found. Creating a new one...")

			// Create bucket
			b, err := tx.CreateBucket([]byte(blocksBucket))
			if err != nil {
				log.Panic("Coudld not create bucket", err)
				return err
			}

			// Create + Insert Genesis block
			genesis := NewGenesisBlock()
			err = b.Put(genesis.Hash, genesis.Serialize())
			if err != nil {
				log.Panic("Could not insert Genesis Block", err)
				return err
			}

			// Mark last block available is Genesis block
			err = b.Put([]byte("l"), genesis.Hash)
			if err != nil {
				log.Panic("Could not insert tail pointer", err)
				return err
			}

			tip = genesis.Hash
		} else {

			tip = b.Get([]byte("l"))
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	return &Blockchain{tip, db}
}

func (bc *Blockchain) AddBlock(data string) {
	var lastHash []byte

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))

		// TODO
		return nil
	})

	newBlock := NewBlock(data, lastHash)

	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket))
		err := b.Put(newBlock.Hash, newBlock.Serialize())
		err = b.Put([]byte("l"), newBlock.Hash)
		bc.tip = newBlock.Hash

		return nil
	})

}
