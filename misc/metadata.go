package misc

import "encoding/json"

func ParseMetadata(comment string) map[string]any {
	if comment == "" {
		return nil
	}
	metadata := make(map[string]any)
	err := json.Unmarshal([]byte(comment), &metadata)
	if err != nil {
		metadata["comment"] = comment
		return metadata
	}
	return metadata
}

func SerializeMetadata(metadata map[string]any) (string, error) {
	if metadata == nil {
		return "", nil
	}
	b, err := json.Marshal(metadata)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
