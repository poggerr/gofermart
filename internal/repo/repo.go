package repo

import (
	"github.com/google/uuid"
	"github.com/poggerr/gophermart/internal/storage"
)

type AccrualRepo struct {
	takeAccrualChan chan storage.SaveOrd
	repository      storage.Storage
}

func NewRepo(strg *storage.Storage) *AccrualRepo {
	return &AccrualRepo{
		takeAccrualChan: make(chan storage.SaveOrd, 10),
		repository:      *strg,
	}
}

func (r *AccrualRepo) TakeAsync(orderNum int, user *uuid.UUID, accrualURL string) error {
	r.takeAccrualChan <- storage.SaveOrd{
		OrderNum:   orderNum,
		User:       user,
		AccrualURL: accrualURL,
	}
	return nil
}

func (r *AccrualRepo) WorkerTakeAccrual() {
	for accrual := range r.takeAccrualChan {
		r.repository.SaveOrder(accrual)
	}

}
