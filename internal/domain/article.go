package domain

type Article struct {
	Id          string   `json:"id" bson:"id"`
	TeamId      string   `json:"teamId" bson:"teamId"`
	OptaMatchId *string  `json:"optaMatchId" bson:"optaMatchId"`
	Title       string   `json:"title" bson:"title"`
	Type        []string `json:"type" bson:"type"`
	Teaser      *string  `json:"teaser" bson:"teaser"`
	Content     string   `json:"content" bson:"content"`
	Url         string   `json:"url" bson:"url"`
	ImageUrl    string   `json:"imageUrl" bson:"imageUrl"`
	GalleryUrls string   `json:"galleryUrls" bson:"galleryUrls"`
	VideoUrl    string   `json:"videoUrl" bson:"videoUrl"`
	Published   string   `json:"published" bson:"published"`
}
