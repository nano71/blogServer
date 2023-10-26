package api

import (
	"blogServer/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log/slog"
	"time"
)

type Article struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Content     string `json:"content"`
	UpdateTime  string `json:"updateTime"`
	CreateTime  string `json:"createTime"`
	Tags        string `json:"tags"`
	CoverImage  string `json:"coverImage"`
	ReadCount   int    `json:"readCount"`
}

func GetArticleList(c *gin.Context) {
	p := &struct {
		Limit int `json:"limit" binding:"required"`
		Page  int `json:"page"`
	}{}
	preprocess(c, p, func(db *gorm.DB) {
		articles := &[]struct {
			Id          int       `json:"id"`
			Title       string    `json:"title"`
			Description string    `json:"description"`
			CreateTime  time.Time `json:"createTime"`
			Tags        string    `json:"tags"`
			CoverImage  string    `json:"coverImage"`
		}{}
		slog.Info("", p)
		var total int64
		db.Model(Article{}).Count(&total).Limit(p.Limit).Offset(p.Page * p.Limit).Order("create_time desc").Find(articles)
		data := gin.H{
			"total": total,
			"list":  articles,
		}
		response.Success(c, data)
	})

}

func SearchArticles(c *gin.Context) {
	p := &struct {
		Query string `json:"query" binding:"required"`
		Limit int    `json:"limit" binding:"required"`
		Page  int    `json:"page"`
	}{}

	preprocess(c, p, func(db *gorm.DB) {
		articles := &[]Article{}
		search := "%" + p.Query + "%"
		where := db.Model(Article{}).Where("title like ?", search).Or("content like ?", search)
		var total int64
		where.Count(&total).Limit(p.Limit).Offset(p.Page * p.Limit).Order("create_time desc").Find(articles)
		data := gin.H{
			"total": total,
			"list":  articles,
		}
		response.Success(c, data)
	})

}

func SearchArticlesByTag(c *gin.Context) {
	p := &struct {
		Tag   string `json:"tag" binding:"required"`
		Limit int    `json:"limit" binding:"required"`
		Page  int    `json:"page"`
	}{}

	preprocess(c, p, func(db *gorm.DB) {
		articles := &[]Article{}
		where := db.Model(Article{}).Where("tags like ?", "%"+p.Tag+"%")
		var total int64
		where.Count(&total).Limit(p.Limit).Offset(p.Page * p.Limit).Order("create_time desc").Find(articles)
		data := gin.H{
			"total": total,
			"list":  articles,
		}
		response.Success(c, data)
	})
}

func GetArticleContent(c *gin.Context) {
	p := &struct {
		ArticleId int `json:"articleId" binding:"required"`
	}{}
	preprocess(c, p, func(db *gorm.DB) {
		article := &Article{
			Id: p.ArticleId,
		}
		db.First(article)
		response.Success(c, article)
	})
}

func PublishArticle(c *gin.Context) {
	p := &struct {
		Title       string `json:"title" binding:"required"`
		Content     string `json:"content" binding:"required"`
		Description string `json:"description" binding:"required"`
		CreateTime  string `json:"createTime" binding:"required"`
		CoverImage  string `json:"coverImage"`
		Tags        string `json:"tags"`
	}{}
	preprocess(c, p, func(db *gorm.DB) {
		article := &Article{
			Title:       p.Title,
			Content:     p.Content,
			Description: p.Description,
			CreateTime:  p.CreateTime,
			CoverImage:  p.CoverImage,
			Tags:        p.Tags,
		}
		result := db.Omit("UpdateTime").Create(article)
		if result.RowsAffected == 1 {
			response.Success(c, true)
			shouldFetchData = true
		} else {
			response.Fail(c, "文章发布失败")
		}
	})
}
