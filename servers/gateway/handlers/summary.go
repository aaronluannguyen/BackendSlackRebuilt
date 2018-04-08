package handlers

import (
	"io"
	"net/http"
	"fmt"
	"strings"
	"golang.org/x/net/html"
	"strconv"
	url2 "net/url"
	"encoding/json"
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
const headerContentType = "Content-Type"
const contentTypeJSON = "application/json"

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
	w.Header().Add(headerContentType, contentTypeJSON)

	url := r.FormValue("url")
	if len(url) == 0 {
		http.Error(w, "status bad request error", http.StatusBadRequest)
		return
	}

	htmlRC, err := fetchHTML(url)
	if err != nil {
		http.Error(w, "error fetching html: " + err.Error(), http.StatusInternalServerError)
		return
	}
	defer htmlRC.Close()
	pgSummary, err := extractSummary(url, htmlRC)
	if err != nil {
		http.Error(w, "extracting summary error: " + err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(pgSummary)
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

	structMap := map[string]string{}
	var imageArray []*PreviewImage
	var insertImage *PreviewImage
	var iconImg *PreviewImage

	tokenizer := html.NewTokenizer(htmlStream)
	for {
		tokenType := tokenizer.Next()
		if tokenType == html.ErrorToken {
			err := tokenizer.Err()
			if err == io.EOF {
				break
			} else {
				fmt.Errorf("error tokenizing HTML: %v", err)
				return nil, err
			}
		}

		if tokenType == html.EndTagToken {
			token := tokenizer.Token()
			if "head" == token.Data {
				break
			}
		}

		if tokenType == html.StartTagToken || tokenType == html.SelfClosingTagToken {
			token := tokenizer.Token()
			switch {
			case "meta" == token.Data:
				k, v, err := handleAttr(token)
				if err != nil {
					return nil, err
				}
				switch {
				case k == "Image":
					absURL, err := absoluteUrl(pageURL, v)
					if err != nil {
						return nil, err
					}
					insertImage = &PreviewImage{}
					insertImage.URL = absURL
					imageArray = append(imageArray, insertImage)
				case k == "Image:Secure_URL":
					absURL, err := absoluteUrl(pageURL, v)
					if err != nil {
						return nil, err
					}
					insertImage.SecureURL = absURL
				case k == "Image:Type": insertImage.Type = v
				case k == "Image:Width":
					width, err := strconv.Atoi(v)
					if err != nil {
						return nil, err
					}
					if v != "" {
						insertImage.Width = width
					}
				case k == "Image:Height":
					height, err := strconv.Atoi(v)
					if err != nil {
						return nil, err
					}
					if v != "" {
						insertImage.Height = height
					}
				case k == "Image:Alt": insertImage.Alt = v
				default:
					if k != "" && v != "" {
						structMap[k] = v
					}
				}
			case "link" == token.Data:
				isIcon, iconUrl, iconType, iconWidth, iconHeight, iconAlt := false, "", "", 0, 0, ""
				for _, a := range token.Attr {
					switch {
					case a.Key == "rel" && a.Val == "icon": isIcon = true
					case a.Key == "href": iconUrl = a.Val
					case a.Key == "type": iconType = a.Val
					case a.Key == "sizes":
						if a.Val != "any" {
							sizes := strings.Split(a.Val, "x")
							height, err := strconv.Atoi(sizes[0])
							if err != nil {
								return nil, err
							}
							iconHeight = height
							width, err := strconv.Atoi(sizes[1])
							if err != nil {
								return nil, err
							}
							iconWidth = width
						}
					case a.Key == "alt": iconAlt = a.Val
					}
					if isIcon {
						absUrl, err := absoluteUrl(pageURL, iconUrl)
						if err != nil {
							return nil, err
						}
						icon := PreviewImage{
							absUrl, "", iconType, iconWidth, iconHeight, iconAlt,
						}
						iconImg = &icon
					}
				}
			case "title" == token.Data:
				tokenType = tokenizer.Next()
				if tokenType == html.TextToken {
					structMap["Title"] = tokenizer.Token().Data
				}
			}
		}
	}

	summary, err := updatePgSUm(structMap)
	if err != nil {
		fmt.Errorf("error updating page summary: %v", err)
		return nil, err
	}
	summary.Icon = iconImg
	summary.Images = imageArray
	return summary, nil
}

func handleAttr(token html.Token) (string, string, error) {
	prop, cont := "", ""
	for _, a := range token.Attr {
		switch {
		case a.Key == "property" :
			switch {
			case a.Val == "og:type": prop = "Type"
			case a.Val == "og:url": prop = "URL"
			case a.Val == "og:title": prop = "OG:Title"
			case a.Val == "og:site_name": prop = "SiteName"
			case a.Val == "og:description": prop = "OG:Description"
			case a.Val == "og:image": prop = "Image"
			case a.Val == "og:image:secure_url": prop = "Image:Secure_URL"
			case a.Val == "og:image:type": prop = "Image:Type"
			case a.Val == "og:image:width": prop = "Image:Width"
			case a.Val == "og:image:height": prop = "Image:Height"
			case a.Val == "og:image:alt": prop = "Image:Alt"
			}
		case a.Key == "name":
			switch {
			case a.Val == "author": prop = "Author"
			case a.Val == "keywords": prop = "Keywords"
			case a.Val == "description": prop = "Description"
			}
		case a.Key == "content": cont = a.Val
		}
	}
	return prop, cont, nil
}

func updatePgSUm(structMap map[string]string,) (*PageSummary, error) {
	pgSum := &PageSummary{}
	pgSum.Type = structMap["Type"]
	pgSum.URL = structMap["URL"]
	pgSum.Title = structMap["Title"]
	if structMap["OG:Title"] != "" {
		pgSum.Title = structMap["OG:Title"]
	}
	pgSum.SiteName = structMap["SiteName"]
	pgSum.Description = structMap["Description"]
	if structMap["OG:Description"] != "" {
		pgSum.Description = structMap["OG:Description"]
	}
	pgSum.Author = structMap["Author"]
	slicedKW := strings.Split(structMap["Keywords"], ",")
	if len(slicedKW) > 1 {
		for i, word := range slicedKW {
			slicedKW[i] = strings.TrimSpace(word)
		}
		pgSum.Keywords = slicedKW
	}
	return pgSum, nil
}

func absoluteUrl (base string, rel string) (string, error) {
	relUrl, err := url2.Parse(rel)
	if err != nil {
		return "", err
	}
	baseUrl, err := url2.Parse(base)
	if err != nil {
		return "", err
	}
	absUrl := baseUrl.ResolveReference(relUrl).String()
	return absUrl, nil
}