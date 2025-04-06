package dtores

type LoginRes struct {
	AccessToken string `json:"accessToken"`
}

func NewLoginRes(accessToken string) LoginRes {
	return LoginRes{
		AccessToken: accessToken,
	}
}
