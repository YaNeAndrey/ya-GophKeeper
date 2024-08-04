package main

import "fmt"

func main() {
	//myClient := client.NewClient(nil, storage.StorageRepo(storage.NewBaseStorage("temp")))
	//myClient.Start()

	reqURL := "efesfsefsf/222/111/1"
	fmt.Println(reqURL)
	buf := []rune(reqURL)
	buf[len(buf)-1] = '2'
	reqURL = string(buf)
	//([]rune(reqURL))[len(reqURL)-1] = '2'
	fmt.Println(reqURL)
}
