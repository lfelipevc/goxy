package main

import (
	"fmt"
)

func gen(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

func sq(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		fmt.Println(<-in)
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out

}

// func readBody(body ioutil) {

// }

func main2() {

	//resp, _ := http.Get("http://geoservicos.pbh.gov.br/geoserver/wfs?service=wfs&request=GetCapabilities")
	// go readBody(resp.Body)
	//body, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println(string(body))
	// Set up the pipeline.
	//c := gen(2, 3)
	//out := sq(c)

	// Consume the output.
	//fmt.Println(<-out) // 4
	//fmt.Println(<-out) // 9
}
