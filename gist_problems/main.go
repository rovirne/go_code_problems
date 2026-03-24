package main

import "fmt"

func main() {
	// 1. Объявлять переменные через var и напрямую.
	var check bool = true
	number := 5
	fmt.Println(check, number)

	// 2. Объявлять структуру с полями разных типов, уметь ее инстанциировать в коде.
	type Person struct {
		name string
		age  int
	}
	p := Person{name: "alex", age: 21}
	fmt.Println(p.name, p.age)

	// 3. Уметь ембедить структуру в структуру и инстанциировать ее в коде.

	type Employe struct {
		Person
		occupation string
	}
	e := Employe{Person: p, occupation: "miner"}
	fmt.Println(e.age, e.occupation)

	// 4. Уметь инстанциировать массив, слайс и мапу, из литералов и через make.
	// То же самое для разных типов данных в них, в том числе структур,
	// в том числе пустых структур.
	array := [3]int{4, 5, 6} // массив
	fmt.Println(array)

	s := []int{1, 2, 3} // слайс
	fmt.Println(s)
	s = make([]int, 6)
	fmt.Println(s)

	m := map[string]int{ // мапа
		"a": 1,
		"b": 2,
	}
	fmt.Println(m["a"])
	m = make(map[string]int)
	m["a"] = 153
	fmt.Println(m["a"])

	type Emty struct{}

}
