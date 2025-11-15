package main

type apiData struct {
	Chart        string `json:"chart"`
	Release      string `json:"release"`
	Values       string `json:"values"`
	RenderValues string `json:"renderValues"`
	Preview      string `json:"preview"`
}
