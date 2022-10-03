package notify

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"text/template"

	"github.com/pdiorio/santashuffle/internal/selection"
	"gopkg.in/yaml.v3"
)

func readSettingsFromFile(filename string) (map[string]string, error) {
	yamlData, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return readSettingsFromYaml(yamlData)
}

func readSettingsFromYaml(yfile []byte) (map[string]string, error) {
	data := make(map[string]string)
	err := yaml.Unmarshal(yfile, &data)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func NotifyPariticpants(matches []*(selection.Match), settingsFilename string, dryrun bool) error {
	settings, _ := readSettingsFromFile(settingsFilename)

	for _, match := range matches {
		t := template.Must(template.New("body").Parse(settings["email_template"]))

		buf := new(bytes.Buffer)
		err := t.Execute(buf, match)
		if err != nil {
			return err
		}

		body := buf.String()

		email := EmailClient{
			EmailAccount:		settings["email_account"],
			AppPassword:		settings["app_password"],
			SmtpServer:      	settings["smtp_server"],
			SmtpPort:       	settings["smtp_port"],
		}

		if dryrun {
			fmt.Printf("--------------------------------------\nTo: %s\nSubject: %s\n%s\n",
				match.Gifter.Email, settings["email_subject"], body)
		} else {
			fmt.Printf("Notiyfing %s of their assignment.\n", match.Gifter.Name)
			err2 := email.SendEmail(match.Gifter.Email, settings["email_subject"], body)
			if err2 != nil {
				log.Fatal(err2)
			}
		}

	}
	return nil
}
