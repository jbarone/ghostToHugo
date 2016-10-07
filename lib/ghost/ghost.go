package ghost

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"regexp"
	"time"
)

// https://github.com/TryGhost/Ghost/wiki/import-format#json-file-structure

// ExportReader is the interface to a reader of the Ghost export
type ExportReader struct {
	io.ReadSeeker
}

func (r *ExportReader) getRuneReader() io.RuneReader {
	return bufio.NewReader(r)
}

func (r *ExportReader) reset() error {
	_, err := r.Seek(0, 0)
	return err
}

type config struct {
	location   *time.Location
	dateformat string
}

var c config

func init() {
	c.location = time.Local
	c.dateformat = time.RFC3339
}

// SetLocation configures ghost to convert times based on the given location
func SetLocation(l *time.Location) {
	c.location = l
}

// SetDateFormat configures the format to use for parsing date strings
func SetDateFormat(format string) {
	c.dateformat = format
}

func parseTime(raw json.RawMessage) time.Time {
	var pt int64
	if err := json.Unmarshal(raw, &pt); err == nil {
		return time.Unix(0, pt*int64(time.Millisecond)).In(c.location)
	}

	var ps string
	if err := json.Unmarshal(raw, &ps); err == nil {
		t, _ := time.ParseInLocation(c.dateformat, ps, c.location)
		return t
	}

	return time.Time{}
}

// Post is a blog post in Ghost
type Post struct {
	// id: {type: 'increments', nullable: false, primary: true}
	ID int `json:"id"`
	// title: {type: 'string', maxlength: 150, nullable: false}
	Title string `json:"title"`
	// slug: {type: 'string', maxlength: 150, nullable: false, unique: true}
	Slug string `json:"slug"`
	// markdown: {type: 'text', maxlength: 16777215, fieldtype: 'medium', nullable: true}
	Content string `json:"markdown"`
	// image: {type: 'text', maxlength: 2000, nullable: true}
	Image string `json:"image"`
	// page: {type: 'bool', nullable: false, defaultTo: false, validations: {isIn: [[0, 1, false, true]]}}
	Page json.RawMessage `json:"page"`
	// status: {type: 'string', maxlength: 150, nullable: false, defaultTo: 'draft'}
	Status string `json:"status"`
	// meta_description: {type: 'string', maxlength: 200, nullable: true}
	MetaDescription string `json:"meta_description"`
	// author_id: {type: 'integer', nullable: false}
	AuthorID int `json:"author_id"`
	// published_at: {type: 'dateTime', nullable: true}
	PublishedAt json.RawMessage `json:"published_at"`
	// created_at: {type: 'dateTime', nullable: false}
	CreatedAt json.RawMessage `json:"created_at"`
}

// Published is the time and date the Post was published
func (p Post) Published() time.Time {
	return parseTime(p.PublishedAt)
}

// Created is the time and date the Post was created
func (p Post) Created() time.Time {
	return parseTime(p.CreatedAt)
}

// IsDraft returns whether or not this Post is a draft (unpublished)
func (p Post) IsDraft() bool {
	return p.Status == "draft"
}

// IsPage returns whether or not this Post is a page
func (p Post) IsPage() bool {
	var b bool
	if err := json.Unmarshal(p.Page, &b); err == nil {
		return b
	}

	var i int
	if err := json.Unmarshal(p.Page, &i); err == nil {
		return i != 0
	}

	return false
}

// Author returns the user that wrote the post
func (p Post) Author(data *ExportData) *User {
	if data != nil {
		for _, user := range data.Users {
			if user.ID == p.AuthorID {
				return &user
			}
		}
	}
	return nil
}

// Tags returns the tags that are associated with the Post
func (p Post) Tags(data *ExportData) []Tag {
	var tags []Tag

	for _, pt := range data.PostTags {
		if pt.PostID == p.ID {
			for _, t := range data.Tags {
				if t.ID == pt.TagID {
					tags = append(tags, t)
					break
				}
			}
		}
	}

	return tags
}

// User is a user in Ghost
type User struct {
	// id: {type: 'increments', nullable: false, primary: true}
	ID int `json:"id"`
	// name: {type: 'string', maxlength: 150, nullable: false}
	Name string `json:"name"`
}

// Tag is a tag that would be found on a post
type Tag struct {
	// id: {type: 'increments', nullable: false, primary: true}
	ID int `json:"id"`
	// name: {type: 'string', maxlength: 150, nullable: false, validations: {matches: /^([^,]|$)/}}
	Name string `json:"name"`
}

// PostTag is the combination of post and tag
type PostTag struct {
	// id: {type: 'increments', nullable: false, primary: true},
	ID int `json:"id"`
	// post_id: {type: 'integer', nullable: false, unsigned: true, references: 'posts.id'},
	PostID int `json:"post_id"`
	// tag_id: {type: 'integer', nullable: false, unsigned: true, references: 'tags.id'},
	TagID int `json:"tag_id"`
	// sort_order: {type: 'integer',  nullable: false, unsigned: true, defaultTo: 0}
	SortOrder int `json:"sort_order,omitempty"`
}

// ExportData is the "data" portion of the Ghost export
type ExportData struct {
	Posts    []Post    `json:"posts"`
	Users    []User    `json:"users"`
	Tags     []Tag     `json:"tags"`
	PostTags []PostTag `json:"posts_tags"`
}

// ExportMeta is the "meta" portion of the Ghost export
type ExportMeta struct {
	ExportedOn int64  `json:"exported_on"`
	Version    string `json:"version"`
}

// ExportEntry is a export entring in the Ghost export (entries in "db")
type ExportEntry struct {
	Meta ExportMeta `json:"meta"`
	Data ExportData `json:"data"`
}

// Export is the top level collection of the Ghost export
type Export struct {
	DB []ExportEntry `json:"db"`
}

func processWrapped(reader ExportReader) ([]ExportEntry, error) {
	var entry ExportEntry
	var entries []ExportEntry

	dec := json.NewDecoder(reader)

	// discard first 3 tokens
	for i := 0; i < 3; i++ {
		if _, err := dec.Token(); err != nil {
			log.Fatal(err)
		}
	}

	for dec.More() {
		if err := dec.Decode(&entry); err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}

	// discard last 2 tokens
	for i := 0; i < 2; i++ {
		if _, err := dec.Token(); err != nil {
			log.Fatal(err)
		}
	}

	return entries, nil
}

func processUnwrapped(reader ExportReader) ([]ExportEntry, error) {
	var entry ExportEntry

	dec := json.NewDecoder(reader)
	if err := dec.Decode(&entry); err != nil {
		return nil, err
	}

	return []ExportEntry{entry}, nil
}

// Process json data to generate Golang objects
func Process(reader ExportReader) ([]ExportEntry, error) {
	isWrapped, err := regexp.MatchReader(
		"\\s*{\\s*\"db\"\\s*:\\s*\\[",
		reader.getRuneReader(),
	)
	if err != nil {
		return nil, err
	}
	reader.reset()

	if isWrapped {
		return processWrapped(reader)
	}
	return processUnwrapped(reader)
}
