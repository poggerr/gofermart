package accrual_service

//import "github.com/poggerr/gophermart/internal/storage"
//
//type AccrualRepo struct {
//	takeAccrualChan chan AccrualRepo
//	repository      storage.Storage
//}
//
//func NewRepo(strg *storage.Storage) *AccrualRepo {
//	return &AccrualRepo{
//		takeAccrualChan: make(chan AccrualRepo, 10),
//		repository:      *strg,
//	}
//}
//
//func (r *AccrualRepo) TakeAsync(ids []string, userID string) error {
//	r.takeAccrualChan <- storage.UserURLs{UserID: userID, URLs: ids}
//	return nil
//}
//
//func (r *AccrualRepo) WorkerTakeAccrual() {
//	for accrual := range r.takeAccrualChan {
//		accrual, err := TakeAccrual(accrual)
//		if err != nil {
//			return
//		}
//	}
//}
