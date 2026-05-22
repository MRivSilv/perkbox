package cli

import (
	"crypto/rand"
	"math/big"
)

func charGen(length int, specialCount int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789@!#$%&?"
	const specialCharSet = "@!#$%&?"

	genSpecialCharset := make([]byte, specialCount)
	genCharset := make([]byte, length-specialCount)
	for i := range genCharset {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "Error generating password", err
		}
		genCharset[i] = charset[num.Int64()]
	}

	for x := range genSpecialCharset {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(specialCharSet))))
		if err != nil {
			return "Error generating special characters", err
		}
		genSpecialCharset[x] = specialCharSet[num.Int64()]
	}

	genPassword := append(genCharset, genSpecialCharset...)

	for i := len(genPassword) - 1; i > 0; i-- {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		if err != nil {
			return "Error mixing characters", err
		}
		j := num.Int64()
		genPassword[i], genPassword[j] = genPassword[j], genPassword[i]
	}
	return string(genPassword), nil
}
