package dockerhub

type ImageList struct {
	Count   int     `json:"count"`
	Next    string  `json:"next"`
	Results []Image `json:"results"`
}
