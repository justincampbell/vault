package command

import (
	"bytes"
	"strings"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/cli"
)

func TestOutputSecret_shell(t *testing.T) {
	var raw = strings.TrimSpace(`{
	  "lease_id": "foo",
	  "renewable": true,
	  "lease_duration": 10,
	  "data": {
		"bar": "456",
		"foo": "123"
	  }
	}`)
	secret, err := api.ParseSecret(strings.NewReader(raw))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	writer := new(bytes.Buffer)
	ui := &cli.BasicUi{
		Writer: writer,
	}

	OutputSecret(ui, "shell", secret)

	if writer.String() != "bar=456\nfoo=123\n" {
		t.Fatalf("bad: %s", writer.String())
	}
}
