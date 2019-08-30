package config

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"strings"
)

type Language struct {
	Code string `yaml:"code"`
	Name string `yaml:"name"`
}

type Template struct {
	Subject      string `yaml:"subject"`
	TemplatePath string `yaml:"template_path,omitempty"`
	Template     string `yaml:"template,omitempty"`
}

type Event struct {
	Name      string              `yaml:"name"`
	Key       string              `yaml:"key"`
	Templates map[string]Template `yaml:"templates"`
}

type Exchange struct {
	ID   string `yaml:"id"`
	Name string `yaml:"name"`
}

// General application configuration.
type Config struct {
	Languages []Language       `yaml:"languages"`
	Exchanges []Exchange       `yaml:"exchanges"`
	Events    map[string]Event `yaml:"events"`
}

func (e *Event) Template(key string) Template {
	return e.Templates[strings.ToUpper(key)]
}

func (t *Template) Content(data interface{}) ([]byte, error) {
	var err error

	buff := new(bytes.Buffer)
	tpl := new(template.Template)

	if strings.TrimSpace(t.Template) != "" {
		tpl, err = template.New(t.Subject).Parse(t.Template)
	} else {
		tpl, err = template.ParseFiles(t.TemplatePath)
	}

	if err != nil {
		return nil, err
	}

	if err := tpl.Execute(buff, &data); err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}

func (config *Config) ContainsLanguage(code string) bool {
	for _, lang := range config.Languages {
		if strings.EqualFold(lang.Code, code) {
			return true
		}
	}

	return false
}

func (lang *Language) Valid() bool {
	notEmpty := len(strings.TrimSpace(lang.Code)) != 0
	isUp := lang.Code == strings.ToUpper(lang.Code)

	return notEmpty && isUp
}

func (config *Config) validateLanguages() (bool, error) {
	for _, lang := range config.Languages {
		if !lang.Valid() {
			return false, fmt.Errorf("language \"%s\" should be uppercased", lang.Code)
		}
	}

	return true, nil
}

// Validate configuration file.
func (config *Config) Validate() (bool, error) {
	if _, err := config.validateLanguages(); err != nil {
		return false, err
	}

	for _, event := range config.Events {
		for lang, tpl := range event.Templates {
			strippedTpl := strings.TrimSpace(tpl.Template)
			strippedTplPath := strings.TrimSpace(tpl.TemplatePath)

			if strippedTpl != "" && strippedTplPath != "" {
				return false, errors.New("template and template path is specified")
			}

			if lang != strings.ToUpper(lang) {
				err := fmt.Errorf("language \"%s\" in event \"%s\" should be uppercased", lang, event.Name)
				return false, err
			}
		}

		for _, lang := range config.Languages {
			if _, exists := event.Templates[lang.Code]; !exists {
				err := fmt.Errorf(
					"language \"%s\" in event \"%s\" is not defined", lang.Code, event.Name)
				return false, err
			}
		}
	}

	return true, nil
}
