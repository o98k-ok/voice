package pkg

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/duke-git/lancet/v2/netutil"
)

var (
	ErrHttpRequest = errors.New("http error happen")
)

func Request[T any](cli *netutil.HttpClient, req *netutil.HttpRequest, validate func(result *T) bool) (*T, error) {
	response, err := cli.SendRequest(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var result T
	if err = json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		fmt.Println(response.StatusCode)
		return nil, ErrHttpRequest
	}

	if !validate(&result) {
		fmt.Println(result)
		return nil, ErrHttpRequest
	}
	return &result, nil
}
