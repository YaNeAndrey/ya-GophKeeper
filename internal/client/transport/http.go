package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"ya-GophKeeper/internal/client/const/clerror"
	"ya-GophKeeper/internal/client/storage"
	"ya-GophKeeper/internal/constants/urlsuff"
	"ya-GophKeeper/internal/content"
)

type TransportHTTP struct {
	srvAddr string
}

func InitTransport(srvAddr string) *TransportHTTP {
	return &TransportHTTP{srvAddr: srvAddr}
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
	err = tr.SyncUpdatedItems(ctx, items, &client)
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
	//add another Headers

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.Body.Close() != nil {
		return err
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

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
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

func (tr *TransportHTTP) SyncUpdatedItems(ctx context.Context, items storage.Collection, client *http.Client) error {
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
		return fmt.Errorf("SyncUpdatedItems(TransportHTTP) %s", clerror.ErrIncorrectType)
	}
	req, _ = http.NewRequest(http.MethodPost, reqURL, bodyReader)
	req.Header.Add("Content-Type", "application/json")
	reqURL, _ = url.JoinPath(reqURL, "1")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&srvAnswer)
	if err != nil {
		return err
	}
	err = items.AddOrUpdateItems(srvAnswer.Items)
	if err != nil {
		return err
	}

	buf := []rune(reqURL)
	buf[len(buf)-1] = '2'
	reqURL = string(buf)

	itemsForServer := items.GetItems(srvAnswer.IDs)
	bodyJSON, err = json.Marshal(itemsForServer)
	if err != nil {
		return err
	}
	bodyReader = bytes.NewReader(bodyJSON)

	([]rune(reqURL))[len(reqURL)-1] = '2'

	req, _ = http.NewRequest(http.MethodPost, reqURL, bodyReader)
	req.Header.Add("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		return err
	}
	return nil
}
