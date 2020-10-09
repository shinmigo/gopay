package kernel

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/url"
)

/**
 * Xml转Map
 */
type XmlToMap url.Values

type XmlToMapRes struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

func (m XmlToMap) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	for {
		res := XmlToMapRes{}
		err := decoder.Decode(&res)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		(m)[res.XMLName.Local] = []string{res.Value}
	}

	return nil
}

// Get gets the first value associated with the given key.
// If there are no values associated with the key, Get returns
// the empty string. To access multiple values, use the map
// directly.
func (m XmlToMap) Get(key string) string {
	if m == nil {
		return ""
	}
	vs := m[key]
	if len(vs) == 0 {
		return ""
	}
	return vs[0]
}

// Set sets the key to value. It replaces any existing
// values.
func (m XmlToMap) Set(key, value string) {
	m[key] = []string{value}
}

// Add adds the value to key. It appends to any existing
// values associated with key.
func (m XmlToMap) Add(key, value string) {
	m[key] = append(m[key], value)
}

// Del deletes the values associated with key.
func (m XmlToMap) Del(key string) {
	delete(m, key)
}

/**
 * Map转XML
 */
func mapToXml(params url.Values) string {
	buffer := &bytes.Buffer{}
	buffer.WriteString("<xml>")

	for paramKey := range params {
		paramValue := params.Get(paramKey)
		if paramKey == "total_fee" || paramKey == "refund_fee" || paramKey == "execute_time_" {
			buffer.WriteString(fmt.Sprintf("<%v>%v</%v>", paramKey, paramValue, paramKey))
		} else {
			buffer.WriteString(fmt.Sprintf("<%v><![CDATA[%v]]></%v>", paramKey, paramValue, paramKey))
		}
	}
	buffer.WriteString("</xml>")

	return buffer.String()
}
