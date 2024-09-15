package storage

import (
	"time"
	"ya-GophKeeper/internal/constants/clerror"
	"ya-GophKeeper/internal/content"
)

/*type Collection interface {
	GetRemovedIDs() []int
	ClearRemovedList()

	GetNewItems() interface{}
	RemoveItemsWithoutID()
	AddOrUpdateItems(interface{}) error

	GetAllIDsWithModtime() map[int]time.Time
	GetItems([]int) interface{}
	RemoveItems([]int)

	Clear()
}*/

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

func (c *CreditCards) GetNewItems() []content.CreditCardInfo {
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
func (c *CreditCards) AddOrUpdateItems(newItemsSlice []content.CreditCardInfo) error {
	if newItemsSlice == nil {
		return clerror.ErrNilSlice
	}
	newItemsCopy := newItemsSlice
	for i := range c.stored {
		for j := range newItemsCopy {
			if c.stored[i].ID == newItemsCopy[j].ID {
				c.stored[i] = newItemsCopy[j]
				newItemsCopy = append(newItemsCopy[:j], newItemsCopy[j+1:]...)
				break
			}
		}
	}
	c.stored = append(c.stored, newItemsCopy...)
	return nil
}

func (c *CreditCards) GetAllIDsWithModtime() map[int]time.Time {
	res := make(map[int]time.Time)
	for _, item := range c.stored {
		res[item.ID] = item.ModificationTime
	}
	return res
}
func (c *CreditCards) GetItems(IDs []int) []content.CreditCardInfo {
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
func (c *CreditCards) Clear() {
	c.stored = nil
	c.removed = nil
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

func (c *Credentials) GetNewItems() []content.CredentialInfo {
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
func (c *Credentials) AddOrUpdateItems(newItemsSlice []content.CredentialInfo) error {
	newItemsCopy := newItemsSlice
	for i := range c.stored {
		for j := range newItemsCopy {
			if c.stored[i].ID == newItemsCopy[j].ID {
				c.stored[i] = newItemsCopy[j]
				newItemsCopy = append(newItemsCopy[:j], newItemsCopy[j+1:]...)
				break
			}
		}
	}
	c.stored = append(c.stored, newItemsCopy...)
	return nil
}

func (c *Credentials) GetAllIDsWithModtime() map[int]time.Time {
	res := make(map[int]time.Time)
	for _, item := range c.stored {
		res[item.ID] = item.ModificationTime
	}
	return res
}
func (c *Credentials) GetItems(IDs []int) []content.CredentialInfo {
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
func (c *Credentials) Clear() {
	c.stored = nil
	c.removed = nil
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

func (c *Texts) GetNewItems() []content.TextInfo {
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
func (c *Texts) AddOrUpdateItems(newItemsSlice []content.TextInfo) error {
	newItemsCopy := newItemsSlice

	for i := range c.stored {
		for j := range newItemsCopy {
			if c.stored[i].ID == newItemsCopy[j].ID {
				c.stored[i] = newItemsCopy[j]
				newItemsCopy = append(newItemsCopy[:j], newItemsCopy[j+1:]...)
				break
			}
		}
	}
	c.stored = append(c.stored, newItemsCopy...)
	return nil
}

func (c *Texts) GetAllIDsWithModtime() map[int]time.Time {
	res := make(map[int]time.Time)
	for _, item := range c.stored {
		res[item.ID] = item.ModificationTime
	}
	return res
}
func (c *Texts) GetItems(IDs []int) []content.TextInfo {
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
func (c *Texts) Clear() {
	c.stored = nil
	c.removed = nil
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

func (c *Files) GetNewItems() []content.BinaryFileInfo {
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
func (c *Files) AddOrUpdateItems(newItemsSlice []content.BinaryFileInfo) error {
	newItemsCopy := newItemsSlice
	for i := range c.stored {
		for j := range newItemsCopy {
			if c.stored[i].ID == newItemsCopy[j].ID {
				c.stored[i] = newItemsCopy[j]
				newItemsCopy = append(newItemsCopy[:j], newItemsCopy[j+1:]...)
				break
			}
		}
	}
	c.stored = append(c.stored, newItemsCopy...)
	return nil
}

func (c *Files) GetAllIDsWithModtime() map[int]time.Time {
	res := make(map[int]time.Time)
	for _, item := range c.stored {
		res[item.ID] = item.ModificationTime
	}
	return res
}
func (c *Files) GetItems(IDs []int) []content.BinaryFileInfo {
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
func (c *Files) Clear() {
	c.stored = nil
	c.removed = nil
}
