package main

import (
	m "github.com/vault-thirteen/SWVD/models"
)

type Task struct {
	s  *m.TaskSettings
	um *m.UrlMaker
}

func NewTask(settings *m.TaskSettings) (t *Task) {
	t = &Task{
		s:  settings,
		um: m.NewUrlMaker(settings.Schema, settings.Domain),
	}

	return t
}
