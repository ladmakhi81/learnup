package zibalv1

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ladmakhi81/learnup/pkg/contracts"
	"github.com/ladmakhi81/learnup/pkg/dtos"
	"net/http"
	"strconv"
)

type ZibalClient struct {
	httpClient contracts.HttpClient
}

func NewZibalClient(
	httpClient contracts.HttpClient,
) *ZibalClient {
	return &ZibalClient{
		httpClient: httpClient,
	}
}

func (svc ZibalClient) CreateRequest(dto dtos.CreatePaymentGatewayDto) (*dtos.CreatePaymentGatewayResDto, error) {
	body := CreateRequestDTO{
		Merchant:    "zibal",
		CallbackURL: dto.CallbackURL,
		Amount:      dto.Amount,
	}
	httpResp, httpRespErr := svc.httpClient.Post(dtos.PostRequestDTO{
		URL:  "https://gateway.zibal.ir/v1/request",
		Body: body,
	})
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
		ID:      strconv.Itoa(resp.TrackID),
		PayLink: fmt.Sprintf("https://gateway.zibal.ir/start/%d", resp.TrackID),
	}, nil
}

func (svc ZibalClient) VerifyTransaction(dto dtos.VerifyTransactionDto) (*dtos.VerifyTransactionResDto, error) {
	parsedID, parsedIDErr := strconv.Atoi(dto.ID)
	if parsedIDErr != nil {
		return nil, parsedIDErr
	}

	body := VerifyRequestDTO{
		Merchant: "zibal",
		TrackID:  parsedID,
	}
	httpResp, httpRespErr := svc.httpClient.Post(dtos.PostRequestDTO{
		URL:  "https://gateway.zibal.ir/v1/verify",
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
		IsSuccess: resp.Message == "success",
		RefCode:   resp.RefNumber,
	}, nil
}
