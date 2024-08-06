package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"ya-GophKeeper/internal/client/storage"
	"ya-GophKeeper/internal/constants/clerror"
	"ya-GophKeeper/internal/constants/urlsuff"
	"ya-GophKeeper/internal/content"
)

type TransportHTTP struct {
	srvAddr  string
	jwtToken string
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
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusConflict {
		return clerror.ErrLoginAlreadyTaken
	}
	if resp.StatusCode != http.StatusOK {
		return BadResponseHandler(resp, "From registration")
	}

	jwtToken := resp.Header.Get("token")
	if jwtToken == "" {
		return clerror.ErrHeaderWithAuthTokenIsEmpty
	}
	tr.jwtToken = jwtToken
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
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return BadResponseHandler(resp, "From login operation")
	}
	jwtToken := resp.Header.Get("token")
	if jwtToken == "" {
		return clerror.ErrHeaderWithAuthTokenIsEmpty
	}
	tr.jwtToken = jwtToken
	return nil
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
	err = tr.SyncChangesFirstChange(ctx, items, &client)
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
	req.Header.Add("token", tr.jwtToken)
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
	reqURL, _ := url.JoinPath(tr.srvAddr, urlsuff.OperationInsertNew)
	srvAnswer := struct {
		Items interface{}
	}{}
	switch items.(type) {
	case *storage.Credentials:
		reqURL, _ = url.JoinPath(reqURL, urlsuff.DatatypeCredential)
		srvAnswer.Items = []content.CredentialInfo{}
	case *storage.CreditCards:
		reqURL, _ = url.JoinPath(reqURL, urlsuff.DatatypeCreditCard)
		srvAnswer.Items = []content.CreditCardInfo{}
	case *storage.Texts:
		reqURL, _ = url.JoinPath(reqURL, urlsuff.DatatypeText)
		srvAnswer.Items = []content.TextInfo{}
	case *storage.Files:
		reqURL, _ = url.JoinPath(reqURL, urlsuff.DatatypeFile)
		//return SyncFiles()
		return nil
	default:
		return fmt.Errorf("SyncNewItems(TransportHTTP) %s", clerror.ErrIncorrectType)
	}

	newItems := items.GetNewItems()
	bodyJSON, err := json.Marshal(newItems)
	if err != nil {
		return err
	}
	bodyReader := bytes.NewReader(bodyJSON)
	var req *http.Request
	req, _ = http.NewRequest(http.MethodPost, reqURL, bodyReader)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("token", tr.jwtToken)
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
	err = items.AddOrUpdateItems(srvAnswer.Items)
	if err != nil {
		return err
	}
	return nil
}

func (tr *TransportHTTP) SyncChangesFirstChange(ctx context.Context, items storage.Collection, client *http.Client) error {
	IDsWithModtime := items.GetAllIDsWithModtime()
	bodyJSON, err := json.Marshal(IDsWithModtime)
	if err != nil {
		return err
	}
	bodyReader := bytes.NewReader(bodyJSON)
	reqURL, _ := url.JoinPath(tr.srvAddr, urlsuff.OperationSync)
	var req *http.Request
	srvAnswer := struct {
		IDs   []int
		Items interface{}
	}{}
	switch items.(type) {
	case *storage.Credentials:
		reqURL, _ = url.JoinPath(reqURL, urlsuff.DatatypeCredential)
		srvAnswer.Items = []content.CredentialInfo{}
	case *storage.CreditCards:
		reqURL, _ = url.JoinPath(reqURL, urlsuff.DatatypeCreditCard)
		srvAnswer.Items = []content.CreditCardInfo{}
	case *storage.Texts:
		reqURL, _ = url.JoinPath(reqURL, urlsuff.DatatypeText)
		srvAnswer.Items = []content.TextInfo{}
	case *storage.Files:
		reqURL, _ = url.JoinPath(reqURL, urlsuff.DatatypeFile)
		srvAnswer.Items = []content.BinaryFileInfo{}
	default:
		return fmt.Errorf("SyncChangesFirstStep(TransportHTTP) %s", clerror.ErrIncorrectType)
	}
	req, _ = http.NewRequest(http.MethodPost, reqURL, bodyReader)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("token", tr.jwtToken)
	reqURL, _ = url.JoinPath(reqURL, "1")

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
	err = items.AddOrUpdateItems(srvAnswer.Items)
	if err != nil {
		return err
	}
	return nil
}

func (tr *TransportHTTP) SyncChangesSecondChange(ctx context.Context, items storage.Collection, sendingItemsIDs []int, client *http.Client) error {
	reqURL, _ := url.JoinPath(tr.srvAddr, urlsuff.OperationSync)

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
		return fmt.Errorf("SyncChangesSecondStep(TransportHTTP) %s", clerror.ErrIncorrectType)
	}
	itemsForServer := items.GetItems(sendingItemsIDs)
	bodyJSON, err := json.Marshal(itemsForServer)
	if err != nil {
		return err
	}
	reqURL, _ = url.JoinPath(tr.srvAddr, "2")
	bodyReader := bytes.NewReader(bodyJSON)
	req, _ := http.NewRequest(http.MethodPost, reqURL, bodyReader)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("token", tr.jwtToken)
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return BadResponseHandler(resp, "From sync changes operation step two (Update)")
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
