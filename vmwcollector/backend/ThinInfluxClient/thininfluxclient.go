package ThinInfluxClient

import (
	"bytes"
	"compress/gzip"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// constants
var precisions = []string{"ns", "u", "ms", "s", "m", "h"}
var maxlines = 5000

// InfluxError message returned on error
// ffjson: noencoder
type InfluxError struct {
	Error string
}

// ThinInfluxClient sends data to influxdb
// ffjson: skip
type ThinInfluxClient struct {
	URL      string
	Username string
	password string
}

// NewThinInlfuxClient creates a new thin influx client
func NewThinInlfuxClient(server string, port int, database, username, password, precision string, ssl bool) (ThinInfluxClient, error) {
	if len(server) == 0 {
		return ThinInfluxClient{}, errors.New("No url indicated")
	}
	if port < 1000 || port > 65535 {
		return ThinInfluxClient{}, errors.New("Port not in acceptable range")
	}
	if len(database) == 0 {
		return ThinInfluxClient{}, errors.New("No database indicated")
	}
	found := false
	for _, p := range precisions {
		if p == precision {
			found = true
			break
		}
	}
	if !found {
		return ThinInfluxClient{}, errors.New("Precision '" + precision + "' not in suppoted presisions " + strings.Join(precisions, ","))
	}
	fullurl := "http"
	if ssl {
		fullurl += "s"
	}
	fullurl += "://" + server + ":" + strconv.Itoa(port) + "/write?db=" + database + "&precision=" + precision

	return ThinInfluxClient{URL: fullurl, Username: username, password: password}, nil
}

// Send data to influx db
// Data is represented by lines of influxdb lineprotocol
// see https://docs.influxdata.com/influxdb/v1.2/write_protocols/line_protocol_tutorial/
// Limiting submits to maxlines (currently 5000) items as in
// https://docs.influxdata.com/influxdb/v1.2/guides/writing_data/
func (client *ThinInfluxClient) Send(lines []string) error {
	if len(lines) > maxlines {
		// split push per maxlines
		for i := 0; i <= len(lines); i += maxlines {
			end := i + maxlines
			if end > len(lines) {
				end = len(lines)
			}
			flush := lines[i:end]
			err := client.Send(flush)
			if err != nil {
				return err
			}
		}
		return nil
	}
	// prepare the content
	push := strings.Join(lines, "\n")
	var buf bytes.Buffer
	g := gzip.NewWriter(&buf)
	if _, err := g.Write([]byte(push)); err != nil {
		return err
	}
	if err := g.Flush(); err != nil {
		return err
	}
	if err := g.Close(); err != nil {
		return err
	}
	// prepare the request
	req, err := http.NewRequest("POST", client.URL, &buf)
	if err != nil {
		return err
	}
	if len(client.Username) >= 0 && len(client.password) >= 0 {
		req.SetBasicAuth(client.Username, client.password)
	}
	req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	req.Header.Set("Content-Encoding", "gzip")
	clt := &http.Client{}
	resp, err := clt.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode == 204 {
		return nil
	}
	jsonerr := InfluxError{}
	if resp.StatusCode == 400 || resp.StatusCode == 404 || resp.StatusCode == 500 {
		defer resp.Body.Close() // nolint: errcheck
		// Check that the server actually sent compressed data
		var reader io.ReadCloser
		switch resp.Header.Get("Content-Encoding") {
		case "gzip":
			reader, err = gzip.NewReader(resp.Body)
			if err != nil {
				return err
			}
			defer reader.Close() // nolint: errcheck
		default:
			reader = resp.Body
		}
		body, err := ioutil.ReadAll(reader)
		if err != nil {
			return err
		}
		err = jsonerr.UnmarshalJSON(body)
		if err != nil {
			return err
		}
	}
	if resp.StatusCode == 400 {
		return errors.New("Influxdb Unacceptable request: " + strings.Trim(jsonerr.Error, " "))
	}
	if resp.StatusCode == 401 {
		return errors.New("Unauthorized access: check credentials and db")
	}
	if resp.StatusCode == 404 {
		return errors.New("Database not found: " + strings.Trim(jsonerr.Error, " "))
	}
	if resp.StatusCode == 500 {
		return errors.New("Server Busy: " + strings.Trim(jsonerr.Error, " "))
	}
	return nil
}
