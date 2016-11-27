package ghostToHugo

import (
	"os"
	"path/filepath"
	"testing"
)

func TestImportGhost(t *testing.T) {
	data := []struct {
		Filename string
	}{
		{"wrapped.json"},
		{"unwrapped.json"},
	}
	for _, d := range data {
		f, err := os.Open(filepath.Join("testdata", d.Filename))
		if err != nil {
			t.Error(err) //TODO: better message
		}
		gth, _ := NewGhostToHugo()
		posts, err := gth.importGhost(f)
		if err != nil {
			t.Error(err)
		}
		for p := range posts {
			t.Log(p)
		}
		err = f.Close()
		if err != nil {
			t.Error(err)
		}
	}
}
