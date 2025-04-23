package entities

type VideoStatus string

const (
	VideoStatus_Pending VideoStatus = "pending"
	VideoStatus_Fail    VideoStatus = "fail"
	VideoStatus_Done    VideoStatus = "done"
)
