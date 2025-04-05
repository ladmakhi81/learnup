package token

type Token interface {
	GenerateToken(userID uint) (string, error)
}
