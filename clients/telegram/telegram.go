package telegram

import (
	"TelegramBot/libs/e"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

const (
	APIMETHOD         = "getUpdates"
	SENDMESSAGEMETHOD = "sendMessage"
)

type Client struct {
	host     string
	basePath string
	client   http.Client
}

func New(host string, token string) *Client {
	return &Client{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
	}
}

func (c *Client) SendMessage(chat_id int, messageText string) error {
	query := url.Values{}
	query.Add("chat_id", strconv.Itoa(chat_id))
	query.Add("text", messageText)
	_, err := c.makeRequest(SENDMESSAGEMETHOD, query)
	if err != nil {
		return e.Wrap("can't send the message", err)
	}
	return nil
}

func newBasePath(token string) string {
	return "bot" + token
}

func (c *Client) Updates(offset int, limit int) (updates []Update, err error) {
	defer func() { err = e.WrapIfNil("can't handle the request", err) }()
	query := url.Values{}
	query.Add("offset", strconv.Itoa(offset))
	query.Add("limit", strconv.Itoa(limit))
	data, err := c.makeRequest(APIMETHOD, query)

	if err != nil {
		return nil, err
	}
	var res UpdatesResponse
	if err := json.Unmarshal(data, &res); err != nil {

		return nil, err
	}

	return res.Result, nil
}

func (c *Client) makeRequest(method string, query url.Values) (data []byte, err error) {
	defer func() { err = e.WrapIfNil("can't handle the request", err) }()
	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}
	request, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	request.URL.RawQuery = query.Encode()
	response, err := c.client.Do(request)

	if err != nil {
		return nil, err
	}

	defer func() { _ = response.Body.Close() }()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
