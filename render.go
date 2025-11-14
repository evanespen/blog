package main

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/niklasfasching/go-org/org"
)

// renderHome renders the home page of the website.
// It processes the index template, executes it with the provided posts and tags,
// and writes the resulting HTML to the build directory.
// Parameters:
//   - posts: A slice of Post structs representing the blog posts.
//   - tags: A slice of strings representing the tags.
//   - css: A string containing the compiled CSS styles.
//
// Returns:
//   - An error if any step of the process fails, otherwise nil.
func renderHome(posts []Post, tags []string, css string) error {
	indexTmpl, _ := template.ParseFiles("templates/parts/index.html")
	var indexContentBuf strings.Builder
	indexData := struct {
		Posts []Post
	}{
		Posts: posts,
	}
	_ = indexTmpl.Execute(&indexContentBuf, indexData)

	// Parse the index.html template
	tmpl, err := template.ParseFiles("templates/layout.html", "templates/parts/header.html")
	if err != nil {
		return fmt.Errorf("error parsing template: %v", err)
	}

	// Create a buffer to hold the template output
	var buf strings.Builder

	// Execute the template with the necessary data
	data := struct {
		Css         template.CSS
		Content     template.HTML
		Hero        template.HTML
		Tags        []string
		ShowSidebar bool
	}{
		Css:         template.CSS(css),
		Content:     template.HTML(indexContentBuf.String()),
		Hero:        template.HTML("<div id=\"hero\"></div>"),
		Tags:        tags,
		ShowSidebar: true,
	}

	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("error executing template: %v", err)
	}

	// Create the build directory if it doesn't exist
	if err := os.MkdirAll("build", os.ModePerm); err != nil {
		return fmt.Errorf("error creating directory: %v", err)
	}

	// Write the HTML content to the index.html file in the build directory
	if err := os.WriteFile("build/index.html", []byte(buf.String()), 0644); err != nil {
		return fmt.Errorf("error writing HTML file: %v", err)
	}

	log.Println("Wrote build/index.html")
	return nil
}

// highlightCodeBlock highlights a code block using the specified language and parameters.
// It uses the chroma library to tokenize and format the code block.
// Parameters:
//   - source: The source code to highlight.
//   - lang: The programming language of the code.
//   - inline: Whether the code block is inline or not.
//   - params: Additional parameters for highlighting, such as highlighted lines.
//
// Returns:
//   - A string containing the highlighted code block in HTML format.
func highlightCodeBlock(source, lang string, inline bool, params map[string]string) string {
	var w strings.Builder
	l := lexers.Get(lang)
	if l == nil {
		l = lexers.Fallback
	}
	l = chroma.Coalesce(l)
	it, _ := l.Tokenise(nil, source)
	options := []html.Option{}
	if params[":hl_lines"] != "" {
		ranges := org.ParseRanges(params[":hl_lines"])
		if ranges != nil {
			options = append(options, html.HighlightLines(ranges))
		}
	}
	_ = html.New(options...).Format(&w, styles.Get("dracula"), it)
	if inline {
		return `<div class="highlight-inline">` + "\n" + w.String() + "\n" + `</div>`
	}
	return `<div class="highlight">` + "\n" + w.String() + "\n" + `</div>`
}

// renderPost renders a single blog post to an HTML file.
// It processes the post content, applies syntax highlighting to code blocks,
// and writes the resulting HTML to the build directory.
// Parameters:
//   - post: The Post struct representing the blog post.
//   - css: A string containing the compiled CSS styles.
//   - tags: A slice of strings representing the tags.
//
// Returns:
//   - An error if any step of the process fails, otherwise nil.
func renderPost(post Post, css string, tags []string) error {
	htmlFilePath := "build/posts/" + post.Slug + ".html"
	render := func(w org.Writer) string {
		out, err := post.Content.Write(w)
		if err != nil {
			log.Fatal(err)
		}
		return out
	}

	renderer := org.NewHTMLWriter()
	renderer.HighlightCodeBlock = highlightCodeBlock
	htmlContent := render(renderer)

	if err := os.MkdirAll("build/posts", os.ModePerm); err != nil {
		log.Fatal("Error creating directory:", err)
		return err
	}

	// Generate the new file path for the HTML output
	tmpl, err := template.ParseFiles("templates/layout.html", "templates/parts/header.html")
	if err != nil {
		log.Fatal("Error parsing template:", err)
		return err
	}

	hero := func(post Post) template.HTML {
		if post.Hero != "/medias/" && post.Hero != "/medias/none" {
			return template.HTML(fmt.Sprintf("<img id=\"hero\" src=\"%s\"/>", post.Hero))
		} else {
			return template.HTML("")
		}
	}

	// Create a buffer to hold the template output
	var buf strings.Builder

	// Execute the template with the necessary data
	data := struct {
		Content     template.HTML
		Css         template.CSS
		Hero        template.HTML
		Tags        []string
		ShowSidebar bool
	}{
		Content:     template.HTML(htmlContent),
		Css:         template.CSS(css),
		Hero:        hero(post),
		Tags:        tags,
		ShowSidebar: false,
	}

	if err := tmpl.Execute(&buf, data); err != nil {
		log.Fatal("Error executing template:", err)
		return err
	}

	// Write the HTML content to the new file
	if err := os.WriteFile(htmlFilePath, []byte(buf.String()), 0644); err != nil {
		log.Fatal("Error writing HTML file:", err)
		return err
	}

	log.Println("Wrote", htmlFilePath)
	return nil
}

// renderTagPage renders a tag page for a specific tag.
// It processes the tag page template, executes it with the provided tag and posts,
// and writes the resulting HTML to the build directory.
// Parameters:
//   - tag: The tag for which the page is being rendered.
//   - posts: A slice of Post structs representing the blog posts associated with the tag.
//   - tags: A slice of strings representing all tags.
//   - css: A string containing the compiled CSS styles.
//
// Returns:
//   - An error if any step of the process fails, otherwise nil.
func renderTagPage(tag string, posts []Post, tags []string, css string) error {
	htmlFilePath := "build/tags/" + tag + ".html"

	if err := os.MkdirAll("build/tags", os.ModePerm); err != nil {
		log.Fatal("Error creating directory:", err)
		return err
	}

	tagPageTmpl, _ := template.ParseFiles("templates/parts/tagPage.html")
	var tagPageContentBuf strings.Builder
	tagPageData := struct {
		Tag   string
		Posts []Post
	}{
		Tag:   tag,
		Posts: posts,
	}
	_ = tagPageTmpl.Execute(&tagPageContentBuf, tagPageData)

	// Generate the new file path for the HTML output
	tmpl, err := template.ParseFiles("templates/layout.html", "templates/parts/header.html")
	if err != nil {
		log.Fatal("Error parsing template:", err)
		return err
	}

	// Create a buffer to hold the template output
	var buf strings.Builder

	// Execute the template with the necessary data
	data := struct {
		Content     template.HTML
		Css         template.CSS
		Hero        template.HTML
		Tags        []string
		ShowSidebar bool
		Tag         string
	}{
		Content:     template.HTML(tagPageContentBuf.String()),
		Css:         template.CSS(css),
		Hero:        template.HTML(""),
		Tags:        tags,
		ShowSidebar: false,
		Tag:         tag,
	}

	if err := tmpl.Execute(&buf, data); err != nil {
		log.Fatal("Error executing template:", err)
		return err
	}

	// Write the HTML content to the new file
	if err := os.WriteFile(htmlFilePath, []byte(buf.String()), 0644); err != nil {
		log.Fatal("Error writing HTML file:", err)
		return err
	}

	log.Println("Wrote", htmlFilePath)
	return nil
}
