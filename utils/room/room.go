package room

import (
	"crypto/rand"
	"math/big"
	"strings"
)

// GenerateCode generate verification code (like 6HNBD)
func GenerateCode(length int) (string, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var code []string
	for i := 0; i < length; i++ {
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		code = append(code, string(charset[index.Int64()]))
	}
	return strings.Join(code, ""), nil
}

// GenerateRoomID generate the roomID like xxx-ddfs-fff
func GenerateRoomID() (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	sections := []int{3, 4, 3}
	var roomID []string

	for _, length := range sections {
		var section []string
		for i := 0; i < length; i++ {
			index, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
			if err != nil {
				return "", err
			}
			section = append(section, string(charset[index.Int64()]))
		}
		roomID = append(roomID, strings.Join(section, ""))
	}

	return strings.Join(roomID, "-"), nil
}
