package client

import (
	"bufio"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
	"time"
	"ya-GophKeeper/internal/content"
)

func RootPage(c *Client) bool {
	fmt.Println("Main page: ")
	fmt.Println("1) Login/Re-login")
	fmt.Println("2) Sync data")
	fmt.Println("3) Add new information")
	fmt.Println("4) Update information")
	fmt.Println("5) Remove some information")
	fmt.Println("6) Print all data")
	fmt.Println("7) Get OTP for new device")
	fmt.Println("8) Help and Information")
	fmt.Println("9) Exit")
	choice := ReadOneLine()
	switch choice {
	case "1":
		RunConsoleFunc(c, LoginPage)
	case "2":
		RunConsoleFunc(c, Synchronization)
	case "3":
		RunConsoleFunc(c, AddPage)
	case "4":
		RunConsoleFunc(c, UpdatePage)
	case "5":
		RunConsoleFunc(c, RemovePage)
	case "6":
		PrintCreditCards(c)
		PrintTexts(c)
		PrintFiles(c)
		PrintCredentials(c)
	case "7":
	case "8":
		PrintHelpAndInformation(c)
	case "9":
		return false
	default:
		fmt.Println("Incorrect input")
	}
	return true
}

func LoginPage(c *Client) bool {
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

func UpdatePage(c *Client) bool {
	fmt.Println("Update page: ")
	fmt.Println("1) Update credential")
	fmt.Println("2) Update credit card")
	fmt.Println("3) Update file")
	fmt.Println("4) Update text")
	fmt.Println("5) Return to previous page")
	answer := ReadOneLine()
	switch answer {
	case "1":
		PrintCredentials(c)
		RunConsoleFunc(c, UpdateCredential)
	case "2":
		PrintCreditCards(c)
		RunConsoleFunc(c, UpdateCreditCards)
	case "3":
		PrintFiles(c)
		RunConsoleFunc(c, UpdateFiles)
	case "4":
		PrintTexts(c)
		RunConsoleFunc(c, UpdateText)
	case "5":
		return false
	}
	return true
}
func AddPage(c *Client) bool {
	fmt.Println("Add page: ")
	fmt.Println("1) Add new credential")
	fmt.Println("2) Add new credit card")
	fmt.Println("3) Add new file")
	fmt.Println("4) Add new text")
	fmt.Println("5) Return to previous page")
	answer := ReadOneLine()
	switch answer {
	case "1":
		RunConsoleFunc(c, AddCredential)
	case "2":
		RunConsoleFunc(c, AddCreditCard)
	case "3":
		RunConsoleFunc(c, AddFile)
	case "4":
		RunConsoleFunc(c, AddText)
	case "5":
		return false
	}
	return true
}
func RemovePage(c *Client) bool {
	fmt.Println("Remove page: ")
	fmt.Println("1) Remove credential")
	fmt.Println("2) Remove credit card")
	fmt.Println("3) Remove file")
	fmt.Println("4) Remove text")
	fmt.Println("5) Return to previous page")
	answer := ReadOneLine()
	switch answer {
	case "1":
		PrintCredentials(c)
	case "2":
		PrintCreditCards(c)
	case "3":
		PrintFiles(c)
	case "4":
		PrintTexts(c)
	case "5":
		return false
	default:
		return true
	}
	sn, err := ReadSerialNumber()
	if err != nil {
		log.Error(err)
		return true
	}
	if sn == -1 {
		return true
	}
	switch answer {
	case "1":
		//Work with storage
	case "2":
		//Work with storage
	case "3":
		//Work with storage
	case "4":
		//Work with storage
	}
	return true
}

func Synchronization(c *Client) bool {
	return false
}

func PrintCreditCards(c *Client) bool {
	return false
}
func PrintFiles(c *Client) bool {
	return false
}
func PrintTexts(c *Client) bool {
	return false
}
func PrintCredentials(c *Client) bool {
	return false
}

func UpdateCredential(c *Client) bool {
	fmt.Println("Update credential: ")
	credInfo, err := ReadCredential()
	if err != nil {
		log.Error(err)
		return true
	}
	if credInfo == nil {
		return false
	}
	return false
}
func UpdateCreditCards(c *Client) bool {
	fmt.Println("Update credential: ")
	creditCardInfo, err := ReadCreditCard()
	if err != nil {
		log.Error(err)
		return true
	}
	if creditCardInfo == nil {
		return false
	}
	return false
}
func UpdateFiles(c *Client) bool {
	fmt.Println("Update credential: ")
	fileInfo, err := ReadBinaryFile()
	if err != nil {
		log.Error(err)
		return true
	}
	if fileInfo == nil {
		return false
	}
	return false
}
func UpdateText(c *Client) bool {
	fmt.Println("Update credential: ")
	textInfo, err := ReadText()
	if err != nil {
		log.Error(err)
		return true
	}
	if textInfo == nil {
		return false
	}
	return false
}

func AddCredential(c *Client) bool {
	fmt.Println("Add new credential data: ")
	credInfo, err := ReadCredential()
	if err != nil {
		log.Error(err)
		return true
	}
	if credInfo == nil {
		return false
	}
	return false
}
func AddCreditCard(c *Client) bool {
	fmt.Println("Add new credit card data: ")
	creditCardInfo, err := ReadCreditCard()
	if err != nil {
		log.Error(err)
		return true
	}
	if creditCardInfo == nil {
		return false
	}
	return false
}
func AddFile(c *Client) bool {
	fmt.Println("Add new file data: ")
	fileInfo, err := ReadBinaryFile()
	if err != nil {
		log.Error(err)
		return true
	}
	if fileInfo == nil {
		return false
	}
	return false
}
func AddText(c *Client) bool {
	fmt.Println("Add new text data: ")
	textInfo, err := ReadText()
	if err != nil {
		log.Error(err)
		return true
	}
	if textInfo == nil {
		return false
	}
	return false
}

/*
	func RemoveCredential() bool {
		fmt.Println("Remove credential data: ")
		credInfo, err := ReadCredential()
		if err != nil {
			log.Error(err)
			return true
		}
		if credInfo == nil {
			return false
		}
		return false
	}

	func RemoveCreditCard() bool {
		fmt.Println("Remove credit card data: ")
		creditCardInfo, err := ReadCreditCard()
		if err != nil {
			log.Error(err)
			return true
		}
		if creditCardInfo == nil {
			return false
		}
		return false
	}

	func RemoveFile() bool {
		fmt.Println("Remove file data: ")
		fileInfo, err := ReadBinaryFile()
		if err != nil {
			log.Error(err)
			return true
		}
		if fileInfo == nil {
			return false
		}
		return false
	}

	func RemoveText() bool {
		fmt.Println("Remove text data: ")
		textInfo, err := ReadText()
		if err != nil {
			log.Error(err)
			return true
		}
		if textInfo == nil {
			return false
		}
		return false
	}
*/
func ReadSerialNumber() (int, error) {
	fmt.Println("Serial number (double Enter for return to previous page):")
	serialNumberStr := ReadOneLine()
	_ = serialNumberStr
	if serialNumberStr == "" {
		return -1, nil
	}

	sn, err := strconv.Atoi(serialNumberStr)
	if err != nil {
		return 0, err
	}
	return sn, nil
}

/*TODO: Add checks for all fields*/
func ReadCredential() (*content.CredentialInfo, error) {
	sn, err := ReadSerialNumber()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if sn == -1 {
		return nil, nil
	}

	fmt.Println("Resource: ")
	resource := ReadOneLine()

	fmt.Println("Login: ")
	login := ReadOneLine()

	fmt.Println("Password: ")
	password := ReadOneLine()

	return &content.CredentialInfo{
		ID:               sn,
		Resource:         resource,
		Login:            login,
		Password:         password,
		ModificationTime: time.Time{},
	}, nil
}
func ReadCreditCard() (*content.CreditCardInfo, error) {
	sn, err := ReadSerialNumber()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if sn == -1 {
		return nil, nil
	}
	fmt.Println("Bank: ")
	bank := ReadOneLine()
	_ = bank
	fmt.Println("Card Number (16 digits with spaces. format: 1111 2222 3333 4444): ")
	cardNumber := ReadOneLine()
	_ = cardNumber
	fmt.Println("CVV (3 digits): ")
	cvv := ReadOneLine()
	_ = cvv
	fmt.Println("ValidThru date (format dd.mm.yyyy/30.12.2001): ")
	validThru := ReadOneLine()
	_ = validThru
	return &content.CreditCardInfo{}, nil
}
func ReadBinaryFile() (*content.BinaryFileInfo, error) {
	sn, err := ReadSerialNumber()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if sn == -1 {
		return nil, nil
	}
	fmt.Println("Base file name (example: file.txt): ")
	fileName := ReadOneLine()
	_ = fileName

	fmt.Println("Description: ")
	description := ReadOneLine()
	_ = description

	fmt.Println("File path for reading: ")
	filePath := ReadOneLine()
	_ = filePath
	//read file
	return &content.BinaryFileInfo{}, nil
}
func ReadText() (*content.TextInfo, error) {
	sn, err := ReadSerialNumber()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if sn == -1 {
		return nil, nil
	}
	fmt.Println("Description: ")
	description := ReadOneLine()
	_ = description

	fmt.Println("Text (use double Enter to stop): ")
	textArray := ReadMultipleLines()
	text := strings.Join(textArray, "\r\n")
	_ = text
	return &content.TextInfo{}, nil
}

func GetOTP(c *Client) string {
	return ""
}

func PrintHelpAndInformation(c *Client) bool {
	return false
}

func RunConsoleFunc(c *Client, f func(*Client) bool) {
	for {
		fmt.Println("*******************************************")
		if !f(c) {
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
