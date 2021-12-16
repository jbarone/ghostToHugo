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
	if !ok || description == nil {
		jww.ERROR.Println("cardBookmark: missing description")
		return ""
	}

	var thumbnail, icon, author, publisher, caption string
	thumb, ok := metadata["thumbnail"]
	if ok && thumb != nil {
		jww.TRACE.Println("cardBookmark: found thumbnail")
		thumbnail = stripContentFolder(thumb.(string))
	}
	iconinfo, ok := metadata["icon"]
	if ok && iconinfo != nil {
		jww.TRACE.Println("cardBookmark: found icon")
		icon = stripContentFolder(iconinfo.(string))
	}
	authorinfo, ok := metadata["author"]
	if ok && authorinfo != nil {
		jww.TRACE.Println("cardBookmark: found author")
		author = authorinfo.(string)
	}
	publisherinfo, ok := metadata["publisher"]
	if ok && publisherinfo != nil {
		jww.TRACE.Println("cardBookmark: found publisher")
		publisher = publisherinfo.(string)
	}
	capt, ok := m["caption"]
	if ok && capt != nil {
		jww.TRACE.Println("cardBookmark: found caption")
		caption = capt.(string)
	}

	return fmt.Sprintf(
		"{{< bookmark url=%q title=%q description=%q icon=%q"+
			" author=%q publisher=%q thumbnail=%q caption=%q >}}",
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
	images, ok := m["images"]
	if !ok {
		jww.ERROR.Println("cardGallery: missing images")
		return ""
	}

	var buf bytes.Buffer
	buf.WriteString("{{< gallery")
	if caption, ok := m["caption"]; ok {
		buf.WriteString(fmt.Sprintf(" caption=\"%s\"", caption.(string)))
	}
	buf.WriteString(" >}}\n")

	for _, img := range images.([]interface{}) {
		image, ok := img.(map[string]interface{})
		if !ok {
			continue
		}
		buf.WriteString("{{< galleryImg ")
		buf.WriteString(fmt.Sprintf(" src=%q", stripContentFolder(image["src"].(string))))
		buf.WriteString(fmt.Sprintf(" width=\"%.0f\"", image["width"].(float64)))
		buf.WriteString(fmt.Sprintf(" height=\"%.0f\"", image["height"].(float64)))
		if alt, ok := image["alt"]; ok {
			buf.WriteString(fmt.Sprintf(" alt=%q", alt.(string)))
		}
		if title, ok := image["title"]; ok {
			buf.WriteString(fmt.Sprintf(" title=%q", title.(string)))
		}
		buf.WriteString(" >}}")
	}

	buf.WriteString("{{< /gallery >}}")

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
			stripContentFolder(src.(string)),
			caption,
		)
	}

	return fmt.Sprintf(
		"{{< figure src=\"%s\" >}}\n",
		stripContentFolder(src.(string)),
	)
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
