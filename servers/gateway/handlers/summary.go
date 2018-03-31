package handlers

import (
	"io"
	"net/http"
	"path"
	"fmt"
	"strings"
	"golang.org/x/net/html"
	"strconv"
)

//PreviewImage represents a preview image for a page
type PreviewImage struct {
	URL       string `json:"url,omitempty"`
	SecureURL string `json:"secureURL,omitempty"`
	Type      string `json:"type,omitempty"`
	Width     int    `json:"width,omitempty"`
	Height    int    `json:"height,omitempty"`
	Alt       string `json:"alt,omitempty"`
}

//PageSummary represents summary properties for a web page
type PageSummary struct {
	Type        string          `json:"type,omitempty"`
	URL         string          `json:"url,omitempty"`
	Title       string          `json:"title,omitempty"`
	SiteName    string          `json:"siteName,omitempty"`
	Description string          `json:"description,omitempty"`
	Author      string          `json:"author,omitempty"`
	Keywords    []string        `json:"keywords,omitempty"`
	Icon        *PreviewImage   `json:"icon,omitempty"`
	Images      []*PreviewImage `json:"images,omitempty"`
}

const headerCORS = "Access-Control-Allow-Origin"

//SummaryHandler handles requests for the page summary API.
//This API expects one query string parameter named `url`,
//which should contain a URL to a web page. It responds with
//a JSON-encoded PageSummary struct containing the page summary
//meta-data.
func SummaryHandler(w http.ResponseWriter, r *http.Request) {
	/*TODO: add code and additional functions to do the following:
	- Add an HTTP header to the response with the name
	 `Access-Control-Allow-Origin` and a value of `*`. This will
	  allow cross-origin AJAX requests to your server.
	- Get the `url` query string parameter value from the request.
	  If not supplied, respond with an http.StatusBadRequest error.
	- Call fetchHTML() to fetch the requested URL. See comments in that
	  function for more details.
	- Call extractSummary() to extract the page summary meta-data,
	  as directed in the assignment. See comments in that function
	  for more details
	- Close the response HTML stream so that you don't leak resources.
	- Finally, respond with a JSON-encoded version of the PageSummary
	  struct. That way the client can easily parse the JSON back into
	  an object. Remember to tell the client that the response content
	  type is JSON.

	Helpful Links:
	https://golang.org/pkg/net/http/#Request.FormValue
	https://golang.org/pkg/net/http/#Error
	https://golang.org/pkg/encoding/json/#NewEncoder
	*/

	w.Header().Add(headerCORS, "*")

	url := path.Base(r.URL.Path)
	if len(url) == 0 {
		fmt.Errorf("status bad request error: %v", http.StatusBadRequest)
		// !!! Check to see which way of error handling is preferred !!!
		// http.Error(w, "Status Bad Request Error", http.StatusBadRequest)
	}

	rc, err := fetchHTML(url)
	if err == nil {
		extractSummary(url, rc)
	}
}

//fetchHTML fetches `pageURL` and returns the body stream or an error.
//Errors are returned if the response status code is an error (>=400),
//or if the content type indicates the URL is not an HTML page.
func fetchHTML(pageURL string) (io.ReadCloser, error) {
	/*TODO: Do an HTTP GET for the page URL. If the response status
	code is >= 400, return a nil stream and an error. If the response
	content type does not indicate that the content is a web page, return
	a nil stream and an error. Otherwise return the response body and
	no (nil) error.

	To test your implementation of this function, run the TestFetchHTML
	test in summary_test.go. You can do that directly in Visual Studio Code,
	or at the command line by running:
		go test -run TestFetchHTML

	Helpful Links:
	https://golang.org/pkg/net/http/#Get
	*/

	response, err := http.Get(pageURL)
	if err != nil {
		return nil, fmt.Errorf("bad request error: %d", err)
	}
	if response.StatusCode >= 400 {
		return nil, fmt.Errorf("error fetching URL: %v", err)
	}
	contentType := response.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "text/html") {
		return nil, fmt.Errorf("response content type was %s not text/html", contentType)
	}
	return response.Body, nil
}

//extractSummary tokenizes the `htmlStream` and populates a PageSummary
//struct with the page's summary meta-data.
func extractSummary(pageURL string, htmlStream io.ReadCloser) (*PageSummary, error) {
	/*TODO: tokenize the `htmlStream` and extract the page summary meta-data
	according to the assignment description.

	To test your implementation of this function, run the TestExtractSummary
	test in summary_test.go. You can do that directly in Visual Studio Code,
	or at the command line by running:
		go test -run TestExtractSummary

	Helpful Links:
	https://drstearns.github.io/tutorials/tokenizing/
	http://ogp.me/
	https://developers.facebook.com/docs/reference/opengraph/
	https://golang.org/pkg/net/url/#URL.ResolveReference
	*/

	pgSum := new(PageSummary)
	pgSum.Images = []*PreviewImage{}

	tokenizer := html.NewTokenizer(htmlStream)
	for {
		tokenType := tokenizer.Next()

		// Error handling
		if tokenType == html.ErrorToken {
			fmt.Errorf("error tokenizing HTML: %v", tokenizer.Err())
			return nil, tokenizer.Err()
		}

		if tokenType == html.EndTagToken {
			token := tokenizer.Token()
			if "/head" == token.Data {
				break
			}
		}

		if tokenType == html.StartTagToken {
			token := tokenizer.Token()
			if "meta" == token.Data {
				structKey := ""
				// Looping over all attributes in "meta"
				for _, a := range token.Attr {
					if a.Key == "property" {
						if a.Val == "og:type" {
							structKey = "Type"
						} else if a.Val == "og:url" {
							structKey = "URL"
						} else if a.Val == "title" {
							structKey = "Title"
						} else if a.Val == "og:site_name" {
							structKey = "SiteName"
						} else if a.Val == "og:description" {
							structKey = "Description"
						} else if a.Val == "og:image" {
							structKey = "Images"
						} else if a.Val == "og:image:url" {
							structKey = "Images:url"
						} else if a.Val == "og:image:secure_url" {
							structKey = "Images:secure_url"
						} else if a.Val == "og:image:type" {
							structKey = "Images:type"
						} else if a.Val == "og:image:width" {
							structKey = "Images:width"
						} else if a.Val == "og:image:height" {
							structKey = "Images:height"
						} else if a.Val == "og:image:alt" {
							structKey = "Images:alt"
						}
					} else if a.Key == "name" {
						if a.Val == "author" {
							structKey = "Author"
						} else if a.Val == "keywords" {
							structKey = "Keywords"
						} else if pgSum.Description == "" && a.Val == "description" {
							structKey = "Description"
						}
					} else if structKey != "" && a.Key == "content" {
						if structKey == "Images" {
							pgSum.Images = append(pgSum.Images, new(PreviewImage))
							pgSum.Images[len(pgSum.Images)].URL = a.Val
						} else if structKey == "Images:secure_url" {
							pgSum.Images[len(pgSum.Images)].SecureURL = a.Val
						} else if structKey == "Images:type" {
							pgSum.Images[len(pgSum.Images)].Type = a.Val
						} else if structKey == "Images:width" {
							width, err := strconv.Atoi(a.Val)
							if err != nil {
								pgSum.Images[len(pgSum.Images)].Width = width
							}
						} else if structKey == "Images:height" {
							height, err := strconv.Atoi(a.Val)
							if err != nil {
								pgSum.Images[len(pgSum.Images)].Height = height
							}
						} else if structKey == "Images:alt" {
							pgSum.Images[len(pgSum.Images)].Alt = a.Val
						} else {
							pgSum.structKey = a.Val
						}
						structKey = ""
					}
				}
			} else if "link" == token.Data {
				isIcon := false
				for _, a := range token.Attr {
					if a.Key == "rel" && a.Val == "icon" {
						isIcon = true
					} else if isIcon && a.Key == "href" {
						pgSum.Icon.URL = a.Val
					}
				}
			} else if pgSum.Title == "" && "title" == token.Data {
				tokenType = tokenizer.Next()
				if tokenType == html.TextToken {
					pgSum.Title = tokenizer.Token().Data
				}
			}
		}
	}
	return pgSum, nil
}