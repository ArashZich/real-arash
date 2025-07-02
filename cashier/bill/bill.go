package bill

import (
	"gitag.ir/armogroup/armo/services/reality/cashier/trait"
	"github.com/google/uuid"
)

// Bill is a struct that holds the invoice data
type Bill struct {
	uUID          string
	amount        uint32
	transactionID string
	description   string
	trait.HasDetail
}

// NewBill creates a new invoice
func NewBill() *Bill {
	return &Bill{
		uUID: uuid.New().String(),
	}
}

// SetAmount sets the amount of the invoice
func (i *Bill) SetAmount(amount uint32) *Bill {
	i.amount = amount
	return i
}

// GetAmount returns the amount of the invoice
func (i *Bill) GetAmount() uint32 {
	return i.amount
}

// SetUUID sets the UUID of the invoice
func (i *Bill) SetUUID(uid ...string) *Bill {
	if len(uid) > 0 {
		i.uUID = uid[0]
	}
	if i.uUID == "" {
		i.uUID = uuid.New().String()
	}
	return i
}

// GetUUID returns the UUID of the invoice
func (i *Bill) GetUUID() string {
	if i.uUID == "" {
		i.SetUUID()
	}

	return i.uUID
}

// SetTransactionID sets the transaction ID of the invoice
func (i *Bill) SetTransactionID(transactionID string) *Bill {
	i.transactionID = transactionID
	return i
}

// GetTransactionID returns the transaction ID of the invoice
func (i *Bill) GetTransactionID() string {
	return i.transactionID
}

// SetDescription sets the description of the invoice
func (i *Bill) SetDescription(description string) *Bill {
	i.description = description
	return i
}

// GetDescription returns the description of the invoice
func (i *Bill) GetDescription() string {
	return i.description
}
