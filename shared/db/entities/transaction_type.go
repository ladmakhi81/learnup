package entities

type TransactionTag string

const (
	TransactionTag_Sell          TransactionTag = "sell"
	TransactionTag_ChargeWallet  TransactionTag = "charge_wallet"
	TransactionTag_DepositWallet TransactionTag = "deposit_wallet"
	TransactionTag_SalaryPayment TransactionTag = "salary_payment"
)
