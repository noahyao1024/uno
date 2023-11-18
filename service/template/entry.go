package template

import (
	"bytes"
	"fmt"
	"html/template"
	"uno/pkg/setting"
)

type Entry struct {
	Sender  string
	Subject string
	Content string
}

// Get ...
func Get(key string, data map[string]string) (*Entry, error) {
	cfg, ok := setting.AppInstance.Template.Email[key]
	if !ok {
		return nil, fmt.Errorf("template key %s not found", key)
	}

	tpl, err := template.ParseFiles(fmt.Sprintf("data/template/%s", cfg.File))
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)

	if data == nil {
		data = make(map[string]string)
	}

	data["subject"] = cfg.Subject

	err = tpl.Execute(buf, data)
	if err != nil {
		return nil, err
	}

	return &Entry{
		Sender:  cfg.Sender,
		Subject: cfg.Subject,
		Content: buf.String(),
	}, nil
}
