package console

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"ya-GophKeeper/internal/client/config"
)

type Client struct {
	config *config.Config
	//storage *storage
}

func (c *Client) Start() {
	fmt.Println("*******************************************")
	fmt.Println("Client settings:")
	fmt.Println(c.config)
	//fmt.Println("*******************************************")
	RunConsoleFunc(RootPage)
}

func RootPage() bool {
	fmt.Println("Main page: ")
	fmt.Println("1) Login/Re-login")
	fmt.Println("2) Sync data")
	fmt.Println("3) Add new information")
	fmt.Println("4) Update information")
	fmt.Println("5) Remove some information")
	fmt.Println("6) Print all data")
	fmt.Println("7) Help")
	fmt.Println("8) Exit")
	choice := ReadOneLine()
	switch choice {
	case "1":
		RunConsoleFunc(LoginPage)
	case "2":
		RunConsoleFunc(Synchronization)
	case "3":
	case "4":
		RunConsoleFunc(UpdatePage)
	case "5":
		RunConsoleFunc(RemovePage)
	case "6":
		PrintCreditCards()
		PrintFilesAndText()
		PrintCredentials()
	case "7":
	case "8":
		return false
	default:
		fmt.Println("Incorrect input")
	}
	return true
}

func LoginPage() bool {
	fmt.Println("Login page: (all old data will be deleted)")
	fmt.Println("1) With password")
	fmt.Println("2) With OTP (an authorized device is required)")
	fmt.Println("3) Return to previous page")
	answer := ReadOneLine()
	switch answer {
	case "1":

	case "2":

	case "3":
		return false
	}
	return true
}

func UpdatePage() bool {
	fmt.Println("Update page: ")
	fmt.Println("1) Update credential")
	fmt.Println("2) Update file")
	fmt.Println("3) Update credit card")
	fmt.Println("4) Update text (store as file)")
	fmt.Println("5) Return to previous page")
	answer := ReadOneLine()
	switch answer {
	case "1":
		PrintCredentials()
		RunConsoleFunc(UpdateCredential)
	case "2":

	case "3":

	case "4":

	case "5":
		return false
	}
	return true
}

func RemovePage() bool {
	fmt.Println("Remove page: ")
	fmt.Println("1) Remove credential")
	fmt.Println("2) Remove file")
	fmt.Println("3) Remove credit card")
	fmt.Println("4) Remove text (store as file)")
	fmt.Println("5) Return to previous page")
	answer := ReadOneLine()
	switch answer {
	case "1":

	case "2":

	case "3":

	case "4":

	case "5":
		return false
	}
	return true
}

func Synchronization() bool {
	return false
}

func PrintCreditCards() bool {
	return false
}
func PrintFilesAndText() bool {
	return false
}
func PrintCredentials() bool {
	return false
}

func UpdateCredential() bool {
	fmt.Println("Update credential: ")
	fmt.Println("Select credential for update (enter zero for exit or credential number):")
	credNumber := ReadOneLine()
	_ = credNumber
	if credNumber == "0" {
		return false
	}
	fmt.Println("Resource: ")
	resource := ReadOneLine()
	_ = resource
	fmt.Println("Login: ")
	login := ReadOneLine()
	_ = login
	fmt.Println("Password: ")
	password := ReadOneLine()
	_ = password
	return false
}

func PrintHelp() bool {
	return false
}

func RunConsoleFunc(f func() bool) {
	for {
		fmt.Println("*******************************************")
		//ok := f()
		if !f() {
			return
		}
	}
}

func ReadMultipleLines() []string {
	scanner := bufio.NewScanner(os.Stdin)

	var lines []string
	for {
		scanner.Scan()
		line := scanner.Text()
		if len(line) == 0 {
			break
		}
		lines = append(lines, line)
	}

	err := scanner.Err()
	if err != nil {
		log.Println(err)
		return nil
	}
	return lines
}

func ReadOneLine() string {
	fmt.Print(">> ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	err := scanner.Err()
	if err != nil {
		log.Fatal(err)
		return ""
	}
	return scanner.Text()
}
