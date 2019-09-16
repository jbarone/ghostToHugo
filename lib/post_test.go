package ghostToHugo

import (
	"reflect"
	"testing"
	"time"
)

func TestMetadata(t *testing.T) {
	data := []struct {
		post     post
		expected map[string]interface{}
	}{
		{
			post{
				Title:           "Title",
				Slug:            "test",
				IsDraft:         true,
				MetaDescription: "test description",
				Created:         time.Now(),
			},
			map[string]interface{}{
				"title":       "Title",
				"draft":       true,
				"slug":        "test",
				"description": "test description",
			},
		},
		{
			post{
				Title:           "Title",
				Slug:            "test",
				IsDraft:         false,
				MetaDescription: "test description",
				Published:       time.Now(),
			},
			map[string]interface{}{
				"title":       "Title",
				"draft":       false,
				"slug":        "test",
				"description": "test description",
			},
		},
		{
			post{
				Title:           "Title",
				Slug:            "test",
				IsDraft:         false,
				MetaDescription: "test description",
				Published:       time.Now(),
				Image:           "/content/test.jpg",
			},
			map[string]interface{}{
				"title":       "Title",
				"draft":       false,
				"slug":        "test",
				"description": "test description",
				"image":       "/test.jpg",
			},
		},
		{
			post{
				Title:           "Title",
				Slug:            "test",
				IsDraft:         false,
				MetaDescription: "test description",
				Published:       time.Now(),
				Tags:            []string{"Test"},
			},
			map[string]interface{}{
				"title":       "Title",
				"draft":       false,
				"slug":        "test",
				"description": "test description",
				"tags":        []string{"Test"},
				"categories":  []string{"Test"},
			},
		},
		{
			post{
				Title:           "Title",
				Slug:            "test",
				IsDraft:         false,
				MetaDescription: "test description",
				Published:       time.Now(),
				Author:          "Aurthor",
			},
			map[string]interface{}{
				"title":       "Title",
				"draft":       false,
				"slug":        "test",
				"description": "test description",
				"author":      "Aurthor",
			},
		},
	}

	for _, d := range data {
		if reflect.DeepEqual(d.expected, d.post.metadata()) {
			t.Errorf(
				"Post.metadata() = %v; want %v",
				d.post.metadata(),
				d.expected,
			)
		}
	}
}

func Test_cardHR(t *testing.T) {
	type args struct {
		payload interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"simple",
			args{""},
			"---\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cardHR(tt.args.payload); got != tt.want {
				t.Errorf("cardHR() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cardImage(t *testing.T) {
	type args struct {
		payload interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"src_no_caption",
			args{
				map[string]interface{}{
					"src": "test.jpg",
				},
			},
			"{{< figure src=\"test.jpg\" >}}\n",
		},
		{
			"src_with_caption",
			args{
				map[string]interface{}{
					"src":     "test.jpg",
					"caption": "this is a test",
				},
			},
			"{{< figure src=\"test.jpg\" caption=\"this is a test\" >}}\n",
		},
		{
			"empty_src",
			args{
				map[string]interface{}{},
			},
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cardImage(tt.args.payload); got != tt.want {
				t.Errorf("cardImage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cardMarkdown(t *testing.T) {
	type args struct {
		payload interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"markdown",
			args{
				map[string]interface{}{
					"markdown": "markdown text",
				},
			},
			"markdown text\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cardMarkdown(tt.args.payload); got != tt.want {
				t.Errorf("cardMarkdown() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cardCode(t *testing.T) {
	type args struct {
		payload interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"code_no_lang",
			args{
				map[string]interface{}{
					"code": "print 'Hello'",
				},
			},
			"```\nprint 'Hello'\n```\n",
		},
		{
			"code_python",
			args{
				map[string]interface{}{
					"code":     "print 'Hello'",
					"language": "python",
				},
			},
			"```python\nprint 'Hello'\n```\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cardCode(tt.args.payload); got != tt.want {
				t.Errorf("cardCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cardEmbed(t *testing.T) {
	type args struct {
		payload interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"html",
			args{
				map[string]interface{}{
					"html": "<h1>Hello</h1>",
				},
			},
			"<h1>Hello</h1>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cardEmbed(tt.args.payload); got != tt.want {
				t.Errorf("cardEmbed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cardGallery(t *testing.T) {
	type args struct {
		payload interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"with_caption",
			args{
				map[string]interface{}{
					"images": []interface{}{
						map[string]interface{}{
							"fileName": "olya-kuzovkina-gAPXLS1LRVE-unsplash.jpg",
							"row":      0,
							"width":    float64(5184),
							"height":   float64(3456),
							"src":      "/content/images/2019/09/olya-kuzovkina-gAPXLS1LRVE-unsplash.jpg",
						},
						map[string]interface{}{
							"fileName": "q-aila-Y_pLBbSAhHI-unsplash.jpg",
							"row":      0,
							"width":    float64(5184),
							"height":   float64(3456),
							"src":      "/content/images/2019/09/q-aila-Y_pLBbSAhHI-unsplash.jpg",
						},
						map[string]interface{}{
							"fileName": "raul-varzar-1l2waV8glIQ-unsplash.jpg",
							"row":      0,
							"width":    float64(3200),
							"height":   float64(2361),
							"src":      "/content/images/2019/09/raul-varzar-1l2waV8glIQ-unsplash.jpg",
						},
					},
					"caption": "Kittens!",
				},
			},
			`<figure>
  <div>
    <div>
      <div><img src="/images/2019/09/olya-kuzovkina-gAPXLS1LRVE-unsplash.jpg" width="5184" height="3456"/></div>
      <div><img src="/images/2019/09/q-aila-Y_pLBbSAhHI-unsplash.jpg" width="5184" height="3456"/></div>
      <div><img src="/images/2019/09/raul-varzar-1l2waV8glIQ-unsplash.jpg" width="3200" height="2361"/></div>
    </div>
  </div>
  <figcaption>
    Kittens!
  </figcaption>
</figure>`,
		},
		{
			"without_caption",
			args{
				map[string]interface{}{
					"images": []interface{}{
						map[string]interface{}{
							"fileName": "olya-kuzovkina-gAPXLS1LRVE-unsplash.jpg",
							"row":      0,
							"width":    float64(5184),
							"height":   float64(3456),
							"src":      "/content/images/2019/09/olya-kuzovkina-gAPXLS1LRVE-unsplash.jpg",
						},
						map[string]interface{}{
							"fileName": "q-aila-Y_pLBbSAhHI-unsplash.jpg",
							"row":      0,
							"width":    float64(5184),
							"height":   float64(3456),
							"src":      "/content/images/2019/09/q-aila-Y_pLBbSAhHI-unsplash.jpg",
						},
						map[string]interface{}{
							"fileName": "raul-varzar-1l2waV8glIQ-unsplash.jpg",
							"row":      0,
							"width":    float64(3200),
							"height":   float64(2361),
							"src":      "/content/images/2019/09/raul-varzar-1l2waV8glIQ-unsplash.jpg",
						},
					},
				},
			},
			`<figure>
  <div>
    <div>
      <div><img src="/images/2019/09/olya-kuzovkina-gAPXLS1LRVE-unsplash.jpg" width="5184" height="3456"/></div>
      <div><img src="/images/2019/09/q-aila-Y_pLBbSAhHI-unsplash.jpg" width="5184" height="3456"/></div>
      <div><img src="/images/2019/09/raul-varzar-1l2waV8glIQ-unsplash.jpg" width="3200" height="2361"/></div>
    </div>
  </div>
</figure>`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cardGallery(tt.args.payload); got != tt.want {
				t.Errorf("cardGallery() = %v, want %v", got, tt.want)
			}
		})
	}
}
