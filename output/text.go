package output

import (
	"io"
	"net/http"

	"github.com/keithpaterson/postal/config"
)

type textOutputter struct {
	template *template
	writer   io.Writer
}

func newTextOutputter(_ config.OutputConfig, template *template, writer io.Writer) *textOutputter {
	return &textOutputter{template: template, writer: writer}
}

func (o *textOutputter) Write(resp *http.Response) error {
	_, err := o.writer.Write([]byte(o.template.Apply(resp)))
	return err
}
