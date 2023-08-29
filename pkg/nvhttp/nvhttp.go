package nvhttp

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
)

type Method string

const (
	GET  Method = "GET"
	POST Method = "POST"
)

func UnmarshalBody[T any](resp *http.Response, response *T) error {

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, &response)
	if err != nil {
		return err
	}

	return nil
}

type Nvhttp struct {
	token string
	debug bool
}

func New(token string) *Nvhttp {
	return &Nvhttp{
		token: token,
		debug: false,
	}
}

func (n *Nvhttp) SetDebug(debug bool) {
	n.debug = debug
}

func (n *Nvhttp) Get(url string) (*http.Response, error) {
	// Criação do cliente HTTP
	client := &http.Client{}

	req, err := http.NewRequest(string(GET), url, nil)
	if err != nil {

		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+n.token)
	req.Header.Add("Content-Type", "application/json") // Define o tipo de conteúdo do corpo

	requestDump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		return nil, err
	}
	fmt.Println("Requisição:")
	fmt.Println(string(requestDump))

	// Faz a requisição HTTP
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
