package config

type GitAccessProtocol string

const (
	SSH  GitAccessProtocol = "ssh"
	HTTP GitAccessProtocol = "http"
)

func (r GitAccessProtocol) String() string {
	return string(r)
}
