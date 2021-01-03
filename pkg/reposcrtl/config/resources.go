package config

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	kyaml "k8s.io/apimachinery/pkg/util/yaml"
)

type resource struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
}

func parseResources(paths []string, resourceFunc func([]byte) error) error {
	for _, path := range paths {
		var bs []byte

		var err error

		switch {
		case path == "-":
			bs, err = ioutil.ReadAll(os.Stdin)
			if err != nil {
				return fmt.Errorf("reading config from stdin: %s", err)
			}
		case strings.HasPrefix(path, "https://"):
			resp, err := http.Get(path)
			if err != nil {
				return fmt.Errorf("reading config from url: %s", err)
			}
			defer resp.Body.Close()

			bs, err = ioutil.ReadAll(resp.Body)

			if err != nil {
				return fmt.Errorf("reading config from url: %s", err)
			}
		default:
			bs, err = ioutil.ReadFile(path)
			if err != nil {
				return fmt.Errorf("reading config '%s': %s", path, err)
			}
		}

		reader := kyaml.NewYAMLReader(bufio.NewReaderSize(bytes.NewReader(bs), 4096))

		for {
			docBytes, err := reader.Read()
			if err == io.EOF {
				break
			}

			if err != nil {
				return fmt.Errorf("parsing config '%s': %s", path, err)
			}

			err = resourceFunc(docBytes)
			if err != nil {
				return fmt.Errorf("parsing resource config '%s': %s", path, err)
			}
		}
	}

	return nil
}
