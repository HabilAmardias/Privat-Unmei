package utils

import (
	"crypto/rand"
	"math/big"
	"privat-unmei/internal/customerrors"
)


type OTPGenUtil struct{}

func CreateOTPUtil() *OTPGenUtil{
	return &OTPGenUtil{}
}

func (ogu *OTPGenUtil) GenerateOTP() (int64, error){
	min := big.NewInt(100000)
	max := big.NewInt(999999)

	rangeExclusive := new(big.Int).Sub(max, min)
	rangeExclusive.Add(rangeExclusive, big.NewInt(1))

	n, err := rand.Int(rand.Reader, rangeExclusive)
	if err != nil {
		return 0, customerrors.NewError(
			"something went wrong",
			err,
			customerrors.CommonErr,
		)
	}
	sixDigitInt := new(big.Int).Add(n, min)
	return sixDigitInt.Int64(), nil
}