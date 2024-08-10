package book

type chapter struct {
	url        string
	body        string
	name        string
	author      string
	content     string
	subChapters []chapter
	config      *ScrapeConfig
}

func NewEmptyChapter() chapter {
	return chapter{"","", "", "", "", []chapter{}, NewScrapeConfigNoInclude()}
}

func NewChapter(url, body, name, author, content string, subChapters []chapter, config *ScrapeConfig) chapter {
	return chapter{url, body, name, author, content, subChapters, config}
}

func (c chapter) Body() string {
	return c.body
}

func (c chapter) Name() string {
	return c.name
}

func (c *chapter) SetName(name string) {
	c.name = name
}

func (c chapter) Author() string {
	return c.author
}

func (c chapter) Url() string {
	return c.url
}

func (c chapter) Content() string {
		return c.content
}

func (c chapter) SubChapters() []chapter {
	return c.subChapters
}

func (c *chapter) AddSubChapter(newChapter chapter) {
	c.subChapters = append(c.subChapters, newChapter)
}
