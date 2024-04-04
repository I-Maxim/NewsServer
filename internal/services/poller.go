package services

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"log"
	"net/http"
	"newsServer/internal/domain"
	"newsServer/internal/repo"
	"time"
)

type Poller struct {
	Period time.Duration
	Count  uint
	repo   repo.Repository
}

func NewPoller(period time.Duration, count uint, repo repo.Repository) Poller {
	return Poller{
		Period: period,
		Count:  count,
		repo:   repo,
	}
}

func (p *Poller) Start(ctx context.Context) {
	ticker := time.NewTicker(p.Period)
	defer ticker.Stop()

	fetchListUrl := fmt.Sprintf("https://www.htafc.com/api/incrowd/getnewlistinformation?count=%d", p.Count)
	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			resp, err := http.Get(fetchListUrl)
			if err != nil {
				log.Println(err)
				continue
			}

			body, err := io.ReadAll(resp.Body)
			if len(body) == 0 || err != nil {
				log.Println(err)
				continue
			}

			list := domain.NewListInformation{}
			err = xml.Unmarshal(body, &list)
			if err := resp.Body.Close(); err != nil {
				log.Println(err)
			}
			if err != nil {
				log.Println(err)
				continue
			}

			articles := []*domain.ArticleDB{}
			for _, v := range list.NewsletterNewsItems.NewsletterNewsItem {
				timeoutCtx, cancel := context.WithTimeout(ctx, time.Second*10)
				existing, err := p.repo.Load(timeoutCtx, v.NewsArticleID)
				cancel()
				if existing != nil || (err != nil && !errors.Is(err, mongo.ErrNoDocuments)) {
					continue
				}

				url := fmt.Sprintf("https://www.htafc.com/api/incrowd/getnewsarticleinformation?id=%s", v.NewsArticleID)
				resp, err := http.Get(url)
				if err != nil {
					log.Println(err)
					continue
				}

				body, err := io.ReadAll(resp.Body)
				if len(body) == 0 || err != nil {
					log.Println(err)
					continue
				}

				extArticle := domain.NewsArticleInformation{}
				err = xml.Unmarshal(body, &extArticle)
				if err := resp.Body.Close(); err != nil {
					log.Println(err)
				}
				if err != nil {
					log.Println(err)
					continue
				}

				articles = append(articles, &domain.ArticleDB{
					//TeamId:      extArticle.NewsArticle.,
					Id:          extArticle.NewsArticle.NewsArticleID,
					OptaMatchId: extArticle.NewsArticle.OptaMatchId,
					Title:       extArticle.NewsArticle.Title,
					Type:        []string{extArticle.NewsArticle.Taxonomies},
					Teaser:      extArticle.NewsArticle.TeaserText,
					Content:     extArticle.NewsArticle.BodyText,
					Url:         extArticle.NewsArticle.ArticleURL,
					ImageUrl:    extArticle.NewsArticle.ThumbnailImageURL,
					GalleryUrls: extArticle.NewsArticle.GalleryImageURLs,
					VideoUrl:    extArticle.NewsArticle.VideoURL,
					Published:   extArticle.NewsArticle.PublishDate,
				})
			}
			if len(articles) == 0 {
				continue
			}
			timeoutCtx, cancel := context.WithTimeout(ctx, time.Second*10)
			if err = p.repo.Save(timeoutCtx, articles...); err != nil {
				log.Println(err)
			}
			cancel()
			log.Printf("saved %d articles\n", len(articles))
		}
	}
}
