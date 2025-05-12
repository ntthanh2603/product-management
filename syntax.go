package main

import "fmt"

/*
Các kiểu dữ liệu: int, float, string, bool
Khi khai báo
var <tên biến> <kiểu dữ liệu>
var <tên biến> = <giá trị>
<tên biến> := <giá trị>
var <tên biến>, <tên biến> = <giá trị>, <giá trị>
<tên biến>, <tên biến> := <giá trị>, <giá trị>
Khi dùng ':=' thì kiểu dữ liệu là kiểu dữ liệu của giá trị trên phái
*/

func sum(a, b int) int {
	return a + b
}

func main() {
	// Khai báo biến không có giá trị ban đầu
	var a int
	var b float32
	var c string
	var d bool

	fmt.Println(a) // 0
	fmt.Println(b) // 0
	fmt.Println(c) // ""
	fmt.Println(d) // false

	// Sự khác nhau giữa var và :=
	// var có thể dùng tron và ngoài các hàm, khai báo và gán có thể riêng biệt
	// := chỉ dùng trong hàm, khai báo và gán nằm chung 1 dòng

	// Ví dụ khai báo nhiều biến trên 1 dòng
	// var a, b, c, d int = 1, 3, 5, 7

	// Cũng có thể khai báo nhiều biến trong 1 khối
	// var (
	// 	a int
	// 	b int = 1
	// 	c string = "hello"
	//   )

	/*
		Quy tắc đặt tên biến:

		- Tên biến phải bắt đầu bằng một chữ cái hoặc ký tự gạch dưới (_)
		- Tên biến không thể bắt đầu bằng một chữ số
		- Tên biến chỉ có thể chứa các ký tự chữ và số và dấu gạch dưới ( a-z, A-Z, 0-9, và _)
		- Tên biến phân biệt chữ hoa chữ thường (age, Age và AGE là ba biến khác nhau)
		- Tên biến không có giới hạn về độ dài của tên biến
		- Tên biến không được chứa khoảng trắng
		- Tên biến không thể là bất kỳ từ khóa Go nào
	*/

	// const không được thay đổi và chỉ có thể đọc: hằng số có thể không cần kiểu dữ liệu

	/*
		In ra màn hình:

		- Print(): In các đối số theo dạng mạc định.
		- Println(): In các đối số có khoảng trắng ở giữa và thêm một dòng mới vào cuối.
		- Printf(): Sẽ định dạng đối số của mình dựa trên động từ định dạng đã cho.
	*/
	arr1 := [5]int{4, 5, 6, 7}
	fmt.Println("arr1:", arr1) // [4 5 6 7 0]

	arr2 := [...]int{4, 5, 6, 7, 8}
	fmt.Println("arr2:", arr2)       // [4 5 6 7 8]
	fmt.Println("arr2[0]:", arr2[1]) // 4

	arr3 := [5]int{1: 10, 2: 40}
	fmt.Println("arr3:", arr3)           // [0 10 40 0 0]
	fmt.Println("len(arr3):", len(arr3)) // 5

	// Slice: mảng động, không cần khai báo độ dài
	slice1 := []int{1, 2, 3, 4, 5}
	fmt.Println("slice1:", slice1)           // [1 2 3 4 5]
	fmt.Println("len(slice1):", len(slice1)) // 5
	fmt.Println("cap(slice1):", cap(slice1)) // 5

	// Lấy slice từ mẳng
	slice2 := arr1[1:3]
	fmt.Println("slice2:", slice2)           // [5 6]
	fmt.Println("len(slice2):", len(slice2)) // 2
	fmt.Println("cap(slice2):", cap(slice2)) // 4

	// Tạo slice bằng make
	slice3 := make([]int, 5, 10)
	fmt.Println("slice3:", slice3)           // [0 0 0 0 0]
	fmt.Println("len(slice3):", len(slice3)) // 5
	fmt.Println("cap(slice3):", cap(slice3)) // 10
	// Thay đổi thành phần của slice
	slice3[0] = 10
	fmt.Println("slice3:", slice3) // [10 0 0 0 0]
	// Hàm append
	slice3 = append(slice3, 20)
	fmt.Println("slice3:", slice3) // [10 0 0 0 0 20]
	// Hàm copy
	slice4 := make([]int, 5)
	copy(slice4, slice3)
	fmt.Println("slice4:", slice4) // [10 0 0 0 0]

	// Các toán tử trong Go: +, -, *, /, %, ++, --
	// Các toán tử so sánh: ==, !=, >, <, >=, <=
	// Các toán tử logic: &&, ||, !
	// Các toán tủ gán: =, +=, -=, *=, /=, %=, <<=, >>=, &=, |=, ^=, &^=

	// Câu lệnh điều kiện if else
	time := 20
	if time < 18 {
		fmt.Println("Good day.")
	} else {
		fmt.Println("Good evening.")
	}

	// Câu lệnh switch case
	switch time {
	case 18, 19, 20:
		fmt.Println("Good day.")
	default:
		fmt.Println("Good evening.")
	}

	// Vòng lặp for
	fmt.Println("Vòng lặp for:")
	for i := 0; i < 5; i++ {
		fmt.Println("i:", i) // 0 1 2 3 4
	}

	// continue dùng để bỏ qua lần lặp hiện tại và tiếp tục vòng lặp
	// break dùng để kết thúc vòng lặp

	// range
	fmt.Println("Vòng lặp range:")
	fruits := [3]string{"apple", "orange", "banana"}
	for index, val := range fruits {
		fmt.Println("index:", index, "val:", val) // 0 apple, 1 orange, 2 banana
	}

	// Hàm trong Go
	// func sum(a, b int) int {
	// 	return a + b
	// }
	fmt.Println("sum(1, 2):", sum(1, 2)) // 3
	// Giá trị trả về của hàm:
	// func myFunction(x int, y int) int {
	// 	return x + y
	// }

	// Giá trị trả về của hàm cũng có thể được đặt tên hoặc có nhiều giá trị:
	// func myFunction(x int, y int) (result int) {
	// 	result = x + y
	// 	return
	// } // hàm này trả về giá trị của biến result

	// Đệ quy
	// func testcount(x int) int {
	// 	if x == 11 {
	// 	  return 0
	// 	}
	// 	fmt.Println(x)
	// 	return testcount(x + 1)
	// }

	//  Struct trong Go
	type person struct {
		name string
		age  int
	}
	p := person{
		name: "John",
		age:  30,
	}
	fmt.Println("p:", p) // {John 30}

	// Map trong Go
	map_test := map[string]int{"one": 1, "two": 2, "three": 3, "four": 4}

	fmt.Println("Map:")
	for k, v := range map_test {
		fmt.Printf("%v : %v,\n", k, v)
	}
}
