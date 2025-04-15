package ffmpeg

import "io"

type Ffmpeg interface {
	EncodeVideo(videoReader io.Reader) (string, error)
	GetVideoDuration(videoReader io.Reader) (string, error)
}
