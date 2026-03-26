package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"sync"
	"time"
)

func (d Dog) Speak() {
	fmt.Println("woof")
}

type Animal interface {
	Speak()
}

type Dog struct{}

func res(n *int) (int, int) {
	*n++
	return *n, 1
}

func res2(n int) (int, int) {
	n++
	return n, 1
}

func check() error {
	return errors.New("errror from func_check")
}

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

	// 8. Передача аргументов в функции по указателю, по значению.
	// Возврат нескольких аргументов из функции.

	{
		fmt.Println("\n8-th example")
		n := 5
		fmt.Println(res(&n))
		fmt.Println(n)

		fmt.Println(res2(n))
		fmt.Println(n)
	}

	// 9. Создание ошибки, возврат ее из функции. Вложенные ошибки.
	{
		fmt.Println("\n9-th example")
		err := fmt.Errorf("something went wrong: %s", "smth")
		fmt.Println(err)
		err = errors.New("problem!")
		fmt.Println(err)
		err = check()
		fmt.Println(err)
	}

	// 10. Анонимные функции - создание, вызов с аргументами, захват значений.
	{
		fmt.Println("\n10-th example")
		res := func(a, b int) int {
			return a + b
		}(5, 6)
		fmt.Println(res)

		res2 := func(a, b int) int {
			return a + b
		}
		fmt.Println(res2(1, 3))

		x := 10

		f := func() int { // closure способность анонимной функции использовать
			//внешие переменные
			x += 5
			return x
		}
		fmt.Println(f())
		fmt.Println(x)
	}

	// 11. Использование defer для закрытия каналов, освобождения ресурсов.
	{
		fmt.Println("\n11-th example")
		check := make(chan int, 5)
		defer close(check)
		check <- 1
		fmt.Println(<-check)
	}
	// 12. Форматирование строки.
	// Вывод инта, флоата, була, кастомной структуры в форматированной строке.
	{
		fmt.Println("\n12-th example")

		type User struct {
			Name string
			Age  int
		}
		n := 42
		f := 3.14159
		b := true
		s := "hello"
		u := User{Name: "Ilya", Age: 25}

		fmt.Printf("its int: %d\n", n)             // int
		fmt.Printf("its float: %f\n", f)           // float
		fmt.Printf("its 2 digts float: %.2f\n", f) // float с точностью
		fmt.Printf("its bool: %t\n", b)            // bool
		fmt.Printf("its string: %s\n", s)          // string
		fmt.Printf("its val?: %v\n", u)            // значение
		fmt.Printf("its struct: %+v\n", u)         // struct с полями
		fmt.Printf("its go look: %#v\n", u)        // Go-вид
	}

	// 13. Переменная с функцией.
	// Вызов такой переменной с аргументами.
	{
		fmt.Println("\n13-th example")
		sum := func(a, b int) int {
			return a + b
		}
		fmt.Println(sum(3, 1))

	}

	// 14. Передача указателя в функцию.
	// Модификация значения в указателе внутри функции.

	{
		fmt.Println("\n14-th example")
		n := 5
		fmt.Println(res(&n))
		fmt.Println(n)
	}

	// 15. Маршал и анмаршал структуры в json.

	{
		fmt.Println("\n15-th example")

		type User struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}

		u := User{
			Name: "Ilia",
			Age:  23,
		}

		// Marshal
		data, err := json.Marshal(u)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(data)) // {"name":"Ilia","age":23}

		// Unmarshal
		var u2 User
		err = json.Unmarshal(data, &u2)
		if err != nil {
			panic(err)
		}
		fmt.Println(u2.Name, u2.Age) // Ilia 23

	}

	//	16. Итерация по строке (через слайс и через ренж).
	// Замена символа в середине строки
	// (в том числе для строк с многобайтовыми символами).

	{
		fmt.Println("\n16-th example")
		s := "héllo"

		for z := range s {
			fmt.Println(string(s[z])) // многобайтовые символы ломаются
		}
		fmt.Println("\n")

		for i, r := range s {
			fmt.Println(i, r) // i — байтовый индекс, r — rune
		}
		fmt.Println("\n")

		b := []byte(s)

		for i := 0; i < len(b); i++ {
			fmt.Println(i, string(b[i])) // по байтам
		} // многобайтовые символы ломаются

		fmt.Println("\n")
		runes := []rune(s)

		for i := 0; i < len(runes); i++ {
			fmt.Println(i, string(runes[i]))
		}

		s = "héllo"
		// s[1] = 'a' // ошибка компиляции нельзя так

		new := []rune(s) // []byte в случае стандартных символов, в других []rune

		new[1] = 'a'

		s = string(new)
		fmt.Println(s) // hallo
	}

	// 17. Конвертация строку в слайс рун и обратно.

	{
		fmt.Println("\n17-th example")

		s := "héllo"

		r := []rune(s)

		fmt.Println(r) // [104 233 108 108 111]

		r = []rune{104, 233, 108, 108, 111}

		s = string(r)

		fmt.Println(s) // héllo
	}

	// 18. Доставание куска строки из середины по индексам
	// (в том числе многобайтовых символов).

	{
		fmt.Println("\n18-th example")
		s := "héllo"

		// может сломать символ WRONG!
		sub := s[1:4]

		fmt.Println(sub)

		r := []rune(s)

		// берем с 1 по 4 символ
		sub = string(r[1:4])

		fmt.Println(sub) // éll

	}

	// 19. Сортировка слайса встроенной функцией.
	{
		fmt.Println("\n19-th example")
		arr := []int{5, 2, 8, 1}

		sort.Ints(arr)

		fmt.Println(arr) // [1 2 5 8]

		arr1 := []string{"banana", "apple", "cherry"}

		sort.Strings(arr1)

		fmt.Println(arr1) // [apple banana cherry]

		arr3 := []int{5, 2, 8, 1}

		sort.Slice(arr3, func(i, j int) bool {
			return arr3[i] > arr3[j]
		})

		fmt.Println(arr3)

	}

	// 20. Остановка исполнения через time.Sleep()
	{
		fmt.Println("\n20-th example")

		fmt.Println("start")

		// time.Sleep(1 * time.Second)
		time.Sleep(1 * time.Millisecond) // 0.5 сек
		// time.Sleep(1 * time.Minute)

		fmt.Println("end")
	}

	// 21. Использование контекста с таймаутом и отменой.
	// Мочь писать код который таймаутится и заканчивает исполнение
	// через н секунд (через контекст с таймаутом и select)
	{
		fmt.Println("\n21-th example")
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		ch := make(chan string)

		go func() {
			time.Sleep(2 * time.Second)
			ch <- "done"
		}()

		select {
		case res := <-ch:
			fmt.Println(res)
		case <-ctx.Done():
			fmt.Println("timeout:", ctx.Err())
		}
	}

	// 22. Использование for (бесконечный цикл, проверка условия,
	// проверка счетчика до значения вперед или назад)

	{
		fmt.Println("\n22-th example")
		i := 0

		// for {
		// 	//do smth or do break
		// }

		for i < 5 {
			fmt.Println(i)
			i++
		}

		fmt.Println("\n")

		for i := 0; i < 5; i++ {
			fmt.Println(i)
		}
		fmt.Println("\n")

		for i := 5; i > 0; i-- {
			fmt.Println(i)
		}
	}

	// 23. Каст обьектов друг к другу. Float -> Int и обратно.
	// Каст обьектов одного интерфейса к конкретным типам.
	// Разные виды кастов.
	{
		fmt.Println("\n23-th example")

		var f float64 = 3.7
		i := int(f)

		fmt.Println(i) // 3 (дробная часть отбрасывается)

		var in int = 3
		fl := float64(in)

		fmt.Println(fl) // 3.0

		var x interface{} = "hello"

		s := x.(string)

		fmt.Println(s) // hello
	}

	{
		var x interface{} = 123

		s, ok := x.(string)

		fmt.Println(s, ok) // "" false
	}

	{
		var x interface{} = 10.42

		switch v := x.(type) {
		case int:
			fmt.Println("int:", v)
		case string:
			fmt.Println("string:", v)
		default:
			fmt.Println("unknown")
		}
	}

	{
		fmt.Println("\n")
		// type Animal interface {
		// 	Speak()
		// }

		// type Dog struct{}

		// func (d Dog) Speak() {
		// 	fmt.Println("woof")
		// }

		var a Animal = Dog{}

		d := a.(Dog)

		d.Speak()
	}

	//24. Конвертация строки в int или float.
	{
		fmt.Println("\n24-th example")
		s := "123"

		i, err := strconv.Atoi(s)
		if err != nil {
			panic(err)
		}

		fmt.Println(i) // 123

		f, err := strconv.ParseFloat("3.14", 32)
		if err != nil {
			panic(err)
		}

		fmt.Println(f) // 3.14
	}

}
