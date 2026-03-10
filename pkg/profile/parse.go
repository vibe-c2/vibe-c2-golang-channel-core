package profile

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"

	coreerrors "github.com/vibe-c2/vibe-c2-golang-channel-core/pkg/errors"
)

// ParseYAML parses one YAML profile and validates semantic requirements.
func ParseYAML(data []byte) (Profile, error) {
	profile, err := parseSingleProfile(data)
	if err != nil {
		return Profile{}, coreerrors.Wrap(coreerrors.CodeProfileInvalid, "unable to parse profile YAML", err)
	}
	if err := Validate(profile); err != nil {
		return Profile{}, err
	}
	return profile, nil
}

// ParseYAMLProfiles parses one or multiple YAML documents split by `---`.
// Each document must contain one profile object.
func ParseYAMLProfiles(data []byte) ([]Profile, error) {
	trimmed := strings.TrimSpace(string(data))
	if trimmed == "" {
		return nil, coreerrors.New(coreerrors.CodeProfileInvalid, "profile YAML is empty")
	}

	segments := splitYAMLDocuments(trimmed)
	profiles := make([]Profile, 0, len(segments))
	for i, segment := range segments {
		p, err := ParseYAML([]byte(segment))
		if err != nil {
			return nil, coreerrors.Wrap(coreerrors.CodeProfileInvalid, fmt.Sprintf("invalid profile document %d", i), err)
		}
		profiles = append(profiles, p)
	}
	return profiles, nil
}

func parseSingleProfile(data []byte) (Profile, error) {
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	result := Profile{}
	section := ""
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()

		if idx := strings.Index(line, "#"); idx >= 0 {
			line = line[:idx]
		}
		if strings.TrimSpace(line) == "" {
			continue
		}
		if strings.TrimSpace(line) == "---" {
			continue
		}

		indent := countLeadingSpaces(line)
		trimmed := strings.TrimSpace(line)

		if indent == 0 && strings.HasSuffix(trimmed, ":") {
			section = strings.TrimSuffix(trimmed, ":")
			continue
		}

		key, value, ok := splitYAMLKeyValue(trimmed)
		if !ok {
			return Profile{}, fmt.Errorf("line %d: expected key: value", lineNumber)
		}

		if section == "mapping" && indent > 0 {
			switch key {
			case "profile_id":
				result.Mapping.ProfileID = value
			case "id":
				result.Mapping.ID = value
			case "encrypted_data":
				result.Mapping.EncryptedData = value
			}
			continue
		}

		section = ""
		switch key {
		case "profile_id":
			result.ProfileID = value
		case "channel_type":
			result.ChannelType = value
		case "enabled":
			parsed, err := strconv.ParseBool(value)
			if err != nil {
				return Profile{}, fmt.Errorf("line %d: invalid bool for enabled", lineNumber)
			}
			result.Enabled = parsed
		case "default_fallback":
			parsed, err := strconv.ParseBool(value)
			if err != nil {
				return Profile{}, fmt.Errorf("line %d: invalid bool for default_fallback", lineNumber)
			}
			result.DefaultFallback = parsed
		case "priority":
			parsed, err := strconv.Atoi(value)
			if err != nil {
				return Profile{}, fmt.Errorf("line %d: invalid integer for priority", lineNumber)
			}
			result.Priority = parsed
		}
	}

	if err := scanner.Err(); err != nil {
		return Profile{}, err
	}

	return result, nil
}

func splitYAMLDocuments(data string) []string {
	parts := strings.Split(data, "\n---")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		segment := strings.TrimSpace(strings.TrimPrefix(part, "---"))
		if segment != "" {
			out = append(out, segment)
		}
	}
	if len(out) == 0 && strings.TrimSpace(data) != "" {
		return []string{strings.TrimSpace(data)}
	}
	return out
}

func splitYAMLKeyValue(line string) (key string, value string, ok bool) {
	idx := strings.Index(line, ":")
	if idx <= 0 {
		return "", "", false
	}
	key = strings.TrimSpace(line[:idx])
	value = strings.TrimSpace(line[idx+1:])
	value = strings.Trim(value, "\"'")
	return key, value, true
}

func countLeadingSpaces(s string) int {
	count := 0
	for _, r := range s {
		if r != ' ' {
			break
		}
		count++
	}
	return count
}
