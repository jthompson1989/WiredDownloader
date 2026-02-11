# Wired Article Downloader

A Go program that downloads Wired articles and saves them as text files in your Documents/Wired folder.

## Usage

```bash
go run main.go <wired-article-url>
```

Example:
```bash
go run main.go https://www.wired.com/story/example-article/
```

The program will:
1. Fetch the Wired article from the provided URL
2. Extract the title and content from the HTML
3. Create a "Wired" folder in your Documents directory (if it doesn't exist)
4. Save the article as a text file named after the article title

## Building

To build an executable:
```bash
go build -o wired-downloader.exe
```

Then run:
```bash
wired-downloader.exe <wired-article-url>
```