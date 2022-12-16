package isbn

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	. "github.com/emad-elsaid/xlog"
	. "github.com/emad-elsaid/xlog/extensions/shortcode"
)

type Books struct {
	Kind       string `json:"kind"`
	TotalItems int    `json:"totalItems"`
	Items      []struct {
		Kind       string `json:"kind"`
		ID         string `json:"id"`
		Etag       string `json:"etag"`
		SelfLink   string `json:"selfLink"`
		VolumeInfo struct {
			Title               string   `json:"title"`
			Authors             []string `json:"authors"`
			Publisher           string   `json:"publisher"`
			PublishedDate       string   `json:"publishedDate"`
			Description         string   `json:"description"`
			IndustryIdentifiers []struct {
				Type       string `json:"type"`
				Identifier string `json:"identifier"`
			} `json:"industryIdentifiers"`
			ReadingModes struct {
				Text  bool `json:"text"`
				Image bool `json:"image"`
			} `json:"readingModes"`
			PageCount           int      `json:"pageCount"`
			PrintType           string   `json:"printType"`
			Categories          []string `json:"categories"`
			AverageRating       float32  `json:"averageRating"`
			RatingsCount        int      `json:"ratingsCount"`
			MaturityRating      string   `json:"maturityRating"`
			AllowAnonLogging    bool     `json:"allowAnonLogging"`
			ContentVersion      string   `json:"contentVersion"`
			PanelizationSummary struct {
				ContainsEpubBubbles  bool `json:"containsEpubBubbles"`
				ContainsImageBubbles bool `json:"containsImageBubbles"`
			} `json:"panelizationSummary"`
			ImageLinks struct {
				SmallThumbnail string `json:"smallThumbnail"`
				Thumbnail      string `json:"thumbnail"`
			} `json:"imageLinks"`
			Language            string `json:"language"`
			PreviewLink         string `json:"previewLink"`
			InfoLink            string `json:"infoLink"`
			CanonicalVolumeLink string `json:"canonicalVolumeLink"`
		} `json:"volumeInfo"`
		SaleInfo struct {
			Country     string `json:"country"`
			Saleability string `json:"saleability"`
			IsEbook     bool   `json:"isEbook"`
		} `json:"saleInfo"`
		AccessInfo struct {
			Country                string `json:"country"`
			Viewability            string `json:"viewability"`
			Embeddable             bool   `json:"embeddable"`
			PublicDomain           bool   `json:"publicDomain"`
			TextToSpeechPermission string `json:"textToSpeechPermission"`
			Epub                   struct {
				IsAvailable bool `json:"isAvailable"`
			} `json:"epub"`
			Pdf struct {
				IsAvailable bool `json:"isAvailable"`
			} `json:"pdf"`
			WebReaderLink       string `json:"webReaderLink"`
			AccessViewStatus    string `json:"accessViewStatus"`
			QuoteSharingAllowed bool   `json:"quoteSharingAllowed"`
		} `json:"accessInfo"`
		SearchInfo struct {
			TextSnippet string `json:"textSnippet"`
		} `json:"searchInfo"`
	} `json:"items"`
}

func init() {
	ShortCode("google-books-isbn", isbnShortCode)
}

func isbnShortCode(n Markdown) Markdown {
	var u url.URL
	u.Scheme = "https"
	u.Host = "www.googleapis.com"
	u.Path = "/books/v1/volumes"
	query := u.Query()
	query.Add("q", fmt.Sprintf("isbn:%s", strings.TrimSpace(string(n))))
	u.RawQuery = query.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return Markdown(err.Error())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Markdown(err.Error())
	}

	var books Books
	err = json.Unmarshal(body, &books)
	if err != nil {
		return Markdown(err.Error())
	}

	if len(books.Items) == 0 {
		return Markdown(err.Error())
	}

	book := books.Items[0]

	out := fmt.Sprintf(`
|||
|----:|----|
| ![](%s) | **[%s](%s)**<br>By %s<br>%d Pages <br>Publisher: %s (%s) |
`,
		book.VolumeInfo.ImageLinks.Thumbnail,
		book.VolumeInfo.Title,
		book.VolumeInfo.InfoLink,
		strings.Join(book.VolumeInfo.Authors, ", "),
		book.VolumeInfo.PageCount,
		book.VolumeInfo.Publisher,
		book.VolumeInfo.PublishedDate,
	)

	return Markdown(out)
}
