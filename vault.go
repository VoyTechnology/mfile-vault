// Package mfilevault establishes a connection to vault to get the data
package mfilevault // import "github.com/voytechnology/mfile-vault"

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hashicorp/vault/api"
	"github.com/voytechnology/mfile"
)

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrMultipleQueries = Error("mfile(vault): multiple query specified")
)

func init() {
	c, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		panic(fmt.Sprintf("mfile(vault): unable to connect: %v", err))
	}

	if err := mfile.Register("vault", handler{c.Logical()}); err != nil {
		panic(fmt.Sprintf("mfile(vault): regisration failed: %v", err))
	}
}

type handler struct {
	l *api.Logical
}

// ReadFile would return the JSON bytes. if a query is provided, the value
// would be looked up and returned.
// eg.
// 		data: foo/bar = {"a": "1", "b": "2"}
//
//      query: foo/bar
//      output: {"a": "1", "b": "2"}
//
//      query: foo/bar?b
//      output: 2
func (h handler) ReadFile(path string) ([]byte, error) {
	u, err := url.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("unable to read path: %w", err)
	}

	sec, err := h.l.Read(u.Path)
	if err != nil {
		return nil, err
	}

	switch len(u.Query()) {
	case 0:
		return json.Marshal(sec.Data)
	case 1:
		return []byte(fmt.Sprint(sec.Data[u.RawQuery])), nil
	default:
		return nil, ErrMultipleQueries
	}
}
