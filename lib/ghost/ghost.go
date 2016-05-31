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
	location *time.Location
}

var c config

func init() {
	c.location = time.Local
}

// SetLocation configures ghost to convert times based on the given location
func SetLocation(l *time.Location) {
	c.location = l
}

// Post is a blog post in Ghost
type Post struct {
	// id: {type: 'increments', nullable: false, primary: true}
	ID int `json:"id"`
	// uuid: {type: 'string', maxlength: 36, nullable: false, validations: {isUUID: true}}
	UUID string `json:"uuid"`
	// title: {type: 'string', maxlength: 150, nullable: false}
	Title string `json:"title"`
	// slug: {type: 'string', maxlength: 150, nullable: false, unique: true}
	Slug string `json:"slug"`
	// markdown: {type: 'text', maxlength: 16777215, fieldtype: 'medium', nullable: true}
	Content string `json:"markdown"`
	// html: {type: 'text', maxlength: 16777215, fieldtype: 'medium', nullable: true}
	HTML string `json:"html"`
	// image: {type: 'text', maxlength: 2000, nullable: true}
	Image string `json:"image"`
	// featured: {type: 'bool', nullable: false, defaultTo: false, validations: {isIn: [[0, 1, false, true]]}}
	Featured json.RawMessage `json:"featured"`
	// page: {type: 'bool', nullable: false, defaultTo: false, validations: {isIn: [[0, 1, false, true]]}}
	Page json.RawMessage `json:"page"`
	// status: {type: 'string', maxlength: 150, nullable: false, defaultTo: 'draft'}
	Status string `json:"status"`
	// language: {type: 'string', maxlength: 6, nullable: false, defaultTo: 'en_US'}
	Language string `json:"language"`
	// meta_title: {type: 'string', maxlength: 150, nullable: true}
	MetaTitle string `json:"meta_title"`
	// meta_description: {type: 'string', maxlength: 200, nullable: true}
	MetaDescription string `json:"meta_description"`
	// author_id: {type: 'integer', nullable: false}
	AuthorID int `json:"author_id"`
	// created_at: {type: 'dateTime', nullable: false}
	CreatedAt int64 `json:"created_at"`
	// created_by: {type: 'integer', nullable: false}
	CreatedBy int `json:"created_by"`
	// updated_at: {type: 'dateTime', nullable: true}
	UpdatedAt int64 `json:"updated_at"`
	// updated_by: {type: 'integer', nullable: true}
	UpdatedBy int `json:"updated_by"`
	// published_at: {type: 'dateTime', nullable: true}
	PublishedAt int64 `json:"published_at"`
	// published_by: {type: 'integer', nullable: true}
	PublishedBy int `json:"published_by"`
}

// Created is the time and date the Post was created
func (p Post) Created() time.Time {
	return time.Unix(0, p.CreatedAt*int64(time.Millisecond)).In(c.location)
}

// Updated is the time and date the Post was updated
func (p Post) Updated() time.Time {
	return time.Unix(0, p.UpdatedAt*int64(time.Millisecond)).In(c.location)
}

// Published is the time and date the Post was published
func (p Post) Published() time.Time {
	return time.Unix(0, p.PublishedAt*int64(time.Millisecond)).In(c.location)
}

// IsFeatured returns whether or not this Post is featured
func (p Post) IsFeatured() bool {
	var b bool
	if err := json.Unmarshal(p.Featured, &b); err == nil {
		return b
	}

	var i int
	if err := json.Unmarshal(p.Featured, &i); err == nil {
		return i != 0
	}

	return false
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
	// uuid: {type: 'string', maxlength: 36, nullable: false, validations: {isUUID: true}}
	UUID string `json:"uuid"`
	// name: {type: 'string', maxlength: 150, nullable: false}
	Name string `json:"name"`
	// slug: {type: 'string', maxlength: 150, nullable: false, unique: true}
	Slug string `json:"slug"`
	// password: {type: 'string', maxlength: 60, nullable: false}
	Password string `json:"password"`
	// email: {type: 'string', maxlength: 254, nullable: false, unique: true, validations: {isEmail: true}}
	Email string `json:"email"`
	// image: {type: 'text', maxlength: 2000, nullable: true}
	Image string `json:"image"`
	// cover: {type: 'text', maxlength: 2000, nullable: true}
	Cover string `json:"cover"`
	// bio: {type: 'string', maxlength: 200, nullable: true}
	Bio string `json:"bio"`
	// website: {type: 'text', maxlength: 2000, nullable: true, validations: {isEmptyOrURL: true}}
	Website string `json:"website"`
	// location: {type: 'text', maxlength: 65535, nullable: true}
	Location string `json:"location"`
	// accessibility: {type: 'text', maxlength: 65535, nullable: true}
	Accessibility string `json:"accessibility"`
	// status: {type: 'string', maxlength: 150, nullable: false, defaultTo: 'active'}
	Status string `json:"status"`
	// language: {type: 'string', maxlength: 6, nullable: false, defaultTo: 'en_US'}
	Language string `json:"language"`
	// meta_title: {type: 'string', maxlength: 150, nullable: true}
	MetaTitle string `json:"meta_title"`
	// meta_description: {type: 'string', maxlength: 200, nullable: true}
	MetaDescription string `json:"meta_description"`
	// tour: {type: 'text', maxlength: 65535, nullable: true}
	Tour string `json:"tour"`
	// last_login: {type: 'dateTime', nullable: true}
	LastLogin int64 `json:"last_login"`
	// created_at: {type: 'dateTime', nullable: false}
	CreatedAt int64 `json:"created_at"`
	// created_by: {type: 'integer', nullable: false}
	CreatedBy int `json:"created_by"`
	// updated_at: {type: 'dateTime', nullable: true}
	UpdatedAt int64 `json:"updated_at"`
	// updated_by: {type: 'integer', nullable: true}
	UpdatedBy int `json:"updated_by"`
}

// Created is the time and date the user was created
func (u User) Created() time.Time {
	return time.Unix(0, u.CreatedAt*int64(time.Millisecond)).In(c.location)
}

// Updated is the time and date the user was updated
func (u User) Updated() time.Time {
	return time.Unix(0, u.UpdatedAt*int64(time.Millisecond)).In(c.location)
}

// Tag is a tag that would be found on a post
type Tag struct {
	// id: {type: 'increments', nullable: false, primary: true}
	ID int `json:"id"`
	// uuid: {type: 'string', maxlength: 36, nullable: false, validations: {isUUID: true}}
	UUID string `json:"uuid"`
	// name: {type: 'string', maxlength: 150, nullable: false, validations: {matches: /^([^,]|$)/}}
	Name string `json:"name"`
	// slug: {type: 'string', maxlength: 150, nullable: false, unique: true}
	Slug string `json:"slug"`
	// description: {type: 'string', maxlength: 200, nullable: true}
	Description string `json:"description"`
	// image: {type: 'text', maxlength: 2000, nullable: true}
	Image string `json:"image"`
	// hidden: {type: 'bool', nullable: false, defaultTo: false, validations: {isIn: [[0, 1, false, true]]}}
	Hidden json.RawMessage `json:"hidden"`
	// parent_id: {type: 'integer', nullable: true}
	ParentID int `json:"parent_id"`
	// meta_title: {type: 'string', maxlength: 150, nullable: true}
	MetaTitle string `json:"meta_title"`
	// meta_description: {type: 'string', maxlength: 200, nullable: true}
	MetaDescription string `json:"meta_description"`
	// created_at: {type: 'dateTime', nullable: false}
	CreatedAt int64 `json:"created_at"`
	// created_by: {type: 'integer', nullable: false}
	CreatedBy int `json:"created_by"`
	// updated_at: {type: 'dateTime', nullable: true}
	UpdatedAt int64 `json:"updated_at"`
	// updated_by: {type: 'integer', nullable: true}
	UpdatedBy int `json:"updated_by"`
}

// IsHidden returns whether or not this Tag is hidden
func (t Tag) IsHidden() bool {
	var b bool
	if err := json.Unmarshal(t.Hidden, &b); err == nil {
		return b
	}

	var i int
	if err := json.Unmarshal(t.Hidden, &i); err == nil {
		return i != 0
	}

	return false
}

// Created is the time and date a tag was added to a post
func (t Tag) Created() time.Time {
	return time.Unix(0, t.CreatedAt*int64(time.Millisecond)).In(c.location)
}

// Updated is the time and date a tag was added to a post
func (t Tag) Updated() time.Time {
	return time.Unix(0, t.UpdatedAt*int64(time.Millisecond)).In(c.location)
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
