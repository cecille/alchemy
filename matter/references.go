package matter

import "strings"

var DisallowedReferenceSuffixes = []string{"Command", "Feature", "Attribute", "Field", "Event"}

func StripReferenceSuffixes(newId string) string {
	for _, suffix := range DisallowedReferenceSuffixes {
		if strings.HasSuffix(newId, suffix) {
			newId = newId[0 : len(newId)-len(suffix)]
			break
		}
	}
	return newId
}
