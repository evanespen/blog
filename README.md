# Evrard Van Espen's Blog

This is a personal blog built with Go, using the go-org library to parse Org-mode files and generate HTML content.

## Features

- Parse Org-mode files from the `posts` directory
- Generate HTML pages for each post
- Create an index page with all posts
- Support for tags and tag pages
- Syntax highlighting for code blocks
- Responsive design with SCSS
- Media handling for images and videos

## Installation

1. Clone the repository:
   ```sh
   git clone https://github.com/evanespen/blog.git
   cd blog
   ```

2. Install dependencies:
   ```sh
   go mod download
   ```

## Usage

1. Add your posts in Org-mode format to the `posts` directory.
2. Add the medias of the articles alongside with them
3. Build the site in dev mode (auto reload):
   ```sh
   ./build-dev.sh
   ```
   Or with `go run .`

4. The generated site will be in the `build` directory.

## Development

- The main entry point is `main.go`.
- Templates are located in the `templates` directory.
- SCSS file is `styles/main.scss`.
- Utility functions are in `utils.go`.
- Media handling is in `medias.go`.
- Org-mode parsing is in `parse.go`.
- HTML rendering is in `render.go`.
- SCSS compilation is in `styles.go`.
- Static file copying is in `static.go`.

## TODO
- [X] code documentation
- [ ] resume page
- [ ] contact page
- [ ] responsive
- [ ] RSS
- [ ] favicon
- [ ] search
- [ ] sitemap.xml
- [ ] robots.txt
- [ ] better error handling

## License

This project is licensed under the MIT License - see the LICENSE file for details.
