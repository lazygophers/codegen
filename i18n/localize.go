package i18n

import (
	"github.com/lazygophers/log"
	"github.com/lazygophers/utils/json"
	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
	"sync"
)

// Localizer is the interface for marshaling and unmarshaling localization files.
// Implementations should support different file formats like JSON, YAML, and TOML.
type Localizer interface {
	MarshalToFile(path string, v any) (err error)
	UnmarshalFromFile(path string, v any) (err error)
}

var (
	localizer      = map[string]Localizer{}
	localizerMutex sync.RWMutex
)

func init() {
	localizer["json"] = jsonLocalizer
	localizer["yaml"] = yamlLocalizer
	localizer["yml"] = yamlLocalizer
	localizer["toml"] = tomlLocalizer
}

// Encoder is the interface for encoding values to a writer
type Encoder interface {
	Encode(v any) error
}

// Decoder is the interface for decoding values from a reader
type Decoder interface {
	Decode(v any) error
}

// EncoderFactory creates an encoder for the given file
type EncoderFactory func(file *os.File) Encoder

// DecoderFactory creates a decoder for the given file
type DecoderFactory func(file *os.File) Decoder

// createFileLocalizer creates a localizer using the given encoder and decoder factories
func createFileLocalizer(encoderFactory EncoderFactory, decoderFactory DecoderFactory) *LocalizerHandle {
	marshal := func(path string, v any) error {
		file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
		if err != nil {
			log.Errorf("failed to open file for writing %s: %v", path, err)
			return err
		}
		defer file.Close()

		encoder := encoderFactory(file)
		err = encoder.Encode(v)
		if err != nil {
			log.Errorf("failed to encode data to %s: %v", path, err)
			return err
		}

		return nil
	}

	unmarshal := func(path string, v any) error {
		file, err := os.Open(path)
		if err != nil {
			log.Errorf("failed to open file for reading %s: %v", path, err)
			return err
		}
		defer file.Close()

		decoder := decoderFactory(file)
		err = decoder.Decode(v)
		if err != nil {
			log.Errorf("failed to decode data from %s: %v", path, err)
			return err
		}

		return nil
	}

	return NewLocalizerHandle(marshal, unmarshal)
}

var (
	jsonLocalizer = createFileLocalizer(
		func(file *os.File) Encoder {
			encoder := json.NewEncoder(file)
			encoder.SetIndent("", "\t")
			return encoder
		},
		func(file *os.File) Decoder {
			return json.NewDecoder(file)
		},
	)

	yamlLocalizer = createFileLocalizer(
		func(file *os.File) Encoder {
			return yaml.NewEncoder(file)
		},
		func(file *os.File) Decoder {
			return yaml.NewDecoder(file)
		},
	)

	tomlLocalizer = createFileLocalizer(
		func(file *os.File) Encoder {
			return toml.NewEncoder(file)
		},
		func(file *os.File) Decoder {
			return toml.NewDecoder(file)
		},
	)
)

// RegisterLocalizer registers a custom localizer for the given format name.
// It's safe to call from multiple goroutines.
func RegisterLocalizer(name string, v Localizer) {
	localizerMutex.Lock()
	defer localizerMutex.Unlock()
	localizer[name] = v
}

// GetLocalizer returns the localizer for the given format name.
// It's safe to call from multiple goroutines.
func GetLocalizer(name string) (Localizer, bool) {
	name = strings.TrimPrefix(name, ".")

	localizerMutex.RLock()
	defer localizerMutex.RUnlock()
	l, ok := localizer[name]
	return l, ok
}

var _ Localizer = (*LocalizerHandle)(nil)

// LocalizerHandle wraps marshal and unmarshal functions to implement Localizer interface
type LocalizerHandle struct {
	marshal   func(path string, v any) (err error)
	unmarshal func(path string, v any) (err error)
}

// MarshalToFile marshals the value to a file
func (p *LocalizerHandle) MarshalToFile(path string, v any) (err error) {
	return p.marshal(path, v)
}

// UnmarshalFromFile unmarshals the value from a file
func (p *LocalizerHandle) UnmarshalFromFile(path string, v any) (err error) {
	return p.unmarshal(path, v)
}

// NewLocalizerHandle creates a new LocalizerHandle with the given marshal and unmarshal functions
func NewLocalizerHandle(marshal func(path string, v any) (err error), unmarshal func(path string, v any) (err error)) *LocalizerHandle {
	return &LocalizerHandle{
		marshal:   marshal,
		unmarshal: unmarshal,
	}
}
