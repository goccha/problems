package problems

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"
)

type Custom struct {
	Name string `xml:"name"`
	Age  int    `xml:"age"`
}

func (c Custom) MarshalXML(e *xml.Encoder, _ xml.StartElement) error {
	if err := e.EncodeElement(c.Name, xml.StartElement{Name: xml.Name{Local: "name"}}); err != nil {
		return err
	}
	if err := e.EncodeElement(c.Age, xml.StartElement{Name: xml.Name{Local: "age"}}); err != nil {
		return err
	}
	return nil
}

func TestMarshalXml(t *testing.T) {

	p := New(Instance("/problems"), ValidationErrors(nil,
		ValidationError{
			Detail:  "required",
			Pointer: "#/name",
		},
		ValidationError{
			Detail:  "required",
			Pointer: "#/nested/0/name",
		}),
		Extension("account", "00001"),
		Extension("accounts", []string{"00002", "00003"}),
		Extension("custom", Custom{Name: "name", Age: 20})).
		BadRequest("bad request")

	if err := RegisterXmlMarshaller("custom", func(e *xml.Encoder, start xml.StartElement, v interface{}) error {
		if c, ok := v.(Custom); ok {
			if err := e.EncodeElement(c, start); err != nil {
				return err
			}
			//if err := e.EncodeToken(xml.StartElement{Name: xml.Name{Local: "name"}}); err != nil {
			//	return err
			//}
			//if err := e.EncodeToken(xml.CharData(c.Name)); err != nil {
			//	return err
			//}
			//if err := e.EncodeToken(xml.EndElement{Name: xml.Name{Local: "name"}}); err != nil {
			//	return err
			//}
			//if err := e.EncodeToken(xml.StartElement{Name: xml.Name{Local: "age"}}); err != nil {
			//	return err
			//}
			//if err := e.EncodeToken(xml.CharData(fmt.Sprintf("%d", c.Age))); err != nil {
			//	return err
			//}
			//if err := e.EncodeToken(xml.EndElement{Name: xml.Name{Local: "age"}}); err != nil {
			//	return err
			//}
			return nil
		}
		return fmt.Errorf("invalid type %T", v)
	}); err != nil {
		t.Errorf("%v", err)
	}
	var pd *ProblemDetails
	errors.As(p, &pd)
	if bin, err := xml.Marshal(pd); err != nil {
		t.Errorf("%v", err)
	} else {
		str := string(bin)
		if !strings.Contains(str, "<errors><error><detail>required</detail><pointer>#/name</pointer></error><error><detail>required</detail><pointer>#/nested/0/name</pointer></error></errors>") {
			t.Errorf(`expected "<errors><error><detail>required</detail><pointer>#/name</pointer></error><error><detail>required</detail><pointer>#/nested/0/name</pointer></error></errors>"`)
		}
		if !strings.Contains(str, "<account>00001</account>") {
			t.Errorf(`expected "<account>00001</account>"`)
		}
		if !strings.Contains(str, "<accounts><i>00002</i><i>00003</i></accounts>") {
			t.Errorf(`expected "<accounts><i>00002</i><i>00003</i></accounts>"`)
		}
		if !strings.Contains(str, "<custom><name>name</name><age>20</age></custom>") {
			t.Errorf(`expected "<custom><name>name</name><age>20</age></custom>"`)
		}
	}
}

func TestUnmarshalXml(t *testing.T) {

	body := `<?xml version="1.0" encoding="UTF-8"?>
<problem xmlns="http://localhost:8080/problems">
	<type>http://localhost:8080/problems</type>
	<status>400</status>
	<title>Bad Request</title>
	<detail>bad request</detail>
	<instance>/problems</instance>
	<code>BadRequest</code>
	<invalid-params>
		<invalid-param>
			<name>name</name>
			<reason>required</reason>
		</invalid-param>
		<invalid-param>
			<name>nested</name>
			<reason>required</reason>
		</invalid-param>
	</invalid-params>
	<custom-attribute>custom value</custom-attribute>
	<int-array>
		<i>1</i>
	</int-array>
	<int8-array>
		<i>1</i>
	</int8-array>
</problem>`

	_ = RegisterXmlUnmarshaler("int-array", IntArray())
	_ = RegisterXmlUnmarshaler("int8-array", Int8Array())

	p, err := UnmarshalXml(t.Context(), http.StatusBadRequest, []byte(body))
	if err != nil {
		t.Errorf("%v", err)
	} else {
		if p.StatusCode() != http.StatusBadRequest {
			t.Errorf("expect = %d, actual = %d", http.StatusBadRequest, p.StatusCode())
		}
		var req *ProblemDetails
		ok := errors.As(p, &req)
		if !ok {
			t.Errorf("invalid struct. %v", p)
		}
	}
}
