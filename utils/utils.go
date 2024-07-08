package utils

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"unicode/utf8"

	"filippo.io/edwards25519"
	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/sha3"
)

const (
	ATOMIC_UNIT                     float64 = 1e12
	CHECKSUM_SIZE                   int     = 4
	KEY_SIZE                        int     = 32
	BASE58_FULL_BLOCK_SIZE          int     = 8
	BASE58_ENCODED_BLOCK_SIZE       int     = 11
	INTEGRATED_ADDRESS_SIZE         int     = 106
	INTEGRATED_ADDRESS_DECODED_SIZE int     = 77
	ADDRESS_SIZE                    int     = 95
	ADDRESS_DECODED_SIZE            int     = 69
)

var (
	invalid_address_err error  = errors.New("invalid Monero address")
	view_tag_prefix     []byte = []byte("view_tag")
	amount_prefix       []byte = []byte("amount")
	subaddr_prefix      []byte = []byte("SubAddr\x00")
	default_extra_tags  []byte = []byte{0x01, 0x02}
)

/********************************************** Cryptography Related Mehtods ***************************************************/

// Derives the public key from the private (either view or spend)
func GetPublicKeyFromPrivate(privKey *PrivateKey) *PublicKey {
	/** Public Key - Point
		A = a * G

		a - Private Key - Scalar
		G - Base Point of Ed25519 Elliptic Curve - Point
	**/
	pubKey := &PublicKey{key: new(edwards25519.Point).ScalarBaseMult(privKey.key)}

	return pubKey
}

// Derives the Private View Key from the Private Spend Key
func GetPrivateViewKeyFromPrivateSpendKey(spend *PrivateKey) (*PrivateKey, error) {
	data, err := Keccak256Hash(spend.Bytes())
	if err != nil {
		return nil, err
	}

	sc, err := keccak256HashToScalar(data)
	if err != nil {
		return nil, err
	}

	return &PrivateKey{key: sc}, nil
}

/********************************************** Tx Related Mehtods ***************************************************/

func calculateSharedKeyHelper(viewKey *PrivateKey, txPub *PublicKey) *edwards25519.Point {
	/** Shared Key - Point
		S = v * R

		v - Private View Key - Scalar
		R - Tx Public Key - Point
	**/
	S := new(edwards25519.Point).ScalarMult(viewKey.key, txPub.key)
	// 8 * S = 8 * v * R - multiplying the Shared Key by cofactor 8
	S = S.MultByCofactor(S)

	return S
}

func concatPrefixSharedKeyOutIndexHelper(prefix []byte, S *edwards25519.Point, index uint64) []byte {
	/** b"view_tag" + S + i
		view_tag_prefix - 8 bytes
		shared key - 32 bytes
		output index - 1-4 bytes
	**/
	indexBytes := uintToLittleEndianBytes(index)

	con := make([]byte, len(prefix)+KEY_SIZE+len(indexBytes))
	copyToSlice(prefix, con, 0)
	copyToSlice(S.Bytes(), con, len(prefix))
	copyToSlice(indexBytes, con, KEY_SIZE+len(prefix))

	return con
}

// Checks whether the output belongs to the specific private view key by comparing view tags
func OutputBelongsViewTag(viewTag string, outIndex uint32, txPub *PublicKey, viewKey *PrivateKey) (bool, error) {
	tag, err := hex.DecodeString(viewTag)
	if err != nil {
		return false, err
	}

	res, err := Keccak256Hash(concatPrefixSharedKeyOutIndexHelper(view_tag_prefix, calculateSharedKeyHelper(viewKey, txPub), uint64(outIndex)))
	if err != nil {
		return false, err
	}

	return res[0] == tag[0], nil
}

func calculateSharedKeyConcatOutIndexHash(S *edwards25519.Point, outIndex uint32) (*edwards25519.Scalar, error) {
	/** Hs(S||i) - Scalar

		Hs - Keccak256Hash - Function
		S - Shared Key - Point
		i - Output index - Scalar
	**/
	hash, err := Keccak256Hash(append(S.Bytes(), uintToLittleEndianBytes(uint64(outIndex))...))
	if err != nil {
		return nil, err
	}

	sc, err := keccak256HashToScalar(hash)
	if err != nil {
		return nil, err
	}

	return sc, nil
}

// Checks whether the output belongs to the specific private view and public spend keys
func OutputBelongsPublicSpendKey(spendKey *PublicKey, outIndex uint32, outKey *PublicKey, txPub *PublicKey, viewKey *PrivateKey) (bool, error) {
	S := calculateSharedKeyHelper(viewKey, txPub)

	Si, err := calculateSharedKeyConcatOutIndexHash(S, outIndex)
	if err != nil {
		return false, err
	}
	/** Si*G + B - Point

		Si - Keccak256Hash of Shared key + Output index - Hs(S||i) - Scalar
		G - Base Point of Ed25519 Elliptic Curve - Point
		B - Public Spend Key - Point
	**/
	outKeyDer := new(edwards25519.Point).ScalarBaseMult(Si)
	outKeyDer.Add(outKeyDer, spendKey.key)

	return bytes.Equal(outKey.Bytes(), outKeyDer.Bytes()), nil
}

func decryptOutputAmountHelper(Si *edwards25519.Scalar, amount string) (uint64, error) {
	amountDecoded, err := hex.DecodeString(amount)
	if err != nil {
		return 0, err
	}

	// keccak("amount"||Hs(8aR||i))
	aSi, err := Keccak256Hash(append(amount_prefix, Si.Bytes()...))
	if err != nil {
		return 0, err
	}
	aSi = aSi[:len(amountDecoded)]

	return binary.LittleEndian.Uint64(xor(amountDecoded, aSi)), nil
}

// Checks whether the output belongs to the specific private view and public spend keys and if so, returns a decrypted amount
func DecryptOutputPublicSpendKey(spendKey *PublicKey, outIndex uint32, outKey *PublicKey, amount string, txPub *PublicKey, viewKey *PrivateKey) (bool, uint64, error) {
	S := calculateSharedKeyHelper(viewKey, txPub)

	Si, err := calculateSharedKeyConcatOutIndexHash(S, outIndex)
	if err != nil {
		return false, 0, err
	}

	/** Si*G + B - Point

		Si - Keccak256Hash of Shared key + Output index - Hs(S||i) - Scalar
		G - Base Point of Ed25519 Elliptic Curve - Point
		B - Public Spend Key - Point
	**/
	outKeyDer := new(edwards25519.Point).ScalarBaseMult(Si)
	outKeyDer.Add(outKeyDer, spendKey.key)

	if !bytes.Equal(outKey.Bytes(), outKeyDer.Bytes()) {
		return false, 0, nil
	}

	amountDec, err := decryptOutputAmountHelper(Si, amount)
	if err != nil {
		return false, 0, err
	}

	return true, amountDec, nil
}

// Checks whether the output belongs to the specific private view key by comparing view tags
// and if so, returns a decrypted amount
func DecryptOutputViewTag(viewTag string, outIndex uint32, amount string, txPub *PublicKey, viewKey *PrivateKey) (bool, uint64, error) {
	tag, err := hex.DecodeString(viewTag)
	if err != nil {
		return false, 0, err
	}

	S := calculateSharedKeyHelper(viewKey, txPub)

	res, err := Keccak256Hash(concatPrefixSharedKeyOutIndexHelper(view_tag_prefix, S, uint64(outIndex)))
	if err != nil {
		return false, 0, err
	}

	if res[0] != tag[0] {
		return false, 0, nil
	}

	Si, err := calculateSharedKeyConcatOutIndexHash(S, outIndex)
	if err != nil {
		return false, 0, err
	}

	amountDec, err := decryptOutputAmountHelper(Si, amount)
	if err != nil {
		return false, 0, err
	}

	return true, amountDec, nil
}

/********************************************** Monero Address Related Mehtods ***************************************************/

func decodeMoneroAddressBase58Helper(addr string) []byte {
	addrSize := len(addr)

	size := ADDRESS_DECODED_SIZE
	if addrSize == INTEGRATED_ADDRESS_SIZE {
		size = INTEGRATED_ADDRESS_DECODED_SIZE
	}
	res := make([]byte, size)

	start := 0
	end := BASE58_ENCODED_BLOCK_SIZE

	for end < addrSize {
		dec := base58.Decode(addr[start:end])
		if len(dec) > BASE58_FULL_BLOCK_SIZE {
			dec = dec[len(dec)-BASE58_FULL_BLOCK_SIZE:]
		}

		copyToSlice(dec, res, BASE58_FULL_BLOCK_SIZE*start/BASE58_ENCODED_BLOCK_SIZE)

		start = end
		end = start + BASE58_ENCODED_BLOCK_SIZE
	}
	copyToSlice(base58.Decode(addr[start:]), res, BASE58_FULL_BLOCK_SIZE*start/BASE58_ENCODED_BLOCK_SIZE)

	return res
}

func verifyMoneroAddressChecksumHelper(addrDec []byte) bool {
	addrDecLen := len(addrDec)

	excsum := addrDec[addrDecLen-CHECKSUM_SIZE:]
	csum, err := Keccak256Hash(addrDec[:addrDecLen-CHECKSUM_SIZE])
	if err != nil {
		return false
	}
	csum = csum[:CHECKSUM_SIZE]

	return bytes.Equal(excsum, csum)
}

// Decodes and validates if the given addr string is a valid Monero integrated/sub-/primary address
func DecodeMoneroAddress(addr string) ([]byte, error) {
	addrSize := utf8.RuneCountInString(addr)
	if addrSize != ADDRESS_SIZE && addrSize != INTEGRATED_ADDRESS_SIZE {
		return nil, invalid_address_err
	}

	addrDec := decodeMoneroAddressBase58Helper(addr)
	if !verifyMoneroAddressChecksumHelper(addrDec) {
		return nil, invalid_address_err
	}

	return addrDec, nil
}

func encodeMoneroAddressBase58Helper(addr []byte) []byte {
	addrSize := len(addr)

	size := ADDRESS_SIZE
	if addrSize == INTEGRATED_ADDRESS_DECODED_SIZE {
		size = INTEGRATED_ADDRESS_SIZE
	}
	res := make([]byte, size)

	start := 0
	end := BASE58_FULL_BLOCK_SIZE

	for end < addrSize {
		dec := []byte(base58.Encode(addr[start:end]))
		if len(dec) < BASE58_ENCODED_BLOCK_SIZE {
			dec = append(bytes.Repeat([]byte{0x31}, BASE58_ENCODED_BLOCK_SIZE-len(dec)), dec...)
		}

		copyToSlice(dec, res, BASE58_ENCODED_BLOCK_SIZE*start/BASE58_FULL_BLOCK_SIZE)

		start = end
		end = start + BASE58_FULL_BLOCK_SIZE
	}
	copyToSlice([]byte(base58.Encode(addr[start:])), res, BASE58_ENCODED_BLOCK_SIZE*start/BASE58_FULL_BLOCK_SIZE)

	return res
}

// Validates if the given addr []byte is a valid Monero integrated/sub-/primary decoded representation and encodes it
func EncodeMoneroAddress(addr []byte) ([]byte, error) {
	addrSize := len(addr)
	if addrSize != ADDRESS_DECODED_SIZE && addrSize != INTEGRATED_ADDRESS_DECODED_SIZE {
		return nil, invalid_address_err
	} else if !verifyMoneroAddressChecksumHelper(addr) {
		return nil, invalid_address_err
	}

	return encodeMoneroAddressBase58Helper(addr), nil
}

// Generates a Monero subaddress base on the primary private view and public spend keys, NetworkType, major and minor indices
func GenerateSubaddress(viewKey *PrivateKey, spendKey *PublicKey, major, minor uint32, nt NetworkType) (*SubAddress, error) {
	index := append(uint32ToLittleEndianBytes(major), uint32ToLittleEndianBytes(minor)...)

	Shash, err := Keccak256Hash(append(append(subaddr_prefix, viewKey.Bytes()...), index...))
	if err != nil {
		return nil, err
	}

	Sscalar, err := keccak256HashToScalar(Shash)
	if err != nil {
		return nil, err
	}

	Si := new(edwards25519.Point).ScalarBaseMult(Sscalar)
	Si.Add(Si, spendKey.key)

	Vi := new(edwards25519.Point).ScalarMult(viewKey.key, Si)

	dec := make([]byte, ADDRESS_DECODED_SIZE)
	pref, err := GetPrefix(nt, Sub)
	if err != nil {
		return nil, err
	}
	dec[0] = pref
	copyToSlice(Si.Bytes(), dec, 1)
	copyToSlice(Vi.Bytes(), dec, 1+KEY_SIZE)

	csum, err := Keccak256Hash(dec[:len(dec)-CHECKSUM_SIZE])
	if err != nil {
		return nil, err
	}
	csum = csum[:CHECKSUM_SIZE]
	copyToSlice(csum, dec, 1+2*KEY_SIZE)

	return &SubAddress{address: address{addr: dec}}, nil
}

/********************************************** Parsing Related Mehtods ***************************************************/

func parseExtraChecksHelper(extra []byte, tags []byte, size int) error {
	tag := []byte{extra[0]}
	if !bytes.Contains(tags, tag) {
		return errors.New("Invalid extra tag: " + hex.EncodeToString(tag))
	}

	if len(extra) < size {
		return errors.New("Invalid extra size: " + string(len(extra)))
	}

	return nil
}

// Parses the extra field and returns a tx pub key
func GetTxPublicKeyFromExtra(extra []byte) (*PublicKey, error) {
	if err := parseExtraChecksHelper(extra, default_extra_tags, 33); err != nil {
		return nil, err
	}

	key, err := new(edwards25519.Point).SetBytes(extra[1:33])
	if err != nil {
		return nil, err
	}

	return &PublicKey{key: key}, nil
}

// Parses the extra field and returns a payment id
func GetPaymentIdFromExtra(extra []byte) ([]byte, error) {
	if err := parseExtraChecksHelper(extra, default_extra_tags, 44); err != nil {
		return nil, err
	}

	return extra[36:44], nil
}

// Parses the extra field and returns a tx pub key (or nil if error) and a payment id (or nil if error)
func ParseExtra(extra []byte) (txKey *PublicKey, payId []byte) {
	txKey, err := GetTxPublicKeyFromExtra(extra)
	if err != nil {
		txKey = nil
	}

	payId, err = GetPaymentIdFromExtra(extra)
	if err != nil {
		payId = nil
	}

	return txKey, payId
}

func parseJson[R any](data []byte) (*R, error) {
	var result R
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func ParseJsonString[R any](str string) (*R, error) {
	return parseJson[R]([]byte(str))
}

func ParseResponse[R any](body io.Reader) (*R, error) {
	data, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}

	return parseJson[R](data)
}

/********************************************** Hash Related Mehtods ***************************************************/

// Returns a hash calculated by the Keccak256 hash algorithm
func Keccak256Hash(data []byte) ([]byte, error) {
	hash := sha3.NewLegacyKeccak256()

	_, err := hash.Write(data)
	if err != nil {
		return nil, err
	}

	return hash.Sum(nil), nil
}

func keccak256HashToScalar(hash []byte) (*edwards25519.Scalar, error) {
	if len(hash) != 32 {
		return nil, errors.New("invalid Keccak256Hash")
	}

	u := make([]byte, 64)
	copyToSlice(hash, u, 0)

	res, err := new(edwards25519.Scalar).SetUniformBytes(u)
	if err != nil {
		return nil, err
	}

	return res, nil
}

/********************************************** Other Mehtods ***************************************************/

// Converts the raw atomic XMR balance to a more human readable format.
func XMRToDecimal(xmr uint64) string {
	str0 := fmt.Sprintf("%013d", xmr)
	l := len(str0)
	return str0[:l-12] + "." + str0[l-12:]
}

// Converts the raw atomic XMR to a float64
func XMRToFloat64(xmr uint64) float64 {
	return float64(xmr) / ATOMIC_UNIT
}

// Converts the float64 to a raw atomic XMR
func Float64ToXMR(xmr float64) uint64 {
	return uint64(xmr * ATOMIC_UNIT)
}

// NewPaymentID64 generates a 64 bit payment ID (hex encoded).
// With 64 bit IDs, there is a non-negligible chance of a collision
// if they are randomly generated. It is up to recipients generating
// them to sanity check for uniqueness.
//
// 1 million IDs at 64-bit (simplified): 1,000,000^2 / (2^64 * 2) = ~1/36,893,488 so
// there is a 50% chance a collision happens around 5.06 billion IDs generated.
func NewPaymentID64() []byte {
	buf := make([]byte, 8)
	rand.Read(buf)
	return buf
}

// NewPaymentID256 generates a 256 bit payment ID (hex encoded).
func NewPaymentID256() []byte {
	buf := make([]byte, 32)
	rand.Read(buf)
	return buf
}

// Returns a NetworkType and an AddressType base on the Monero address prefix
func GetNetworkTypeAndAddressType(prefix byte) (NetworkType, AddressType, error) {
	switch prefix {
	case 0x12:
		return Mainnet, Primary, nil
	case 0x2a:
		return Mainnet, Sub, nil
	case 0x13:
		return Mainnet, Integrated, nil

	case 0x18:
		return Stagenet, Primary, nil
	case 0x24:
		return Stagenet, Sub, nil
	case 0x19:
		return Stagenet, Integrated, nil

	case 0x35:
		return Testnet, Primary, nil
	case 0x3f:
		return Testnet, Sub, nil
	case 0x36:
		return Testnet, Integrated, nil

	default:
		return 255, 255, errors.New("invalid prefix: " + hex.EncodeToString([]byte{prefix}))
	}
}

// Returns a NetworkType and an AddressType base on the Monero address prefix
func GetPrefix(nt NetworkType, at AddressType) (byte, error) {
	switch nt {
	case Mainnet:
		switch at {
		case Primary:
			return 0x12, nil

		case Sub:
			return 0x2a, nil

		case Integrated:
			return 0x13, nil

		default:
			return 255, errors.New("invalid AddressType: " + string(at))
		}

	case Stagenet:
		switch at {
		case Primary:
			return 0x18, nil

		case Sub:
			return 0x24, nil

		case Integrated:
			return 0x19, nil

		default:
			return 255, errors.New("invalid AddressType: " + string(at))
		}

	case Testnet:
		switch at {
		case Primary:
			return 0x35, nil

		case Sub:
			return 0x3f, nil

		case Integrated:
			return 0x36, nil

		default:
			return 255, errors.New("invalid AddressType: " + string(at))
		}

	default:
		return 255, errors.New("invalid NetworkType: " + string(nt))
	}
}

func copyToSlice[T any](src, dst []T, from int) {
	dstSize := len(dst)
	srcSize := len(src)

	if from < 0 || from >= dstSize {
		return
	}

	for i := 0; from+i < dstSize && i < srcSize; i++ {
		dst[from+i] = src[i]
	}
}

func uint32ToLittleEndianBytes(v uint32) []byte {
	r := make([]byte, 4)
	binary.LittleEndian.PutUint32(r, v)

	return r
}

// Non fixed sized
func uintToLittleEndianBytes(v uint64) []byte {
	size := 1
	if (v >> 56) != 0 {
		size = 8
	} else if (v >> 48) != 0 {
		size = 7
	} else if (v >> 40) != 0 {
		size = 6
	} else if (v >> 32) != 0 {
		size = 5
	} else if (v >> 24) != 0 {
		size = 4
	} else if (v >> 16) != 0 {
		size = 3
	} else if (v >> 8) != 0 {
		size = 2
	}

	r := make([]byte, size)
	for i := 0; i < size; i++ {
		r[size-1-i] = byte(v >> (8 * (size - 1 - i)))
	}

	return r
}

func xor(a, b []byte) []byte {
	r := make([]byte, len(a))
	for i := 0; i < len(a); i++ {
		r[i] = a[i] ^ b[i]
	}

	return r
}

func mod(a, m int) int {
	if a < 0 {
		a = -a
		if a < m {
			return m - a
		}
		return a % m
	}

	return a % m
}
