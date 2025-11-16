package main

type apiData struct {
	Chart        string        `json:"chart"`
	Release      string        `json:"release"`
	Values       string        `json:"values"`
	FullValues   string        `json:"fullValues"`
	RenderValues string        `json:"renderValues"`
	PreviewFiles []apiDataFile `json:"previewFiles"`
}

type apiDataFile struct {
	Filename string `json:"filename"`
	Preview  string `json:"preview"`
}
