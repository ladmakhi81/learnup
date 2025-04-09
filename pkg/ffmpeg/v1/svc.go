package ffmpegv1

import (
	"fmt"
	"github.com/google/uuid"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"io"
	"os"
)

type FfmpegSvc struct {
}

func NewFfmpegSvc() *FfmpegSvc {
	return &FfmpegSvc{}
}

func (svc FfmpegSvc) EncodeVideo(videoReader io.Reader) (string, error) {
	tmpDirID := uuid.NewString()
	tmpDir := fmt.Sprintf("/tmp/%s", tmpDirID)
	if err := os.MkdirAll(tmpDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("Error in creating directory: %s", err.Error())
	}
	playlistLocation := tmpDir + "/playlist.m3u8"
	segmentsLocation := tmpDir + "/segment%d.ts"
	err := ffmpeg.Input("pipe:0").
		Output(playlistLocation,
			ffmpeg.KwArgs{
				"c:v":                  "h264",
				"b:v":                  "1500k",
				"c:a":                  "aac",
				"preset":               "fast",
				"profile:v":            "baseline",
				"level":                "3.0",
				"hls_time":             "10",
				"hls_list_size":        "0",
				"hls_segment_filename": segmentsLocation,
				"hls_flags":            "independent_segments",
				"hls_playlist_type":    "vod",
			}).
		WithInput(videoReader).
		Run()
	if err != nil {
		return "", fmt.Errorf("Error in encode video by ffmpeg: %s", err.Error())
	}
	return tmpDir, nil
}
