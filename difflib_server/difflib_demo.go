// A demo for difflib. This program accepts the paths to two files and
// launches a web server at port 8080 that serves the diff results.
package main

import (
	"fmt"
	"github.com/aryann/difflib"
	"html"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var hostPort = "localhost:8080"

var templateString = `
<!doctype html>
<html>
<head>
  <meta charset="utf-8" />
  <title>difflib results</title>
  <style type="text/css">
    table {
      background-color: lightgrey;
      border-spacing: 1px;
    }

    tr {
      background-color: white;
      border-bottom: 1px solid black;
    }

    .line-num {
      width: 50px;
    }

    .added {
      background-color: green;
    }

    .deleted {
      background-color: red;
    }
  </style>
</head>
<body>
  <table>
    <tr>
      <th></th>
      <th><em>{{.Filename1}}</em></th>
      <th><em>{{.Filename2}}</em></th>
      <th></th>
    </tr>
    {{.Diff}}
  </table>
</body>
</html>
`

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "USAGE: %s <file-1> <file-2>\n", os.Args[0])
		os.Exit(1)
	}
	http.HandleFunc("/", diffHandler(os.Args[1], os.Args[2]))
	fmt.Printf("Starting server at %s.\n", hostPort)
	err := http.ListenAndServe(hostPort, nil)
	if err != nil {
		panic(err)
	}
}

// diffHandler returns an http.HandlerFunc that serves the diff of the
// two given files.
func diffHandler(filename1, filename2 string) http.HandlerFunc {
	diff := difflib.HTMLDiff(fileToLines(filename1), fileToLines(filename2))
	tmpl, _ := template.New("diffTemplate").Parse(templateString)
	return func(w http.ResponseWriter, r *http.Request) {
		err := tmpl.Execute(w, map[string]interface{}{
			"Diff":      template.HTML(diff),
			"Filename1": filename1,
			"Filename2": filename2,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// filesToLines takes a path to a file and returns a string array of
// the lines in the file. Any HTML in the file is escaped.
func fileToLines(filename string) []string {
	contents, _ := ioutil.ReadFile(filename)
	return strings.Split(html.EscapeString(string(contents)), "\n")
}
