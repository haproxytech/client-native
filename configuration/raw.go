package configuration

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/haproxytech/models"
)

// GetRawConfiguration returns a struct with configuration version and a
// string containing raw config file
func (c *Client) GetRawConfiguration() (*models.GetHAProxyConfigurationOKBody, error) {
	file, err := os.Open(c.ConfigurationFile)
	if err != nil {
		return nil, NewConfError(ErrCannotReadConfFile, err.Error())
	}
	defer file.Close()

	dataStr := ""
	ondiskV := int64(0)
	scanner := bufio.NewScanner(file)
	// parse out version
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "# _version=") {
			w := strings.Split(line, "=")
			if len(w) == 2 {
				ondiskV, err = strconv.ParseInt(w[1], 10, 64)
				if err != nil {
					ondiskV = int64(0)
				}
			}
		} else {
			dataStr += line + "\n"
		}
	}
	if err = scanner.Err(); err != nil {
		return nil, NewConfError(ErrCannotReadConfFile, err.Error())

	}

	data := &models.GetHAProxyConfigurationOKBody{
		Configuration: dataStr,
		Version:       ondiskV,
	}
	return data, nil
}

// PostRawConfiguration pushes given string to the config file if the version
// matches
func (c *Client) PostRawConfiguration(config *string, version int64) error {
	ondiskV, _ := c.GetVersion("")
	if ondiskV != version {
		return NewConfError(ErrVersionMismatch, fmt.Sprintf("Version in configuration file is %v, given version is %v", ondiskV, version))
	}

	tmp, e := ioutil.TempFile(filepath.Dir(c.ConfigurationFile), filepath.Base(c.ConfigurationFile))
	defer tmp.Close()

	if e != nil {
		return NewConfError(ErrGeneralError, e.Error())
	}

	w := bufio.NewWriter(tmp)
	w.WriteString(fmt.Sprintf("# _version=%v\n%v", ondiskV+1, *config))
	w.Flush()

	tmpPath, err := filepath.Abs(filepath.Dir(tmp.Name()))
	if err != nil {
		return NewConfError(ErrGeneralError, e.Error())
	}

	if err = c.validateConfigFile(tmpPath); err != nil {
		_ = os.Remove(tmpPath)
		return err
	}

	err = os.Rename(tmpPath, c.ConfigurationFile)
	if err != nil {
		os.Remove(tmpPath)
		return NewConfError(ErrGeneralError, e.Error())
	}

	if c.Cache.Enabled() {
		c.Cache.InvalidateCache()
	}
	return nil
}

func (c *Client) validateConfigFile(confFile string) error {
	cmd := exec.Command(c.Haproxy)
	cmd.Args = append(cmd.Args, "-c")

	if confFile != "" {
		cmd.Args = append(cmd.Args, "-f")
		cmd.Args = append(cmd.Args, confFile)
	} else {
		cmd.Args = append(cmd.Args, "-f")
		cmd.Args = append(cmd.Args, c.ConfigurationFile)
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return NewConfError(ErrValidationError, err.Error())
	}
	return nil
}
