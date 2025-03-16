package output

import (
	"io"
	"net/http"

	"github.com/keithpaterson/postal/config"
)

type rawOutputter struct {
	writer io.Writer
}

func newRawOutputter(_ config.OutputConfig, writer io.Writer) *rawOutputter {
	return &rawOutputter{writer: writer}
}

func (o *rawOutputter) Write(resp *http.Response) error {
	_, err := io.Copy(o.writer, resp.Body)
	return err
}
