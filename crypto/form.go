package crypto

type DecryptRequest struct {
	FileName   string `json:"fileName"`
	PassPhrase string `json:"passPhrase"`
}
