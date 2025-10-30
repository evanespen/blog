package main

import (
	"log"
	"slices"
)

func main() {
	// main is the entry point of the application.
	// It orchestrates the process of generating the static website by:
	// 1. Listing all posts.
	// 2. Compiling SCSS styles.
	// 3. Rendering each post.
	// 4. Rendering the home page.
	// 5. Copying static files.
	// 6. Rendering tag pages.
	// 7. Copying media files.
	posts, _ := listPosts()
	var tags []string
	postsByTag := make(map[string][]Post)

	for _, p := range posts {
		for _, t := range p.Tags {
			if !slices.Contains(tags, t) {
				tags = append(tags, t)
			}

			if postsByTag[t] == nil {
				postsByTag[t] = []Post{}
			}

			postsByTag[t] = append(postsByTag[t], p)
		}
	}

	css, _ := compileSCSS()

	log.Println(len(posts), "posts to handle")

	for _, p := range posts {
		_ = renderPost(p, css, tags)
	}

	// Process the index.html template
	if err := renderHome(posts, tags, css); err != nil {
		log.Fatal("Error processing index template:", err)
	}

	// Copy the "static" folder to the "build" folder
	if err := copyDir("static/", "build"); err != nil {
		log.Fatal("Error copying static files:", err)
	}

	for _, t := range tags {
		renderTagPage(t, postsByTag[t], tags, css)
	}

	if err := copyMedias(); err != nil {
		log.Fatal("Erro copying media files:", err)
	}
}
