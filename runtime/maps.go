package runtime

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/haproxytech/models"
)

// ShowMaps returns map files description from runtime
func (s *SingleRuntime) ShowMaps() (models.Maps, error) {
	response, err := s.ExecuteWithResponse("show map")
	if err != nil {
		return nil, err
	}
	return s.parseMaps(response), nil
}

// CreateMap creates a new map file with its entries. Returns an error if file already exists
func (s *SingleRuntime) CreateMap(name string, file multipart.File) (models.MapEntries, error) {
	ext := filepath.Ext(name)
	if ext != ".map" {
		return nil, fmt.Errorf("provided file with %s extension, but supported .map", ext)
	}

	if _, err := os.Stat(name); err == nil {
		return nil, fmt.Errorf("file %s already exists. You should delete an existing file first", name)
	}

	dst, err := os.Create(name)
	if err != nil {
		return nil, fmt.Errorf("file could not be created. %s", err)
	}
	defer dst.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, file)
	if err != nil {
		return nil, err
	}

	entries := buf.String()
	err = ioutil.WriteFile(name, []byte(entries), 0644)
	if err != nil {
		return nil, err
	}

	me := ParseMapEntries(entries, false)
	buf.Reset()
	return me, nil
}

// parseMaps parses output from `show map` command and return array of map files
// First line in output represents format and is ignored
// Sample output format:
// # id (file) description
// -1 (/etc/haproxy/maps/hosts.map) pattern loaded from file '/etc/haproxy/maps/hosts.map' used by map at file '/etc/haproxy/haproxy.cfg' line 26
func (s *SingleRuntime) parseMaps(output string) models.Maps {
	if output == "" {
		return nil
	}
	maps := models.Maps{}

	lines := strings.Split(strings.TrimSpace(output), "\n")
	for _, line := range lines {
		m := s.parseMap(line)
		if m != nil {
			maps = append(maps, m)
		}
	}
	return maps
}

// parseMap parses one line from map files array and return it structured
func (s *SingleRuntime) parseMap(line string) *models.Map {
	if line == "" || strings.HasPrefix(strings.TrimSpace(line), "# id") {
		return nil
	}

	parts := strings.Fields(line)
	if len(parts) < 3 {
		return nil
	}

	m := &models.Map{
		ID:          parts[0],
		File:        strings.TrimSuffix(strings.TrimPrefix(parts[1], "("), ")"),
		Description: strings.Join(parts[2:], " "),
	}

	return m
}

// GetMap returns one structured runtime map file
func (s *SingleRuntime) GetMap(name string) (*models.Map, error) {
	maps, err := s.ShowMaps()
	if err != nil {
		return nil, err
	}

	for _, m := range maps {
		if m.File == name {
			return m, nil
		}
	}
	return nil, fmt.Errorf("%s doesn't exists", name)
}

// ClearMap removes all map entries from the map file.
func (s *SingleRuntime) ClearMap(name string) error {
	cmd := fmt.Sprintf("clear map %s", name)
	err := s.Execute(cmd)
	if err != nil {
		return err
	}
	return nil
}

// ShowMapEntries returns one map runtime entries
func (s *SingleRuntime) ShowMapEntries(name string) (models.MapEntries, error) {
	cmd := fmt.Sprintf("show map %s", name)
	response, err := s.ExecuteWithResponse(cmd)
	if err != nil {
		return nil, err
	}
	return ParseMapEntries(response, true), nil
}

// ParseMapEntries parses array of entries in one map file
// One line sample entry:
// ID			  Key                Value
// 0x55d155c6fbf0 static.example.com be_static
func ParseMapEntries(output string, hasId bool) models.MapEntries {
	if output == "" {
		return nil
	}
	me := models.MapEntries{}

	lines := strings.Split(strings.TrimSpace(output), "\n")
	for _, line := range lines {
		e := parseMapEntry(line, hasId)
		if e != nil {
			me = append(me, e)
		}
	}
	return me
}

// parseMapEntry parses one entry in one map file/runtime and returns it structured
func parseMapEntry(line string, hasId bool) *models.MapEntry {
	if line == "" || strings.HasPrefix(strings.TrimSpace(line), "#") {
		return nil
	}

	parts := strings.Fields(line)
	if len(parts) < 2 {
		return nil
	}

	m := &models.MapEntry{}
	if hasId {
		m.ID = parts[0] //map entries from runtime have ID
		m.Key = parts[1]
		m.Value = parts[2]
	} else {
		m.Key = parts[0] //map entries from file
		m.Value = parts[1]
	}
	return m
}

// AddMapEntry adds an entry into the map file
func (s *SingleRuntime) AddMapEntry(name, key, value string) error {
	cmd := fmt.Sprintf("add map %s %s %s", name, key, value)
	err := s.Execute(cmd)
	if err != nil {
		return err
	}
	return nil
}

// GetMapEntry returns one map runtime setting
func (s *SingleRuntime) GetMapEntry(name, id string) (*models.MapEntry, error) {
	cmd := fmt.Sprintf("get map %s %s", name, id)
	response, err := s.ExecuteWithResponse(cmd)
	if err != nil {
		return nil, err
	}

	m := &models.MapEntry{}
	parts := strings.Split(response, ",")
	for _, p := range parts {
		kv := strings.Split(p, "=")
		switch key := strings.TrimSpace(kv[0]); {
		case key == "key":
			m.Key = strings.TrimPrefix(strings.TrimSuffix(kv[1], "\""), "\"")
		case key == "value":
			m.Value = strings.TrimPrefix(strings.TrimSuffix(kv[1], "\""), "\"")
		}
	}
	//safe guard m.Key != id:
	//when id doesn't exists in runtime maps,
	//but any existing key is substring of id
	//get map command returns wrong result(BUG in HAProxy)
	//so we need to check it
	if m.Key == "" || m.Value == "" || m.Key != id {
		return nil, fmt.Errorf("%s entry not found", id)
	}
	return m, nil
}

// SetMapEntry replaces the value corresponding to each id in a map
func (s *SingleRuntime) SetMapEntry(name, id, value string) error {
	cmd := fmt.Sprintf("set map %s %s %s", name, id, value)
	err := s.Execute(cmd)
	if err != nil {
		return err
	}
	return nil
}

// DeleteMapEntry deletes all the map entries from the map by its id
func (s *SingleRuntime) DeleteMapEntry(name, id string) error {
	cmd := fmt.Sprintf("del map %s %s", name, id)
	err := s.Execute(cmd)
	if err != nil {
		return err
	}
	return nil
}
