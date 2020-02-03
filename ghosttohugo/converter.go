package ghosttohugo

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"strings"
	"time"

	"github.com/gohugoio/hugo/hugolib"
	"github.com/gohugoio/hugo/parser/metadecoders"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
)

// Converter is responsible for importing a Ghost blog export and converting
// it to a Hugo site.
type Converter struct {
	location   *time.Location
	dateformat string
	path       string
	force      bool
	info       info
	site       *hugolib.Site
	kind       metadecoders.Format
}

// WithLocation sets the location used when working with timestamps
func WithLocation(location *time.Location) func(*Converter) {
	return func(c *Converter) {
		c.location = location
	}
}

// WithDateFormat sets the date format to use for ghost imports
func WithDateFormat(format string) func(*Converter) {
	return func(c *Converter) {
		c.dateformat = format
	}
}

// WithHugoPath sets the path to the Hugo project being created
func WithHugoPath(path string) func(*Converter) {
	return func(c *Converter) {
		c.path = path
		viper.AddConfigPath(path)
	}
}

// WithForce sets the converter to forcefully overwrite a Hugo site.
func WithForce() func(*Converter) {
	return func(c *Converter) {
		c.force = true
	}
}

// New creates a new Converter configured with optional settings.
func New(options ...func(*Converter)) (*Converter, error) {
	c := &Converter{
		dateformat: time.RFC3339,
		location:   time.Local,
		path:       "newhugosite",
		kind:       metadecoders.TOML,
	}

	for _, option := range options {
		option(c)
	}

	return c, nil
}

func (c Converter) parseTime(raw json.RawMessage) time.Time {
	var pt int64
	if err := json.Unmarshal(raw, &pt); err == nil {
		return time.Unix(0, pt*int64(time.Millisecond)).In(c.location)
	}

	var ps string
	err := json.Unmarshal(raw, &ps)
	if err != nil {
		jww.ERROR.Printf("error unmarshalling time: %v\n", err)
		return time.Time{}
	}
	t, err := time.ParseInLocation(c.dateformat, ps, c.location)
	if err != nil {
		jww.ERROR.Printf("error parsing time: %v\n", err)
		return time.Time{}
	}
	return t

}

func (c Converter) populatePost(p *post) {
	p.Published = c.parseTime(p.PublishedAt)
	p.Created = c.parseTime(p.CreatedAt)

	for _, user := range c.info.Data.Users {
		if bytes.Equal(user.ID, p.AuthorID) {
			p.Author = user.Name
			break
		}
	}

	for _, posttag := range c.info.Data.PostTags {
		if !bytes.Equal(posttag.PostID, p.ID) {
			continue
		}

		for _, tag := range c.info.Data.Tags {
			if bytes.Equal(tag.ID, posttag.TagID) {
				p.Tags = append(p.Tags, strings.TrimPrefix(tag.Name, "#"))
				break
			}
		}
	}
}

func (c *Converter) Convert(r io.ReadSeeker) (int, error) {
	if err := c.decodeInfo(r); err != nil {
		return 0, err
	}

	if _, err := r.Seek(0, io.SeekStart); err != nil {
		return 0, err
	}

	if err := c.createSite(); err != nil {
		return 0, err
	}

	decoder := json.NewDecoder(r)
	err := seekTo(decoder, "posts")
	if err != nil {
		return 0, err
	}
	_, err = decoder.Token() // Strip Token
	if err != nil {
		return 0, err
	}

	var count int
	for decoder.More() {
		var p post
		err = decoder.Decode(&p)
		if err != nil {
			break
		}
		c.populatePost(&p)
		if err := c.writePost(p); err != nil {
			log.Fatal(err)
		}
		count++
	}

	return count, nil
}
