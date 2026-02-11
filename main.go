package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: wired-downloader <wired-article-url>")
		os.Exit(1)
	}

	url := os.Args[1]

	if !strings.Contains(url, "wired.com") {
		log.Fatal("Please provide a valid Wired.com article URL")
	}

	article, err := fetchArticle(url)
	if err != nil {
		log.Fatalf("Error fetching article: %v", err)
	}

	wiredPath, err := getWiredFolderPath()
	if err != nil {
		log.Fatalf("Error creating Wired folder: %v", err)
	}

	filename := sanitizeFilename(article.Title) + ".txt"
	filepath := filepath.Join(wiredPath, filename)

	err = saveArticleToFile(article, filepath)
	if err != nil {
		log.Fatalf("Error saving article: %v", err)
	}

	fmt.Printf("Article saved to: %s\n", filepath)
}

type Article struct {
	Title   string
	Content string
}

func fetchArticle(url string) (*Article, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error: %s", resp.Status)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	title := extractTitle(url)
	content := extractContent(doc)

	return &Article{
		Title:   title,
		Content: content,
	}, nil
}

func extractTitle(url string) string {
	var title string

	var splitUrl = strings.Split(url, "/")
	title = splitUrl[len(strings.Split(url, "/"))-1]

	return title
}

func extractContent(n *html.Node) string {
	var content strings.Builder

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "p", "h1", "h2", "h3", "h4", "h5", "h6":
				text := getTextContent(n)
				if strings.TrimSpace(text) != "" {
					content.WriteString(text)
					content.WriteString("\n\n")
				}
			case "div":
				if class := getAttr(n, "class"); class != "" {
					if strings.Contains(class, "article-body") ||
						strings.Contains(class, "content") ||
						strings.Contains(class, "post-body") {
						for c := n.FirstChild; c != nil; c = c.NextSibling {
							f(c)
						}
					}
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(n)

	return strings.TrimSpace(content.String())
}

func getTextContent(n *html.Node) string {
	var text strings.Builder

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.TextNode {
			text.WriteString(n.Data)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(n)

	return strings.TrimSpace(text.String())
}

func getAttr(n *html.Node, attrName string) string {
	for _, attr := range n.Attr {
		if attr.Key == attrName {
			return attr.Val
		}
	}
	return ""
}

func getWiredFolderPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	documentsPath := filepath.Join(home, "Documents")
	wiredPath := filepath.Join(documentsPath, "Wired")

	err = os.MkdirAll(wiredPath, 0755)
	if err != nil {
		return "", fmt.Errorf("failed to create Wired directory: %w", err)
	}

	return wiredPath, nil
}

func sanitizeFilename(filename string) string {
	re := regexp.MustCompile(`[<>:"/\\|?*]`)
	filename = re.ReplaceAllString(filename, "")

	re = regexp.MustCompile(`[\s]+`)
	filename = re.ReplaceAllString(filename, "_")

	if len(filename) > 100 {
		filename = filename[:100]
	}

	return strings.TrimSpace(filename)
}

func saveArticleToFile(article *Article, filepath string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	_, err = io.WriteString(file, article.Title+"\n\n")
	if err != nil {
		return fmt.Errorf("failed to write title: %w", err)
	}

	_, err = io.WriteString(file, article.Content)
	if err != nil {
		return fmt.Errorf("failed to write content: %w", err)
	}

	return nil
}
