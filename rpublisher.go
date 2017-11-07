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
	conn, err := connectToDb(connStr)
	if err != nil {
		panic(err)
	}
	err = conn.Listen(pgChannel)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	go updateHtmlOnNotification(conn, ctx, rmarkdownPath)

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

func connectToDb(connStr string) (*pgx.Conn, error) {
	config, err := pgx.ParseConnectionString(connStr)
	if err != nil {
		return nil, err
	}
	return pgx.Connect(config)
}

func updateHtmlOnNotification(conn *pgx.Conn, ctx context.Context, rMarkDownFilePath string) error {
	notifications := make(chan *pgx.Notification)
	go channelNotifications(conn, ctx, notifications)
	for {
		select {
		case <-ctx.Done():
			return nil
		case <- notifications:
			updateHtml(rMarkDownFilePath)
		}
	}
}

func updateHtml(rMarkDownFilePath string) error {
	command := exec.Command("Rscript", "-e", "'rmarkdown::render(\"" + rMarkDownFilePath + "\")' ")
	err := command.Start()
	return err
}

func channelNotifications(conn *pgx.Conn, ctx context.Context, notifications chan<- *pgx.Notification) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			if notification, err := conn.WaitForNotification(ctx); err != nil {
				return err
			} else if notification != nil {
				 notifications <- notification
			}
		}
	}
}
