package blockchain


type Transaction struct {
	Id   []byte
	Ins  []TxIn
	Outs []TxOut
}

func NewTransaction() *Transaction {
	tx := &Transaction{}

	// 模拟id是整比交易的hash值
	tx.Id = getRandomBit(32)
	return tx
}

type TxIn struct {
}

type TxOut struct {
}


