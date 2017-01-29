package ghostToHugo

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/spf13/hugo/hugolib"
	"github.com/spf13/hugo/parser"
	"github.com/spf13/viper"
)

// GhostToHugo handles the import of a Ghot blog export and outputting to
// hugo static blog
type GhostToHugo struct {
	location   *time.Location
	dateformat string
	path       string
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

// WithLocation sets the location used when working with timestamps
func WithLocation(location *time.Location) func(*GhostToHugo) {
	return func(gth *GhostToHugo) {
		gth.location = location
	}
}

// WithDateFormat sets the date format to use for ghost imports
func WithDateFormat(format string) func(*GhostToHugo) {
	return func(gth *GhostToHugo) {
		gth.dateformat = format
	}
}

// WithHugoPath sets the path to the hugo project being written too
func WithHugoPath(path string) func(*GhostToHugo) {
	return func(gth *GhostToHugo) {
		gth.path = path
		viper.AddConfigPath(path)
	}
}

// NewGhostToHugo returns a new instance of GhostToHugo
func NewGhostToHugo(options ...func(*GhostToHugo)) (*GhostToHugo, error) {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("hugo")
	gth := new(GhostToHugo)

	// set defaults
	gth.dateformat = time.RFC3339
	gth.location = time.Local
	gth.path = "."

	for _, option := range options {
		option(gth)
	}

	viper.AddConfigPath(gth.path)
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigParseError); ok {
			return nil, err
		}
		return nil, fmt.Errorf("Unable to locate Config file. Perhaps you need to create a new site. (%v)", err)
	}
	viper.SetDefault("MetaDataFormat", "toml")
	viper.SetDefault("ContentDir", filepath.Join(gth.path, "content"))
	viper.SetDefault("LayoutDir", "layouts")
	viper.SetDefault("StaticDir", "static")
	viper.SetDefault("ArchetypeDir", "archetypes")
	viper.SetDefault("PublishDir", "public")
	viper.SetDefault("DataDir", "data")
	viper.SetDefault("ThemesDir", "themes")
	viper.SetDefault("DefaultLayout", "post")
	viper.SetDefault("Verbose", false)
	viper.SetDefault("Taxonomies", map[string]string{"tag": "tags", "category": "categories"})

	return gth, nil
}

func seekTo(d *json.Decoder, token json.Token) error {
	var tok json.Token
	var err error
	for err == nil && tok != token {
		tok, err = d.Token()
	}
	return err
}

func (gth *GhostToHugo) importGhost(r io.ReadSeeker) (<-chan post, error) {

	gi, err := decodeGhostInfo(r)
	if err != nil {
		return nil, err
	}

	_, err = r.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(r)
	err = seekTo(decoder, "posts")
	if err != nil {
		return nil, err
	}
	_, err = decoder.Token() // Strip Token
	if err != nil {
		return nil, err
	}

	posts := make(chan post)
	go func(decoder *json.Decoder, posts chan post) {
		for decoder.More() {
			var p post
			err = decoder.Decode(&p)
			if err != nil {
				break
			}
			p.populate(&gi, gth)
			posts <- p
		}
		close(posts)
	}(decoder, posts)

	return posts, nil
}

func (gth *GhostToHugo) parseTime(raw json.RawMessage) time.Time {
	var pt int64
	if err := json.Unmarshal(raw, &pt); err == nil {
		return time.Unix(0, pt*int64(time.Millisecond)).In(gth.location)
	}

	var ps string
	if err := json.Unmarshal(raw, &ps); err == nil {
		t, err := time.ParseInLocation(gth.dateformat, ps, gth.location)
		if err != nil {
			return time.Time{}
		}
		return t
	}

	return time.Time{}
}

func (gth *GhostToHugo) exportPosts(posts <-chan post) {
	throttle := make(chan struct{}, goMaxProcs()*5)
	var wg sync.WaitGroup
	var site = hugolib.NewSiteDefaultLang()
	for p := range posts {
		wg.Add(1)
		go func(p post) {
			throttle <- struct{}{}
			defer func() { <-throttle }()
			defer wg.Done()
			var name = p.path()
			log.Println("saving file", name)
			page, err := site.NewPage(name)
			if err != nil {
				fmt.Printf("ERROR writing %s: %v\n", name, err)
				return
			}
			err = page.SetSourceMetaData(
				p.metadata(),
				parser.FormatToLeadRune(viper.GetString("MetaDataFormat")))
			if err != nil {
				fmt.Printf("ERROR writing %s: %v\n", name, err)
				return
			}
			page.SetSourceContent([]byte(p.Content))
			err = page.SafeSaveSourceAs(name)
			if err != nil {
				fmt.Printf("ERROR writing %s: %v\n", name, err)
				return
			}
		}(p)
	}
	wg.Wait()
}

func (gth *GhostToHugo) Export(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Error opening export: %v", err)
	}
	defer file.Close()

	posts, err := gth.importGhost(file)
	if err != nil {
		log.Fatalf("Error processing Ghost export: %v", err)
	}

	gth.exportPosts(posts)
}

func goMaxProcs() int {
	if gmp := os.Getenv("GOMAXPROCS"); gmp != "" {
		if p, err := strconv.Atoi(gmp); err != nil {
			return p
		}
	}
	return 1
}
