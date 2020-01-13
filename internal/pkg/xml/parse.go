package xml

import (
	"bufio"
	"net/http"

	"golang.org/x/net/html/charset"
	"gopkg.in/xmlpath.v2"
)

func Parse(res *http.Response) (*Node, error) {
	r := bufio.NewReader(res.Body)
	data, err := r.Peek(1024)
	if err != nil {
		return nil, err
	}

	enc, _, ok := charset.DetermineEncoding(data, res.Header.Get("Content-Type"))
	if ok {
		return xmlpath.ParseHTML(enc.NewDecoder().Reader(r))
	} else {
		return xmlpath.ParseHTML(res.Body)
	}
}
