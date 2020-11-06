package baidulogin

import (
	"log"
	"net/http"
	"text/template"
)

func indexPage(w http.ResponseWriter, r *http.Request) {
	sess, _ := globalSessions.SessionStart(w, r) // session start
	registerBaiduClient(&sess)                   // 如果没有 baiduClient , 就添加

	// get file contents as string
	contents, err := templateFilesBox.String("index.html")
	if err != nil {
		log.Println(err)
		return
	}
	tmpl, err := template.New("index.html").Parse(contents)
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(w, Version)
	if err != nil {
		panic(err)
	}
}

func favicon(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Location", "//www.baidu.com/favicon.ico")
	http.Error(w, "", 302)
}
