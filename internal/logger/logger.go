package logger

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
)

type Fields map[string]any
type Options struct {
	Directory string
	FileName  string
}

var log zerolog.Logger

func Init(o *Options) error {
	if o == nil {
		return errors.New("logger options is nil")
	}
	_ = os.MkdirAll(o.Directory, 0755)
	file, err := os.OpenFile(
		filepath.Join(o.Directory, o.FileName),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		return err
	}

	log = zerolog.New(file).
		With().
		Timestamp().
		Logger()

	return nil
}

func Debug(msg string, fields Fields) {
	event := log.Debug()

	if fields != nil {
		event.Interface("input_data", fields)
	}
	event.Msg(msg)
}

func Info(msg string, fields Fields) {
	event := log.Info()

	if fields != nil {
		event.Interface("input_data", fields)
	}
	event.Msg(msg)
}

func Warning(msg string, fields Fields) {
	event := log.Warn()
	if fields != nil {
		event.Interface("input_data", fields)
	}
	event.Msg(msg)
}

func Error(msg string, fields Fields, err error) {
	event := log.Error()

	if err != nil {
		event.Err(err)
	}

	if fields != nil {
		event.Interface("input_data", fields)
	}
	event.Msg(msg)
}
