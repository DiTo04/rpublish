package main

import (
	"github.com/jackc/pgx"
	"context"
)

type PgListener struct {
	conn *pgx.Conn
}

func NewPgListener(connStr string, pgChannel string) (*PgListener, error){
	conn, err := connectToDb(connStr)
	if err != nil {
		return nil, err
	}
	err = conn.Listen(pgChannel)
	if err != nil {
		return nil, err
	}
	listener := &PgListener{}
	listener.conn = conn
	return listener, nil
}

func (p *PgListener) NotifyOnNotification(ctx context.Context,channel chan<- *pgx.Notification) {
	go channelNotifications(p.conn, ctx, channel)
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

func connectToDb(connStr string) (*pgx.Conn, error) {
	config, err := pgx.ParseConnectionString(connStr)
	if err != nil {
		return nil, err
	}
	return pgx.Connect(config)
}