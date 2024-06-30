package utils

type NetworkType uint8

const (
	Mainnet NetworkType = iota
	Stagenet
	Testnet
)
