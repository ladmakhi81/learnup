package zarinpalv1

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/pkg/dtos"
	"net/http"
	"strconv"
)

type ZarinpalClient struct {
	httpClient contracts.HttpClient
	config     *dtos.EnvConfig
}

func NewZarinpalClient(
	httpClient contracts.HttpClient,
	config *dtos.EnvConfig,
) *ZarinpalClient {
	return &ZarinpalClient{
		httpClient: httpClient,
		config:     config,
	}
}

func (svc ZarinpalClient) CreateRequest(dto dtos.CreatePaymentGatewayDto) (*dtos.CreatePaymentGatewayResDto, error) {
	body := CreateRequestDTO{
		Amount:      dto.Amount,
		CallbackURL: svc.config.Zarinpal.CallbackURL,
		Description: "...",
		MerchantID:  svc.config.Zarinpal.Merchant,
	}
	httpResp, httpRespErr := svc.httpClient.Post(dtos.PostRequestDTO{
		URL:  "https://sandbox.zarinpal.com/pg/v4/payment/request.json",
		Body: body,
	})
	fmt.Println(string(httpResp.Result), httpResp.StatusCode, httpRespErr)
	if httpRespErr != nil {
		return nil, httpRespErr
	}
	if httpResp.StatusCode != http.StatusOK {
		return nil, errors.New("status code is not okay")
	}
	var resp CreateRequestResDTO
	if err := json.Unmarshal(httpResp.Result, &resp); err != nil {
		return nil, err
	}
	return &dtos.CreatePaymentGatewayResDto{
		ID:      resp.Data.Authority,
		PayLink: fmt.Sprintf("https://sandbox.zarinpal.com/pg/StartPay/%s", resp.Data.Authority),
	}, nil
}

func (svc ZarinpalClient) VerifyTransaction(dto dtos.VerifyTransactionDto) (*dtos.VerifyTransactionResDto, error) {
	body := VerifyRequestDTO{
		Amount:     dto.Amount,
		Authority:  dto.ID,
		MerchantID: svc.config.Zarinpal.Merchant,
	}
	httpResp, httpRespErr := svc.httpClient.Post(dtos.PostRequestDTO{
		URL:  "https://sandbox.zarinpal.com/pg/v4/payment/verify.json",
		Body: body,
	})
	if httpRespErr != nil {
		return nil, httpRespErr
	}
	if httpResp.StatusCode != http.StatusOK {
		return nil, errors.New("status code is not okay")
	}
	var resp VerifyRequestResDTO
	if err := json.Unmarshal(httpResp.Result, &resp); err != nil {
		return nil, err
	}
	return &dtos.VerifyTransactionResDto{
		IsSuccess: resp.Data.Code == 100,
		RefCode:   strconv.Itoa(resp.Data.RefID),
	}, nil
}
