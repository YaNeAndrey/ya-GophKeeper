package collection

import (
	"sync"
	"time"
	"ya-GophKeeper/internal/content"
)

type CreditCards struct {
	New     []content.CreditCardInfo
	Stored  []content.CreditCardInfo
	Removed []int
	mutex   sync.Mutex
}

func (c *CreditCards) GetRemovedIDs() []int {
	return c.Removed
}
func (c *CreditCards) ClearRemovedList() {
	c.Removed = nil
}

func (c *CreditCards) GetNewItems() []content.CreditCardInfo {
	var items []content.CreditCardInfo
	c.mutex.Lock()
	for _, item := range c.Stored {
		if item.ID == 0 {
			items = append(items, item)
		}
	}
	c.mutex.Unlock()
	return items
}
func (c *CreditCards) RemoveItemsWithoutID() {
	var newStored []content.CreditCardInfo
	c.mutex.Lock()
	for _, item := range c.Stored {
		if item.ID != 0 {
			newStored = append(newStored, item)
		}
	}
	c.Stored = newStored
	c.mutex.Unlock()
}
func (c *CreditCards) AddOrUpdateItems(newItemsSlice []content.CreditCardInfo) error {
	newItemsCopy := newItemsSlice
	c.mutex.Lock()
	for i := range c.Stored {
		for j := range newItemsCopy {
			if c.Stored[i].ID == newItemsCopy[j].ID {
				c.Stored[i] = newItemsCopy[j]
				newItemsCopy = append(newItemsCopy[:j], newItemsCopy[j+1:]...)
				break
			}
		}
	}
	c.Stored = append(c.Stored, newItemsCopy...)
	c.mutex.Unlock()
	return nil
}

func (c *CreditCards) GetAllIDsWithModtime() map[int]time.Time {
	res := make(map[int]time.Time)
	c.mutex.Lock()
	for _, item := range c.Stored {
		res[item.ID] = item.ModificationTime
	}
	c.mutex.Unlock()
	return res
}
func (c *CreditCards) GetItems(IDs []int) []content.CreditCardInfo {
	if IDs == nil {
		return c.Stored
	}
	var items []content.CreditCardInfo
	c.mutex.Lock()
	for _, item := range c.Stored {
		for _, id := range IDs {
			if item.ID == id {
				items = append(items, item)
				break
			}
		}
	}
	c.mutex.Unlock()
	return items
}
func (c *CreditCards) RemoveItems(IDs []int) {
	var newStored []content.CreditCardInfo
	c.mutex.Lock()
	for _, item := range c.Stored {
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
	c.Stored = newStored
	c.mutex.Unlock()
}

func (c *CreditCards) Clear() {
	c.mutex.Lock()
	c.Stored = nil
	c.Removed = nil
	c.mutex.Unlock()
}
