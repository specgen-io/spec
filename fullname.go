package spec

import "fmt"

type FullName struct {
	Name Name
	Version Name
}

func (self *FullName) String() string {
	if self.Version.Source != "" {
		return fmt.Sprintf("%s.%s", self.Version.Source, self.Name.Source)
	} else {
		return self.Name.Source
	}
}
