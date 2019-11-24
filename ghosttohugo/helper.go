package ghosttohugo

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	jww "github.com/spf13/jwalterweatherman"
)

func stripContentFolder(original string) string {
	return strings.TrimPrefix(original, "/content")
}

func parseBool(rm json.RawMessage) bool {
	var b bool
	if err := json.Unmarshal(rm, &b); err == nil {
		return b
	}

	var i int
	if err := json.Unmarshal(rm, &i); err == nil {
		return i != 0
	}

	return false
}

func seekTo(d *json.Decoder, token json.Token) error {
	var tok json.Token
	var err error
	for err == nil && tok != token {
		tok, err = d.Token()
	}
	return err
}

func mkdir(x ...string) {
	p := filepath.Join(x...)

	err := os.MkdirAll(p, 0777) // before umask
	if err != nil {
		jww.FATAL.Fatalln(err)
	}
}
