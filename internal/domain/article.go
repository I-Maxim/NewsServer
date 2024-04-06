package domain

type ArticleResponse struct {
	Id          string   `json:"id"`
	TeamId      string   `json:"teamId"`
	OptaMatchId *string  `json:"optaMatchId"`
	Title       string   `json:"title"`
	Type        []string `json:"type"`
	Teaser      *string  `json:"teaser"`
	Content     string   `json:"content"`
	Url         string   `json:"url"`
	ImageUrl    string   `json:"imageUrl"`
	GalleryUrls string   `json:"galleryUrls"`
	VideoUrl    string   `json:"videoUrl"`
	Published   string   `json:"published"`
}

func MapDBtoResponseModel(articleDB *ArticleDB) *ArticleResponse {
	return &ArticleResponse{
		Id:          articleDB.Id,
		TeamId:      articleDB.TeamId,
		OptaMatchId: articleDB.OptaMatchId,
		Title:       articleDB.Title,
		Type:        articleDB.Type,
		Teaser:      articleDB.Teaser,
		Content:     articleDB.Content,
		Url:         articleDB.Url,
		ImageUrl:    articleDB.ImageUrl,
		GalleryUrls: articleDB.GalleryUrls,
		VideoUrl:    articleDB.VideoUrl,
		Published:   articleDB.Published,
	}
}

type ArticleDB struct {
	Id          string   `bson:"id"`
	TeamId      string   `bson:"teamId"`
	OptaMatchId *string  `bson:"optaMatchId"`
	Title       string   `bson:"title"`
	Type        []string `bson:"type"`
	Teaser      *string  `bson:"teaser"`
	Content     string   `bson:"content"`
	Url         string   `bson:"url"`
	ImageUrl    string   `bson:"imageUrl"`
	GalleryUrls string   `bson:"galleryUrls"`
	VideoUrl    string   `bson:"videoUrl"`
	Published   string   `bson:"published"`
}
