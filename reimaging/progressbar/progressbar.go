package progressbar

import (
	pbar "github.com/schollz/progressbar/v3"
)

type ProgressBarHandler interface {
	Add(value int) error
	Finish() error
}

type progressBar struct {
	pb *pbar.ProgressBar
}

func NewProgressBar(max int, title string) *progressBar {
	theme := pbar.Theme{
		Saucer:        "=",
		SaucerHead:    ">",
		SaucerPadding: " ",
		BarStart:      "[",
		BarEnd:        "]",
	}
	return &progressBar{
		pb: pbar.NewOptions64(
			int64(max),
			pbar.OptionShowIts(),
			pbar.OptionSetItsString("photos"),
			pbar.OptionSetDescription(title),
			pbar.OptionSetPredictTime(true),
			pbar.OptionShowCount(),
			pbar.OptionSetTheme(theme),
		),
	}
}

func (p *progressBar) Add(value int) error {
	max := p.pb.GetMax()
	current := int(p.pb.State().CurrentBytes)
	if current+value >= max {
		return p.pb.Add(max - current)
	}
	return p.pb.Add(value)
}

func (p *progressBar) Finish() error {
	return p.pb.Finish()
}
