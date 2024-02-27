package converters

import "fmt"

type OnlineOfflineConverter struct {
	*BaseConverter
}

func (c *OnlineOfflineConverter) Convert(item string) ([]byte, error) {
	return nil, fmt.Errorf("Convert method is not implemented")
}
