package xml

import (
	"net/url"
	"strings"

	"gopkg.in/xmlpath.v2"

	"github.com/muniere/glean/internal/pkg/urls"
)

type Node = xmlpath.Node

func Title(doc *Node) string {
	xpath := xmlpath.MustCompile("//title")
	iter := xpath.Iter(doc)

	if iter.Next() {
		return iter.Node().String()
	} else {
		return ""
	}
}

func Collect(doc *Node, test func(*Node) bool) []*url.URL {
	hrefs := collect(doc, "//a/@href", test)
	srcs := collect(doc, "//img/@src", test)
	return urls.Unique(append(hrefs, srcs...))
}

func collect(doc *Node,path string, test func(*Node) bool) []*url.URL {
	var result []*url.URL

	xpath := xmlpath.MustCompile(path)
	iter := xpath.Iter(doc)

	for iter.Next() {
		n := iter.Node()

		if !test(n) {
			continue
		}

		s := strings.Replace(n.String(), " ", "+", -1)
		u, err := url.Parse(s)
		if err != nil {
			continue
		}

		result = append(result, u)
	}

	return result
}
