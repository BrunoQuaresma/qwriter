package qwriter

type Profile struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

var DefaultProfile = Profile{
	Name:        "Default",
	Description: "Correct any grammar and spelling errors and enhance the clarity and tone of the text.",
}
