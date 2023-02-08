package tg_send

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func Send(to string, msg string, token string) error {
	if to == "" {
		panic(errors.New("no chat_id provided"))
	} else if msg == "" {
		panic(errors.New("no message provided"))
	}

	data := url.Values{
		"parse_mode": {"MarkdownV2"},
		"chat_id":    {to},
		"text":       {msg},
	}
	if resp, err := http.PostForm(fmt.Sprintf(`https://api.telegram.org/bot%s/sendMessage`, token), data); err != nil {
		return err
	} else if bodyString, err := io.ReadAll(resp.Body); err != nil {
		return err
	} else if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("request error, code: %d\n%s", resp.StatusCode, bodyString))
	} else {
		return nil
	}
}
