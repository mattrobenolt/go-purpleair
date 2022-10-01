package purpleair

import (
	"context"
	"net/http"
)

func (c *Client) Keys(ctx context.Context) (*KeysResponse, error) {
	b, resp, err := c.get(ctx, "/keys", 201)
	if err != nil {
		return nil, err
	}

	kr := &KeysResponse{
		body:         b,
		httpResponse: resp,
	}
	err = decodeBody(b, kr)
	return kr, err
}

type KeysResponse struct {
	ApiVersion string `json:"api_version"`
	TimeStamp  int64  `json:"time_stamp"`
	ApiKeyType string `json:"api_key_type"`

	body         []byte
	httpResponse *http.Response
}

func (kr *KeysResponse) Body() []byte                 { return kr.body }
func (kr *KeysResponse) HttpResponse() *http.Response { return kr.httpResponse }
