package command

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/cli"
	"github.com/ryanuber/columnize"
)

func OutputSecret(ui cli.Ui, format string, secret *api.Secret) int {
	switch format {
	case "json":
		return outputFormatJSON(ui, secret)
	case "shell":
		return outputFormatShell(ui, secret)
	case "table":
		fallthrough
	default:
		return outputFormatTable(ui, secret, true)
	}
}

func outputFormatJSON(ui cli.Ui, s *api.Secret) int {
	b, err := json.Marshal(s)
	if err != nil {
		ui.Error(fmt.Sprintf(
			"Error formatting secret: %s", err))
		return 1
	}

	var out bytes.Buffer
	json.Indent(&out, b, "", "\t")
	ui.Output(out.String())
	return 0
}

func outputFormatShell(ui cli.Ui, s *api.Secret) int {
	input := []string{}

	for k, v := range s.Data {
		input = append(input, fmt.Sprintf("%s=%v", k, v))
	}

	ui.Output(strings.Join(input, "\n"))
	return 0
}

func outputFormatTable(ui cli.Ui, s *api.Secret, whitespace bool) int {
	config := columnize.DefaultConfig()
	config.Delim = "♨"
	config.Glue = "\t"
	config.Prefix = ""

	input := make([]string, 0, 5)
	input = append(input, fmt.Sprintf("Key %s Value", config.Delim))

	if s.LeaseID != "" && s.LeaseDuration > 0 {
		input = append(input, fmt.Sprintf("lease_id %s %s", config.Delim, s.LeaseID))
		input = append(input, fmt.Sprintf(
			"lease_duration %s %d", config.Delim, s.LeaseDuration))
	}

	for k, v := range s.Data {
		input = append(input, fmt.Sprintf("%s %s %v", k, config.Delim, v))
	}

	ui.Output(columnize.Format(input, config))
	return 0
}
