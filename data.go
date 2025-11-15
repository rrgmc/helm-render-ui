package main

type apiData struct {
	Chart        string        `json:"chart"`
	Release      string        `json:"release"`
	Values       string        `json:"values"`
	RenderValues string        `json:"renderValues"`
	Preview      string        `json:"preview"`
	PreviewFiles []apiDataFile `json:"previewFiles"`
}

type apiDataFile struct {
	Filename string `json:"filename"`
	Preview  string `json:"preview"`
}
