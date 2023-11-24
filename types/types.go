package types

type Product struct {
	ID        int    `json:"id"`
	Brand     string `json:"brand"`
	Model     string `json:"model"`
	Name      string `json:"name"`
	Color     string `json:"color"`
	Size      string `json:"size"`
	Processor string `json:"processor"`
	FrontMp   string `json:"frontMp"`
	BackMp    string `json:"backMp"`
	Ram       string `json:"ram"`
	Storage   string `json:"storage"`
	MadeIn    string `json:"madeIn"`
	ImageLoc  string `json:"ImageLoc"`
}
