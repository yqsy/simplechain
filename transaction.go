package main

import (
	"fmt"
	"log"
	"encoding/hex"
)

const (
	Subsidy = 10
)

type Transaction struct {
	Id   []byte
	VIn  []TxInput
	VOut []TxOutput
}

func (tx Transaction) IsCoinbase() bool {
	return len(tx.VIn) == 1 && len(tx.VIn[0].TxId) == 0 && tx.VIn[0].Prevout == -1
}

type TxInput struct {
	// 该笔交易的ID
	TxId []byte

	// 存储了output的index ? TODO
	Prevout int

	// 解锁TXOutput的ScriptPubKey,
	// 如果正确output可以解锁,随后可以用来生成新的outputs
	// 如果不正确output不能在input中引用
	// 机制确保使用者不能花费属于别人的coin
	ScripSig string
}

func (in *TxInput) CanUnlockOutputWith(unlockingData string) bool {
	return in.ScripSig == unlockingData
}

type TxOutput struct {
	// coins
	Value int

	// 锁?
	ScriptPubKey string
}

func (out *TxOutput) CanBeUnlockedWith(unlockingData string) bool {
	return out.ScriptPubKey == unlockingData
}

func NewCoinbaseTx(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}

	txin := TxInput{
		TxId:     []byte{},
		Prevout:  -1,
		ScripSig: data}

	txout := TxOutput{
		Value:        Subsidy, // int bitcoin, every 21000 blocks the reward is halved
		ScriptPubKey: to}

	tx := Transaction{
		Id:   nil,
		VIn:  []TxInput{txin},
		VOut: []TxOutput{txout}}

	// TODO
	// tx.SetId()
	return &tx
}

func NewUTXOTransaction(from, to string, amount int, bc *BlockChain) *Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	acc, validOutputs := bc.FindSpendableOutputs(from, amount)

	if acc < amount {
		log.Panic("ERROR: Not enough funds")
	}

	// Build a list of inputs
	for txid, outs := range validOutputs {
		txID, _ := hex.DecodeString(txid)

		for _, out := range outs {
			// 一个input来源于多个output, TxID 是上一笔的交易号,Prevout是上一笔交易号中的OutIdx
			// 组成了这个in
			input := TxInput{txID, out, from}
			inputs = append(inputs, input)
		}
	}

	// Build a list of outputs
	// 给to amount
	outputs = append(outputs, TxOutput{amount, to})

	// 如果还有余额 把钱还给amount (自己给自己转账)
	if acc > amount {
		outputs = append(outputs, TxOutput{acc - amount, from}) // a change
	}

	tx := Transaction{nil, inputs, outputs}
	//tx.SetID()

	return &tx
}
