package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

type Post struct {
	CreatedAt   int64  `json:"created_at"`
	Title       string `json:"title"`
	Slug        string `json:"slug"`
	Content     string `json:"markdown"`
	PublishedAt int64  `json:"published_at"`
}

type ExportData struct {
	Posts []Post `json:"posts"`
}

type ExportEntry struct {
	Data ExportData `json:"data"`
}

type Export struct {
	DB []ExportEntry `json:"db"`
}

func main() {
	output_dir := "./content/posts/"
	file, e := ioutil.ReadFile("./GhostData.json")
	if e != nil {
		fmt.Printf("File error %v\n", e)
		os.Exit(1)
	}

	var export Export
	json.Unmarshal(file, &export)
	if len(export.DB) >= 1 {
		os.MkdirAll(output_dir, 0777)
		done := make(chan bool)

		for _, post := range export.DB[0].Data.Posts {
			go func(post Post) {
				postFile, err := os.Create(output_dir + post.Slug + ".md")
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				defer postFile.Close()

				w := bufio.NewWriter(postFile)

				fmt.Fprintln(w, "+++")
				fmt.Fprint(w, "date = \"%s\"\n", time.Unix(0, post.PublishedAt*int64(time.Millisecond)).UTC().String())
				fmt.Fprintln(w, "draft = false")
				fmt.Fprint(w, "title = \"%s\"\n", post.Title)
				fmt.Fprint(w, "slug = \"%s\"\n", post.Slug)
				fmt.Fprintln(w, "+++")
				fmt.Fprint(w, "%s", post.Content)
				done <- true
			}(post)
		}

		for _ = range export.DB[0].Data.Posts {
			<-done
		}
	}
}
