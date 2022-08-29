package model

type LiveState struct {
	Title  string `json:"title"`
	IsLive bool   `json:"is_live"`
}

func (l *LiveState) GoLive(title string) {
	l.Title = title
	l.IsLive = true
}

func (l *LiveState) GoOffline() {
	l.IsLive = false
	l.Title = ""
}
