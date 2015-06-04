package sdk

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

const sdkIndexURL = "https://developer.android.com/sdk/index.html"

var sdkFilenameRegex = regexp.MustCompile("android-sdk_r([0-9.]+)-([a-zA-Z]+).([a-z]+)")

// VersionInfo SDK version info and download link
type VersionInfo struct {
	Version   string
	OS        string
	Extension string
	Size      int64
	SHA1      string
	URL       string
}

func (v VersionInfo) String() string {
	return fmt.Sprintf(
		"Filename: android-sdk_r%s-%s.%s, Size: %d, SHA1: %s, URL: %s",
		v.Version,
		v.OS,
		v.Extension,
		v.Size,
		v.SHA1,
		v.URL,
	)
}

func parseVersion(filename string) (info *VersionInfo) {
	matches := sdkFilenameRegex.FindStringSubmatch(filename)
	if len(matches) != 4 {
		return
	}
	info = &VersionInfo{
		Version:   matches[1],
		OS:        matches[2],
		Extension: matches[3],
	}
	return
}

// LatestVersions fetch latest sdk version info
func LatestVersions() (links []*VersionInfo, err error) {
	resp, err := http.Get(sdkIndexURL)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return
	}

	for node := doc; node != nil; node = depthFirst(node) {
		if node.Type == html.ElementNode && node.Data == "a" {
			inner := node.FirstChild
			if inner != nil && inner.Type == html.TextNode {
				if info := parseVersion(inner.Data); info != nil {
					if url, ok := getAttr(node.Attr, "href"); ok {
						info.URL = url
						info.SHA1 = findSha1(node)
						info.Size = findSize(node)
						links = append(links, info)
					}
				}
			}
		}
	}

	return
}

func depthFirst(node *html.Node) *html.Node {
	switch {
	case node == nil:
		return nil
	case node.FirstChild != nil:
		return node.FirstChild
	case node.NextSibling != nil:
		return node.NextSibling
	case node.Parent != nil:
		return nextParentSibling(node.Parent)
	default:
		return nil
	}
}

func nextParentSibling(node *html.Node) *html.Node {
	switch {
	case node == nil:
		return nil
	case node.NextSibling != nil:
		return node.NextSibling
	case node.Parent != nil:
		return nextParentSibling(node.Parent)
	default:
		return nil
	}
}

func getAttr(attrs []html.Attribute, key string) (value string, found bool) {
	for _, attr := range attrs {
		if strings.ToLower(attr.Key) == strings.ToLower(key) {
			return attr.Val, true
		}
	}
	return
}

type nextNode func(node *html.Node) *html.Node

func findSha1(node *html.Node) (sha1 string) {
	if node == nil {
		return
	}

	path := []nextNode{
		nextNode(func(n *html.Node) *html.Node { return n.Parent }),
		nextNode(func(n *html.Node) *html.Node { return n.NextSibling }),
		nextNode(func(n *html.Node) *html.Node { return n.NextSibling }),
		nextNode(func(n *html.Node) *html.Node { return n.NextSibling }),
		nextNode(func(n *html.Node) *html.Node { return n.NextSibling }),
		nextNode(func(n *html.Node) *html.Node { return n.FirstChild }),
	}

	for _, next := range path {
		if node = next(node); node == nil {
			return
		}
	}

	if node.Type != html.TextNode || len(node.Data) != 40 {
		return
	}
	return node.Data
}

func findSize(node *html.Node) (size int64) {
	if node == nil {
		return
	}

	path := []nextNode{
		nextNode(func(n *html.Node) *html.Node { return n.Parent }),
		nextNode(func(n *html.Node) *html.Node { return n.NextSibling }),
		nextNode(func(n *html.Node) *html.Node { return n.NextSibling }),
		nextNode(func(n *html.Node) *html.Node { return n.FirstChild }),
	}

	for _, next := range path {
		if node = next(node); node == nil {
			return
		}
	}

	if node.Type != html.TextNode || !strings.HasSuffix(node.Data, " bytes") {
		return
	}
	if maybeSize, err := strconv.ParseInt(strings.TrimSuffix(node.Data, " bytes"), 10, 64); err == nil {
		return maybeSize
	}
	return
}
