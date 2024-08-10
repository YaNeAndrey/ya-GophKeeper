package client

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
	"time"
	"ya-GophKeeper/internal/client/storage"
	"ya-GophKeeper/internal/client/transport"
	"ya-GophKeeper/internal/constants/clerror"
	"ya-GophKeeper/internal/constants/urlsuff"
	"ya-GophKeeper/internal/content"
)

func StartPage(c *Client) bool {
	c.storage.Clear()
	fmt.Println("1) Sign in")
	fmt.Println("2) Sign up")
	fmt.Println("3) Exit")
	choice := ReadOneLine()
	var authOk bool
	switch choice {
	case "1":
		authOk = LoginPage(c)
	case "2":
		authOk = RegistrationPage(c)
	case "3":
		return false
	//case "4":
	//	RunConsoleFunc(c, MainPage)
	default:
		fmt.Println("Incorrect input")
	}
	if authOk {
		RunConsoleFunc(c, MainPage)
	}
	return true
}

func MainPage(c *Client) bool {
	fmt.Println("Main page: ")
	fmt.Println("1) Add information page")
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
		GetOTP(c)
	case "7":
		return false
	case "8":
		ChangePassword(c)
	case "9":
		PrintHelpAndInformation(c)
	case "10":
		os.Exit(0)
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
		RunConsoleFunc(c, AuthorizationPassword)
	case "2":
		RunConsoleFunc(c, AuthorizationOTP)
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

	err := c.transport.Registration(context.Background(), transport.UserInfo{
		Login:    newUserLogin,
		Password: newUserPaswd,
	})
	if err != nil {
		log.Println(err)
		return true
	}
	return false
}
func ChangePassword(c *Client) bool {
	fmt.Println("New password: ")
	passwd := ReadOneLine()
	_ = passwd
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
	fmt.Println("Print page: ")
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

	var errorsList []error
	fmt.Println("Bank: ")
	bank := ReadOneLine()
	fmt.Println("Card Number (16 digits without spaces): ")
	cardNumber := ReadOneLine()
	if cardNumber != "" {
		err = CheckCardNumberWithLuhn(cardNumber)
		if err != nil {
			errorsList = append(errorsList, err)
		}
	}
	fmt.Println("CVV (3 digits): ")
	cvv := ReadOneLine()
	if cvv != "" {
		err = CheckCVV(cvv)
		if err != nil {
			errorsList = append(errorsList, err)
		}
	}
	fmt.Println("ValidThru date (format dd.mm.yyyy/30.12.2001): ")
	validThruStr := ReadOneLine()
	var validThru time.Time
	if validThruStr != "" {
		validThru, err = time.Parse("02.01.2006", validThruStr)
		if err != nil {
			errorsList = append(errorsList, err)
		}
	}
	if len(errorsList) != 0 {
		return 0, nil, errors.Join(errorsList...)
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

	fmt.Println("File path for reading. Max file size: 4GB. : ")
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

	fmt.Println("Text. Max size: 1000 symbols. (use double Enter to stop): ")
	textArray := ReadMultipleLines()
	text := strings.Join(textArray, "\r\n")
	if len(text) > 1000 {
		return 0, nil, clerror.ErrMaxTextSizeExceeded
	}
	return index, &content.TextInfo{
		Content:          text,
		Description:      description,
		ModificationTime: time.Now(),
	}, nil
}

func Synchronization(c *Client) bool {
	fmt.Println("Synchronization page: ")
	fmt.Println("1) Sync credential")
	fmt.Println("2) Sync credit card")
	fmt.Println("3) Sync texts")
	fmt.Println("4) Sync files")
	fmt.Println("5) Full sync")
	fmt.Println("6) Return to previous page")
	answer := ReadOneLine()

	var collection storage.Collection
	switch answer {
	case "1":
		collection = c.storage.GetCredentialsData()
	case "2":
		collection = c.storage.GetCreditCardsData()
	case "3":
		collection = c.storage.GetTextsData()
	case "4":
		collection = c.storage.GetFilesData()
	case "5":
		SyncCollection(c, c.storage.GetCredentialsData())
		SyncCollection(c, c.storage.GetCreditCardsData())
		SyncCollection(c, c.storage.GetTextsData())
		SyncCollection(c, c.storage.GetFilesData())

	case "6":
		return false
	default:
		fmt.Println("Incorrect input")
		return true
	}
	err := SyncCollection(c, collection)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func SyncCollection(c *Client, collection storage.Collection) error {
	err := c.transport.Sync(context.Background(), collection)
	if err != nil {
		return err
	}
	return nil
}

func PrintCreditCards(c *Client) bool {
	fmt.Println("*******************************************")
	fmt.Println("Credit cards:")
	cardsInfo := c.storage.GetCreditCardsData()
	cards := cardsInfo.GetItems(nil)
	switch t := cards.(type) {
	case []content.CreditCardInfo:
		for i, card := range t {
			fmt.Printf("%d) %s", i, card.String())
		}
	}
	return false
}

func PrintFiles(c *Client) bool {
	fmt.Println("*******************************************")
	fmt.Println("Files:")

	filesInfo := c.storage.GetFilesData()
	files := filesInfo.GetItems(nil)
	switch t := files.(type) {
	case []content.BinaryFileInfo:
		for i, file := range t {
			fmt.Printf("%d) %s", i, file.String())
		}
	}
	return false
}
func PrintTexts(c *Client) bool {
	fmt.Println("*******************************************")
	fmt.Println("Texts:")

	textsInfo := c.storage.GetTextsData()
	texts := textsInfo.GetItems(nil)
	switch t := texts.(type) {
	case []content.TextInfo:
		for i, text := range t {
			fmt.Printf("%d) %s", i, text.String())
		}
	}
	return false
}
func PrintCredentials(c *Client) bool {
	fmt.Println("*******************************************")
	fmt.Println("Credentials:")
	credsInfo := c.storage.GetCredentialsData()
	creds := credsInfo.GetItems(nil)
	switch t := creds.(type) {
	case []content.CredentialInfo:
		for i, cred := range t {
			fmt.Printf("%d) %s", i, cred.String())
		}
	}
	return false
}

func GetOTP(c *Client) {
	otp, err := c.transport.GetOTP(context.Background())
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("OTP: %d\r\n", otp)
}

func AuthorizationPassword(c *Client) bool {
	fmt.Println("Login: ")
	login := ReadOneLine()
	fmt.Println("Password: ")
	passwd := ReadOneLine()

	err := c.transport.Login(context.Background(), transport.UserInfo{
		Login:    login,
		Password: passwd,
	}, urlsuff.LoginTypePasswd)
	if err != nil {
		fmt.Printf("Authorization problem. Try again! Err: %s\r\n", err.Error())
		return true
	}

	return false
}

func AuthorizationOTP(c *Client) bool {
	fmt.Println("Login: ")
	login := ReadOneLine()
	fmt.Println("OTP: ")
	otp := ReadOneLine()
	err := c.transport.Login(context.Background(), transport.UserInfo{
		Login:    login,
		Password: otp,
	}, urlsuff.LoginTypeOTP)
	if err != nil {
		fmt.Printf("Authorization problem. Try again! Err: %s", err.Error())
		return true
	}
	return false
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
