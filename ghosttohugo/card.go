package ghosttohugo

import (
	"bytes"
	"fmt"

	jww "github.com/spf13/jwalterweatherman"
)

func cardBookmark(payload interface{}) string {
	m, ok := payload.(map[string]interface{})
	if !ok {
		jww.ERROR.Println("cardBookmark: payload not correct type")
		return ""
	}

	meta, ok := m["metadata"]
	if !ok {
		jww.ERROR.Println("cardBookmark: payload does not contain metadata")
		return ""
	}

	metadata, ok := meta.(map[string]interface{})
	if !ok {
		jww.ERROR.Println("cardBookmark: payload metadata was not correct type")
		return ""
	}

	url, ok := metadata["url"]
	if !ok {
		jww.ERROR.Println("cardBookmark: missing url")
		return ""
	}
	title, ok := metadata["title"]
	if !ok {
		jww.ERROR.Println("cardBookmark: missing title")
		return ""
	}
	description, ok := metadata["description"]
	if !ok {
		jww.ERROR.Println("cardBookmark: missing description")
		return ""
	}

	var thumbnail, icon, author, publisher, caption string
	thumb, ok := metadata["thumbnail"]
	if ok {
		jww.TRACE.Println("cardBookmark: found thumbnail")
		thumbnail = fmt.Sprintf("<div><img src=%q></div>", thumb.(string))
	}
	iconinfo, ok := metadata["icon"]
	if ok {
		jww.TRACE.Println("cardBookmark: found icon")
		icon = fmt.Sprintf("<img src=%q>", iconinfo.(string))
	}
	authorinfo, ok := metadata["author"]
	if ok && authorinfo != nil {
		jww.TRACE.Println("cardBookmark: found author")
		author = fmt.Sprintf("<span>%s</span>", authorinfo.(string))
	}
	publisherinfo, ok := metadata["publisher"]
	if ok && publisherinfo != nil {
		jww.TRACE.Println("cardBookmark: found publisher")
		publisher = fmt.Sprintf("<span>%s</span>", publisherinfo.(string))
	}
	capt, ok := m["caption"]
	if ok && capt != nil {
		jww.TRACE.Println("cardBookmark: found caption")
		caption = fmt.Sprintf("<figcaption>%s</figcaption>", capt.(string))
	}

	return fmt.Sprintf(
		`<figure>
	     <a href=%q>
	       <div>
	         <div>%s</div>
	         <div>%s</div>
	         <div>
	           %s
	           %s
	           %s
	         </div>
	       </div>
	       %s
	     </a>
	     %s
	   </figure>`,
		url.(string),
		title.(string),
		description.(string),
		icon,
		author,
		publisher,
		thumbnail,
		caption,
	)
}

func cardCode(payload interface{}) string {
	m, ok := payload.(map[string]interface{})
	if !ok {
		jww.ERROR.Println("cardCode: payload not correct type")
		return ""
	}

	code, ok := m["code"]
	if !ok {
		jww.ERROR.Println("cardCode: missing code")
		return ""
	}

	var buf bytes.Buffer

	buf.WriteString("```")
	if lang, ok := m["language"]; ok {
		buf.WriteString(lang.(string))
	}
	buf.WriteString("\n")
	buf.WriteString(code.(string))
	buf.WriteString("\n```\n")

	return buf.String()
}

func cardEmbed(payload interface{}) string {
	m, ok := payload.(map[string]interface{})
	if !ok {
		jww.ERROR.Println("cardEmbed: payload not correct type")
		return ""
	}

	html, ok := m["html"]
	if !ok {
		jww.ERROR.Println("cardEmbed: missing html")
		return ""
	}

	return html.(string)
}
func cardGallery(payload interface{}) string {
	m, ok := payload.(map[string]interface{})
	if !ok {
		jww.ERROR.Println("cardGallery: payload not correct type")
		return ""
	}

	var buf bytes.Buffer
	buf.WriteString("<figure>\n")
	buf.WriteString("  <div>\n")   // gallery container
	buf.WriteString("    <div>\n") // start of first row
	images, ok := m["images"]
	if !ok {
		jww.ERROR.Println("cardGallery: missing images")
		return ""
	}

	for i, img := range images.([]interface{}) {
		if i > 0 && i%3 == 0 {
			buf.WriteString("    </div>") // end of current row
			buf.WriteString("    <div>")  // start of next row
		}
		image, ok := img.(map[string]interface{})
		if !ok {
			continue
		}
		buf.WriteString("      <div>")
		buf.WriteString("<img")
		buf.WriteString(fmt.Sprintf(" src=%q", stripContentFolder(image["src"].(string))))
		buf.WriteString(fmt.Sprintf(" width=\"%.0f\"", image["width"].(float64)))
		buf.WriteString(fmt.Sprintf(" height=\"%.0f\"", image["height"].(float64)))
		if alt, ok := image["alt"]; ok {
			buf.WriteString(fmt.Sprintf(" alt=%q", alt.(string)))
		}
		if title, ok := image["title"]; ok {
			buf.WriteString(fmt.Sprintf(" title=%q", title.(string)))
		}
		buf.WriteString("/>")
		buf.WriteString("</div>\n")
	}
	buf.WriteString("    </div>\n") // end of current row
	buf.WriteString("  </div>\n")   // end of gallery container

	if caption, ok := m["caption"]; ok {
		buf.WriteString(fmt.Sprintf("  <figcaption>\n    %s\n  </figcaption>\n", caption.(string)))
	}

	buf.WriteString("</figure>")

	return buf.String()
}

func cardHR(payload interface{}) string {
	return "---\n"
}

func cardHTML(payload interface{}) string {
	m, ok := payload.(map[string]interface{})
	if !ok {
		jww.ERROR.Println("cardHTML: payload not correct type")
		return ""
	}

	if html, ok := m["html"]; ok {
		return html.(string)
	}

	return ""
}

func cardImage(payload interface{}) string {
	m, ok := payload.(map[string]interface{})
	if !ok {
		jww.ERROR.Println("cardImage: payload not correct type")
		return ""
	}

	src, ok := m["src"]
	if !ok {
		jww.ERROR.Println("cardImage: missing src")
		return ""
	}

	if caption, ok := m["caption"]; ok {
		return fmt.Sprintf(
			"{{< figure src=\"%s\" caption=\"%s\" >}}\n",
			src,
			caption,
		)
	}

	return fmt.Sprintf("{{< figure src=\"%s\" >}}\n", src)
}

func cardMarkdown(payload interface{}) string {
	m, ok := payload.(map[string]interface{})
	if !ok {
		jww.ERROR.Println("cardMarkdown: payload not correct type")
		return ""
	}
	if markdown, ok := m["markdown"]; ok {
		return fmt.Sprintf("%s\n", markdown.(string))
	}
	return ""
}
