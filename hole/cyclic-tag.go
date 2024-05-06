package hole

import (
	"math/rand/v2"
	"strings"
)

func generate_case() string {
	var j int
	var choice = []byte{';', '0', '0', '0', '0', '0', '1', '1', '1', '1', '1'}
	var t strings.Builder
	result.WriteByte('1')
	var len = rand.IntN(6) +1 
	for j = 0; j < len; j++ {
		t.WriteByte(choice[rand.IntN(11)])
	}
	t.WriteByte(';')
	len = rand.IntN(6) + 1
	for j = 0; j < len; j++ {
		t.WriteByte(choice[rand.IntN(11)])
	}
	t.WriteByte(';')
	return t.String()
}
func pop[T any](a *[]T)         { *a = (*a)[1:] }
func push[T any](a *[]T, val T) { *a = append(*a, val) }
func interpret_ct(prog string) string {
	var data = []byte{'1'}
	var id = 0
	var count = 0
	var result strings.Builder
	for {
		var c = prog[id]
		var popped = data[0]
		if c == ';' {
			pop(&data)
		} else {
			popped = '2'
		}
		if c == '0' && '1' == data[0] {
			push(&data, '0')
		}
		if c == '1' && '1' == data[0] {
			push(&data, '1')
		}
		if popped != '2' {
			result.WriteByte(popped)
			count = count + 1
		}
		id = (id + 1) % len(prog)
		if count == 200 || 0 == len(data) { //We can halt after 200 removals, as the hole only requires us to print the first 200 removals
			break
		}
	}
	return result.String()
}
func cyclicTag() {
	var i int
	args := []string{ //Fixed halting testcases
		"1010;;1;;",
		"10000;010;1;;;",
		"1;0011;0;0;00",
		"111;000;000;111;000;000;000;000;000;000;000;000;000;000;000",
		"10;100;0;01;10000;;",
	}

	outs := []string{
		"11010",
		"1100000101",
		"1100110010011001000100",
		"1111000000111000000000",
		"110100011000",
	}
	for range 95 {
		cur_case := generate_case()
		result := interpret_ct(cur_case)
		for len(result) < 200 && i >= 30 { //Make at least 70 nonhalting cases
			cur_case = generate_case()
			result = interpret_ct(cur_case)
		}
		args = append(args, cur_case)
		outs = append(outs, result)
	}
	rand.Shuffle(len(args), func(i, j int) {
		args[i], args[j] = args[j], args[i]
		outs[i], outs[j] = outs[j], outs[i]
	})

	return []Run{{Args: args, Answer: strings.Join(outs, "\n")}}
}
