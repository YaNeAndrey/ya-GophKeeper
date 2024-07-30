package client

import (
	"bufio"
	"container/list"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
	"time"
	"ya-GophKeeper/internal/client/clerror"
	"ya-GophKeeper/internal/content"
)

func StartPage(c *Client) bool {
	fmt.Println("1) Sign in")
	fmt.Println("2) Sign up")
	fmt.Println("3) Exit")
	choice := ReadOneLine()
	switch choice {
	case "1":
		RunConsoleFunc(c, LoginPage)
	case "2":
		RunConsoleFunc(c, RegistrationPage)
	case "3":
		return false
	case "4":
		RunConsoleFunc(c, MainPage)
	default:
		fmt.Println("Incorrect input")
	}
	return true
}

func MainPage(c *Client) bool {
	fmt.Println("Main page: ")
	fmt.Println("1) Add new information page")
	fmt.Println("2) Update information page")
	fmt.Println("3) Remove information page")
	fmt.Println("4) Print information page")
	fmt.Println("5) Sync data")
	fmt.Println("6) Get OTP for new device")
	fmt.Println("7) Re-login")
	fmt.Println("8) Change password")
	fmt.Println("9) Help")
	fmt.Println("10) Exit")
	choice := ReadOneLine()
	switch choice {
	case "1":
		RunConsoleFunc(c, AddPage)
	case "2":
		RunConsoleFunc(c, UpdatePage)
	case "3":
		RunConsoleFunc(c, RemovePage)
	case "4":
		RunConsoleFunc(c, PrintPage)
	case "5":
		RunConsoleFunc(c, Synchronization)
	case "6":
		//Get OTP
	case "7":
		RunConsoleFunc(c, LoginPage)
	case "8":
		//Change password
	case "9":
		PrintHelpAndInformation(c)
	case "10":
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
func RegistrationPage(c *Client) bool {
	fmt.Println("Login: ")
	newUserLogin := ReadOneLine()
	fmt.Println("Password: ")
	newUserPaswd := ReadOneLine()
	_ = newUserPaswd
	_ = newUserLogin
	return false
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
	default:
		fmt.Println("Incorrect input")
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
	default:
		fmt.Println("Incorrect input")
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
		fmt.Println("Incorrect input")
		return true
	}
	index, err := ReadDataIndex()
	if err != nil {
		log.Errorf("RemovePage : %w", err)
		return true
	}
	if index == -1 {
		return true
	}
	switch answer {
	case "1":
		err = c.storage.RemoveCredential(index)
		if err != nil {
			log.Errorf("RemovePage : %w", err)
		}
	case "2":
		err = c.storage.RemoveCreditCard(index)
		if err != nil {
			log.Errorf("RemovePage : %w", err)
		}
	case "3":
		err = c.storage.RemoveFile(index)
		if err != nil {
			log.Errorf("RemovePage : %w", err)
		}
	case "4":
		err = c.storage.RemoveText(index)
		if err != nil {
			log.Errorf("RemovePage : %w", err)
		}
	}
	return true
}
func PrintPage(c *Client) bool {
	fmt.Println("Remove page: ")
	fmt.Println("1) Print credential")
	fmt.Println("2) Print credit card")
	fmt.Println("3) Print file")
	fmt.Println("4) Print text")
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
		fmt.Println("Incorrect input")
	}
	return true
}

func UpdateCredential(c *Client) bool {
	fmt.Println("Update credential: ")
	index, credInfo, err := ReadCredential()
	if err != nil {
		log.Error(err)
		return true
	}
	if credInfo == nil {
		return false
	}
	err = c.storage.UpdateCredentials(index, credInfo)
	if err != nil {
		log.Error(err)
		return true
	}
	return false
}
func UpdateCreditCards(c *Client) bool {
	fmt.Println("Update credential: ")
	index, creditCardInfo, err := ReadCreditCard()
	if err != nil {
		log.Error(err)
		return true
	}
	if creditCardInfo == nil {
		return false
	}
	err = c.storage.UpdateCreditCards(index, creditCardInfo)
	if err != nil {
		log.Error(err)
		return true
	}
	return false
}
func UpdateFiles(c *Client) bool {
	fmt.Println("Update credential: ")
	index, fileInfo, err := ReadBinaryFile()
	if err != nil {
		log.Error(err)
		return true
	}
	if fileInfo == nil {
		return false
	}
	err = c.storage.UpdateFiles(index, fileInfo)
	if err != nil {
		log.Error(err)
		return true
	}
	return false
}
func UpdateText(c *Client) bool {
	fmt.Println("Update credential: ")
	index, textInfo, err := ReadText()
	if err != nil {
		log.Error(err)
		return true
	}
	if textInfo == nil {
		return false
	}
	err = c.storage.UpdateTexts(index, textInfo)
	if err != nil {
		log.Error(err)
		return true
	}
	return false
}

func AddCredential(c *Client) bool {
	fmt.Println("Add new credential data: ")
	_, credInfo, err := ReadCredential()
	if err != nil {
		log.Error(err)
		return true
	}
	if credInfo == nil {
		return false
	}
	err = c.storage.AddNewCredential(credInfo)
	if err != nil {
		log.Error(err)
		return true
	}
	return false
}
func AddCreditCard(c *Client) bool {
	fmt.Println("Add new credit card data: ")
	_, creditCardInfo, err := ReadCreditCard()
	if err != nil {
		log.Error(err)
		return true
	}
	if creditCardInfo == nil {
		return false
	}
	err = c.storage.AddNewCreditCard(creditCardInfo)
	if err != nil {
		log.Error(err)
		return true
	}
	return false
}
func AddFile(c *Client) bool {
	fmt.Println("Add new file data: ")
	_, fileInfo, err := ReadBinaryFile()
	if err != nil {
		log.Error(err)
		return true
	}
	if fileInfo == nil {
		return false
	}
	err = c.storage.AddNewFile(fileInfo)
	if err != nil {
		log.Error(err)
		return true
	}
	return false
}
func AddText(c *Client) bool {
	fmt.Println("Add new text data: ")
	_, textInfo, err := ReadText()
	if err != nil {
		log.Error(err)
		return true
	}
	if textInfo == nil {
		return false
	}
	err = c.storage.AddNewText(textInfo)
	if err != nil {
		log.Error(err)
		return true
	}
	return false
}

func ReadDataIndex() (int, error) {
	fmt.Println("Select index (leave empty for return to previous page. Use 0 for Add operation.):")
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

func ReadCredential() (int, *content.CredentialInfo, error) {
	index, err := ReadDataIndex()
	if err != nil {
		return 0, nil, err
	}
	if index == -1 {
		return 0, nil, nil
	}
	fmt.Println("Resource: ")
	resource := ReadOneLine()

	fmt.Println("Login: ")
	login := ReadOneLine()

	fmt.Println("Password: ")
	password := ReadOneLine()

	return index, &content.CredentialInfo{
		Resource:         resource,
		Login:            login,
		Password:         password,
		ModificationTime: time.Now(),
	}, nil
}
func ReadCreditCard() (int, *content.CreditCardInfo, error) {
	index, err := ReadDataIndex()
	if err != nil {
		return 0, nil, err
	}
	if index == -1 {
		return 0, nil, nil
	}

	errorsList := list.List{}
	fmt.Println("Bank: ")
	bank := ReadOneLine()
	fmt.Println("Card Number (16 digits without spaces): ")
	cardNumber := ReadOneLine()
	if cardNumber != "" {
		err = CheckCardNumberWithLuhn(cardNumber)
		if err != nil {
			errorsList.PushBack(err)
		}
	}
	fmt.Println("CVV (3 digits): ")
	cvv := ReadOneLine()
	if cvv != "" {
		err = CheckCVV(cvv)
		if err != nil {
			errorsList.PushBack(err)
		}
	}
	fmt.Println("ValidThru date (format dd.mm.yyyy/30.12.2001): ")
	validThruStr := ReadOneLine()
	var validThru time.Time
	if validThruStr != "" {
		validThru, err = time.Parse("01/02/2006", validThruStr)
		if err != nil {
			errorsList.PushBack(err)
		}
	}
	if errorsList.Len() != 0 {
		return 0, nil, fmt.Errorf("%v", errorsList)
	} else {
		return index, &content.CreditCardInfo{
			CardNumber:       cardNumber,
			CVV:              cvv,
			ValidThru:        validThru,
			Bank:             bank,
			ModificationTime: time.Now(),
		}, nil
	}
}
func ReadBinaryFile() (int, *content.BinaryFileInfo, error) {
	index, err := ReadDataIndex()
	if err != nil {
		log.Error(err)
		return 0, nil, err
	}
	if index == -1 {
		return 0, nil, nil
	}
	fmt.Println("Base file name (example: file.txt): ")
	fileName := ReadOneLine()

	fmt.Println("Description: ")
	description := ReadOneLine()

	fmt.Println("File path for reading: ")
	filePath := ReadOneLine()
	return index, &content.BinaryFileInfo{
		FileName:         fileName,
		FilePath:         filePath,
		Description:      description,
		ModificationTime: time.Now(),
	}, nil
}
func ReadText() (int, *content.TextInfo, error) {
	index, err := ReadDataIndex()
	if err != nil {
		log.Error(err)
		return 0, nil, err
	}
	if index == -1 {
		return 0, nil, nil
	}
	fmt.Println("Description: ")
	description := ReadOneLine()
	_ = description

	fmt.Println("Text (use double Enter to stop): ")
	textArray := ReadMultipleLines()
	text := strings.Join(textArray, "\r\n")
	_ = text
	return index, &content.TextInfo{}, nil
}

func Synchronization(c *Client) bool {
	return false
}

func PrintCreditCards(c *Client) bool {
	fmt.Println("*******************************************")
	fmt.Println("Credit cards:")
	for i, card := range c.storage.GetCreditCardData() {
		fmt.Printf("%d) %s", i, card.String())
	}
	return false
}
func PrintFiles(c *Client) bool {
	fmt.Println("*******************************************")
	fmt.Println("Files:")
	for i, file := range c.storage.GetFilesData() {
		fmt.Printf("%d) %s", i, file.String())
	}
	return false
}
func PrintTexts(c *Client) bool {
	fmt.Println("*******************************************")
	fmt.Println("Texts:")
	for i, text := range c.storage.GetTextData() {
		fmt.Printf("%d) %s", i, text.String())
	}
	return false
}
func PrintCredentials(c *Client) bool {
	fmt.Println("*******************************************")
	fmt.Println("Credentials:")
	for i, cred := range c.storage.GetCredentials() {
		fmt.Printf("%d) %s", i, cred.String())
	}
	return false
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

func CheckCardNumberWithLuhn(cardNumberStr string) error {
	total := 0
	isSecondDigit := false

	for i := len(cardNumberStr) - 1; i >= 0; i-- {
		digit, err := strconv.Atoi(string(cardNumberStr[i]))
		if err != nil {
			return clerror.ErrCreditCardIncorrectChar
		}
		if isSecondDigit {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		total += digit
		isSecondDigit = !isSecondDigit
	}
	if total%10 == 0 {
		return nil
	} else {
		return clerror.ErrCreditCardLuhn
	}
}

func CheckCVV(cvvStr string) error {
	cvv, err := strconv.Atoi(cvvStr)
	if err != nil {
		return err
	}
	if cvv < 1000 && cvv > 99 {
		return nil
	}
	return clerror.ErrIncorrectValueCVV
}
