package hole

import (
	"fmt"
	"math/rand/v2"
	"strconv"
	"strings"
)

func computeZeckendorf(n uint64) string {
	// a list of fibonacci values less than 2^64
	fib := []uint64{12200160415121876738, 7540113804746346429, 4660046610375530309, 2880067194370816120, 1779979416004714189, 1100087778366101931, 679891637638612258, 420196140727489673, 259695496911122585, 160500643816367088, 99194853094755497, 61305790721611591, 37889062373143906, 23416728348467685, 14472334024676221, 8944394323791464, 5527939700884757, 3416454622906707, 2111485077978050, 1304969544928657, 806515533049393, 498454011879264, 308061521170129, 190392490709135, 117669030460994, 72723460248141, 44945570212853, 27777890035288, 17167680177565, 10610209857723, 6557470319842, 4052739537881, 2504730781961, 1548008755920, 956722026041, 591286729879, 365435296162, 225851433717, 139583862445, 86267571272, 53316291173, 32951280099, 20365011074, 12586269025, 7778742049, 4807526976, 2971215073, 1836311903, 1134903170, 701408733, 433494437, 267914296, 165580141, 102334155, 63245986, 39088169, 24157817, 14930352, 9227465, 5702887, 3524578, 2178309, 1346269, 832040, 514229, 317811, 196418, 121393, 75025, 46368, 28657, 17711, 10946, 6765, 4181, 2584, 1597, 987, 610, 377, 233, 144, 89, 55, 34, 21, 13, 8, 5, 3, 2, 1}

	out := make([]uint64, 0)

	for _, f := range fib {
		if f > n {
			continue
		}
		n -= f
		out = append(out, f)
	}

	return strings.Trim(strings.Replace(fmt.Sprint(out), " ", " + ", -1), "[]")
}

func zeckendorfRepresentation() []Run {
	fixedCases := []test{
		{"64", "55 + 8 + 1"},
		{"89", "89"},
		{"144", "144"},
		{"701408733", "701408733"},
		{"7", "5 + 2"},
		{"27777890035290", "27777890035288 + 2"},
		{"44945838127149", "44945570212853 + 267914296"},
		{"165580140", "102334155 + 39088169 + 14930352 + 5702887 + 2178309 + 832040 + 317811 + 121393 + 46368 + 17711 + 6765 + 2584 + 987 + 377 + 144 + 55 + 21 + 8 + 3 + 1"},
	}

	tests := make([]test, 40)

	for i := uint64(0); i < 20; i++ {
		n := rand.Uint64()
		tests[i] = test{
			strconv.FormatUint(n, 10),
			computeZeckendorf(n),
		}
	}

	for i := uint64(20); i < 40; i++ {
		n := uint64(rand.Uint32())
		tests[i] = test{
			strconv.FormatUint(n, 10),
			computeZeckendorf(n),
		}
	}

	tests = append(fixedCases, tests...)

	return outputTests(shuffle(tests))
}
