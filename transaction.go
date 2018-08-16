package main

import "fmt"

const (
	Subsidy = 10
)

type Transaction struct {
	Id   []byte
	Vin  []TXInput
	Vout []TXOutput
}

type TXInput struct {
	// 该笔交易的ID
	Txid []byte

	// 存储了output的index ? TODO
	Vout int

	// 解锁TXOutput的ScriptPubKey,
	// 如果正确output可以解锁,随后可以用来生成新的outputs
	// 如果不正确output不能在input中引用
	// 机制确保使用者不能花费属于别人的coin
	ScripSig string
}

type TXOutput struct {
	// coins
	Value int

	// 锁?
	ScriptPubKey string
}

func NewCoinbaseTx(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}

	txin := TXInput{
		Txid:     []byte{},
		Vout:     -1,
		ScripSig: data}

	txout := TXOutput{
		Value:        Subsidy, // int bitcoin, every 21000 blocks the reward is halved
		ScriptPubKey: to}

	tx := Transaction{
		Id:   nil,
		Vin:  []TXInput{txin},
		Vout: []TXOutput{txout}}

	// TODO
	// tx.SetId()
	return &tx
}
