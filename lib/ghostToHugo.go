package ghostToHugo

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/spf13/hugo/helpers"
	"github.com/spf13/hugo/hugolib"
	"github.com/spf13/hugo/parser"
	"github.com/spf13/viper"
)

// GhostToHugo handles the imprt of a Ghot blog export and outputting to
// hugo static blog
type GhostToHugo struct {
	location   *time.Location
	dateformat string
	path       string
}

// Post is a blog post in Ghost
type Post struct {
	ID              int             `json:"id"`
	Title           string          `json:"title"`
	Slug            string          `json:"slug"`
	Content         string          `json:"markdown"`
	Image           string          `json:"image"`
	Page            json.RawMessage `json:"page"`
	Status          string          `json:"status"`
	MetaDescription string          `json:"meta_description"`
	AuthorID        int             `json:"author_id"`
	PublishedAt     json.RawMessage `json:"published_at"`
	CreatedAt       json.RawMessage `json:"created_at"`

	Published time.Time
	Created   time.Time
	IsDraft   bool
	IsPage    bool
	Author    string
	Tags      []string
}

func (p *Post) populate(gi *ghostInfo, gth *GhostToHugo) {
	p.Published = gth.parseTime(p.PublishedAt)
	p.Created = gth.parseTime(p.CreatedAt)
	p.IsDraft = p.Status == "draft"
	p.IsPage = parseBool(p.Page)

	for _, user := range gi.users {
		if user.ID == p.AuthorID {
			p.Author = user.Name
			break
		}
	}

	for _, pt := range gi.posttags {
		if pt.PostID == p.ID {
			for _, t := range gi.tags {
				if t.ID == pt.TagID {
					p.Tags = append(p.Tags, t.Name)
					break
				}
			}
		}
	}
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

type meta struct {
	ExportedOn int64  `json:"exported_on"`
	Version    string `json:"version"`
}

type user struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type tag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type posttag struct {
	ID        int `json:"id"`
	PostID    int `json:"post_id"`
	TagID     int `json:"tag_id"`
	SortOrder int `json:"sort_order,omitempty"`
}

type ghostInfo struct {
	m        meta
	users    []user
	tags     []tag
	posttags []posttag
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
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigParseError); ok {
			return nil, err
		}
		return nil, fmt.Errorf("Unable to locate Config file. Perhaps you need to create a new site. (%s)\n", err)
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

func stripContentFolder(originalString string) string {
	return strings.Replace(originalString, "/content/", "/", -1)
}

func seekTo(d *json.Decoder, token json.Token) error {
	var tok json.Token
	var err error
	for err == nil && tok != token {
		tok, err = d.Token()
	}
	return err
}

func decodeGhostInfo(r io.Reader) (ghostInfo, error) {
	var gi ghostInfo
	var decoder = json.NewDecoder(r)
	var doneCount int

	for doneCount < 4 {
		tok, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return gi, err
		}

		switch tok {
		case "meta":
			err = decoder.Decode(&gi.m)
			doneCount++
		case "users":
			err = decoder.Decode(&gi.users)
			doneCount++
		case "tags":
			err = decoder.Decode(&gi.tags)
			doneCount++
		case "posts_tags":
			err = decoder.Decode(&gi.posttags)
			doneCount++
		}

		if err != nil {
			return gi, err
		}
	}

	return gi, nil
}

func (gth *GhostToHugo) importGhost(r io.ReadSeeker) (<-chan Post, error) {

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

	posts := make(chan Post)
	go func(decoder *json.Decoder, posts chan Post) {
		for decoder.More() {
			var p Post
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

func (p Post) getPath() string {
	if p.IsPage {
		return helpers.AbsPathify(
			path.Join(viper.GetString("contentDir"), p.Slug+".md"))
	}

	return helpers.AbsPathify(
		path.Join(viper.GetString("contentDir"), "post", p.Slug+".md"))
}

func (p Post) getMetadata() map[string]interface{} {
	metadata := make(map[string]interface{})

	switch p.IsDraft {
	case true:
		metadata["date"] = p.Created
	case false:
		metadata["date"] = p.Published
	}
	metadata["title"] = p.Title
	metadata["draft"] = p.IsDraft
	metadata["slug"] = p.Slug
	metadata["description"] = p.MetaDescription
	if p.Image != "" {
		metadata["image"] = stripContentFolder(p.Image)
	}
	if len(p.Tags) > 0 {
		metadata["tags"] = p.Tags
		metadata["categories"] = p.Tags
	}
	if p.Author != "" {
		metadata["author"] = p.Author
	}

	return metadata
}

func (gth *GhostToHugo) exportPosts(posts <-chan Post) {
	throttle := make(chan struct{}, 500)
	var wg sync.WaitGroup
	var site = hugolib.NewSiteDefaultLang()
	for post := range posts {
		wg.Add(1)
		go func(p Post) {
			throttle <- struct{}{}
			defer func() { <-throttle }()
			defer wg.Done()
			var name = p.getPath()
			// log.Println("saving file", name)
			page, err := site.NewPage(name)
			if err != nil {
				fmt.Printf("ERROR writing %s: %v\n", name, err)
				return
			}
			err = page.SetSourceMetaData(
				p.getMetadata(),
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
		}(post)
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
