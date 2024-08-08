package book

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
	epub "github.com/bmaupin/go-epub"
)

func Filename(name string) string {
	filename := name

	filename = strings.ReplaceAll(filename, " ", "_")
	filename = strings.ReplaceAll(filename, "/", "")

	return filename
}

// TODO this saves it in markdown file with title and contents
// TODO: edit this to save it in a file
func ToMarkdownString(c chapter) string {
	markdown := ""

	// chapter content
	if c.config.Include {
		// title
		markdown += fmt.Sprintf("%s\n", c.Name())
		markdown += fmt.Sprintf("%s\n\n", strings.Repeat("=", len(c.Name())))

		// convert content to markdown
		content, err := md.NewConverter("", true, nil).ConvertString(c.Content())
		if err != nil {
			log.Fatal(err)
		}
		markdown += fmt.Sprintf("%s\n\n\n", content)
	}

	// subchapters content
	for _, sc := range c.SubChapters() {
		markdown += fmt.Sprintf("%s\n\n\n", ToMarkdownString(sc))
	}

	return markdown
}


// HandleSubChapter calls ToMarkdown and saves a subchapter to a file.
// The file is saved in the root directory with the filename of the subchapter.
//
// Parameters:
// - sc: the subchapter to save
// - rootDirStr: the root directory where the file will be saved
//
// Returns:
// - the filename of the saved file
func HandleSubChapter(sc chapter, rootDirStr string) string {
	// create the filename for the file
	filename := fmt.Sprintf("%s.md", Filename(sc.Name()))

	// create the path to the file
	pathToFile := filepath.Join(rootDirStr, filename)

	// print the path to the file
	fmt.Println("[-] pathToFile: ", pathToFile)

	// save the subchapter to the file
	filename = ToMarkdown(sc, pathToFile)

	// print the filename of the saved file
	fmt.Println("Inside HandleChapter!: ", filename)

	return filename
}


func ToMarkdown(c chapter, filename string) string {
	if len(filename) == 0 {
		filename = fmt.Sprintf("%s.md", Filename(c.Name()))
	}

	markdown := ToMarkdownString(c)

	// write to file
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	_, err2 := f.WriteString(markdown)
	if err2 != nil {
		log.Fatal(err2)
	}
	f.Close()

	return filename
}

func ToHtmlString(c chapter) string {
	html := ""

	// chapter content
	if c.config.Include {
		html += fmt.Sprintf("<h1>%s</h1>", c.Name())
		html += c.Content()
	}

	// subchapters content
	for _, sc := range c.SubChapters() {
		html += ToHtmlString(sc)
	}

	return html
}

func ToHtml(c chapter, filename string) string {
	if len(filename) == 0 {
		filename = fmt.Sprintf("%s.html", Filename(c.Name()))
	}

	html := fmt.Sprintf("<html><head></head><body>%s</body></html>", ToHtmlString(c))

	// write to file
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	_, err2 := f.WriteString(html)
	if err2 != nil {
		log.Fatal(err2)
	}
	f.Close()

	return filename
}

func ToEpub(c chapter, filename string) string {
	if len(filename) == 0 {
		filename = fmt.Sprintf("%s.epub", Filename(c.Name()))
	}

	// init ebook
	e := epub.NewEpub(c.Name())
	e.SetAuthor(c.Author())

	AppendToEpub(e, c)

	err := e.Write(filename)
	if err != nil {
		log.Fatal(err)
	}

	return filename
}

func AppendToEpub(e *epub.Epub, c chapter) {
	content := ""

	// append table of content
	if len(c.SubChapters()) > 1 {
		html := "<h1>Table of Contents</h1>"

		html += "<ol>"
		for _, sc := range c.SubChapters() {
			html += fmt.Sprintf("<li>%s</li>", sc.Name())
		}
		html += "</ol>"

		_, err := e.AddSection(html, "Table of Contents", "", "")
		if err != nil {
			log.Fatal(err)
		}
	}

	// chapter content
	if c.config.Include {

		if c.config.ImagesOnly == false {
			content = c.Content()
		}

		// parse content
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(c.Content()))
		if err != nil {
			log.Fatal(err)
		}

		// download images and replace src in img tags of content
		doc.Find("img").Each(func(i int, s *goquery.Selection) {
			src, _ := s.Attr("src")
			src = strings.Split(src, "?")[0] // remove query part
			imagePath, _ := e.AddImage(src, "")

			if c.config.ImagesOnly {
				imageTag, _ := goquery.OuterHtml(s)
				content += strings.Replace(imageTag, src, imagePath, 1)
			} else {
				content = strings.Replace(content, src, imagePath, 1)
			}
		})

		html := ""
		// add title only if ImagesOnly = false
		if c.config.ImagesOnly == false {
			html += fmt.Sprintf("<h1>%s</h1>", c.Name())
		}
		html += content

		//  write to epub file
		_, err = e.AddSection(html, c.Name(), "", "")
		if err != nil {
			log.Fatal(err)
		}

	}

	// subchapters content
	for _, sc := range c.SubChapters() {
		AppendToEpub(e, sc)
	}
}

func ToMobi(c chapter, filename string) string {
	if len(filename) == 0 {
		filename = fmt.Sprintf("%s.mobi", Filename(c.Name()))
	} else {

		// add .mobi extension if not specified
		if strings.HasSuffix(filename, ".mobi") == false {
			filename = fmt.Sprintf("%s.mobi", filename)
		}

	}

	filenameEPUB := strings.ReplaceAll(filename, ".mobi", ".epub")
	ToEpub(c, filenameEPUB)

	exec.Command("kindlegen", filenameEPUB).Run()
	// exec command always returns status 1 even if it succeed
	// if err != nil {
	// 	log.Fatal(err)
	// }

	err := os.Remove(filenameEPUB)
	if err != nil {
		log.Fatal(err)
	}

	return filename
}
