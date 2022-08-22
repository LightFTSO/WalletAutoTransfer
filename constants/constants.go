package constants

type Network struct {
	Name    string
	ChainId int
	RpcUrl  string
	Nat     string
}

type Wallet struct {
	Address    string
	PrivateKey string
	Mnemonic   string
}
