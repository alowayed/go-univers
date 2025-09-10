package mattermost

const Name = "mattermost"

type Ecosystem struct{}

func (e *Ecosystem) Name() string {
	return Name
}
