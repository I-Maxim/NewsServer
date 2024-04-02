package domain

type ClubHeader struct {
	ClubName       string `xml:"ClubName"`
	ClubWebsiteURL string `xml:"ClubWebsiteURL"`
}

type ExternalArticle struct {
	//Text              string `xml:",chardata"`
	ArticleURL        string  `xml:"ArticleURL"`
	NewsArticleID     string  `xml:"NewsArticleID"`
	PublishDate       string  `xml:"PublishDate"`
	Taxonomies        string  `xml:"Taxonomies"`
	TeaserText        *string `xml:"TeaserText"`
	Subtitle          string  `xml:"Subtitle"` //
	ThumbnailImageURL string  `xml:"ThumbnailImageURL"`
	Title             string  `xml:"Title"`
	BodyText          string  `xml:"BodyText"`         //
	GalleryImageURLs  string  `xml:"GalleryImageURLs"` //
	VideoURL          string  `xml:"VideoURL"`         //
	OptaMatchId       *string `xml:"OptaMatchId"`
	LastUpdateDate    string  `xml:"LastUpdateDate"`
	IsPublished       string  `xml:"IsPublished"`
}

type ArticleList struct {
	//Text               string `xml:",chardata"`
	NewsletterNewsItem []ExternalArticle `xml:"NewsletterNewsItem"`
}

type NewsArticleInformation struct {
	//XMLName        xml.Name `xml:"NewsArticleInformation"`
	//Text           string   `xml:",chardata"`
	ClubHeader
	NewsArticle ExternalArticle `xml:"NewsArticle"`
}

type NewListInformation struct {
	//XMLName             xml.Name `xml:"NewListInformation"`
	//Text                string   `xml:",chardata"`
	ClubHeader
	NewsletterNewsItems ArticleList `xml:"NewsletterNewsItems"`
}
