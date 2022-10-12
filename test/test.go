package main

import "fmt"

func main()  {
	a := 1
	t:
		for {
			if a > 2 {
				fmt.Println(a)
				continue t
			}
			a++
		}
}

