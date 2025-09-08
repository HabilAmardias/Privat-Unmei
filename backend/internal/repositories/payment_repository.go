package repositories

import "privat-unmei/internal/db"

type PaymentRepositoryImpl struct {
	DB *db.CustomDB
}

func CreatePaymentRepository(db *db.CustomDB) *PaymentRepositoryImpl {
	return &PaymentRepositoryImpl{db}
}
