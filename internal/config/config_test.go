package config

import (
	"os"
	"strings"
	"testing"

	"github.com/go-yaml/yaml"
	"github.com/stretchr/testify/assert"
)

var (
	configPath   = "../../config/postmaster.yml"
	templatePath = "../../test/test.tpl"
)

type FakeData struct {
	Test string
}

func NewFakeData(text string) *FakeData {
	return &FakeData{
		Test: text,
	}
}

func FakeConfig() Config {
	return Config{
		Languages: []Language{
			{
				Code: "EN",
				Name: "English",
			},
			{
				Code: "FR",
				Name: "French",
			},
		},
	}
}

func FakeTemplate() Template {
	return Template{
		Subject: "Fake",
	}
}

func DefaultConfig() Config {
	config := Config{}

	file, _ := os.Open(configPath)
	if err := yaml.NewDecoder(file).Decode(&config); err != nil {
		panic(err)
	}

	return config
}

func ExampleConfig() Config {
	return Config{
		Languages: []Language{
			{
				Code: "EN",
				Name: "English",
			},
		},
		Events: []Event{
			{
				Name: "Example",
				Key:  "example",
				Templates: map[string]Template{
					"EN": {
						Subject:  "Example",
						Template: "Yo",
					},
				},
			},
		},
	}
}

func TestConfig_ContainsLanguage(t *testing.T) {
	config := FakeConfig()

	t.Run("missing language code", func(t *testing.T) {
		contains := config.ContainsLanguage("UA")
		assert.False(t, contains)
	})

	t.Run("valid upper cased language code", func(t *testing.T) {
		contains := config.ContainsLanguage("EN")
		assert.True(t, contains)
	})

	t.Run("valid lower cased language code", func(t *testing.T) {
		contains := config.ContainsLanguage("en")
		assert.True(t, contains)
	})
}

func TestEvent_Template(t *testing.T) {
	code := "RU"

	config := DefaultConfig()

	file, err := os.Open(configPath)
	assert.NoError(t, err)

	err = yaml.NewDecoder(file).Decode(&config)
	assert.NoError(t, err)

	assert.Equal(t,
		config.Events[0].Templates[strings.ToUpper(code)],
		config.Events[0].Template(strings.ToUpper(code)),
	)

	assert.Equal(t,
		config.Events[0].Templates[strings.ToUpper(code)],
		config.Events[0].Template(strings.ToLower(code)),
	)
}

func TestTemplate_Content(t *testing.T) {
	t.Run("has only template", func(t *testing.T) {
		tmp := FakeTemplate()
		tmp.Template = "{{ .Test }}"

		data := NewFakeData("OpenWare")
		result, err := tmp.Content(&data)
		assert.NoError(t, err)

		assert.Equal(t, "OpenWare", string(result))
	})

	t.Run("has only template path", func(t *testing.T) {
		tmp := FakeTemplate()
		tmp.TemplatePath = templatePath

		data := NewFakeData("OpenWare")
		result, err := tmp.Content(&data)
		assert.NoError(t, err)

		assert.Equal(t, "OpenWare", string(result))
	})

	t.Run("has both template and template path", func(t *testing.T) {
		tmp := FakeTemplate()
		tmp.TemplatePath = templatePath
		tmp.Template = "Nothing"

		data := NewFakeData("OpenWare")
		result, err := tmp.Content(&data)
		assert.NoError(t, err)

		assert.Equal(t, "Nothing", string(result))
	})
}

func TestLanguage_Valid(t *testing.T) {
	name := "French"
	code := "FR"

	t.Run("has empty language code", func(t *testing.T) {
		lang := Language{
			Name: name,
		}

		assert.False(t, lang.Valid())
	})

	t.Run("has lower cased language code", func(t *testing.T) {
		lang := Language{
			Code: strings.ToLower(code),
			Name: name,
		}
		assert.False(t, lang.Valid())
	})

	t.Run("has upper cased language code", func(t *testing.T) {
		lang := Language{
			Code: strings.ToUpper(code),
			Name: name,
		}
		assert.True(t, lang.Valid())
	})
}

func TestValidate(t *testing.T) {
	t.Run("lower cased language codes", func(t *testing.T) {
		tmp := ExampleConfig()
		tmp.Languages[0].Code = strings.ToLower(tmp.Languages[0].Code)

		err := tmp.Validate()
		assert.Error(t, err)
	})

	t.Run("event without templates", func(t *testing.T) {
		tmp := ExampleConfig()
		tmp.Events[0].Templates = make(map[string]Template, 0)

		err := tmp.Validate()
		assert.Error(t, err)
	})

	t.Run("default should be valid", func(t *testing.T) {
		tmp := DefaultConfig()
		err := tmp.Validate()

		assert.NoError(t, err)
	})
}
