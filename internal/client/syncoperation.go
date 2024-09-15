package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gosuri/uilive"
	"time"
	"ya-GophKeeper/internal/constants/urlsuff"
	"ya-GophKeeper/internal/content"

	log "github.com/sirupsen/logrus"
)

func SyncCollection(c *Client, dataType string) error {
	finishCh := make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())
	go SynchronizationPrinter(ctx, finishCh)
	var err error
	switch dataType {
	case urlsuff.DatatypeCredential:
		err = SyncCredentials(c)
	case urlsuff.DatatypeCreditCard:
		err = SyncCreditCards(c)
	case urlsuff.DatatypeText:
		err = SyncTexts(c)
	case urlsuff.DatatypeFile:
		err = SyncFiles(c)
	}
	cancel()
	<-finishCh
	if err != nil {
		return err
	}
	return nil
}

func SyncCreditCards(c *Client) error {
	creditCards := c.storage.GetCreditCardsData()
	rem := creditCards.GetRemovedIDs()
	ctx := context.Background()
	if rem != nil {
		err := c.transport.SyncRemovedItems(ctx, rem, urlsuff.DatatypeCreditCard)
		if err != nil {
			return err
		}
		creditCards.ClearRemovedList()
	}
	newCreditCards := creditCards.GetNewItems()
	if newCreditCards != nil {
		bodyJSON, err := json.Marshal(newCreditCards)
		if err != nil {
			return err
		}
		srvAns, err := c.transport.SyncNewItems(ctx, bodyJSON, urlsuff.DatatypeCreditCard)
		if err != nil {
			return err
		}
		var newCreditCardsFromSrv []content.CreditCardInfo
		err = json.Unmarshal(srvAns, &newCreditCardsFromSrv)
		if err != nil {
			return err
		}
		creditCards.RemoveItemsWithoutID()
		err = creditCards.AddOrUpdateItems(newCreditCardsFromSrv)
		if err != nil {
			log.Println(err)
		}
	}

	IDsWithModtime := creditCards.GetAllIDsWithModtime()
	//if IDsWithModtime == nil {
	bodyJSON, err := json.Marshal(IDsWithModtime)
	if err != nil {
		return err
	}
	srvAnsBytes, err := c.transport.SyncChangesFirstStep(ctx, bodyJSON, urlsuff.DatatypeCreditCard)
	if err != nil {
		return err
	}
	srvAnswer := struct {
		DataForSrv    []int                    `json:",omitempty"`
		RemoveFromCli []int                    `json:",omitempty"`
		DataForCli    []content.CreditCardInfo `json:",omitempty"`
	}{}
	err = json.Unmarshal(srvAnsBytes, &srvAnswer)
	if err != nil {
		return err
	}

	err = creditCards.AddOrUpdateItems(srvAnswer.DataForCli)
	if err != nil {
		log.Println(err)
	}

	if srvAnswer.RemoveFromCli != nil {
		creditCards.RemoveItems(srvAnswer.RemoveFromCli)
	}

	if srvAnswer.DataForSrv != nil {
		itemsForServer := creditCards.GetItems(srvAnswer.DataForSrv)
		bodyJSON, err = json.Marshal(itemsForServer)
		if err != nil {
			return err
		}
		err = c.transport.SyncChangesSecondStep(ctx, bodyJSON, urlsuff.DatatypeCreditCard)
		if err != nil {
			return err
		}
	}
	//}
	return nil
}
func SyncCredentials(c *Client) error {
	credentials := c.storage.GetCredentialsData()
	rem := credentials.GetRemovedIDs()
	ctx := context.Background()
	if rem != nil {
		err := c.transport.SyncRemovedItems(ctx, rem, urlsuff.DatatypeCredential)
		if err != nil {
			return err
		}
		credentials.ClearRemovedList()
	}
	newCredentials := credentials.GetNewItems()
	if newCredentials != nil {
		bodyJSON, err := json.Marshal(newCredentials)
		if err != nil {
			return err
		}
		srvAns, err := c.transport.SyncNewItems(ctx, bodyJSON, urlsuff.DatatypeCredential)
		if err != nil {
			return err
		}
		var newCredentialsFromSrv []content.CredentialInfo
		err = json.Unmarshal(srvAns, &newCredentialsFromSrv)
		if err != nil {
			return err
		}
		credentials.RemoveItemsWithoutID()
		err = credentials.AddOrUpdateItems(newCredentialsFromSrv)
		if err != nil {
			log.Println(err)
		}
	}

	IDsWithModtime := credentials.GetAllIDsWithModtime()
	//if IDsWithModtime == nil {
	bodyJSON, err := json.Marshal(IDsWithModtime)
	if err != nil {
		return err
	}
	srvAnsBytes, err := c.transport.SyncChangesFirstStep(ctx, bodyJSON, urlsuff.DatatypeCredential)
	if err != nil {
		return err
	}
	srvAnswer := struct {
		DataForSrv    []int                    `json:",omitempty"`
		RemoveFromCli []int                    `json:",omitempty"`
		DataForCli    []content.CredentialInfo `json:",omitempty"`
	}{}
	err = json.Unmarshal(srvAnsBytes, &srvAnswer)
	if err != nil {
		return err
	}

	err = credentials.AddOrUpdateItems(srvAnswer.DataForCli)
	if err != nil {
		log.Println(err)
	}

	if srvAnswer.RemoveFromCli != nil {
		credentials.RemoveItems(srvAnswer.RemoveFromCli)
	}

	if srvAnswer.DataForSrv != nil {
		itemsForServer := credentials.GetItems(srvAnswer.DataForSrv)
		bodyJSON, err = json.Marshal(itemsForServer)
		if err != nil {
			return err
		}
		err = c.transport.SyncChangesSecondStep(ctx, bodyJSON, urlsuff.DatatypeCredential)
		if err != nil {
			return err
		}
	}
	//}
	return nil
}
func SyncTexts(c *Client) error {
	texts := c.storage.GetTextsData()
	rem := texts.GetRemovedIDs()
	ctx := context.Background()
	if rem != nil {
		err := c.transport.SyncRemovedItems(ctx, rem, urlsuff.DatatypeText)
		if err != nil {
			return err
		}
		texts.ClearRemovedList()
	}
	newTexts := texts.GetNewItems()
	if newTexts != nil {
		bodyJSON, err := json.Marshal(newTexts)
		if err != nil {
			return err
		}
		srvAns, err := c.transport.SyncNewItems(ctx, bodyJSON, urlsuff.DatatypeText)
		if err != nil {
			return err
		}
		var newTextsFromSrv []content.TextInfo
		err = json.Unmarshal(srvAns, &newTextsFromSrv)
		if err != nil {
			return err
		}
		texts.RemoveItemsWithoutID()
		err = texts.AddOrUpdateItems(newTextsFromSrv)
		if err != nil {
			log.Println(err)
		}
	}

	IDsWithModtime := texts.GetAllIDsWithModtime()
	//if IDsWithModtime == nil {
	bodyJSON, err := json.Marshal(IDsWithModtime)
	if err != nil {
		return err
	}
	srvAnsBytes, err := c.transport.SyncChangesFirstStep(ctx, bodyJSON, urlsuff.DatatypeText)
	if err != nil {
		return err
	}
	srvAnswer := struct {
		DataForSrv    []int              `json:",omitempty"`
		RemoveFromCli []int              `json:",omitempty"`
		DataForCli    []content.TextInfo `json:",omitempty"`
	}{}
	err = json.Unmarshal(srvAnsBytes, &srvAnswer)
	if err != nil {
		return err
	}

	err = texts.AddOrUpdateItems(srvAnswer.DataForCli)
	if err != nil {
		log.Println(err)
	}

	if srvAnswer.RemoveFromCli != nil {
		texts.RemoveItems(srvAnswer.RemoveFromCli)
	}

	if srvAnswer.DataForSrv != nil {
		itemsForServer := texts.GetItems(srvAnswer.DataForSrv)
		bodyJSON, err = json.Marshal(itemsForServer)
		if err != nil {
			return err
		}
		err = c.transport.SyncChangesSecondStep(ctx, bodyJSON, urlsuff.DatatypeText)
		if err != nil {
			return err
		}
	}
	//}
	return nil
}
func SyncFiles(c *Client) error {
	files := c.storage.GetFilesData()
	rem := files.GetRemovedIDs()
	ctx := context.Background()
	if rem != nil {
		err := c.transport.SyncRemovedItems(ctx, rem, urlsuff.DatatypeFile)
		if err != nil {
			return err
		}
		files.ClearRemovedList()
	}
	newFiles := files.GetNewItems()
	if newFiles != nil {
		bodyJSON, err := json.Marshal(newFiles)
		if err != nil {
			return err
		}
		srvAns, err := c.transport.SyncNewItems(ctx, bodyJSON, urlsuff.DatatypeFile)
		if err != nil {
			return err
		}
		var newFilesFromSrv []content.BinaryFileInfo
		err = json.Unmarshal(srvAns, &newFilesFromSrv)
		if err != nil {
			return err
		}
		files.RemoveItemsWithoutID()
		err = files.AddOrUpdateItems(newFilesFromSrv)
		if err != nil {
			log.Println(err)
		}
	}

	IDsWithModtime := files.GetAllIDsWithModtime()
	//if IDsWithModtime == nil {
	bodyJSON, err := json.Marshal(IDsWithModtime)
	if err != nil {
		return err
	}
	srvAnsBytes, err := c.transport.SyncChangesFirstStep(ctx, bodyJSON, urlsuff.DatatypeFile)
	if err != nil {
		return err
	}
	srvAnswer := struct {
		DataForSrv    []int                    `json:",omitempty"`
		RemoveFromCli []int                    `json:",omitempty"`
		DataForCli    []content.BinaryFileInfo `json:",omitempty"`
	}{}
	err = json.Unmarshal(srvAnsBytes, &srvAnswer)
	if err != nil {
		return err
	}

	err = files.AddOrUpdateItems(srvAnswer.DataForCli)
	if err != nil {
		log.Println(err)
	}

	if srvAnswer.RemoveFromCli != nil {
		files.RemoveItems(srvAnswer.RemoveFromCli)
	}

	if srvAnswer.DataForSrv != nil {
		itemsForServer := files.GetItems(srvAnswer.DataForSrv)
		bodyJSON, err = json.Marshal(itemsForServer)
		if err != nil {
			return err
		}
		err = c.transport.SyncChangesSecondStep(ctx, bodyJSON, urlsuff.DatatypeFile)
		if err != nil {
			return err
		}
	}
	//}
	return nil
}

func SyncMonitor(c *Client, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return

		case <-time.After(c.config.SyncInterval):
			err := FullSync(c)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func SynchronizationPrinter(ctx context.Context, finishCh chan struct{}) {
	writer := uilive.New()
	writer.Start()
	defer writer.Stop()

	for i := 1; ; i++ {
		select {
		case <-ctx.Done():
			fmt.Fprintf(writer, "Synchronization complete\n")
			finishCh <- struct{}{}
			return
		default:
			str := ""
			for j := 0; j < i%20; j++ {
				str += "*"
			}
			fmt.Fprintf(writer, "Synchronization: %s\n", str)
		}
		time.Sleep(time.Millisecond * 10)
	}
}

func FullSync(c *Client) error {
	fmt.Print("Credentials: ")
	err := SyncCollection(c, urlsuff.DatatypeCredential)
	if err != nil {
		return err
	}
	fmt.Print("Credit Cards: ")
	err = SyncCollection(c, urlsuff.DatatypeCreditCard)
	if err != nil {
		return err
	}
	fmt.Print("Texts: ")
	err = SyncCollection(c, urlsuff.DatatypeText)
	if err != nil {
		return err
	}

	fmt.Print("Files: ")
	err = SyncCollection(c, urlsuff.DatatypeFile)
	if err != nil {
		return err
	}

	return nil
}
