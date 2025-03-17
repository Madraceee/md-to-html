# Markdown to HTML generator
This project can convert markdown to HTML. Given a markdown file, it generates a HTML file.
This does not fully implement the grammar of markdown rather only a subset and cannot be used for production. Error handling is minimal.

#### Note
The project is a playground for me to play around with scanners, parsers, etc. More features will be added slowly.
If you have any recommendation, please raise an issue.

### How to run
Requires [Go](https://go.dev/).
```
make build
```
This generates a binary for your system.

To run the binary
```
./md-to-html [FILE]
```
