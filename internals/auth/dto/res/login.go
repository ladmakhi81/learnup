package dtores

type LoginResDto struct {
	AccessToken string `json:"accessToken"`
}

func NewLoginResDto(accessToken string) LoginResDto {
	return LoginResDto{
		AccessToken: accessToken,
	}
}
