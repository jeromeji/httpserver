package main

import (
	"bufio"
	"fmt"
	"net/http"
)

func main()  {
	res,err := http.Get("http://127.0.0.1:8080/healthz")
	defer res.Body.Close()
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
	fmt.Println(res.Status)
	scanner := bufio.NewScanner(res.Body)
	for i := 0; scanner.Scan() && i < 5; i++ {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

}