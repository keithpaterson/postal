/*
package output contains tools for outputting the response.

For supported output types, refer to the [config] package

This package will resolve the template string before writing 'text' output.
The template is ignored when using the 'raw' output format.

# Output Templates:

The output template is a block of text that can contain token strings which will be expanded
before the output is generated.

Supported tokens include:
  - ${date:...}, ${time:...}, ${datetime:...} and ${epoch:...} tokens
  - ${env:...} tokens
  - ${prop:...} tokens
  - ${response:...} tokens

Response tokens allow you specify where in the template text you want the response information
to appear.

Supported response token resplacements are:
  - ${response:body}: the entire response body (partial/custom replacement is unsupported)
  - ${response:headers}: all the response headers as a single-line string separated by semicolons,
    e.g.
    "Access-Control-Allow-Origin=*; Access-Control-Allow-Credentials=true; Content-Type=application/json; Content-Length=429"
  - ${response:header=xxx}: the header value specified by "xxx"
  - ${response:status}: the status string for the response as supplied by golang
  - ${response:status-code}: the status code number for the response
  - ${response:content-length}: shortcut for extracting just the content-length value from the headers
*/
package output
