package storage

import (
	"encoding/json"
	"time"
	"ya-GophKeeper/internal/content"
)

type Collection interface {
	GetRemovedIDs() []int
	ClearRemovedList()

	GetNewItems() interface{}
	RemoveItemsWithoutID()
	AddOrUpdateItems(interface{}) error

	GetAllIDsWithModtime() map[int]time.Time
	GetItems([]int) interface{}
	RemoveItems([]int)
}

type CreditCards struct {
	stored  []content.CreditCardInfo
	removed []int
}

func (c *CreditCards) GetRemovedIDs() []int {
	return c.removed
}
func (c *CreditCards) ClearRemovedList() {
	c.removed = nil
}

func (c *CreditCards) GetNewItems() interface{} {
	var items []content.CreditCardInfo
	for _, item := range c.stored {
		if item.ID == 0 {
			items = append(items, item)
		}
	}

	return items
}
func (c *CreditCards) RemoveItemsWithoutID() {
	var newStored []content.CreditCardInfo
	for _, item := range c.stored {
		if item.ID != 0 {
			newStored = append(newStored, item)
		}
	}
	c.stored = newStored
}
func (c *CreditCards) AddOrUpdateItems(newItemsSlice interface{}) error {
	var newItemsWithType []content.CreditCardInfo
	jsonbody, err := json.Marshal(newItemsSlice)
	if err != nil {
		// do error check
		return err
	}
	if err = json.Unmarshal(jsonbody, &newItemsWithType); err != nil {
		return err
	}

	for i := range c.stored {
		for j := range newItemsWithType {
			if c.stored[i].ID == newItemsWithType[j].ID {
				c.stored[i] = newItemsWithType[j]
				newItemsWithType = append(newItemsWithType[:j], newItemsWithType[j+1:]...)
				break
			}
		}
	}
	c.stored = append(c.stored, newItemsWithType...)
	return nil
}

func (c *CreditCards) GetAllIDsWithModtime() map[int]time.Time {
	res := make(map[int]time.Time)
	for _, item := range c.stored {
		res[item.ID] = item.ModificationTime
	}
	return res
}
func (c *CreditCards) GetItems(IDs []int) interface{} {
	if IDs == nil {
		return c.stored
	}
	var items []content.CreditCardInfo
	for _, item := range c.stored {
		for _, id := range IDs {
			if item.ID == id {
				items = append(items, item)
				break
			}
		}
	}
	return items
}
func (c *CreditCards) RemoveItems(IDs []int) {
	var newStored []content.CreditCardInfo
	for _, item := range c.stored {
		ok := true
		for _, id := range IDs {
			if item.ID == id {
				ok = false
				break
			}
		}
		if ok {
			newStored = append(newStored, item)
		}
	}
	c.stored = newStored
}

type Credentials struct {
	stored  []content.CredentialInfo
	removed []int
}

func (c *Credentials) GetRemovedIDs() []int {
	return c.removed
}
func (c *Credentials) ClearRemovedList() {
	c.removed = nil
}

func (c *Credentials) GetNewItems() interface{} {
	var items []content.CredentialInfo
	for _, item := range c.stored {
		if item.ID == 0 {
			items = append(items, item)
		}
	}

	return items
}
func (c *Credentials) RemoveItemsWithoutID() {
	var newStored []content.CredentialInfo
	for _, item := range c.stored {
		if item.ID != 0 {
			newStored = append(newStored, item)
		}
	}
	c.stored = newStored
}
func (c *Credentials) AddOrUpdateItems(newItemsSlice interface{}) error {
	var newItemsWithType []content.CredentialInfo
	jsonbody, err := json.Marshal(newItemsSlice)
	if err != nil {
		// do error check
		return err
	}
	if err = json.Unmarshal(jsonbody, &newItemsWithType); err != nil {
		return err
	}

	for i := range c.stored {
		for j := range newItemsWithType {
			if c.stored[i].ID == newItemsWithType[j].ID {
				c.stored[i] = newItemsWithType[j]
				newItemsWithType = append(newItemsWithType[:j], newItemsWithType[j+1:]...)
				break
			}
		}
	}
	c.stored = append(c.stored, newItemsWithType...)
	return nil
}

func (c *Credentials) GetAllIDsWithModtime() map[int]time.Time {
	res := make(map[int]time.Time)
	for _, item := range c.stored {
		res[item.ID] = item.ModificationTime
	}
	return res
}
func (c *Credentials) GetItems(IDs []int) interface{} {
	if IDs == nil {
		return c.stored
	}
	var items []content.CredentialInfo
	for _, item := range c.stored {
		for _, id := range IDs {
			if item.ID == id {
				items = append(items, item)
				break
			}
		}
	}
	return items
}
func (c *Credentials) RemoveItems(IDs []int) {
	var newStored []content.CredentialInfo
	for _, item := range c.stored {
		ok := true
		for _, id := range IDs {
			if item.ID == id {
				ok = false
				break
			}
		}
		if ok {
			newStored = append(newStored, item)
		}
	}
	c.stored = newStored
}

type Texts struct {
	stored  []content.TextInfo
	removed []int
}

func (c *Texts) GetRemovedIDs() []int {
	return c.removed
}
func (c *Texts) ClearRemovedList() {
	c.removed = nil
}

func (c *Texts) GetNewItems() interface{} {
	var items []content.TextInfo
	for _, item := range c.stored {
		if item.ID == 0 {
			items = append(items, item)
		}
	}

	return items
}
func (c *Texts) RemoveItemsWithoutID() {
	var newStored []content.TextInfo
	for _, item := range c.stored {
		if item.ID != 0 {
			newStored = append(newStored, item)
		}
	}
	c.stored = newStored
}
func (c *Texts) AddOrUpdateItems(newItemsSlice interface{}) error {
	var newItemsWithType []content.TextInfo
	jsonbody, err := json.Marshal(newItemsSlice)
	if err != nil {
		// do error check
		return err
	}
	if err = json.Unmarshal(jsonbody, &newItemsWithType); err != nil {
		return err
	}

	for i := range c.stored {
		for j := range newItemsWithType {
			if c.stored[i].ID == newItemsWithType[j].ID {
				c.stored[i] = newItemsWithType[j]
				newItemsWithType = append(newItemsWithType[:j], newItemsWithType[j+1:]...)
				break
			}
		}
	}
	c.stored = append(c.stored, newItemsWithType...)
	return nil
}

func (c *Texts) GetAllIDsWithModtime() map[int]time.Time {
	res := make(map[int]time.Time)
	for _, item := range c.stored {
		res[item.ID] = item.ModificationTime
	}
	return res
}
func (c *Texts) GetItems(IDs []int) interface{} {
	if IDs == nil {
		return c.stored
	}
	var items []content.TextInfo
	for _, item := range c.stored {
		for _, id := range IDs {
			if item.ID == id {
				items = append(items, item)
				break
			}
		}
	}
	return items
}
func (c *Texts) RemoveItems(IDs []int) {
	var newStored []content.TextInfo
	for _, item := range c.stored {
		ok := true
		for _, id := range IDs {
			if item.ID == id {
				ok = false
				break
			}
		}
		if ok {
			newStored = append(newStored, item)
		}
	}
	c.stored = newStored
}

type Files struct {
	tempDir string
	stored  []content.BinaryFileInfo
	removed []int
}

func (c *Files) GetRemovedIDs() []int {
	return c.removed
}
func (c *Files) ClearRemovedList() {
	c.removed = nil
}

func (c *Files) GetNewItems() interface{} {
	var items []content.BinaryFileInfo
	for _, item := range c.stored {
		if item.ID == 0 {
			items = append(items, item)
		}
	}

	return items
}
func (c *Files) RemoveItemsWithoutID() {
	var newStored []content.BinaryFileInfo
	for _, item := range c.stored {
		if item.ID != 0 {
			newStored = append(newStored, item)
		}
	}
	c.stored = newStored
}
func (c *Files) AddOrUpdateItems(newItemsSlice interface{}) error {
	var newItemsWithType []content.BinaryFileInfo
	jsonbody, err := json.Marshal(newItemsSlice)
	if err != nil {
		// do error check
		return err
	}
	if err = json.Unmarshal(jsonbody, &newItemsWithType); err != nil {
		return err
	}

	for i := range c.stored {
		for j := range newItemsWithType {
			if c.stored[i].ID == newItemsWithType[j].ID {
				c.stored[i] = newItemsWithType[j]
				newItemsWithType = append(newItemsWithType[:j], newItemsWithType[j+1:]...)
				break
			}
		}
	}
	c.stored = append(c.stored, newItemsWithType...)
	return nil
}

func (c *Files) GetAllIDsWithModtime() map[int]time.Time {
	res := make(map[int]time.Time)
	for _, item := range c.stored {
		res[item.ID] = item.ModificationTime
	}
	return res
}
func (c *Files) GetItems(IDs []int) interface{} {
	if IDs == nil {
		return c.stored
	}
	var items []content.BinaryFileInfo
	for _, item := range c.stored {
		for _, id := range IDs {
			if item.ID == id {
				items = append(items, item)
				break
			}
		}
	}
	return items
}
func (c *Files) RemoveItems(IDs []int) {
	var newStored []content.BinaryFileInfo
	for _, item := range c.stored {
		ok := true
		for _, id := range IDs {
			if item.ID == id {
				ok = false
				break
			}
		}
		if ok {
			newStored = append(newStored, item)
		}
	}
	c.stored = newStored
}
