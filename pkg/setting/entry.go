package setting

import (
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

type Server struct {
	RunMode  string
	HttpPort int
}

var ServerInstance = &Server{
	RunMode:  gin.ReleaseMode,
	HttpPort: 80,
}

type App struct {
	APIs map[string]struct {
		Host      string `yaml:"host,omitempty"`
		Path      string `yaml:"path,omitempty"`
		AppID     string `yaml:"app_id,omitempty"`
		AppSecret string `yaml:"app_secret,omitempty"`
	} `yaml:"apis,omitempty"`

	Provider struct {
		EmailSES struct {
			Region string `yaml:"region,omitempty"`
			Access string `yaml:"access,omitempty"`
			Secret string `yaml:"secret,omitempty"`
		} `yaml:"email_ses,omitempty"`
	} `yaml:"provider,omitempty"`
	Template struct {
		Email map[string]struct {
			Sender  string `json:"sender,omitempty"`
			File    string `json:"file,omitempty"`
			Subject string `json:"subject,omitempty"`
		} `json:"email,omitempty"`
	} `json:"template,omitempty"`
}

var AppInstance = &App{}

func init() {
	yamlFile, err := ioutil.ReadFile("conf/app.yaml")
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(yamlFile, AppInstance)
	if err != nil {
		panic(err)
	}
}
