package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/niklasfasching/go-org/org"
)

// Post represents a blog post with metadata and content.
type Post struct {
	Title       string        // Title of the post
	Slug        string        // URL-friendly identifier for the post
	Tags        []string      // Tags associated with the post
	Description string        // Brief description of the post
	Date        time.Time     // Date when the post was published
	DateStr     string        // Date when the post was published (YYYY-MM-DD)
	Timestamp   int64         // Unix timestamp of the publication date
	Path        string        // File path to the original .org file
	PathHtml    string        // URL path to the rendered HTML file
	Content     *org.Document // Parsed content of the post
	ReadTime    uint8         // Estimated reading time in minutes
	Hero        string        // URL path to the hero image for the post
}

// listPosts reads the posts directory and returns a slice of Post structs.
// It filters out non-.org files, parses each .org file, and sorts the posts by date in descending order.
// Returns the slice of posts and any error encountered during the process.
func listPosts() ([]Post, error) {
	entries, err := os.ReadDir("posts")
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return nil, err
	}

	entries = filter(entries, func(e os.DirEntry) bool { return filepath.Ext(e.Name()) == ".org" })

	var posts []Post
	for _, entry := range entries {
		filePath := filepath.Join("posts", entry.Name())

		post, err := parseOrg(filePath)
		if err != nil {
			log.Println("[!] Unable to parse ", filePath)
		} else {
			posts = append(posts, post)
		}
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Timestamp > posts[j].Timestamp
	})

	return posts, nil
}

// handleImages processes image and video links in the org document.
// It updates the URL of the link to point to the media directory.
// Parameters:
//   - protocol: The protocol of the link (e.g., "file", "http").
//   - description: The description of the link.
//   - link: The URL of the link.
//
// Returns:
//   - The processed link node.
func handleImages(protocol string, description []org.Node, link string) org.Node {
	linked := org.RegularLink{protocol, description, link, false}
	if linked.Kind() == "image" || linked.Kind() == "video" {
		linked.URL = path.Join("/medias/", linked.URL)
	}
	return linked
}

// parseOrg parses an org file and returns a Post struct.
// It reads the file, extracts metadata, and calculates the reading time.
// Parameters:
//   - filePath: The path to the org file.
//
// Returns:
//   - The parsed Post struct.
//   - Any error encountered during the process.
func parseOrg(filePath string) (Post, error) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Error reading file")
		return Post{}, err
	}

	config := org.New()
	config.ResolveLink = handleImages

	orgData := config.Parse(file, filePath)

	title := orgData.Get("TITLE")
	description := orgData.Get("DESCRIPTION")
	dateStr := strings.Split(orgData.Get("DATE"), "T")[0]
	slug := orgData.Get("SLUG")
	tags := strings.Split(orgData.Get("TAGS"), ", ")
	hero := path.Join("/medias", orgData.Get("HERO"))

	date, _ := time.Parse("2006-01-02", dateStr)
	ts := date.Unix()

	raw, _ := os.ReadFile(filePath)
	readTime := len(strings.Split(string(raw), " ")) / 200

	return Post{
		Title:       title,
		Slug:        slug,
		Tags:        tags,
		Description: description,
		Date:        date,
		DateStr:     date.Format("2006-01-02"),
		Timestamp:   ts,
		Path:        filePath,
		PathHtml:    "/posts/" + slug + ".html",
		Content:     orgData,
		ReadTime:    uint8(readTime),
		Hero:        hero,
	}, nil
}
