package main

import (
	"fmt"
	"sync"
)

func main() {
	// 1. Объявлять переменные через var и напрямую.
	{
		fmt.Println("1-th example")
		var check bool = true
		number := 5
		fmt.Println(check, number)
	}
	// 2. Объявлять структуру с полями разных типов, уметь ее инстанциировать в коде.
	{
		fmt.Println("\n2-th example")
		type Person struct {
			name string
			age  int
		}
		p := Person{name: "alex", age: 21}
		fmt.Println(p.name, p.age)
	}

	// 3. Уметь ембедить структуру в структуру и инстанциировать ее в коде.
	{
		fmt.Println("\n3-th example")
		type Person struct {
			name string
			age  int
		}
		p := Person{name: "alex", age: 21}

		type Employe struct {
			Person
			occupation string
		}
		e := Employe{Person: p, occupation: "miner"}
		fmt.Println(e.age, e.occupation)
	}

	// 4. Уметь инстанциировать массив, слайс и мапу, из литералов и через make.
	// То же самое для разных типов данных в них, в том числе структур,
	// в том числе пустых структур.
	{
		fmt.Println("\n4-th example")
		array := [3]int{4, 5, 6} // массив порядок сохр
		fmt.Println(array)

		s := []int{1, 2, 3} // слайс порядок сохр
		fmt.Println(s)
		s = make([]int, 6)
		fmt.Println(s)

		m := map[string]int{ // мапа порядок не сохр
			"a": 1,
			"b": 2,
		}
		fmt.Println(m["a"])
		m = make(map[string]int)
		m["a"] = 153
		fmt.Println(m["a"])

		type Emty struct{}
	}
	// 5. Итерация по слайсу (с индексом),
	// мапе (итерация по ключам, значениям, использование ok).
	// Доставание элементов по индексу или ключу. Добавление элементов.
	{
		fmt.Println("\n5-th example")
		s := []int{45, 46, 47} // slice
		for index, val := range s {
			fmt.Println(index, val)
		}
		s = append(s, 48)
		fmt.Println(s[3])

		m := map[string]int{ // мапа
			"a": 1,
			"b": 2,
		}
		for v, k := range m {
			fmt.Println(v, k)
		}
		value, ok := m["c"]
		fmt.Println(value, ok)

		m["c"] = 3
		fmt.Println(m["c"])
	}

	// 6. Создание каналов, с буфером и без. Закрытие канала.
	//  Запись и чтение из канала (range, select, напрямую, использование ok).
	// Канал пустых структур, кастомных структур. Паттерн Fan in.
	{
		fmt.Println("\n6-th example")
		// var wg sync.WaitGroup
		check := make(chan int, 3) // 3 - буфер
		check <- 1
		check <- 2
		close(check)
		for i := range check {
			fmt.Println(i)
		}
		_, ok := <-check
		if !ok {
			fmt.Println(ok)
		}

		ch1 := make(chan int)
		ch2 := make(chan int)

		go func() {
			ch1 <- 1
		}()

		go func() {
			ch2 <- 2
		}()

		select {
		case v := <-ch1:
			fmt.Println("ch1:", v)
		case v := <-ch2:
			fmt.Println("ch2:", v)
		}

		done := make(chan struct{})
		go func() {
			//do some work
			done <- struct{}{}
		}()
		<-done
		fmt.Println("work is done")

		type User struct {
			Name string
			Age  int
		}

		ch := make(chan User)

		go func() {
			ch <- User{Name: "Ilia", Age: 25}
		}()

		user := <-ch
		fmt.Println(user.Name, user.Age)

		// fan in обьединение каналов
		ch1 = make(chan int)
		ch2 = make(chan int)
		result := make(chan int)

		go func() {
			ch1 <- 1
		}()

		go func() {
			ch2 <- 2
		}()

		// fan-in
		go func() {
			result <- <-ch1
		}()

		go func() {
			result <- <-ch2
		}()

		fmt.Println(<-result)
		fmt.Println(<-result)

	}
	// 7. Создание горутины, возврат данных из горутины в главный поток.
	// Использование mutex, wait group.
	{
		fmt.Println("\n7-th example")
		ch := make(chan int)

		go func() {
			result := 42
			ch <- result
		}()

		res := <-ch
		fmt.Println(res)

		counter := 0
		var wg sync.WaitGroup
		var mu sync.Mutex

		for i := 0; i < 1000; i++ {
			wg.Add(1)

			go func() {

				defer wg.Done() // wg гарантирует что каждая горутина закончится

				mu.Lock() // лок гарантирует выполненение всех добавлений раздельно
				counter++
				mu.Unlock()
			}()
		}

		// чтобы горутины успели выполниться
		wg.Wait()

		fmt.Println(counter)
	}

}
