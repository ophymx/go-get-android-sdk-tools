package sdk

// ProgressRenderer renders progress of downloads
type ProgressRenderer interface {
	Init()
	StartPhase(phase string)
	Progress(percent float32)
	Complete()
}

type nullRenderer struct {
}

func (r nullRenderer) Init() {
}

func (r nullRenderer) StartPhase(phase string) {
}

func (r nullRenderer) Progress(percent float32) {
}

func (r nullRenderer) Complete() {
}

// TODO(abic): convert this to be an io.Reader instead of an io.Writer
type progressWriter struct {
	renderer ProgressRenderer
	size     int64
	total    int64
}

func (w *progressWriter) Write(p []byte) (n int, err error) {
	length := len(p)

	w.size += int64(length)
	percent := float32(w.size) / float32(w.total)
	w.renderer.Progress(percent)
	return length, nil
}
