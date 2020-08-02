package domain

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/buger/jsonparser"
)

type EntryList []*Entry

func (e *EntryList) UnmarshalJSON(data []byte) (err error) {
	*e = make([]*Entry, 0)

	_, err = jsonparser.ArrayEach(data, e.parseArray)

	return err
}

func (e EntryList) SetRemoteIP(remoteIP string) {
	for _, entry := range e {
		entry.SetRemoteIP(remoteIP)
	}
}

func (e *EntryList) parseArray(value []byte, dataType jsonparser.ValueType, offset int, err error) {
	if err != nil {
		return
	}

	if dataType != jsonparser.Object {
		return
	}

	entry := &Entry{}

	if err = entry.UnmarshalJSON(value); err == nil {
		*e = append(*e, entry)
	}
}

type Entry struct {
	Time        time.Time
	Namespace   string
	Source      string
	Host        string
	Level       string
	TraceID     string
	Message     string
	BuildCommit string
	ConfigHash  string
	RemoteIP    string
	Params      json.RawMessage
	StringKey   []string
	StringVal   []string
	FloatKey    []string
	FloatVal    []float64
}

func (e *Entry) UnmarshalJSON(data []byte) (err error) {
	*e = Entry{Params: data}

	return jsonparser.ObjectEach(data, e.parseRootObject)
}

func (e *Entry) SetRemoteIP(remoteIP string) {
	e.RemoteIP = remoteIP
}

func (e *Entry) parseRootObject(key, value []byte, dataType jsonparser.ValueType, offset int) (err error) {
	switch string(key) {
	case "time":
		e.Time, err = e.parseTime(value)
	case "namespace":
		e.Namespace = e.parseString(value)
	case "source":
		e.Source = e.parseString(value)
	case "host":
		e.Host = e.parseString(value)
	case "level":
		e.Level = e.parseString(value)
	case "trace_id":
		e.TraceID = e.parseString(value)
	case "message":
		e.Message = e.parseString(value)
	case "build_commit":
		e.BuildCommit = e.parseString(value)
	case "config_hash":
		e.ConfigHash = e.parseString(value)
	default:
		return e.parseOtherObject(key, value, dataType, offset)
	}

	return err
}

func (e *Entry) parseOtherObject(key, value []byte, dataType jsonparser.ValueType, _ int) (err error) {
	switch dataType {
	case jsonparser.Array:
		return e.parseArray(key, value)
	case jsonparser.Number:
		return e.appendFloat(key, value)
	case jsonparser.Object:
		return jsonparser.ObjectEach(value, e.parseOtherObject)
	case jsonparser.Boolean, jsonparser.NotExist, jsonparser.Null, jsonparser.String, jsonparser.Unknown:
		e.appendString(key, value)
	}

	return nil
}

func (e *Entry) parseArray(key, value []byte) (err error) {
	_, err = jsonparser.ArrayEach(value, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if err != nil {
			log.Printf("[critical] error in callback function: %v", err)
			return
		}

		if err = e.parseOtherObject(key, value, dataType, offset); err != nil {
			log.Printf("[critical] parse object failed: %v", err)
			return
		}
	})

	return err
}

func (e *Entry) appendFloat(key, value []byte) error {
	f, err := strconv.ParseFloat(string(value), 64)
	if err != nil {
		return err
	}

	e.FloatKey = append(e.FloatKey, e.parseString(key))
	e.FloatVal = append(e.FloatVal, f)

	return nil
}

func (e *Entry) appendString(key, value []byte) {
	e.StringKey = append(e.StringKey, e.parseString(key))
	e.StringVal = append(e.StringVal, e.parseString(value))
}

func (e *Entry) parseTime(data []byte) (time.Time, error) {
	str, err := jsonparser.ParseString(data)
	if err != nil {
		return time.Time{}, err
	}

	return time.Parse(time.RFC3339Nano, str)
}

func (e *Entry) parseString(data []byte) string {
	return strings.ToLower(string(data))
}
