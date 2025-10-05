package driver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
)

type HttpDriver interface {
	Get(ctx context.Context, endpoint string, params any, header map[string]string) ([]byte, error)
	Post(ctx context.Context, endpoint string, reqBody any, header map[string]string) ([]byte, error)
}

type httpDriver struct {
	httpClient *http.Client
}

func NewClient(c *http.Client) HttpDriver {
	return &httpDriver{
		httpClient: c,
	}
}

func (c *httpDriver) Get(ctx context.Context, endpoint string, params any, header map[string]string) ([]byte, error) {
	parsedURL, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	if params != nil {
		queryParams, err := buildQueryParams(params)
		if err != nil {
			return nil, err
		}
		parsedURL.RawQuery = queryParams.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, parsedURL.String(), nil)
	if err != nil {
		return nil, err
	}

	for k, v := range header {
		req.Header.Set(k, v)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (c *httpDriver) Post(ctx context.Context, endpoint string, reqBody any, header map[string]string) ([]byte, error) {
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	for k, v := range header {
		req.Header.Set(k, v)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err

	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err

	}
	return body, nil
}

func buildQueryParams(params any) (url.Values, error) {
	values := url.Values{}

	if params == nil {
		return values, nil
	}

	v := reflect.ValueOf(params)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("params must be a struct or pointer to struct")
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		tag := fieldType.Tag.Get("param")
		if tag == "" || tag == "-" {
			continue
		}

		if !field.IsValid() || (field.Kind() == reflect.Ptr && field.IsNil()) {
			continue
		}

		var value string
		switch field.Kind() {
		case reflect.String:
			value = field.String()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			value = strconv.FormatInt(field.Int(), 10)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			value = strconv.FormatUint(field.Uint(), 10)
		case reflect.Float32, reflect.Float64:
			value = strconv.FormatFloat(field.Float(), 'f', -1, 64)
		case reflect.Bool:
			value = strconv.FormatBool(field.Bool())
		case reflect.Ptr:
			if !field.IsNil() {
				elem := field.Elem()
				switch elem.Kind() {
				case reflect.String:
					value = elem.String()
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					value = strconv.FormatInt(elem.Int(), 10)
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					value = strconv.FormatUint(elem.Uint(), 10)
				case reflect.Float32, reflect.Float64:
					value = strconv.FormatFloat(elem.Float(), 'f', -1, 64)
				case reflect.Bool:
					value = strconv.FormatBool(elem.Bool())
				default:
					continue
				}
			}
		default:
			continue
		}

		if value != "" {
			values.Add(tag, value)
		}
	}

	return values, nil
}
