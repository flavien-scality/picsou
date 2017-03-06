package settings

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

type User struct {
	Name string `yaml:"name"`
	Email string `yaml:"email"`
}

type Settings struct {
	Users []User `yaml:",flow"`
	Path string `yaml:"path"`
}

func New(path string) *Settings {
	var settings Settings
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal([]byte(data), &settings)
	if err != nil {
		panic(err)
	}
	return &settings
}

func (s *Settings) GetUsersEmail() []string {
	var result []string
	for _, user := range s.Users {
		result = append(result, user.Email)
	}
	return result
}

func (s *Settings) GetName(email string) string {
	for _, user := range s.Users {
		if user.Email == email {
			return user.Name
		}
	}
	return ""
}
