package utils

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func LookupValueByRef(ref string) (string, error) {
	re := regexp.MustCompile(`ref\+(.*?)://(.*)`)
	match := re.FindStringSubmatch(ref)
	refLookup := match[1]
	lookup := strings.TrimSpace(match[2])

	switch refLookup {
	case "env":
		refValue := os.Getenv(lookup)
		if refValue == "" {
			return "", fmt.Errorf("no value exists for ref '%s'", ref)
		}

		return refValue, nil
	case "cmd":
		cmd := exec.Command("sh", "-c", fmt.Sprintf("echo $(echo \"%s\"); 1>&2", lookup)) // noqa

		var out bytes.Buffer

		cmd.Stdout = &out

		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}

		return strings.TrimSpace(out.String()), nil
	default:
		return "", fmt.Errorf("not valid ref strategy" + refLookup)
	}
}
