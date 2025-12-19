package misc

import "encoding/json"

func ParseMetadata(comment string) map[string]interface{} {
	if comment == "" {
		return nil
	}
	metadata := make(map[string]interface{})
	err := json.Unmarshal([]byte(comment), &metadata)
	if err != nil {
		metadata["comment"] = comment
		return metadata
	}
	return metadata
}

func SerializeMetadata(metadata map[string]interface{}) (string, error) {
	if metadata == nil {
		return "", nil
	}
	b, err := json.Marshal(metadata)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
