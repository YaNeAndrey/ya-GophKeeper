package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"ya-GophKeeper/internal/constants/clerror"
	"ya-GophKeeper/internal/constants/urlsuff"
	"ya-GophKeeper/internal/content"
)

type TransportHTTP struct {
	srvAddr   string
	chunkSize int64
	jwtToken  string
}

type UserInfo struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func InitTransport(srvAddr string, chunkSize int64) *TransportHTTP {
	return &TransportHTTP{srvAddr: srvAddr, chunkSize: chunkSize}
}
func (tr *TransportHTTP) Registration(ctx context.Context, userAuthData UserInfo) error {
	client := http.Client{}
	bodyJSON, err := json.Marshal(userAuthData)
	if err != nil {
		return err
	}
	bodyReader := bytes.NewReader(bodyJSON)
	reqURL, _ := url.JoinPath(tr.srvAddr, urlsuff.OperationRegistration)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bodyReader)
	//req, _ := http.NewRequest(http.MethodPost, reqURL, bodyReader)
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

func (tr *TransportHTTP) Login(ctx context.Context, userAuthData UserInfo, loginType string) error {
	client := http.Client{}
	bodyJSON, err := json.Marshal(userAuthData)
	if err != nil {
		return err
	}
	bodyReader := bytes.NewReader(bodyJSON)
	reqURL, _ := url.JoinPath(tr.srvAddr, urlsuff.OperationLogin, loginType)

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bodyReader)
	//req, _ := http.NewRequest(http.MethodPost, reqURL, bodyReader)
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
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, strings.NewReader(newPasswd))
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
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
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

func (tr *TransportHTTP) SyncRemovedItems(ctx context.Context, removedIDs []int, dataType string) error {
	client := http.Client{}
	rem := removedIDs
	if rem == nil {
		return nil
	}
	bodyJSON, err := json.Marshal(rem)
	if err != nil {
		return err
	}
	bodyReader := bytes.NewReader(bodyJSON)
	reqURL, _ := url.JoinPath(tr.srvAddr, urlsuff.OperationRemove, dataType)
	var req *http.Request

	req, _ = http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bodyReader)
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
		return BadResponseHandler(resp, "From sync removed operation (Remove)")
	}
	return nil
}

func (tr *TransportHTTP) SyncNewItems(ctx context.Context, bodyWithNewItems []byte, dataType string) ([]byte, error) {
	client := http.Client{}
	reqURL, _ := url.JoinPath(tr.srvAddr, urlsuff.OperationInsertNew, dataType)
	bodyReader := bytes.NewReader(bodyWithNewItems)
	var req *http.Request
	req, _ = http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bodyReader)
	req.Header.Add("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: tr.jwtToken,
	})

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, BadResponseHandler(resp, "From sync new items operation (Add)")
	}
	srvAns, err := io.ReadAll(resp.Body)

	return srvAns, err

	/*	err = json.NewDecoder(resp.Body).Decode(&srvAnswer)
		if err != nil {
			return err
		}
		//items.RemoveItemsWithoutID()
		//	err = items.AddOrUpdateItems(srvAnswer)
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
	*/

}

func (tr *TransportHTTP) UploadFiles(ctx context.Context, files []content.BinaryFileInfo) ([]content.BinaryFileInfo, error) {
	client := http.Client{}
	var returnErr error
	var wg sync.WaitGroup
	errorCh := make(chan error)
	filesCopy := make([]content.BinaryFileInfo, len(files))
	copy(filesCopy, files)

	go func() {
		for {
			newErr, ok := <-errorCh
			if ok {
				returnErr = fmt.Errorf("%s%s\r\n", returnErr, newErr.Error())
			} else {
				return
			}
		}
	}()

	for i, fileInfo := range filesCopy {
		//fileInfo := filesCopy[i]
		wg.Add(1)
		func() {
			defer wg.Done()
			downloadFilePath, err := tr.UploadFile(ctx, fileInfo.FilePath, fileInfo.ID, &client)
			if err != nil {
				errorCh <- fmt.Errorf("%s\r\n%s\r\n", fileInfo.FileName, err.Error())
				return
			}
			filesCopy[i].FilePath = downloadFilePath
		}()
	}

	wg.Wait()
	close(errorCh)
	return filesCopy, returnErr
}

func (tr *TransportHTTP) UploadFile(ctx context.Context, filepath string, fileID int, client *http.Client) (string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}

	defer file.Close()

	fileInfo, _ := file.Stat()

	fileSize := fileInfo.Size()

	totalPartsNum := int64(math.Ceil(float64(fileSize) / float64(tr.chunkSize)))

	for i := int64(0); i < totalPartsNum; i++ {
		partSize := int(math.Min(float64(tr.chunkSize), float64(fileSize-int64(i*tr.chunkSize))))
		partBuffer := make([]byte, partSize)
		_, err = file.Read(partBuffer)
		if err != nil {
			return "", err
		}

		resBody, err := tr.SendChunk(ctx, partBuffer, fileID, i, totalPartsNum, fileSize, client)
		if err != nil {
			return "", err
		}
		if i == totalPartsNum-1 && resBody != nil {
			return string(resBody), err
		}
	}
	return "", nil
}

func (tr *TransportHTTP) SyncChangesFirstStep(ctx context.Context, bodyIDsWithModtime []byte, dataType string) ([]byte, error) {
	client := http.Client{}
	bodyReader := bytes.NewReader(bodyIDsWithModtime)
	reqURL, _ := url.JoinPath(tr.srvAddr, urlsuff.OperationSync, urlsuff.SyncFirstStep, dataType)
	var req *http.Request
	/*	srvAnswer := struct {
			DataForSrv    []int       `json:",omitempty"`
			RemoveFromCli []int       `json:",omitempty"`
			DataForCli    interface{} `json:",omitempty"`
		}{}
	*/
	req, _ = http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bodyReader)
	req.Header.Add("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: tr.jwtToken,
	})

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, BadResponseHandler(resp, "From sync changes operation step one (Update)")
	}
	/*
	   srvAnswer := struct {
	   		DataForSrv    []int                    `json:",omitempty"`
	   		RemoveFromCli []int                    `json:",omitempty"`
	   		DataForCli    []content.CredentialInfo `json:",omitempty"`
	   	}{}
	   	err = json.NewDecoder(resp.Body).Decode(&srvAnswer)
	*/
	srvAns, err := io.ReadAll(resp.Body)
	return srvAns, err
	/*
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
	*/
}

func (tr *TransportHTTP) SyncChangesSecondStep(ctx context.Context, bodyDataForSrv []byte, dataType string) error {
	client := http.Client{}
	reqURL, _ := url.JoinPath(tr.srvAddr, urlsuff.OperationSync, urlsuff.SyncSecondStep, dataType)

	bodyReader := bytes.NewReader(bodyDataForSrv)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bodyReader)
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
	/*
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
	*/
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

func (tr *TransportHTTP) SendChunk(ctx context.Context, chunk []byte, fileID int, chunkNumber int64, chunkCount int64, totalFileSize int64, client *http.Client) ([]byte, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	metadata := struct {
		ChunkNumber   int64
		TotalChunks   int64
		FileID        int
		TotalFileSize int64
	}{
		ChunkNumber:   chunkNumber,
		TotalChunks:   chunkCount,
		FileID:        fileID,
		TotalFileSize: totalFileSize,
	}
	bodyJSON, err := json.Marshal(metadata)
	if err != nil {
		return nil, err
	}

	metadataHeader := textproto.MIMEHeader{}
	metadataHeader.Set("Content-Disposition", "form-data; name=\"metadata\"")
	metadataHeader.Set("Content-Type", "application/json")

	part, _ := writer.CreatePart(metadataHeader)
	_, err = part.Write(bodyJSON)
	if err != nil {
		return nil, err
	}

	mediaHeader := textproto.MIMEHeader{}
	mediaHeader.Set("Content-Disposition", "form-data; name=\"chunk\"")
	mediaPart, _ := writer.CreatePart(mediaHeader)

	_, err = io.Copy(mediaPart, bytes.NewReader(chunk))
	if err != nil {
		return nil, err
	}
	writer.Close()
	reqURL, _ := url.JoinPath(tr.srvAddr, urlsuff.DatatypeFile, urlsuff.FileOperationUpload)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, body)
	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: tr.jwtToken,
	})
	if err != nil {
		return nil, err
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
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		//BadResponseHandler(res, "")
		return nil, BadResponseHandler(res, "")
	}
	if res.ContentLength != 0 {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		return bodyBytes, nil
	}
	return nil, nil
}

func (tr *TransportHTTP) Clear() {
	tr.jwtToken = ""
}
