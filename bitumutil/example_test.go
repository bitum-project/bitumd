package bitumutil_test

import (
	"fmt"
	"math"

	"github.com/bitum-project/bitumd/bitumutil"
)

func ExampleAmount() {

	a := bitumutil.Amount(0)
	fmt.Println("Zero Atom:", a)

	a = bitumutil.Amount(1e8)
	fmt.Println("100,000,000 Atoms:", a)

	a = bitumutil.Amount(1e5)
	fmt.Println("100,000 Atoms:", a)
	// Output:
	// Zero Atom: 0 BITUM
	// 100,000,000 Atoms: 1 BITUM
	// 100,000 Atoms: 0.001 BITUM
}

func ExampleNewAmount() {
	amountOne, err := bitumutil.NewAmount(1)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(amountOne) //Output 1

	amountFraction, err := bitumutil.NewAmount(0.01234567)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(amountFraction) //Output 2

	amountZero, err := bitumutil.NewAmount(0)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(amountZero) //Output 3

	amountNaN, err := bitumutil.NewAmount(math.NaN())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(amountNaN) //Output 4

	// Output: 1 BITUM
	// 0.01234567 BITUM
	// 0 BITUM
	// invalid coin amount
}

func ExampleAmount_unitConversions() {
	amount := bitumutil.Amount(44433322211100)

	fmt.Println("Atom to kCoin:", amount.Format(bitumutil.AmountKiloCoin))
	fmt.Println("Atom to Coin:", amount)
	fmt.Println("Atom to MilliCoin:", amount.Format(bitumutil.AmountMilliCoin))
	fmt.Println("Atom to MicroCoin:", amount.Format(bitumutil.AmountMicroCoin))
	fmt.Println("Atom to Atom:", amount.Format(bitumutil.AmountAtom))

	// Output:
	// Atom to kCoin: 444.333222111 kBITUM
	// Atom to Coin: 444333.222111 BITUM
	// Atom to MilliCoin: 444333222.111 mBITUM
	// Atom to MicroCoin: 444333222111 μBITUM
	// Atom to Atom: 44433322211100 Atom
}
