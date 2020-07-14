package main

import (
	"os"
	"strings" 
	"bufio"
	"log"
	"net/http"
    "html/template"
)


type Main struct {
	Key 	string
	Value 	string
}

type Context struct {
	Env 		[]Main
	Distro 		[]Main
	Hostname    string
	Heder       http.Header
}


func getEnv() []Main {
	var m []Main
	for _, element := range os.Environ() {
        variable := strings.Split(element, "=")
        key := variable[0]
        val := variable[1] 
        if len(key) > 0 && len(val) > 0 {
            m = append(m, Main {Key: variable[0], Value: variable[1] })
        }
	}
	return m
}

func distroInfo() []Main {
	file, _ := os.Open("/etc/os-release")
    defer file.Close()

	scanner := bufio.NewScanner(file)
	var info []Main
    for scanner.Scan() {
		variable := strings.Split(scanner.Text(), "=")
        key := variable[0]
		val := variable[1] 
		if len(key) > 0 && len(val) > 0 {
            info = append(info, Main {Key: key, Value: val})
        }
    }
    return info
}

func getHostname() string {
	name, _ := os.Hostname()
	return name
}

func createHatml() string {
	val := `
	<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <link
      rel="stylesheet"
      href="https://stackpath.bootstrapcdn.com/bootstrap/4.5.0/css/bootstrap.min.css"
    />
    <title>{{ .Hostname }}</title>
  </head>

  <body class="bg-dark text-white">
    <div class="container-fluid">
      <h1 class="mx-auto text-center p-3 m-3">Hostname: {{ .Hostname }}</h1>
      <div class="row">
        <div class="col-md-6">
          <h2 class="m-3 p-3 mx-auto text-center">OS info</h2>
          <table class="table table-dark">
            <thead>
              <tr>
                <th scope="col">Name</th>
                <th scope="col">Value</th>
              </tr>
            </thead>
            <tbody>
              {{ range .Distro }}
              <tr>
                <td>{{ .Key }}</td>
                <td>{{ .Value }}</td>
              </tr>
              {{ end }}
            </tbody>
          </table>
        </div>
        <div class="col-md-6">
          <h2 class="m-3 p-3 mx-auto text-center">Headers</h2>
          <table class="table table-dark">
            <thead>
              <tr>
                <th scope="col">Name</th>
                <th scope="col">Value</th>
              </tr>
            </thead>
            <tbody>
              {{ range $key, $value := .Heder }}
              <tr>
                <td>{{ $key }}</td>
                <td>{{ $value }}</td>
              </tr>
              {{ end }}
            </tbody>
          </table>
        </div>
        <div class="col-md-12">
          <h2 class="m-3 p-3 mx-auto text-center">Enviroments</h2>
          <table class="table table-dark">
            <thead>
              <tr>
                <th scope="col">Name</th>
                <th scope="col">Value</th>
              </tr>
            </thead>
            <tbody>
              {{ range .Env }}
              <tr>
                <td>{{ .Key }}</td>
                <td>{{ .Value }}</td>
              </tr>
              {{ end }}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </body>
</html>
	`
	return val
}


func handler(w http.ResponseWriter, r *http.Request) { 
	values := r.URL.Query()
	for key, val := range values {
		if len(key) != 0 && len(val) != 0  {
            os.Setenv(key, val[0])
            http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	}
	context := Context{ Env: getEnv(), Distro: distroInfo(), Hostname: getHostname(), Heder: r.Header}
	tmpl, _ := template.New("test").Parse(createHatml())
	tmpl.Execute(w, context)
}


func main()  {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":80", nil))
}

