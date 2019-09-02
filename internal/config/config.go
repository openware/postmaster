package config

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"strings"

	"github.com/openware/postmaster/pkg/eventapi"
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
	Name       string              `yaml:"name"`
	Key        string              `yaml:"key"`
	Exchange   string              `yaml:"exchange"`
	Templates  map[string]Template `yaml:"templates"`
	Expression string              `yaml:"expression"`
}

// General application configuration.
type Config struct {
	Languages []Language                    `yaml:"languages"`
	Keychain  map[string]eventapi.Validator `yaml:"keychain"`
	Exchanges map[string]string             `yaml:"exchanges"`
	Events    []Event                       `yaml:"events"`
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

func (config *Config) ContainsExchange(id string) bool {
	_, ok := config.Exchanges[id]
	return ok
}

func (config *Config) ContainsKey(id string) bool {
	_, ok := config.Keychain[id]
	return ok
}

func (lang *Language) Valid() bool {
	notEmpty := len(strings.TrimSpace(lang.Code)) != 0
	isUp := lang.Code == strings.ToUpper(lang.Code)

	return notEmpty && isUp
}

func (config *Config) validateLanguages() (bool, error) {
	for _, lang := range config.Languages {
		if !lang.Valid() {
			return false, fmt.Errorf("language %s should be uppercased", lang.Code)
		}
	}

	return true, nil
}

func (config *Config) ValidateExchanges() error {
	if len(config.Exchanges) < 1 {
		return errors.New("no exchanges was specified")
	}

	for k, v := range config.Exchanges {
		if v == "" {
			return fmt.Errorf("exchange %s can not have empty value", k)
		}

	}

	return nil
}

func (config *Config) ValidateKeychain() error {
	for id := range config.Exchanges {
		if !config.ContainsKey(id) {
			return fmt.Errorf("exchange %s doesn't have a key", id)
		}
	}

	for k, v := range config.Keychain {
		if v.Value == "" {
			return fmt.Errorf("key for %s has an empty value", k)
		}

		if v.Algorithm == "" {
			return fmt.Errorf("key for %s has an empty algorithm", k)
		}
	}

	return nil
}

// Validate configuration file.
func (config *Config) Validate() error {
	if _, err := config.validateLanguages(); err != nil {
		return err
	}

	if err := config.ValidateExchanges(); err != nil {
		return err
	}

	if err := config.ValidateKeychain(); err != nil {
		return err
	}

	for _, event := range config.Events {
		for lang, tpl := range event.Templates {
			strippedTpl := strings.TrimSpace(tpl.Template)
			strippedTplPath := strings.TrimSpace(tpl.TemplatePath)

			if strippedTpl != "" && strippedTplPath != "" {
				return errors.New("template and template path is specified")
			}

			if lang != strings.ToUpper(lang) {
				err := fmt.Errorf("language %s in event %s should be uppercased", lang, event.Name)
				return err
			}
		}

		// Check, that at least one language exist.
		if len(event.Templates) == 0 {
			return fmt.Errorf("templates can not be empty")
		}

		// Check, that exchange was declared.
		if !config.ContainsExchange(event.Exchange) {
			err := fmt.Errorf("exchange %s in event %s is not defined", event.Exchange, event.Name)
			return err
		}
	}

	return nil
}
