package main

import (
	"log"
	"net/http"
	"testing"

	"gotest.tools/assert"
)

func Test_handlePublishTask(t *testing.T) {
	req, err := http.NewRequest("POST", "http://127.0.0.1:8080/publishTask", nil)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header["X-Appengine-Cron"] = []string{"Valid Cron Header"}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	assert.Equal(t, resp.StatusCode, http.StatusOK)
}

func Test_handleUploadImage(t *testing.T) {
	req, err := http.NewRequest("POST", "http://127.0.0.1:8080/uploadImage", nil)
	req.Header["X-Appengine-Cron"] = []string{"Valid Cron Header"}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	assert.Equal(t, resp.StatusCode, http.StatusOK)
}
