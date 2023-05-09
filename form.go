package jhttp

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/textproto"
	"path/filepath"
	"strings"
)

//from.go is used to post form-data requests

type FieldType uint8
type ContentType string

const (
	Text FieldType = 1
	File FieldType = 2
)

const (
	Zip  ContentType = "application/zip"
	Byte ContentType = "application/octet-stream"
	Json ContentType = "application/json"
	XML  ContentType = "text/xml"
)

type FormOption = func(form *form)
type form struct {
	fields     []string
	values     []string
	fieldTypes []FieldType
}
type FormData struct {
	buf    *bytes.Buffer
	writer *multipart.Writer
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

func createFormFile(fieldName, fileName, contentType string, w *multipart.Writer) (io.Writer, error) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			escapeQuotes(fieldName), escapeQuotes(fileName)))
	h.Set("Content-Type", contentType)
	return w.CreatePart(h)
}

func NewFormParams(opts ...FormOption) (FormData, error) {
	var form form
	for _, opt := range opts {
		opt(&form)
	}
	formData, err := form.build()
	if err != nil {
		return FormData{}, err
	}
	return *formData, nil
}

func AddFormParams(field string, value string, fieldType FieldType) FormOption {
	return func(form *form) {
		form.fields = append(form.fields, field)
		form.values = append(form.values, value)
		form.fieldTypes = append(form.fieldTypes, fieldType)
	}
}

func (f form) build() (*FormData, error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for i := 0; i < len(f.fields); i++ {
		field := f.fields[i]
		value := f.values[i]
		fieldType := f.fieldTypes[i]
		switch fieldType {
		case File:
			file, err := ioutil.ReadFile(value)
			if err != nil {
				return nil, err
			}
			fileName := filepath.Base(value)
			contentType := getContentType(fileName)
			part, err := createFormFile(field, fileName, string(contentType), w)
			if err != nil {
				return nil, err
			}
			_, err = part.Write(file)
			if err != nil {
				return nil, err
			}
		case Text:
			err := w.WriteField(field, value)
			if err != nil {
				return nil, err
			}
		}
	}
	err := w.Close()
	if err != nil {
		return nil, err
	}
	return &FormData{
		buf:    &b,
		writer: w,
	}, nil
}

func getContentType(filename string) ContentType {
	str := strings.Split(filename, ".")
	if len(str) != 2 {
		return Byte
	}
	switch str[1] {
	case "zip":
		return Zip
	case "json":
		return Json
	case "xml":
		return XML
	default:
		return Byte
	}
}
