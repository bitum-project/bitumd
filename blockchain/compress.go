// Copyright (c) 2015-2016 The btcsuite developers
// Copyright (c) 2015-2018 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package blockchain

import (
	"fmt"

	"github.com/bitum-project/bitumd/blockchain/stake"
	"github.com/bitum-project/bitumd/bitumec/secp256k1"
	"github.com/bitum-project/bitumd/txscript"
)

// currentCompressionVersion is the current script compression version of the
// database.
const currentCompressionVersion = 1

// -----------------------------------------------------------------------------
// A variable length quantity (VLQ) is an encoding that uses an arbitrary number
// of binary octets to represent an arbitrarily large integer.  The scheme
// employs a most significant byte (MSB) base-128 encoding where the high bit in
// each byte indicates whether or not the byte is the final one.  In addition,
// to ensure there are no redundant encodings, an offset is subtracted every
// time a group of 7 bits is shifted out.  Therefore each integer can be
// represented in exactly one way, and each representation stands for exactly
// one integer.
//
// Another nice property of this encoding is that it provides a compact
// representation of values that are typically used to indicate sizes.  For
// example, the values 0 - 127 are represented with a single byte, 128 - 16511
// with two bytes, and 16512 - 2113663 with three bytes.
//
// While the encoding allows arbitrarily large integers, it is artificially
// limited in this code to an unsigned 64-bit integer for efficiency purposes.
//
// Example encodings:
//           0 -> [0x00]
//         127 -> [0x7f]                 * Max 1-byte value
//         128 -> [0x80 0x00]
//         129 -> [0x80 0x01]
//         255 -> [0x80 0x7f]
//         256 -> [0x81 0x00]
//       16511 -> [0xff 0x7f]            * Max 2-byte value
//       16512 -> [0x80 0x80 0x00]
//       32895 -> [0x80 0xff 0x7f]
//     2113663 -> [0xff 0xff 0x7f]       * Max 3-byte value
//   270549119 -> [0xff 0xff 0xff 0x7f]  * Max 4-byte value
//      2^64-1 -> [0x80 0xfe 0xfe 0xfe 0xfe 0xfe 0xfe 0xfe 0xfe 0x7f]
//
// References:
//   https://en.wikipedia.org/wiki/Variable-length_quantity
//   http://www.codecodex.com/wiki/Variable-Length_Integers
// -----------------------------------------------------------------------------

// serializeSizeVLQ returns the number of bytes it would take to serialize the
// passed number as a variable-length quantity according to the format described
// above.
func serializeSizeVLQ(n uint64) int {
	size := 1
	for ; n > 0x7f; n = (n >> 7) - 1 {
		size++
	}

	return size
}

// putVLQ serializes the provided number to a variable-length quantity according
// to the format described above and returns the number of bytes of the encoded
// value.  The result is placed directly into the passed byte slice which must
// be at least large enough to handle the number of bytes returned by the
// serializeSizeVLQ function or it will panic.
func putVLQ(target []byte, n uint64) int {
	offset := 0
	for ; ; offset++ {
		// The high bit is set when another byte follows.
		highBitMask := byte(0x80)
		if offset == 0 {
			highBitMask = 0x00
		}

		target[offset] = byte(n&0x7f) | highBitMask
		if n <= 0x7f {
			break
		}
		n = (n >> 7) - 1
	}

	// Reverse the bytes so it is MSB-encoded.
	for i, j := 0, offset; i < j; i, j = i+1, j-1 {
		target[i], target[j] = target[j], target[i]
	}

	return offset + 1
}

// deserializeVLQ deserializes the provided variable-length quantity according
// to the format described above.  It also returns the number of bytes
// deserialized.
func deserializeVLQ(serialized []byte) (uint64, int) {
	var n uint64
	var size int
	for _, val := range serialized {
		size++
		n = (n << 7) | uint64(val&0x7f)
		if val&0x80 != 0x80 {
			break
		}
		n++
	}

	return n, size
}

// -----------------------------------------------------------------------------
// In order to reduce the size of stored scripts, a domain specific compression
// algorithm is used which recognizes standard scripts and stores them using
// less bytes than the original script.  The compression algorithm used here was
// obtained from Bitcoin Core, so all credits for the algorithm go to it.
//
// The general serialized format is:
//
//   <script size or type><script data>
//
//   Field                 Type     Size
//   script size or type   VLQ      variable
//   script data           []byte   variable
//
// The specific serialized format for each recognized standard script is:
//
// - Pay-to-pubkey-hash: (21 bytes) - <0><20-byte pubkey hash>
// - Pay-to-script-hash: (21 bytes) - <1><20-byte script hash>
// - Pay-to-pubkey**:    (33 bytes) - <2, 3, 4, or 5><32-byte pubkey X value>
//   2, 3 = compressed pubkey with bit 0 specifying the y coordinate to use
//   4, 5 = uncompressed pubkey with bit 0 specifying the y coordinate to use
//   ** Only valid public keys starting with 0x02, 0x03, and 0x04 are supported.
//
// Any scripts which are not recognized as one of the aforementioned standard
// scripts are encoded using the general serialized format and encode the script
// size as the sum of the actual size of the script and the number of special
// cases.
// -----------------------------------------------------------------------------

// The following constants specify the special constants used to identify a
// special script type in the domain-specific compressed script encoding.
//
// NOTE: This section specifically does not use iota since these values are
// serialized and must be stable for long-term storage.
const (
	// cstPayToPubKeyHash identifies a compressed pay-to-pubkey-hash script.
	cstPayToPubKeyHash = 0

	// cstPayToScriptHash identifies a compressed pay-to-script-hash script.
	cstPayToScriptHash = 1

	// cstPayToPubKeyCompEven identifies a compressed pay-to-pubkey script to
	// a compressed pubkey whose y coordinate is not odd.
	cstPayToPubKeyCompEven = 2

	// cstPayToPubKeyCompOdd identifies a compressed pay-to-pubkey script to
	// a compressed pubkey whose y coordinate is odd.
	cstPayToPubKeyCompOdd = 3

	// cstPayToPubKeyUncompEven identifies a compressed pay-to-pubkey script to
	// an uncompressed pubkey whose y coordinate is not odd when compressed.
	cstPayToPubKeyUncompEven = 4

	// cstPayToPubKeyUncompOdd identifies a compressed pay-to-pubkey script to
	// an uncompressed pubkey whose y coordinate is odd when compressed.
	cstPayToPubKeyUncompOdd = 5

	// numSpecialScripts is the number of special scripts possibly recognized
	// by the domain-specific script compression algorithm. It is one more
	// than half the number required to overflow a single byte in VLQ format
	// (127). All scripts prefixed 64 and higher for their size are considered
	// uncompressed scripts that are stored uncompressed. Because only 5
	// special script types are currently stored by Bitum, there is a large
	// amount of room for future upgrades to the compression algorithm with
	// scripts that are common, such as those for the staking system.
	numSpecialScripts = 64
)

// extractPubKeyHash extracts a pubkey hash that is being paid from the passed
// public key script if it is a standard pay-to-pubkey-hash script.  It will
// return nil otherwise.
func extractPubKeyHash(script []byte) []byte {
	if len(script) == 25 && script[0] == txscript.OP_DUP &&
		script[1] == txscript.OP_HASH160 &&
		script[2] == txscript.OP_DATA_20 &&
		script[23] == txscript.OP_EQUALVERIFY &&
		script[24] == txscript.OP_CHECKSIG {

		return script[3:23]
	}

	return nil
}

// isPubKeyHash returns whether or not the passed public key script is a
// standard pay-to-pubkey-hash script.
func isPubKeyHash(script []byte) bool {
	return extractPubKeyHash(script) != nil
}

// extractScriptHash extracts a script hash that is being paid from the passed
// public key script if it is a standard pay-to-script-hash script.  It will
// return nil otherwise.
func extractScriptHash(script []byte) []byte {
	if len(script) == 23 && script[0] == txscript.OP_HASH160 &&
		script[1] == txscript.OP_DATA_20 &&
		script[22] == txscript.OP_EQUAL {

		return script[2:22]
	}

	return nil
}

// isScriptHash returns whether or not the passed public key script is a
// standard pay-to-script-hash script.
func isScriptHash(script []byte) bool {
	return extractScriptHash(script) != nil
}

// isPubKey returns whether or not the passed public key script is a standard
// pay-to-pubkey script that pays to a valid compressed or uncompressed public
// key along with the serialized pubkey it is paying to if it is.
//
// NOTE: This function ensures the public key is actually valid since the
// compression algorithm requires valid pubkeys.  It does not support hybrid
// pubkeys.  This means that even if the script has the correct form for a
// pay-to-pubkey script, this function will only return true when it is paying
// to a valid compressed or uncompressed pubkey.
func isPubKey(script []byte) (bool, []byte) {
	// Pay-to-compressed-pubkey script.
	if len(script) == 35 && script[0] == txscript.OP_DATA_33 &&
		script[34] == txscript.OP_CHECKSIG && (script[1] == 0x02 ||
		script[1] == 0x03) {

		// Ensure the public key is valid.
		serializedPubKey := script[1:34]
		_, err := secp256k1.ParsePubKey(serializedPubKey)
		if err == nil {
			return true, serializedPubKey
		}
	}

	// Pay-to-uncompressed-pubkey script.
	if len(script) == 67 && script[0] == txscript.OP_DATA_65 &&
		script[66] == txscript.OP_CHECKSIG && script[1] == 0x04 {

		// Ensure the public key is valid.
		serializedPubKey := script[1:66]
		_, err := secp256k1.ParsePubKey(serializedPubKey)
		if err == nil {
			return true, serializedPubKey
		}
	}

	return false, nil
}

// compressedScriptSize returns the number of bytes the passed script would take
// when encoded with the domain specific compression algorithm described above.
func compressedScriptSize(scriptVersion uint16, pkScript []byte,
	compressionVersion uint32) int {
	// Pay-to-pubkey-hash or pay-to-script-hash script.
	if isPubKeyHash(pkScript) || isScriptHash(pkScript) {
		return 21
	}

	// Pay-to-pubkey (compressed or uncompressed) script.
	if valid, _ := isPubKey(pkScript); valid {
		return 33
	}

	// When none of the above special cases apply, encode the script as is
	// preceded by the sum of its size and the number of special cases
	// encoded as a variable length quantity.
	return serializeSizeVLQ(uint64(len(pkScript)+numSpecialScripts)) +
		len(pkScript)
}

// decodeCompressedScriptSize treats the passed serialized bytes as a compressed
// script, possibly followed by other data, and returns the number of bytes it
// occupies taking into account the special encoding of the script size by the
// domain specific compression algorithm described above.
func decodeCompressedScriptSize(serialized []byte, compressionVersion uint32) int {
	scriptSize, bytesRead := deserializeVLQ(serialized)
	if bytesRead == 0 {
		return 0
	}

	switch scriptSize {
	case cstPayToPubKeyHash:
		return 21

	case cstPayToScriptHash:
		return 21

	case cstPayToPubKeyCompEven, cstPayToPubKeyCompOdd,
		cstPayToPubKeyUncompEven, cstPayToPubKeyUncompOdd:
		return 33
	}

	scriptSize -= numSpecialScripts
	scriptSize += uint64(bytesRead)
	return int(scriptSize)
}

// putCompressedScript compresses the passed script according to the domain
// specific compression algorithm described above directly into the passed
// target byte slice.  The target byte slice must be at least large enough to
// handle the number of bytes returned by the compressedScriptSize function or
// it will panic.
func putCompressedScript(target []byte, scriptVersion uint16, pkScript []byte,
	compressionVersion uint32) int {
	if len(target) == 0 {
		target[0] = 0x00
		return 1
	}

	// Pay-to-pubkey-hash script.
	if hash := extractPubKeyHash(pkScript); hash != nil {
		target[0] = cstPayToPubKeyHash
		copy(target[1:21], hash)
		return 21
	}

	// Pay-to-script-hash script.
	if hash := extractScriptHash(pkScript); hash != nil {
		target[0] = cstPayToScriptHash
		copy(target[1:21], hash)
		return 21
	}

	// Pay-to-pubkey (compressed or uncompressed) script.
	if valid, serializedPubKey := isPubKey(pkScript); valid {
		pubKeyFormat := serializedPubKey[0]
		switch pubKeyFormat {
		case 0x02, 0x03:
			if pubKeyFormat == 0x02 {
				target[0] = cstPayToPubKeyCompEven
			}
			if pubKeyFormat == 0x03 {
				target[0] = cstPayToPubKeyCompOdd
			}
			copy(target[1:33], serializedPubKey[1:33])
			return 33
		case 0x04:
			// Encode the oddness of the serialized pubkey into the
			// compressed script type.
			target[0] = cstPayToPubKeyUncompEven
			if (serializedPubKey[64] & 0x01) == 0x01 {
				target[0] = cstPayToPubKeyUncompOdd
			}
			copy(target[1:33], serializedPubKey[1:33])
			return 33
		}
	}

	// When none of the above special cases apply, encode the unmodified
	// script preceded by the script version, the sum of its size and
	// the number of special cases encoded as a variable length quantity.
	encodedSize := uint64(len(pkScript) + numSpecialScripts)
	vlqSizeLen := putVLQ(target, encodedSize)
	copy(target[vlqSizeLen:], pkScript)
	return vlqSizeLen + len(pkScript)
}

// decompressScript returns the original script obtained by decompressing the
// passed compressed script according to the domain specific compression
// algorithm described above.
//
// NOTE: The script parameter must already have been proven to be long enough
// to contain the number of bytes returned by decodeCompressedScriptSize or it
// will panic.  This is acceptable since it is only an internal function.
func decompressScript(compressedPkScript []byte,
	compressionVersion uint32) []byte {
	// Empty scripts, specified by 0x00, are considered nil.
	if len(compressedPkScript) == 0 {
		return nil
	}

	// Decode the script size and examine it for the special cases.
	encodedScriptSize, bytesRead := deserializeVLQ(compressedPkScript)
	switch encodedScriptSize {
	// Pay-to-pubkey-hash script.  The resulting script is:
	// <OP_DUP><OP_HASH160><20 byte hash><OP_EQUALVERIFY><OP_CHECKSIG>
	case cstPayToPubKeyHash:
		pkScript := make([]byte, 25)
		pkScript[0] = txscript.OP_DUP
		pkScript[1] = txscript.OP_HASH160
		pkScript[2] = txscript.OP_DATA_20
		copy(pkScript[3:], compressedPkScript[bytesRead:bytesRead+20])
		pkScript[23] = txscript.OP_EQUALVERIFY
		pkScript[24] = txscript.OP_CHECKSIG
		return pkScript

	// Pay-to-script-hash script.  The resulting script is:
	// <OP_HASH160><20 byte script hash><OP_EQUAL>
	case cstPayToScriptHash:
		pkScript := make([]byte, 23)
		pkScript[0] = txscript.OP_HASH160
		pkScript[1] = txscript.OP_DATA_20
		copy(pkScript[2:], compressedPkScript[bytesRead:bytesRead+20])
		pkScript[22] = txscript.OP_EQUAL
		return pkScript

	// Pay-to-compressed-pubkey script.  The resulting script is:
	// <OP_DATA_33><33 byte compressed pubkey><OP_CHECKSIG>
	case cstPayToPubKeyCompEven, cstPayToPubKeyCompOdd:
		pkScript := make([]byte, 35)
		pkScript[0] = txscript.OP_DATA_33
		oddness := byte(0x02)
		if encodedScriptSize == cstPayToPubKeyCompOdd {
			oddness = 0x03
		}
		pkScript[1] = oddness
		copy(pkScript[2:], compressedPkScript[bytesRead:bytesRead+32])
		pkScript[34] = txscript.OP_CHECKSIG
		return pkScript

	// Pay-to-uncompressed-pubkey script.  The resulting script is:
	// <OP_DATA_65><65 byte uncompressed pubkey><OP_CHECKSIG>
	case cstPayToPubKeyUncompEven, cstPayToPubKeyUncompOdd:
		// Change the leading byte to the appropriate compressed pubkey
		// identifier (0x02 or 0x03) so it can be decoded as a
		// compressed pubkey.  This really should never fail since the
		// encoding ensures it is valid before compressing to this type.
		compressedKey := make([]byte, 33)
		oddness := byte(0x02)
		if encodedScriptSize == cstPayToPubKeyUncompOdd {
			oddness = 0x03
		}
		compressedKey[0] = oddness
		copy(compressedKey[1:], compressedPkScript[1:])
		key, err := secp256k1.ParsePubKey(compressedKey)
		if err != nil {
			return nil
		}

		pkScript := make([]byte, 67)
		pkScript[0] = txscript.OP_DATA_65
		copy(pkScript[1:], key.SerializeUncompressed())
		pkScript[66] = txscript.OP_CHECKSIG
		return pkScript
	}

	// When none of the special cases apply, the script was encoded using
	// the general format, so reduce the script size by the number of
	// special cases and return the unmodified script.
	scriptSize := int(encodedScriptSize - numSpecialScripts)
	pkScript := make([]byte, scriptSize)
	copy(pkScript, compressedPkScript[bytesRead:bytesRead+scriptSize])
	return pkScript
}

// -----------------------------------------------------------------------------
// In order to reduce the size of stored amounts, a domain specific compression
// algorithm is used which relies on there typically being a lot of zeroes at
// end of the amounts.  The compression algorithm used here was obtained from
// Bitcoin Core, so all credits for the algorithm go to it.
//
// While this is simply exchanging one uint64 for another, the resulting value
// for typical amounts has a much smaller magnitude which results in fewer bytes
// when encoded as variable length quantity.  For example, consider the amount
// of 0.1 BITUM which is 10000000 atoms.  Encoding 10000000 as a VarInt would take
// 4 bytes while encoding the compressed value of 8 as a VarInt only takes 1 byte.
//
// Essentially the compression is achieved by splitting the value into an
// exponent in the range [0-9] and a digit in the range [1-9], when possible,
// and encoding them in a way that can be decoded.  More specifically, the
// encoding is as follows:
// - 0 is 0
// - Find the exponent, e, as the largest power of 10 that evenly divides the
//   value up to a maximum of 9
// - When e < 9, the final digit can't be 0 so store it as d and remove it by
//   dividing the value by 10 (call the result n).  The encoded value is thus:
//   1 + 10*(9*n + d-1) + e
// - When e==9, the only thing known is the amount is not 0.  The encoded value
//   is thus:
//   1 + 10*(n-1) + e   ==   10 + 10*(n-1)
//
// Example encodings:
// (The numbers in parenthesis are the number of bytes when serialized as a VarInt)
//            0 (1) -> 0        (1)           *  0.00000000 BTC
//         1000 (2) -> 4        (1)           *  0.00001000 BTC
//        10000 (2) -> 5        (1)           *  0.00010000 BTC
//     12345678 (4) -> 111111101(4)           *  0.12345678 BTC
//     50000000 (4) -> 47       (1)           *  0.50000000 BTC
//    100000000 (4) -> 9        (1)           *  1.00000000 BTC
//    500000000 (5) -> 49       (1)           *  5.00000000 BTC
//   1000000000 (5) -> 10       (1)           * 10.00000000 BTC
// -----------------------------------------------------------------------------

// compressTxOutAmount compresses the passed amount according to the domain
// specific compression algorithm described above.
func compressTxOutAmount(amount uint64) uint64 {
	// No need to do any work if it's zero.
	if amount == 0 {
		return 0
	}

	// Find the largest power of 10 (max of 9) that evenly divides the
	// value.
	exponent := uint64(0)
	for amount%10 == 0 && exponent < 9 {
		amount /= 10
		exponent++
	}

	// The compressed result for exponents less than 9 is:
	// 1 + 10*(9*n + d-1) + e
	if exponent < 9 {
		lastDigit := amount % 10
		amount /= 10
		return 1 + 10*(9*amount+lastDigit-1) + exponent
	}

	// The compressed result for an exponent of 9 is:
	// 1 + 10*(n-1) + e   ==   10 + 10*(n-1)
	return 10 + 10*(amount-1)
}

// decompressTxOutAmount returns the original amount the passed compressed
// amount represents according to the domain specific compression algorithm
// described above.
func decompressTxOutAmount(amount uint64) uint64 {
	// No need to do any work if it's zero.
	if amount == 0 {
		return 0
	}

	// The decompressed amount is either of the following two equations:
	// x = 1 + 10*(9*n + d - 1) + e
	// x = 1 + 10*(n - 1)       + 9
	amount--

	// The decompressed amount is now one of the following two equations:
	// x = 10*(9*n + d - 1) + e
	// x = 10*(n - 1)       + 9
	exponent := amount % 10
	amount /= 10

	// The decompressed amount is now one of the following two equations:
	// x = 9*n + d - 1  | where e < 9
	// x = n - 1        | where e = 9
	var n uint64
	if exponent < 9 {
		lastDigit := amount%9 + 1
		amount /= 9
		n = amount*10 + lastDigit
	} else {
		n = amount + 1
	}

	// Apply the exponent.
	for ; exponent > 0; exponent-- {
		n *= 10
	}

	return n
}

// -----------------------------------------------------------------------------
// Compressed transaction outputs for UTXOS consist of an amount and a public
// key script both compressed using the domain specific compression algorithms
// previously described.
//
// The serialized format is:
//
//   <compressed amount><compressed script>
//
//   Field                 Type     Size
//     compressed amount   VLQ      variable
//     compressed script   []byte   variable
// -----------------------------------------------------------------------------

// compressedTxOutSize returns the number of bytes the passed transaction output
// fields would take when encoded with the format described above.  The
// preCompressed flag indicates the provided amount and script are already
// compressed.  This is useful since loaded utxo entries are not decompressed
// until the output is accessed.
func compressedTxOutSize(amount uint64, scriptVersion uint16, pkScript []byte,
	compressionVersion uint32, preCompressed bool, hasAmount bool) int {
	scriptVersionSize := serializeSizeVLQ(uint64(scriptVersion))
	if preCompressed && !hasAmount {
		return scriptVersionSize + len(pkScript)
	}
	if preCompressed && hasAmount {
		return scriptVersionSize + serializeSizeVLQ(compressTxOutAmount(amount)) +
			len(pkScript)
	}
	if !preCompressed && !hasAmount {
		return scriptVersionSize + compressedScriptSize(scriptVersion,
			pkScript, compressionVersion)
	}

	// if !preCompressed && hasAmount
	return scriptVersionSize + serializeSizeVLQ(compressTxOutAmount(amount)) +
		compressedScriptSize(scriptVersion, pkScript, compressionVersion)
}

// putCompressedTxOut potentially compresses the passed amount and script
// according to their domain specific compression algorithms and encodes them
// directly into the passed target byte slice with the format described above.
// The preCompressed flag indicates the provided amount and script are already
// compressed in which case the values are not modified.  This is useful since
// loaded utxo entries are not decompressed until the output is accessed.  The
// target byte slice must be at least large enough to handle the number of bytes
// returned by the compressedTxOutSize function or it will panic.
func putCompressedTxOut(target []byte, amount uint64, scriptVersion uint16,
	pkScript []byte, compressionVersion uint32, preCompressed bool,
	hasAmount bool) int {
	if preCompressed && hasAmount {
		offset := putVLQ(target, compressTxOutAmount(amount))
		offset += putVLQ(target[offset:], uint64(scriptVersion))
		copy(target[offset:], pkScript)
		return offset + len(pkScript)
	}
	if preCompressed && !hasAmount {
		offset := putVLQ(target, uint64(scriptVersion))
		copy(target[offset:], pkScript)
		return offset + len(pkScript)
	}
	if !preCompressed && !hasAmount {
		offset := putVLQ(target, uint64(scriptVersion))
		offset += putCompressedScript(target[offset:], scriptVersion, pkScript,
			compressionVersion)
		return offset
	}

	// if !preCompressed && hasAmount
	offset := putVLQ(target, compressTxOutAmount(amount))
	offset += putVLQ(target[offset:], uint64(scriptVersion))
	offset += putCompressedScript(target[offset:], scriptVersion, pkScript,
		compressionVersion)
	return offset
}

// decodeCompressedTxOut decodes the passed compressed txout, possibly followed
// by other data, into its compressed amount and compressed script and returns
// them along with the number of bytes they occupied.
func decodeCompressedTxOut(serialized []byte, compressionVersion uint32,
	hasAmount bool) (int64, uint16, []byte, int, error) {
	var amount int64
	var bytesRead int
	var offset int
	if hasAmount {
		// Deserialize the compressed amount and ensure there are bytes
		// remaining for the compressed script.
		var compressedAmount uint64
		compressedAmount, bytesRead = deserializeVLQ(serialized)
		if bytesRead >= len(serialized) {
			return 0, 0, nil, bytesRead, errDeserialize("unexpected end of " +
				"data after compressed amount")
		}
		amount = int64(decompressTxOutAmount(compressedAmount))
		offset += bytesRead
	}

	// Decode the script version.
	var scriptVersion uint64
	scriptVersion, bytesRead = deserializeVLQ(serialized[offset:])
	offset += bytesRead

	// Decode the compressed script size and ensure there are enough bytes
	// left in the slice for it.
	scriptSize := decodeCompressedScriptSize(serialized[offset:],
		compressionVersion)
	if scriptSize < 0 {
		return 0, 0, nil, offset, errDeserialize("negative script size")
	}
	if len(serialized[offset:]) < scriptSize {
		return 0, 0, nil, offset, errDeserialize(fmt.Sprintf("unexpected end of "+
			"data after script size (got %v, need %v)", len(serialized[offset:]),
			scriptSize))
	}

	// Make a copy of the compressed script so the original serialized data
	// can be released as soon as possible.
	compressedScript := make([]byte, scriptSize)
	copy(compressedScript, serialized[offset:offset+scriptSize])

	return amount, uint16(scriptVersion), compressedScript,
		offset + scriptSize, nil
}

// -----------------------------------------------------------------------------
// Bitum specific transaction encoding flags
//
// Details about a transaction needed to determine how it may be spent
// according to consensus rules are given by these flags.
//
// The following details are encoded into a single byte, where the index
// of the bit is given in zeroeth order:
//     0: Is coinbase
//     1: Has an expiry
//   2-3: Transaction type
//     4: Fully spent
//   5-7: Unused
//
// 0, 1, and 4 are bit flags, while the transaction type is encoded with a bitmask
// and used to describe the underlying int.
//
// The fully spent flag should always come as the *last* flag (highest bit index)
// in this data type should flags be updated to include more rules in the future,
// such as rules governing new script OP codes. This ensures that we may still use
// these flags in the UTX serialized data without consequence, where the last flag
// indicating fully spent will always be zeroed.
//
// -----------------------------------------------------------------------------

const (
	// txTypeBitmask describes the bitmask that yields the 3rd and 4th bits
	// from the flags byte.
	txTypeBitmask = 0x0c

	// txTypeShift is the number of bits to shift falgs to the right to yield the
	// correct integer value after applying the bitmask with AND.
	txTypeShift = 2
)

// encodeFlags encodes transaction flags into a single byte.
func encodeFlags(isCoinBase bool, hasExpiry bool, txType stake.TxType, fullySpent bool) byte {
	b := uint8(txType)
	b <<= txTypeShift

	if isCoinBase {
		b |= 0x01 // Set bit 0
	}
	if hasExpiry {
		b |= 0x02 // Set bit 1
	}
	if fullySpent {
		b |= 0x10 // Set bit 4
	}

	return b
}

// decodeFlags decodes transaction flags from a single byte into their respective
// data types.
func decodeFlags(b byte) (bool, bool, stake.TxType, bool) {
	isCoinBase := b&0x01 != 0
	hasExpiry := b&(1<<1) != 0
	fullySpent := b&(1<<4) != 0
	txType := stake.TxType((b & txTypeBitmask) >> txTypeShift)

	return isCoinBase, hasExpiry, txType, fullySpent
}

// decodeFlagsFullySpent decodes whether or not a transaction was fully spent.
func decodeFlagsFullySpent(b byte) bool {
	return b&(1<<4) != 0
}