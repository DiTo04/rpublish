package main

import (
	"testing"
	"os"
	//"github.com/jackc/pgx/pgmock"
)

func TestDeliverRequestFactory(t *testing.T) {
	htmlPath := "test"
	_ , err := deliverRequestFactory(htmlPath)
	if err == nil {
		t.Error("Should return nil with faulty page!")
	}
	htmlPath = "test_page.html"
	if _, err := os.Stat(htmlPath); os.IsNotExist(err) {
		t.Error("Could not locate " + htmlPath)
	}
	_, err = deliverRequestFactory(htmlPath)
	if err != nil {
		t.Error("Should not generate error with " + htmlPath)
	}
}

func TestChannelNotifications(t *testing.T) {
	//steps := pgmock.AcceptUnauthenticatedConnRequestSteps()
	//script := &pgmock.Script{Steps: steps}
	//server, _ := pgmock.NewServer(script) TODO(David): Implement the rest of all tests.

}