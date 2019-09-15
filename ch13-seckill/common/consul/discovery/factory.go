package discovery

import (
	"context"
	_ "encoding/json"
	_ "errors"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"
	kithttp "github.com/go-kit/kit/transport/http"
	"io"
	_ "net/http"
	"net/url"
	_ "strconv"
	"strings"
)

func Factory(_ context.Context, method, path string, encodeRequest, decodeReponse) sd.Factory {
	return func(instance string) (endpoint endpoint.Endpoint, closer io.Closer, err error) {
		if !strings.HasPrefix(instance, "http") {
			instance = "http://" + instance
		}

		tgt, err := url.Parse(instance)
		if err != nil {
			return nil, nil, err
		}
		tgt.Path = path

		var (
			enc kithttp.EncodeRequestFunc
			dec kithttp.DecodeResponseFunc
		)
		enc, dec = encodeRequest, decodeReponse

		return kithttp.NewClient(method, tgt, enc, dec).Endpoint(), nil, nil
	}
}