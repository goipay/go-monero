package utils

type AddressType uint8

const (
	Primary AddressType = iota
	Sub
	Integrated
)

type MoneroAddress interface {
	PublicSpendKey() *PublicKey
	PublicViewKey() *PublicKey
	AddressType() AddressType
	NetworkType() NetworkType
	Address() string
}

// Base struct that impelments MoneroAddress interface
type address struct {
	addr []byte
}

// Returns a Public Spend Key from the Monero integrated/sub-/primary address
func (a *address) PublicSpendKey() *PublicKey {
	s := 1
	key, _ := newPublicKeyHelper(a.addr[s : s+KEY_SIZE])
	return key
}

// Returns a Public View Key from the Monero integrated/sub-/primary address
func (a *address) PublicViewKey() *PublicKey {
	s := KEY_SIZE + 1
	key, _ := newPublicKeyHelper(a.addr[s : s+KEY_SIZE])
	return key
}

// Returns an AddressType from the Monero integrated/sub-/primary address
func (a *address) AddressType() AddressType {
	_, at, _ := GetNetworkTypeAndAddressType(a.addr[0])
	return at
}

// Returns a NetworkType from the Monero integrated/sub-/primary address
func (a *address) NetworkType() NetworkType {
	nt, _, _ := GetNetworkTypeAndAddressType(a.addr[0])
	return nt
}

// Returns a string representation of the Monero integrated/sub-/primary address
func (a *address) Address() string {
	addr, _ := EncodeMoneroAddress(a.addr)
	return string(addr)
}

type PrimaryAddress struct {
	address
}

type SubAddress struct {
	address
}

type IntegratedAddress struct {
	address
}

// Returns a PaymentId from the Monero integrated address
func (a *IntegratedAddress) PaymentId() []byte {
	s := 2*KEY_SIZE + 1
	return a.addr[s : s+8]
}

// Creates a Monero integrated/sub-/primary address
func NewAddress(addr string) (MoneroAddress, error) {
	addrDec, err := DecodeMoneroAddress(addr)
	if err != nil {
		return nil, err
	}

	_, at, err := GetNetworkTypeAndAddressType(addrDec[0])
	if err != nil {
		return nil, err
	}

	switch at {
	case Primary:
		return &PrimaryAddress{address: address{addr: addrDec}}, nil
	case Sub:
		return &SubAddress{address: address{addr: addrDec}}, nil
	case Integrated:
		return &IntegratedAddress{address: address{addr: addrDec}}, nil

	default:
		return nil, err
	}

}
