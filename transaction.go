package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

const subsidy int = 20000

type Transaction struct {
	ID   []byte
	VIn  []TxInput
	VOut []TxOutput
}

type TxInput struct {
	TxID      []byte
	VOut      int
	ScriptSig string
}

type TxOutput struct {
	Value        int
	ScriptPubKey string
}

func NewUTXTransaction(from, to string, amount int, bc *Blockchain) *Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	acc, validOutputs := bc.FindSpendableOutputs(from, amount)

	if acc < amount {
		log.Panic("Error: Insufficient funds")
	}

	// Build a list of inputs
	for txIndex, outs := range validOutputs {
		txID, err := hex.DecodeString(txIndex)
		if err != nil {
			log.Panic("Could not decode string", err)
		}

		for _, out := range outs {
			input := TxInput{txID, out, from}
			inputs = append(inputs, input)
		}
	}

	outputs = append(outputs, TxOutput{amount, to})

	if acc > amount {
		outputs = append(outputs, TxOutput{acc - amount, from})
	}

	tx := Transaction{nil, inputs, outputs}
	tx.SetID()

	return &tx
}

func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}

	txin := TxInput{[]byte{}, -1, data}
	txout := TxOutput{subsidy, to}

	tx := Transaction{nil, []TxInput{txin}, []TxOutput{txout}}
	tx.SetID()

	return &tx
}

func (tx *Transaction) IsCoinbase() bool {
	if len(tx.VIn) != 1 {
		return false
	}

	if len(tx.VIn[0].TxID) != 0 {
		return false
	}

	return tx.VIn[0].VOut == -1
}

func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic("Could not encode transaction", err)
	}

	hash = sha256.Sum256(encoded.Bytes())

	tx.ID = hash[:]
}

func (in *TxInput) CanUnlockOutputWith(unlockingData string) bool {
	return in.ScriptSig == unlockingData
}

func (out *TxOutput) CanBeUnlockedWith(unlockingData string) bool {
	return out.ScriptPubKey == unlockingData
}
