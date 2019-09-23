package config

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"strings"

	"github.com/openware/postmaster/pkg/eventapi"
)

// Language represents configuration for every language registered for further usage.
type Language struct {
	Code string `yaml:"code"`
	Name string `yaml:"name"`
}

// Template represents email massage content and subject.
type Template struct {
	Subject      string `yaml:"subject"`
	TemplatePath string `yaml:"template_path,omitempty"`
	Template     string `yaml:"template,omitempty"`
}

// Event represent configuration for listening an message from RabbitMQ.
type Event struct {
	Name       string              `yaml:"name"`
	Key        string              `yaml:"key"`
	Exchange   string              `yaml:"exchange"`
	Templates  map[string]Template `yaml:"templates"`
	Expression string              `yaml:"expression"`
}

// Exchange contains exchange name and signer unique identifier.
type Exchange struct {
	Name   string `yaml:"name"`
	Signer string `yaml:"signer"`
}

// Config represents application configuration model.
type Config struct {
	Languages []Language                    `yaml:"languages"`
	Keychain  map[string]eventapi.Validator `yaml:"keychain"`
	Exchanges map[string]Exchange           `yaml:"exchanges"`
	Events    []Event                       `yaml:"events"`
}

// Template returns Template model for given unique key.
func (e *Event) Template(key string) Template {
	return e.Templates[strings.ToUpper(key)]
}

// Content returns ready to go message with specified data.
// Note: "template" has bigger priority, than "template_path".
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

// ContainsLanguage reports whether the language code is known.
func (config *Config) ContainsLanguage(code string) bool {
	for _, lang := range config.Languages {
		if strings.EqualFold(lang.Code, code) {
			return true
		}
	}

	return false
}

// ContainsExchange reports whether the exchange with specified key exist.
func (config *Config) ContainsExchange(id string) bool {
	_, ok := config.Exchanges[id]
	return ok
}

// ContainsKey reports whether the keychain key exist.
func (config *Config) ContainsKey(id string) bool {
	_, ok := config.Keychain[id]
	return ok
}

// Valid reports whether configuration is valid or not.
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

// ValidateExchanges validates exchanges config.
func (config *Config) ValidateExchanges() error {
	if len(config.Exchanges) < 1 {
		return errors.New("no exchanges was specified")
	}

	for k, v := range config.Exchanges {
		if v.Name == "" {
			return fmt.Errorf("exchange name can not be empty: %s", k)
		}

		// Check, that signer is not empty and exist in keychain.
		if v.Signer == "" {
			return fmt.Errorf("signer %s of exchange %s can not be empty", v.Signer, k)
		} else if _, ok := config.Keychain[v.Signer]; !ok {
			return fmt.Errorf("signer %s is not registered", v.Signer)
		}
	}

	return nil
}

// ValidateKeychain validates keychain config.
func (config *Config) ValidateKeychain() error {
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
