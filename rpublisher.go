package main

import (
	"os"
	"net/http"
	"strings"
	"github.com/jackc/pgx"
	"context"
	"os/exec"
	"html/template"
)

func main() {
	port := os.Getenv("PORT")
	rmarkdownPath := "lol.Rmd" //os.Args[1]
	pgHost := os.Getenv("PG_HOST")
	pgPort := os.Getenv("PG_PORT")
	pgDbname := os.Getenv("PG_DBNAME")
	pgUser := os.Getenv("PG_USER")
	pgPassword := os.Getenv("PG_PASSWORD")
	pgChannel := os.Getenv("PG_CHANNEL")
	connStr :=
		" user=" + pgUser +
		" password=" + pgPassword +
		" dbname=" + pgDbname +
		" sslmode=disable" +
		" port=" + pgPort +
		" host=" + pgHost
	listener, err := NewPgListener(connStr, pgChannel)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	notifications := make(chan *pgx.Notification)
	listener.NotifyOnNotification(ctx, notifications)
	go updateHtmlOnNotification(notifications, rmarkdownPath)

	htmlPath := strings.Replace(rmarkdownPath, ".Rmd", ".html", 1)
	fn, err := deliverRequestFactory(htmlPath)
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/", fn)
	http.ListenAndServe("0.0.0.0:"+port, nil)
}

func deliverRequestFactory(htmlPath string) (func(http.ResponseWriter, *http.Request), error) {
	htmlName := strings.Trim(htmlPath, ".html")
	tmpl, err := template.ParseFiles(htmlPath)
	if err != nil {
		return nil, err
	}
	return func(writer http.ResponseWriter, request *http.Request) {
		tmpl.ExecuteTemplate(writer, htmlName, nil)
	}, nil
}


func updateHtmlOnNotification(notifications <-chan *pgx.Notification, rMarkDownFilePath string) error {
	for {
		<-notifications
		updateHtml(rMarkDownFilePath)
	}
}

func updateHtml(rMarkDownFilePath string) error {
	command := exec.Command("Rscript", "-e", "'rmarkdown::render(\"" + rMarkDownFilePath + "\")' ")
	err := command.Start()
	return err
}