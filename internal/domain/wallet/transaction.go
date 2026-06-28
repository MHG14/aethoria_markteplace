package wallet

type TxType string

const (
	TxReserve TxType = "reserve"
	TxRelease TxType = "release"
	TxDebit   TxType = "debit"
	TxCredit  TxType = "credit"
)
