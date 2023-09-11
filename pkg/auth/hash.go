package auth

import "golang.org/x/crypto/bcrypt"

type Hasher struct {
	cost int
}

func NewHasher(cost int) *Hasher {
	return &Hasher{cost: cost}
}

func (h Hasher) Compare(expected string, actual string) error {
	byteExpected := []byte(expected)
	byteActual := []byte(actual)
	return bcrypt.CompareHashAndPassword(byteExpected, byteActual)
}

func (h Hasher) Hash(element string) (string, error) {
	byteElement := []byte(element)
	result, err := bcrypt.GenerateFromPassword(byteElement, h.cost)
	return string(result), err
}
