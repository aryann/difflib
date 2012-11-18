

package main

import (
	"fmt"
	"html"
	"html/template"
	"strings"
	"os"
	"io/ioutil"
	"difflib"
	"net/http"
)

var template1 = `
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

    .added, .deleted {
      background-color: black;
      width: 100px;
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
	http.ListenAndServe(":8080", nil)
}

func diffHandler(filename1, filename2 string) http.HandlerFunc {
	diff := difflib.HTMLDiff(fileToLines(filename1), fileToLines(filename2))
	fmt.Println(diff)
	tmpl, _ := template.New("diffTemplate").Parse(template1)
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, map[string]interface{}{
			"Diff": template.HTML(diff),
			"Filename1": filename1,
			"Filename2": filename2,
		})
	}
}

func fileToLines(filename string) []string {
	contents, _ := ioutil.ReadFile(filename)
	return strings.Split(html.EscapeString(string(contents)), "\n")
}