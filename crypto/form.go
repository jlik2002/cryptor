package crypto

type FileRequest struct {
	FileName   string `json:"fileName"`
	PassPhrase string `json:"passPhrase"`
}
