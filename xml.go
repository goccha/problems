package problems

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/goccha/logging/log"
)

const (
	Ns9457 = "urn:ietf:rfc:9457"
	Ns7807 = "urn:ietf:rfc:7807"
)

type xmlInvalidParams []InvalidParam

func (ips *xmlInvalidParams) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = "invalid-params"
	root := []xml.Token{start}
	for _, ip := range *ips {
		root = append(root, ip)
	}
	for _, t := range root {
		err := e.EncodeToken(t)
		if err != nil {
			return err
		}
	}
	return e.EncodeToken(start.End())
}

func (ips *xmlInvalidParams) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	params := make(xmlInvalidParams, 0)
	var ip InvalidParam
	for {
		if err := d.Decode(&ip); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}
		params = append(params, ip)
	}
	*ips = params
	return nil
}

type xmlValidationErrors []ValidationError

func (vel *xmlValidationErrors) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = "validation-errors"
	root := []xml.Token{start}
	for _, ve := range *vel {
		root = append(root, ve)
	}
	for _, t := range root {
		err := e.EncodeToken(t)
		if err != nil {
			return err
		}
	}
	return e.EncodeToken(start.End())
}

func (vel *xmlValidationErrors) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	errs := make(xmlValidationErrors, 0)
	var ve ValidationError
	for {
		if err := d.Decode(&ve); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}
		errs = append(errs, ve)
	}
	*vel = errs
	return nil
}

func EncodeMap(e *xml.Encoder, ex map[string]interface{}) error {
	if ex == nil {
		return nil
	}
	var start xml.StartElement
	for k, v := range ex {
		start = xml.StartElement{Name: xml.Name{Local: k}}
		if err := e.EncodeToken(start); err != nil {
			return err
		}
		switch t := v.(type) {
		case string:
			if err := e.EncodeToken(xml.CharData(t)); err != nil {
				return err
			}
		case int, int8, int16, int32, int64:
			if err := e.EncodeToken(xml.CharData(fmt.Sprintf("%d", t))); err != nil {
				return err
			}
		case float32, float64:
			if err := e.EncodeToken(xml.CharData(fmt.Sprintf("%f", t))); err != nil {
				return err
			}
		case bool:
			if t {
				if err := e.EncodeToken(xml.CharData("true")); err != nil {
					return err
				}
			} else {
				if err := e.EncodeToken(xml.CharData("false")); err != nil {
					return err
				}
			}
		case []byte:
			for _, b := range t {
				if err := e.EncodeToken(xml.StartElement{Name: xml.Name{Local: "i"}}); err != nil {
					return err
				}
				if err := e.EncodeToken(xml.CharData([]byte{b})); err != nil {
					return err
				}
				if err := e.EncodeToken(xml.EndElement{Name: xml.Name{Local: "i"}}); err != nil {
					return err
				}
			}
		case []string:
			for _, s := range t {
				if err := e.EncodeToken(xml.StartElement{Name: xml.Name{Local: "i"}}); err != nil {
					return err
				}
				if err := e.EncodeToken(xml.CharData(s)); err != nil {
					return err
				}
				if err := e.EncodeToken(xml.EndElement{Name: xml.Name{Local: "i"}}); err != nil {
					return err
				}
			}
		case []int:
			for _, i := range t {
				if err := e.EncodeToken(xml.StartElement{Name: xml.Name{Local: "i"}}); err != nil {
					return err
				}
				if err := e.EncodeToken(xml.CharData(fmt.Sprintf("%d", i))); err != nil {
					return err
				}
				if err := e.EncodeToken(xml.EndElement{Name: xml.Name{Local: "i"}}); err != nil {
					return err
				}
			}
		case []int8:
			for _, i := range t {
				if err := e.EncodeToken(xml.StartElement{Name: xml.Name{Local: "i"}}); err != nil {
					return err
				}
				if err := e.EncodeToken(xml.CharData(fmt.Sprintf("%d", i))); err != nil {
					return err
				}
				if err := e.EncodeToken(xml.EndElement{Name: xml.Name{Local: "i"}}); err != nil {
					return err
				}
			}
		case []int16:
			for _, i := range t {
				if err := e.EncodeToken(xml.StartElement{Name: xml.Name{Local: "i"}}); err != nil {
					return err
				}
				if err := e.EncodeToken(xml.CharData(fmt.Sprintf("%d", i))); err != nil {
					return err
				}
				if err := e.EncodeToken(xml.EndElement{Name: xml.Name{Local: "i"}}); err != nil {
					return err
				}
			}
		case []int32:
			for _, i := range t {
				if err := e.EncodeToken(xml.StartElement{Name: xml.Name{Local: "i"}}); err != nil {
					return err
				}
				if err := e.EncodeToken(xml.CharData(fmt.Sprintf("%d", i))); err != nil {
					return err
				}
				if err := e.EncodeToken(xml.EndElement{Name: xml.Name{Local: "i"}}); err != nil {
					return err
				}
			}
		case []int64:
			for _, i := range t {
				if err := e.EncodeToken(xml.StartElement{Name: xml.Name{Local: "i"}}); err != nil {
					return err
				}
				if err := e.EncodeToken(xml.CharData(fmt.Sprintf("%d", i))); err != nil {
					return err
				}
				if err := e.EncodeToken(xml.EndElement{Name: xml.Name{Local: "i"}}); err != nil {
					return err
				}
			}
		case []float32:
			for _, f := range t {
				if err := e.EncodeToken(xml.StartElement{Name: xml.Name{Local: "i"}}); err != nil {
					return err
				}
				if err := e.EncodeToken(xml.CharData(fmt.Sprintf("%f", f))); err != nil {
					return err
				}
				if err := e.EncodeToken(xml.EndElement{Name: xml.Name{Local: "i"}}); err != nil {
					return err
				}
			}
		case []float64:
			for _, f := range t {
				if err := e.EncodeToken(xml.StartElement{Name: xml.Name{Local: "i"}}); err != nil {
					return err
				}
				if err := e.EncodeToken(xml.CharData(fmt.Sprintf("%f", f))); err != nil {
					return err
				}
				if err := e.EncodeToken(xml.EndElement{Name: xml.Name{Local: "i"}}); err != nil {
					return err
				}
			}
		case []bool:
			for _, b := range t {
				if err := e.EncodeToken(xml.StartElement{Name: xml.Name{Local: "i"}}); err != nil {
					return err
				}
				if b {
					if err := e.EncodeToken(xml.CharData("true")); err != nil {
						return err
					}
				} else {
					if err := e.EncodeToken(xml.CharData("false")); err != nil {
						return err
					}
				}
				if err := e.EncodeToken(xml.EndElement{Name: xml.Name{Local: "i"}}); err != nil {
					return err
				}
			}
		default:
			if marshal, ok := marshaller[k]; ok {
				if err := marshal(e, start, v); err != nil {
					return err
				}
			}
		}
		if err := e.EncodeToken(xml.EndElement{Name: xml.Name{Local: k}}); err != nil {
			return err
		}
	}
	return nil
}

type xmlProblemDetails struct {
	XMLName       xml.Name               `xml:"problem"`
	Xmlns         string                 `xml:"xmlns,attr"`
	Type          string                 `xml:"type"`
	Title         string                 `xml:"title"`
	Detail        string                 `xml:"detail"`
	Instance      string                 `xml:"instance"`
	Status        int                    `xml:"status,omitempty"`
	Code          string                 `xml:"code,omitempty"`
	InvalidParams xmlInvalidParams       `xml:"invalid-params,omitempty"`
	Errors        xmlValidationErrors    `xml:"errors,omitempty"`
	Extensions    map[string]interface{} `xml:",any"`
}

func (xo *xmlProblemDetails) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	tagName := start.Name.Local
	token, err := d.Token()
	for ; err == nil; token, err = d.Token() {
		switch el := token.(type) {
		case xml.StartElement:
			start = el
			tagName = el.Name.Local
		case xml.CharData:
			switch tagName {
			case "":
				continue
			case "problem":
				println(string(el))
				continue
			case "type":
				xo.Type = string(el)
			case "title":
				xo.Title = string(el)
			case "detail":
				xo.Detail = string(el)
			case "instance":
				xo.Instance = string(el)
			case "status":
				if status, err := strconv.Atoi(string(el)); err == nil {
					xo.Status = status
				}
			case "code":
				xo.Code = string(el)
			case "invalid-params":
				xo.InvalidParams = make(xmlInvalidParams, 0, 1)
				if err = d.DecodeElement(&xo.InvalidParams, &start); err != nil {
					return err
				}
				tagName = ""
			case "errors":
				xo.Errors = make(xmlValidationErrors, 0, 1)
				if err = d.DecodeElement(&xo.Errors, &start); err != nil {
					return err
				}
				tagName = ""
			default:
				if xo.Extensions == nil {
					xo.Extensions = make(map[string]interface{})
				}
				if unmarshal, ok := unmarshaler[tagName]; ok {
					if res, err := unmarshal(d, start); err != nil {
						return err
					} else {
						xo.Extensions[tagName] = res
					}
				} else {
					xo.Extensions[tagName] = string(el)
				}
				tagName = ""
			}
		case xml.EndElement:
			tagName = ""
		default:
			log.Warn(context.TODO()).Interface("token", token).Msg("xmlOverlay.UnmarshalXML default")
		}
	}
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}
	return nil
}

type XmlMarshaler func(e *xml.Encoder, start xml.StartElement, v interface{}) error

var marshaller = make(map[string]XmlMarshaler)

func RegisterXmlMarshaller(name string, fn XmlMarshaler) error {
	if fn == nil {
		return errors.New("marshaller func is nil")
	}
	if _, ok := marshaller[name]; ok {
		return errors.New("marshaller already registered")
	}
	marshaller[name] = fn
	return nil
}

type XmlUnmarshaler func(d *xml.Decoder, start xml.StartElement) (interface{}, error)

var unmarshaler = make(map[string]XmlUnmarshaler)

func RegisterXmlUnmarshaler(name string, fn XmlUnmarshaler) error {
	if fn == nil {
		return errors.New("unmarshaler func is nil")
	}
	if _, ok := unmarshaler[name]; ok {
		return errors.New("unmarshaler already registered")
	}
	unmarshaler[name] = fn
	return nil
}

func unmarshalIntArray[T int | int8 | int16 | int32 | int64](arr []T, bitSize int, d *xml.Decoder, start xml.StartElement) ([]T, error) {
	tagName := start.Name.Local
	token, err := d.Token()
	for ; err == nil; token, err = d.Token() {
		switch el := token.(type) {
		case xml.StartElement:
			start = el
			tagName = el.Name.Local
		case xml.CharData:
			if tagName == "" {
				continue
			}
			if tagName == "i" {
				switch bitSize {
				case 8:
					if i, err := strconv.ParseInt(string(el), 10, 8); err == nil {
						arr = append(arr, T(i))
					}
				case 16:
					if i, err := strconv.ParseInt(string(el), 10, 16); err == nil {
						arr = append(arr, T(i))
					}
				case 32:
					if i, err := strconv.ParseInt(string(el), 10, 32); err == nil {
						arr = append(arr, T(i))
					}
				case 64:
					if i, err := strconv.ParseInt(string(el), 10, 64); err == nil {
						arr = append(arr, T(i))
					}
				default:
					if i, err := strconv.Atoi(string(el)); err == nil {
						arr = append(arr, T(i))
					}
				}
			}
		case xml.EndElement:
			tagName = ""
		default:
			continue
		}
	}
	return arr, nil
}

func unmarshalFloatArray[T float32 | float64](arr []T, bitSize int, d *xml.Decoder, start xml.StartElement) ([]T, error) {
	tagName := start.Name.Local
	token, err := d.Token()
	for ; err == nil; token, err = d.Token() {
		switch el := token.(type) {
		case xml.StartElement:
			start = el
			tagName = el.Name.Local
		case xml.CharData:
			if tagName == "" {
				continue
			}
			if tagName == "i" {
				switch bitSize {
				case 32:
					if i, err := strconv.ParseFloat(string(el), 32); err == nil {
						arr = append(arr, T(i))
					}
				case 64:
					if i, err := strconv.ParseFloat(string(el), 64); err == nil {
						arr = append(arr, T(i))
					}
				}
			}
		case xml.EndElement:
			tagName = ""
		default:
			continue
		}
	}
	return arr, nil
}

type intArray []int

func (ia *intArray) UnmarshalXML(d *xml.Decoder, start xml.StartElement) (err error) {
	var arr []int
	arr, err = unmarshalIntArray(arr, 0, d, start)
	if err != nil {
		return err
	}
	*ia = arr
	return nil
}

func IntArray() XmlUnmarshaler {
	return func(d *xml.Decoder, start xml.StartElement) (interface{}, error) {
		var arr intArray
		if err := d.DecodeElement(&arr, &start); err != nil {
			return nil, err
		}
		return arr, nil
	}
}

type int8Array []int8

func (ia *int8Array) UnmarshalXML(d *xml.Decoder, start xml.StartElement) (err error) {
	var arr []int8
	arr, err = unmarshalIntArray(arr, 8, d, start)
	if err != nil {
		return err
	}
	*ia = arr
	return nil
}
func Int8Array() XmlUnmarshaler {
	return func(d *xml.Decoder, start xml.StartElement) (interface{}, error) {
		var arr int8Array
		if err := d.DecodeElement(&arr, &start); err != nil {
			return nil, err
		}
		return arr, nil
	}
}

type int16Array []int16

func (ia *int16Array) UnmarshalXML(d *xml.Decoder, start xml.StartElement) (err error) {
	var arr []int16
	arr, err = unmarshalIntArray(arr, 16, d, start)
	if err != nil {
		return err
	}
	*ia = arr
	return nil
}
func Int16Array() XmlUnmarshaler {
	return func(d *xml.Decoder, start xml.StartElement) (interface{}, error) {
		var arr int16Array
		if err := d.DecodeElement(&arr, &start); err != nil {
			return nil, err
		}
		return arr, nil
	}
}

type int32Array []int32

func (ia *int32Array) UnmarshalXML(d *xml.Decoder, start xml.StartElement) (err error) {
	var arr []int32
	arr, err = unmarshalIntArray(arr, 32, d, start)
	if err != nil {
		return err
	}
	*ia = arr
	return nil
}
func Int32Array() XmlUnmarshaler {
	return func(d *xml.Decoder, start xml.StartElement) (interface{}, error) {
		var arr int32Array
		if err := d.DecodeElement(&arr, &start); err != nil {
			return nil, err
		}
		return arr, nil
	}
}

type int64Array []int64

func (ia *int64Array) UnmarshalXML(d *xml.Decoder, start xml.StartElement) (err error) {
	var arr []int64
	arr, err = unmarshalIntArray(arr, 64, d, start)
	if err != nil {
		return err
	}
	*ia = arr
	return nil
}
func Int64Array() XmlUnmarshaler {
	return func(d *xml.Decoder, start xml.StartElement) (interface{}, error) {
		var arr int64Array
		if err := d.DecodeElement(&arr, &start); err != nil {
			return nil, err
		}
		return arr, nil
	}
}

type float32Array []float32

func (fa *float32Array) UnmarshalXML(d *xml.Decoder, start xml.StartElement) (err error) {
	var arr []float32
	arr, err = unmarshalFloatArray(arr, 32, d, start)
	if err != nil {
		return err
	}
	*fa = arr
	return nil
}
func Float32Array() XmlUnmarshaler {
	return func(d *xml.Decoder, start xml.StartElement) (interface{}, error) {
		var arr float32Array
		if err := d.DecodeElement(&arr, &start); err != nil {
			return nil, err
		}
		return arr, nil
	}
}

type float64Array []float64

func (fa *float64Array) UnmarshalXML(d *xml.Decoder, start xml.StartElement) (err error) {
	var arr []float64
	arr, err = unmarshalFloatArray(arr, 64, d, start)
	if err != nil {
		return err
	}
	*fa = arr
	return nil
}
func Float64Array() XmlUnmarshaler {
	return func(d *xml.Decoder, start xml.StartElement) (interface{}, error) {
		var arr float64Array
		if err := d.DecodeElement(&arr, &start); err != nil {
			return nil, err
		}
		return arr, nil
	}
}

type boolArray []bool

func (ba *boolArray) UnmarshalXML(d *xml.Decoder, start xml.StartElement) (err error) {
	var arr []bool
	tagName := start.Name.Local
	token, err := d.Token()
	for ; err == nil; token, err = d.Token() {
		switch el := token.(type) {
		case xml.StartElement:
			start = el
			tagName = el.Name.Local
		case xml.CharData:
			if tagName == "" {
				continue
			}
			if tagName == "i" {
				if string(el) == "true" {
					arr = append(arr, true)
				} else if string(el) == "false" {
					arr = append(arr, false)
				}
			}
		case xml.EndElement:
			tagName = ""
		default:
			continue
		}
	}
	*ba = arr
	return nil
}
func BoolArray() XmlUnmarshaler {
	return func(d *xml.Decoder, start xml.StartElement) (interface{}, error) {
		var arr boolArray
		if err := d.DecodeElement(&arr, &start); err != nil {
			return nil, err
		}
		return arr, nil
	}
}

type stringArray []string

func (sa *stringArray) UnmarshalXML(d *xml.Decoder, start xml.StartElement) (err error) {
	var arr []string
	tagName := start.Name.Local
	token, err := d.Token()
	for ; err == nil; token, err = d.Token() {
		switch el := token.(type) {
		case xml.StartElement:
			start = el
			tagName = el.Name.Local
		case xml.CharData:
			if tagName == "" {
				continue
			}
			if tagName == "i" {
				arr = append(arr, string(el))
			}
		case xml.EndElement:
			tagName = ""
		default:
			continue
		}
	}
	*sa = arr
	return nil
}
func StringArray() XmlUnmarshaler {
	return func(d *xml.Decoder, start xml.StartElement) (interface{}, error) {
		var arr stringArray
		if err := d.DecodeElement(&arr, &start); err != nil {
			return nil, err
		}
		return arr, nil
	}
}
