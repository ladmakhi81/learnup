package reqdto

type TusHookType string

const (
	TusHookType_PreCreate   TusHookType = "pre-create"
	TusHookType_PostCreate  TusHookType = "post-create"
	TusHookType_PostReceive TusHookType = "post-receive"
	TusHookType_PostFinish  TusHookType = "post-finish"
)

type TusWebhookDTO struct {
	Type  TusHookType `json:"Type"`
	Event struct {
		Upload struct {
			ID             string         `json:"ID"`
			Size           int64          `json:"Size"`
			SizeIsDeferred bool           `json:"SizeIsDeferred"`
			Offset         int64          `json:"Offset"`
			MetaData       map[string]any `json:"MetaData"`
			IsPartial      bool           `json:"IsPartial"`
			IsFinal        bool           `json:"IsFinal"`
			PartialUploads interface{}    `json:"PartialUploads"` // or []string if you expect a list
			Storage        map[string]any `json:"Storage"`        // could define this later if needed
		} `json:"Upload"`
		HTTPRequest struct {
			Method     string              `json:"Method"`
			URI        string              `json:"URI"`
			RemoteAddr string              `json:"RemoteAddr"`
			Header     map[string][]string `json:"Header"`
		} `json:"HTTPRequest"`
	} `json:"Event"`
}
