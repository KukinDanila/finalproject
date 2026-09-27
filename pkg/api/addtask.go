package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"finalproject/pkg/db"
)

func checkDate(task *db.Task) error {
	now := time.Now()

	if task.Date == "" {
		task.Date = now.Format("20060102")
	}

	t, err := time.Parse("20060102", task.Date)
	if err != nil {
		return err
	}

	if task.Repeat != "" {
		next, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}

		if afterNow(now, t) {
			task.Date = next
		}
	} else {
		if afterNow(now, t) {
			task.Date = now.Format("20060102")
		}
	}

	return nil
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	if task.Title == "" {
		writeJson(w, map[string]string{"error": "не указан заголовок задачи"})
		return
	}

	if err = checkDate(&task); err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	writeJson(w, map[string]string{"id": strconv.FormatInt(id, 10)})
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	if task.Title == "" {
		writeJson(w, map[string]string{"error": "не указан заголовок задачи"})
		return
	}

	if err = checkDate(&task); err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	if err = db.UpdateTask(&task); err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	writeJson(w, map[string]string{})
}
