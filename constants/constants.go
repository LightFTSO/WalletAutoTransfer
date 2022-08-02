package constants

type Network struct {
	Name    string
	ChainId int
	RpcUrl  string
}

type NetworkTypes struct {
	Coston   Network
	Songbird Network
	Flare    Network
}

type Wallet struct {
	Address    string
	PrivateKey string
	Mnemonic   string
}

var Coston = &Network{Name: "Coston", ChainId: 16}
var Songbird = &Network{Name: "Songbird", ChainId: 19}
var Flare = &Network{Name: "Flare", ChainId: 14}

var Networks = map[string]*Network{
	"Coston":   Coston,
	"Songbird": Songbird,
	"Flare":    Flare,
}

var Nat = map[string]string{
	"Coston":   "CFLR",
	"Songbird": "SGB",
	"Flare":    "FLR",
}
