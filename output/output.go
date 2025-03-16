package output

import (
	"io"
	"net/http"
	"os"

	"github.com/keithpaterson/postal/config"
	"github.com/keithpaterson/postal/logging"

	ulog "github.com/keithpaterson/resweave-utils/logging"
)

type Outputter interface {
	Write(response *http.Response) error
}

func NewOutputter(cfg *config.Config) Outputter {
	log := logging.NamedLogger("output")

	writer, err := openWriter(cfg.Output)
	if err != nil {
		log.Errorw("NewOutputter", ulog.LogKeyError, err)
	}

	template := newTemplate(cfg, log)

	switch cfg.Output.OutFormat() {
	case config.OutFmtRaw:
		return newRawOutputter(cfg.Output, writer)
	case config.OutFmtText:
		return newTextOutputter(cfg.Output, template, writer)
	default:
		log.Warnw("NewOutputter", "warning", "unsupported outputter", "outputter", cfg.Output.Format)
		return newTextOutputter(cfg.Output, template, writer)
	}
}

func openWriter(cfg config.OutputConfig) (io.Writer, error) {
	switch cfg.Filename {
	case "", "stdout":
		return os.Stdout, nil
	case "stderr":
		return os.Stderr, nil
	default:
		return os.Create(cfg.Filename)
	}
}
