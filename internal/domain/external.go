package domain

type ClubHeader struct {
	ClubName       string `xml:"ClubName"`
	ClubWebsiteURL string `xml:"ClubWebsiteURL"`
}

type ExternalArticle struct {
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

type NewsArticleInformation struct {
	ClubHeader
	NewsArticle ExternalArticle `xml:"NewsArticle"`
}

type NewListInformation struct {
	ClubHeader
	NewsletterNewsItems struct {
		NewsletterNewsItem []ExternalArticle `xml:"NewsletterNewsItem"`
	} `xml:"NewsletterNewsItems"`
}
