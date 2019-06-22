package main

import (
	"fmt"
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

func (tx *Transaction) SetID() {

}
