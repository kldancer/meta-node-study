package __4_golang_gin

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"time"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Email    string `gorm:"unique;not null"`
}

type Post struct {
	gorm.Model
	Title    string `gorm:"not null"`
	Content  string `gorm:"not null"`
	UserID   uint
	User     User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Comments []Comment `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type PostInfo struct {
	Title       string
	Content     string
	UserID      uint
	CommentInfo []CommentInfo
}

type Comment struct {
	gorm.Model
	Content string `gorm:"not null"`
	UserID  uint
	User    User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	PostID  uint
	Post    Post `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type CommentInfo struct {
	Content string `gorm:"not null"`
	UserID  uint
}

var (
	db        *gorm.DB
	jwtSecret = []byte("your_secret_key")
)

func initDB() {
	var err error
	db, err = gorm.Open(sqlite.Open("blog.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}
	// 自动迁移模型
	if err := db.AutoMigrate(&User{}, &Post{}, &Comment{}); err != nil {
		panic("failed to migrate database: " + err.Error())
	}
}

// ErrorResponse 统一错误返回
func ErrorResponse(c *gin.Context, code int, err error) {
	c.JSON(code, gin.H{"error": err.Error()})
	c.Abort()
}

func Register(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = string(hashedPassword)

	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

// Login 用户登录，返回 JWT
func Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, err)
		return
	}
	var user User
	if err := db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		ErrorResponse(c, http.StatusUnauthorized, errors.New("invalid credentials"))
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		ErrorResponse(c, http.StatusUnauthorized, errors.New("invalid credentials"))
		return
	}
	// 生成 JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

// AuthMiddleware 验证 JWT
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			ErrorResponse(c, http.StatusUnauthorized, errors.New("authorization header required"))
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			ErrorResponse(c, http.StatusUnauthorized, errors.New("authorization header format must be Bearer {token}"))
			return
		}
		tokenStr := parts[1]
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			ErrorResponse(c, http.StatusUnauthorized, errors.New("invalid token"))
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			ErrorResponse(c, http.StatusUnauthorized, errors.New("invalid token claims"))
			return
		}
		// 保存当前用户 ID 到 Context
		c.Set("user_id", uint(claims["user_id"].(float64)))
		c.Next()
	}
}

// CreatePost 发表文章
func CreatePost(c *gin.Context) {
	var req struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, err)
		return
	}
	userID := c.GetUint("user_id")
	post := Post{
		Title:   req.Title,
		Content: req.Content,
		UserID:  userID,
	}
	if err := db.Create(&post).Error; err != nil {
		ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusCreated, post)
}

// GetPosts 获取所有文章
func GetPosts(c *gin.Context) {
	var posts []Post
	if err := db.Preload("User").Find(&posts).Error; err != nil {
		ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	var postInfos []PostInfo
	for _, post := range posts {
		pi := parsePostInfo(post)
		postInfos = append(postInfos, pi)
	}

	c.JSON(http.StatusOK, postInfos)
}

// GetPost 获取单个文章及其评论
func GetPost(c *gin.Context) {
	id := c.Param("id")
	var post Post
	if err := db.Preload("User").Preload("Comments.User").First(&post, id).Error; err != nil {
		ErrorResponse(c, http.StatusNotFound, errors.New("post not found"))
		return
	}
	c.JSON(http.StatusOK, parsePostInfo(post))
}

func parsePostInfo(p Post) PostInfo {
	return PostInfo{
		Title:       p.Title,
		Content:     p.Content,
		UserID:      p.UserID,
		CommentInfo: nil,
	}
}

// UpdatePost 更新文章（仅作者）
func UpdatePost(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetUint("user_id")
	var post Post
	if err := db.First(&post, id).Error; err != nil {
		ErrorResponse(c, http.StatusNotFound, errors.New("post not found"))
		return
	}
	if post.UserID != userID {
		ErrorResponse(c, http.StatusForbidden, errors.New("you are not the author"))
		return
	}
	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, err)
		return
	}
	post.Title = req.Title
	post.Content = req.Content
	if err := db.Save(&post).Error; err != nil {
		ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, post)
}

// DeletePost 删除文章（仅作者）
func DeletePost(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetUint("user_id")
	var post Post
	if err := db.First(&post, id).Error; err != nil {
		ErrorResponse(c, http.StatusNotFound, errors.New("post not found"))
		return
	}
	if post.UserID != userID {
		ErrorResponse(c, http.StatusForbidden, errors.New("you are not the author"))
		return
	}
	if err := db.Delete(&post).Error; err != nil {
		ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "post deleted"})
}

// CreateComment 为文章添加评论
func CreateComment(c *gin.Context) {
	postID := c.Param("id")
	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, err)
		return
	}
	userID := c.GetUint("user_id")
	comment := Comment{
		Content: req.Content,
		UserID:  userID,
		PostID:  0,
	}
	// 验证 post 是否存在
	var post Post
	if err := db.First(&post, postID).Error; err != nil {
		ErrorResponse(c, http.StatusNotFound, errors.New("post not found"))
		return
	}
	comment.PostID = post.ID
	if err := db.Create(&comment).Error; err != nil {
		ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusCreated, comment)
}

// GetComments 获取某篇文章的所有评论
func GetComments(c *gin.Context) {
	postID := c.Param("id")
	var comments []Comment
	if err := db.Preload("User").
		Where("post_id = ?", postID).
		Find(&comments).Error; err != nil {
		ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, comments)
}

func ginBlog() {
	initDB()

	r := gin.Default() // 内置 Logger 和 Recovery

	// 公共路由
	r.POST("/register", Register)
	r.POST("/login", Login)
	r.GET("/posts", GetPosts)
	r.GET("/posts/:id", GetPost)
	r.GET("/posts/:id/comments", GetComments)

	// 需要认证的路由
	auth := r.Group("/")
	auth.Use(AuthMiddleware())
	{
		auth.POST("/posts", CreatePost)
		auth.PUT("/posts/:id", UpdatePost)
		auth.DELETE("/posts/:id", DeletePost)
		auth.POST("/posts/:id/comments", CreateComment)
	}

	// 启动
	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
