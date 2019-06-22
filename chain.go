package main

import (
	"fmt"
	"log"
	"os"

	"github.com/boltdb/bolt"
)

const dbFile = "blockchain.db"
const blocksBucket = "blocks"
const genesisCoinbaseData = "Son Luong Ngoc created this"

type Blockchain struct {
	tip []byte
	db  *bolt.DB
}

func dbExists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}

func NewBlockchain(address string) *Blockchain {
	if !dbExists() {
		fmt.Println("No existing blockchain found. Create one first.")
		os.Exit(1)
	}

	var tip []byte

	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic("could not open DB file", err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		tip = b.Get([]byte("l"))

		return nil
	})

	if err != nil {
		log.Panic("could not get update db", err)
	}

	return &Blockchain{tip, db}
}

func CreateBlockchain(address string) *Blockchain {
	if dbExists() {
		fmt.Println("Blockchain already exist")
		os.Exit(1)
	}

	var tip []byte

	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic("Could not open dbFile", err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		cbtx := NewCoinbaseTX(address, genesisCoinbaseData)
		genesis := NewGenesisBlock(cbtx)

		b, err := tx.CreateBucket([]byte(blocksBucket))
		if err != nil {
			log.Panic("could not create bucket", err)
			return err
		}

		err = b.Put(genesis.Hash, genesis.Serialize())
		if err != nil {
			log.Panic("could not persist genesis block", err)
			return err
		}

		err = b.Put([]byte("l"), genesis.Hash)
		if err != nil {
			log.Panic("could not update last index", err)
			return err
		}

		tip = genesis.Hash

		return nil
	})

	if err != nil {
		log.Panic("could not update DB", err)
	}

	return &Blockchain{tip, db}
}

func (bc *Blockchain) AddBlock(transactions []*Transaction) {
	var lastHash []byte

	// Get Last hash
	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))

		return nil
	})
	if err != nil {
		log.Panic("Failed to get last hash", err)
	}

	// Create and store new block
	newBlock := NewBlock(transactions, lastHash)
	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		err := b.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			log.Panic("Failed to store new block", err)
			return err
		}

		err = b.Put([]byte("l"), newBlock.Hash)
		if err != nil {
			log.Panic("Failed to update last hash", err)
			return err
		}

		bc.tip = newBlock.Hash

		return nil
	})

}
