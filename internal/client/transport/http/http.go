package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"math"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"ya-GophKeeper/internal/client/storage"
	"ya-GophKeeper/internal/constants/clerror"
	"ya-GophKeeper/internal/constants/urlsuff"
	"ya-GophKeeper/internal/content"
)

type TransportHTTP struct {
	srvAddr   string
	chunkSize uint64
	jwtToken  string
}

type UserInfo struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func InitTransport(srvAddr string) *TransportHTTP {
	return &TransportHTTP{srvAddr: srvAddr}
}
func (tr *TransportHTTP) Registration(ctx context.Context, userAutData UserInfo) error {
	client := http.Client{}
	bodyJSON, err := json.Marshal(userAutData)
	if err != nil {
		return err
	}
	bodyReader := bytes.NewReader(bodyJSON)
	reqURL, _ := url.JoinPath(tr.srvAddr, urlsuff.OperationRegistration)
	req, _ := http.NewRequest(http.MethodPost, reqURL, bodyReader)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusConflict {
		return clerror.ErrLoginAlreadyTaken
	}
	if resp.StatusCode != http.StatusOK {
		return BadResponseHandler(resp, "From registration")
	}
	cookies := resp.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "token" {
			tr.jwtToken = cookie.Value
		}
	}
	if tr.jwtToken == "" {
		return clerror.ErrAuthTokenIsEmpty
	}
	return nil
}

func (tr *TransportHTTP) Login(ctx context.Context, userAutData UserInfo, loginType string) error {
	client := http.Client{}
	bodyJSON, err := json.Marshal(userAutData)
	if err != nil {
		return err
	}
	bodyReader := bytes.NewReader(bodyJSON)
	reqURL, _ := url.JoinPath(tr.srvAddr, urlsuff.OperationLogin, loginType)
	req, _ := http.NewRequest(http.MethodPost, reqURL, bodyReader)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized && resp.ContentLength == 0 {
		return clerror.ErrIncorrectPassword
	}
	if resp.StatusCode != http.StatusOK {
		return BadResponseHandler(resp, "From login operation")
	}
	cookies := resp.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "token" {
			tr.jwtToken = cookie.Value
		}
	}
	if tr.jwtToken == "" {
		return clerror.ErrAuthTokenIsEmpty
	}
	return nil
}

func (tr *TransportHTTP) ChangePassword(ctx context.Context, newPasswd string) error {
	client := http.Client{}
	reqURL, _ := url.JoinPath(tr.srvAddr, urlsuff.OperationChangPassword)
	strings.NewReader(newPasswd)
	req, _ := http.NewRequest(http.MethodPost, reqURL, strings.NewReader(newPasswd))
	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: tr.jwtToken,
	})
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return BadResponseHandler(resp, "From Change Password operation")
	}
	return nil
}

func (tr *TransportHTTP) GetOTP(ctx context.Context) (int, error) {
	client := http.Client{}
	reqURL, _ := url.JoinPath(tr.srvAddr, urlsuff.OperationGenerateOTP)
	req, _ := http.NewRequest(http.MethodGet, reqURL, nil)
	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: tr.jwtToken,
	})
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	otpBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	OTP, err := strconv.Atoi(string(otpBytes))
	if err != nil {
		return 0, err
	}
	return OTP, nil
}

func (tr *TransportHTTP) Sync(ctx context.Context, items storage.Collection) error {
	client := http.Client{}
	err := tr.SyncRemovedItems(ctx, items, &client)
	if err != nil {
		return err
	}
	err = tr.SyncNewItems(ctx, items, &client)
	if err != nil {
		return err
	}
	err = tr.SyncChangesFirstStep(ctx, items, &client)
	if err != nil {
		return err
	}

	return nil
}

func (tr *TransportHTTP) SyncRemovedItems(ctx context.Context, items storage.Collection, client *http.Client) error {
	rem := items.GetRemovedIDs()
	if rem == nil {
		return nil
	}
	bodyJSON, err := json.Marshal(rem)
	if err != nil {
		return err
	}
	bodyReader := bytes.NewReader(bodyJSON)
	reqURL, _ := url.JoinPath(tr.srvAddr, urlsuff.OperationRemove)
	var req *http.Request
	switch items.(type) {
	case *storage.Credentials:
		reqURL, _ = url.JoinPath(reqURL, urlsuff.DatatypeCredential)
	case *storage.CreditCards:
		reqURL, _ = url.JoinPath(reqURL, urlsuff.DatatypeCreditCard)
	case *storage.Texts:
		reqURL, _ = url.JoinPath(reqURL, urlsuff.DatatypeText)
	case *storage.Files:
		reqURL, _ = url.JoinPath(reqURL, urlsuff.DatatypeFile)
	default:
		return fmt.Errorf("SyncRemovedItems(TransportHTTP) %s", clerror.ErrIncorrectType)
	}
	req, _ = http.NewRequest(http.MethodPost, reqURL, bodyReader)
	req.Header.Add("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: tr.jwtToken,
	})
	//add another Headers

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return BadResponseHandler(resp, "From sync removed operation (Remove)")
	}
	items.ClearRemovedList()
	return nil
}

func (tr *TransportHTTP) SyncNewItems(ctx context.Context, items storage.Collection, client *http.Client) error {
	newItems := items.GetNewItems()
	if newItems == nil || (reflect.ValueOf(newItems).Kind() == reflect.Ptr && reflect.ValueOf(newItems).IsNil()) {
		return nil
	}
	/*
		if reflect.ValueOf(newItems).IsNil() {
			return nil
		}
	*/
	reqURL, _ := url.JoinPath(tr.srvAddr, urlsuff.OperationInsertNew)
	var srvAnswer interface{}
	datatype := ""
	switch items.(type) {
	case *storage.Credentials:
		datatype = urlsuff.DatatypeCredential
		//reqURL, _ = url.JoinPath(reqURL, urlsuff.DatatypeCredential)
	case *storage.CreditCards:
		//reqURL, _ = url.JoinPath(reqURL, urlsuff.DatatypeCreditCard)
		datatype = urlsuff.DatatypeCreditCard
	case *storage.Texts:
		//reqURL, _ = url.JoinPath(reqURL, urlsuff.DatatypeText)
		datatype = urlsuff.DatatypeText
	case *storage.Files:
		//reqURL, _ = url.JoinPath(reqURL, urlsuff.DatatypeFile)
		datatype = urlsuff.DatatypeFile
		return nil
	default:
		return fmt.Errorf("SyncNewItems(TransportHTTP) %s", clerror.ErrIncorrectType)
	}

	reqURL, _ = url.JoinPath(reqURL, datatype)
	bodyJSON, err := json.Marshal(newItems)
	if err != nil {
		return err
	}
	bodyReader := bytes.NewReader(bodyJSON)
	var req *http.Request
	req, _ = http.NewRequest(http.MethodPost, reqURL, bodyReader)
	req.Header.Add("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: tr.jwtToken,
	})

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return BadResponseHandler(resp, "From sync new items operation (Add)")
	}
	err = json.NewDecoder(resp.Body).Decode(&srvAnswer)
	if err != nil {
		return err
	}
	items.RemoveItemsWithoutID()
	err = items.AddOrUpdateItems(srvAnswer)
	if err != nil {
		return err
	}

	if datatype == urlsuff.DatatypeFile {
		var newItemsWithType []content.BinaryFileInfo
		jsonbody, err := json.Marshal(srvAnswer)
		if err != nil {
			return err
		}
		if err = json.Unmarshal(jsonbody, &newItemsWithType); err != nil {
			return err
		}

		err = tr.UploadFiles(ctx, newItemsWithType, client)
		if err != nil {
			log.Println(err)
			return nil
		}
	}
	return nil
}

func (tr *TransportHTTP) UploadFiles(ctx context.Context, files []content.BinaryFileInfo, client *http.Client) error {
	var returnErr error
	for _, file := range files {
		fileInfo := content.BinaryFileInfo{
			FileName: file.FileName,
			FilePath: file.FilePath,
		}
		go func() {
			err := tr.UploadFile(ctx, fileInfo.FilePath, fileInfo.FileName, client)
			if err != nil {
				returnErr = fmt.Errorf("%s%s\r\n", returnErr.Error(), err.Error())
			}
		}()
	}
	return returnErr
}

func (tr *TransportHTTP) UploadFile(ctx context.Context, filepath string, fileName string, client *http.Client) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}

	defer file.Close()

	fileInfo, _ := file.Stat()

	fileSize := fileInfo.Size()

	totalPartsNum := uint64(math.Ceil(float64(fileSize) / float64(tr.chunkSize)))

	for i := uint64(0); i < totalPartsNum; i++ {
		partSize := int(math.Min(float64(tr.chunkSize), float64(fileSize-int64(i*tr.chunkSize))))
		partBuffer := make([]byte, partSize)
		_, err = file.Read(partBuffer)
		if err != nil {
			return err
		}

		err = tr.SendChunk(ctx, partBuffer, fileName, i, totalPartsNum, fileSize, client)
		if err != nil {
			return err
		}
	}
	return nil
}

func (tr *TransportHTTP) SyncChangesFirstStep(ctx context.Context, items storage.Collection, client *http.Client) error {
	IDsWithModtime := items.GetAllIDsWithModtime()
	bodyJSON, err := json.Marshal(IDsWithModtime)
	if err != nil {
		return err
	}
	bodyReader := bytes.NewReader(bodyJSON)
	reqURL, _ := url.JoinPath(tr.srvAddr, urlsuff.OperationSync)
	var req *http.Request
	srvAnswer := struct {
		DataForSrv    []int       `json:",omitempty"`
		RemoveFromCli []int       `json:",omitempty"`
		DataForCli    interface{} `json:",omitempty"`
	}{}
	datatype := ""
	switch items.(type) {
	case *storage.Credentials:
		datatype = urlsuff.DatatypeCredential
		//reqURL, _ = url.JoinPath(reqURL, urlsuff.DatatypeCredential)
	case *storage.CreditCards:
		datatype = urlsuff.DatatypeCreditCard
		//reqURL, _ = url.JoinPath(reqURL, urlsuff.DatatypeCreditCard)
	case *storage.Texts:
		datatype = urlsuff.DatatypeText
		//reqURL, _ = url.JoinPath(reqURL, urlsuff.DatatypeText)
	case *storage.Files:
		datatype = urlsuff.DatatypeFile
		//reqURL, _ = url.JoinPath(reqURL, urlsuff.DatatypeFile)
	default:
		return fmt.Errorf("SyncChangesFirstStep(TransportHTTP) %s", clerror.ErrIncorrectType)
	}

	reqURL, _ = url.JoinPath(reqURL, datatype)
	reqURL, _ = url.JoinPath(reqURL, "1")
	req, _ = http.NewRequest(http.MethodPost, reqURL, bodyReader)
	req.Header.Add("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: tr.jwtToken,
	})

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return BadResponseHandler(resp, "From sync changes operation step one (Update)")
	}
	err = json.NewDecoder(resp.Body).Decode(&srvAnswer)
	if err != nil {
		return err
	}
	err = items.AddOrUpdateItems(srvAnswer.DataForCli)
	if err != nil {
		return err
	}
	if srvAnswer.RemoveFromCli != nil {
		items.RemoveItems(srvAnswer.RemoveFromCli)
	}

	if srvAnswer.DataForSrv == nil {
		return nil
	}
	err = tr.SyncChangesSecondStep(ctx, items, srvAnswer.DataForSrv, datatype, client)
	if err != nil {
		return err
	}
	return nil
}

func (tr *TransportHTTP) SyncChangesSecondStep(ctx context.Context, items storage.Collection, sendingItemsIDs []int, datatype string, client *http.Client) error {
	reqURL, _ := url.JoinPath(tr.srvAddr, urlsuff.OperationSync, datatype)

	itemsForServer := items.GetItems(sendingItemsIDs)
	bodyJSON, err := json.Marshal(itemsForServer)
	if err != nil {
		return err
	}
	reqURL, _ = url.JoinPath(reqURL, "2")
	bodyReader := bytes.NewReader(bodyJSON)
	req, _ := http.NewRequest(http.MethodPost, reqURL, bodyReader)
	req.Header.Add("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: tr.jwtToken,
	})
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return BadResponseHandler(resp, "From sync changes operation step two (Update)")
	}

	if datatype == urlsuff.DatatypeFile {
		var newItemsWithType []content.BinaryFileInfo
		jsonbody, err := json.Marshal(itemsForServer)
		if err != nil {
			// do error check
			return err
		}
		if err = json.Unmarshal(jsonbody, &newItemsWithType); err != nil {
			log.Println(err)
			return nil
		}

		err = tr.UploadFiles(ctx, newItemsWithType, client)
		if err != nil {
			log.Println(err)
			return nil
		}
	}

	return nil
}

func BadResponseHandler(r *http.Response, message string) error {
	var bodyBytes []byte
	if r.ContentLength != 0 {
		bodyBytes, _ = io.ReadAll(r.Body)
	}
	code := r.StatusCode
	switch code {
	case http.StatusUnauthorized:
		return clerror.ErrNotAuthorized
	case http.StatusBadRequest:
		if bodyBytes != nil {
			return fmt.Errorf("%s (%w : %s)", message, clerror.ErrBadRequest, bodyBytes)
		}
		return fmt.Errorf("%s (%w)", message, clerror.ErrBadRequest)
	case http.StatusInternalServerError:
		if bodyBytes != nil {
			return fmt.Errorf("%s (%w : %s)", message, clerror.ErrInternalServerError, bodyBytes)
		}
		return fmt.Errorf("%s (%w)", message, clerror.ErrInternalServerError)
	case http.StatusNotFound:
		if bodyBytes != nil {
			return fmt.Errorf("%s (%w : %s)", message, clerror.ErrStatusNotFound, bodyBytes)
		}
		return fmt.Errorf("%s (%w)", message, clerror.ErrStatusNotFound)
	default:
		if bodyBytes != nil {
			return fmt.Errorf("%s (status code: %d; server message : %s)", message, code, bodyBytes)
		}
		return fmt.Errorf("%s (status code: %d; server message : %s)", message, code)
	}
}

func (tr *TransportHTTP) SendChunk(ctx context.Context, chunk []byte, fileName string, chunkNumber uint64, chunkCount uint64, totalFileSize int64, client *http.Client) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	metadata := struct {
		ChunkNumber   uint64
		TotalChunks   uint64
		FileName      string
		TotalFileSize int64
	}{
		ChunkNumber:   chunkNumber,
		TotalChunks:   chunkCount,
		FileName:      fileName,
		TotalFileSize: totalFileSize,
	}
	bodyJSON, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	metadataHeader := textproto.MIMEHeader{}
	metadataHeader.Set("Content-Disposition", "form-data; name=\"metadata\"")
	metadataHeader.Set("Content-Type", "application/json")
	part, _ := writer.CreatePart(metadataHeader)
	_, err = part.Write(bodyJSON)
	if err != nil {
		return err
	}

	mediaHeader := textproto.MIMEHeader{}
	mediaHeader.Set("Content-Disposition", "form-data; name=\"chunk\"")
	mediaPart, _ := writer.CreatePart(mediaHeader)

	_, err = io.Copy(mediaPart, bytes.NewReader(chunk))
	if err != nil {
		return err
	}
	reqURL, _ := url.JoinPath(tr.srvAddr, urlsuff.DatatypeFile, urlsuff.FileOperationUpload)

	req, err := http.NewRequest(http.MethodPost, reqURL, body)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", writer.FormDataContentType())

	/*
		proxyURL, _ := url.Parse("http://localhost:8888")
		proxy := http.ProxyURL(proxyURL)
		transport := &http.Transport{Proxy: proxy}
		client := &http.Client{Transport: transport}
	*/
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}
