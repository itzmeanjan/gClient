package utils

import "strings"

type TopicList []string

func (t *TopicList) String() string {
	if t == nil {
		return ""
	}

	return strings.Join(*t, "")
}

func (t *TopicList) Set(val string) error {
	*t = append(*t, val)
	return nil
}
