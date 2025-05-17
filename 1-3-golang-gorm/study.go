package __3_golang_gorm

import (
	"github.com/jmoiron/sqlx"
	"gorm.io/gorm"
)

/*
Sqlx入门,
题目1：使用SQL扩展库进行查询,
假设你已经使用Sqlx连接到一个数据库，并且有一个 employees 表，包含字段 id 、 name 、 department 、 salary 。
要求 ：
编写Go代码，使用Sqlx查询 employees 表中所有部门为 "技术部" 的员工信息，并将结果映射到一个自定义的 Employee 结构体切片中。,
编写Go代码，使用Sqlx查询 employees 表中工资最高的员工信息，并将结果映射到一个 Employee 结构体中。
*/

type Employee struct {
	ID         int     `db:"id"`
	Name       string  `db:"name"`
	Department string  `db:"department"`
	Salary     float64 `db:"salary"`
}

func QueryTechEmployees(db *sqlx.DB) ([]Employee, error) {
	var employees []Employee
	err := db.Select(&employees, "SELECT id, name, department, salary FROM employees WHERE department = ?", "技术部")
	if err != nil {
		return nil, err
	}
	return employees, nil
}

func QueryHighestPaidEmployee(db *sqlx.DB) (Employee, error) {
	var employee Employee
	err := db.Get(&employee, "SELECT id, name, department, salary FROM employees ORDER BY salary DESC LIMIT 1")
	if err != nil {
		return Employee{}, err
	}
	return employee, nil
}

/*
题目2：实现类型安全映射,
假设有一个 books 表，包含字段 id 、 title 、 author 、 price 。
要求 ：
定义一个 Book 结构体，包含与 books 表对应的字段。,
编写Go代码，使用Sqlx执行一个复杂的查询，例如查询价格大于 50 元的书籍，并将结果映射到 Book 结构体切片中，确保类型安全。
*/

type Book struct {
	ID     int     `db:"id"`
	Title  string  `db:"title"`
	Author string  `db:"author"`
	Price  float64 `db:"price"`
}

func QueryBooksGreaterThan50(db *sqlx.DB) ([]Book, error) {
	var books []Book
	err := db.Select(&books, "SELECT id, title, author, price FROM books WHERE price > ?", 50)
	if err != nil {
		return nil, err
	}
	return books, nil
}

/*


 */

type User struct {
	ID        int    `gorm:"primaryKey"`
	Name      string `gorm:"column:name"`
	Email     string `gorm:"column:email"`
	Password  string `gorm:"column:password"`
	Posts     []Post `gorm:"foreignKey:UserID"` // 一对多关系，外键关联到Post表的UserID字段
	PostCount int    `gorm:"column:postCount"`
}

type Post struct {
	ID            int       `gorm:"primaryKey"`
	Title         string    `gorm:"column:title"`
	Content       string    `gorm:"column:content"`
	UserID        int       `gorm:"column:user_id"`    // 外键，指向User表的ID
	User          User      `gorm:"foreignkey:UserID"` // 关联User模型
	Comments      []Comment `gorm:"foreignKey:PostID"` // 一对多关系，外键关联到Comment表的PostID字段
	CommentStatus string    `gorm:"column:comment_status"`
}

type Comment struct {
	ID      int    `gorm:"primaryKey"`
	Content string `gorm:"column:content"`
	PostID  int    `gorm:"column:post_id"`    // 外键，指向Post表的ID
	Post    Post   `gorm:"foreignkey:PostID"` // 关联Post模型
}

func CreateTables(db *gorm.DB) error {
	err := db.AutoMigrate(&User{}, &Post{}, &Comment{})
	if err != nil {
		return err
	}
	return nil
}

/*
题目2：关联查询,
基于上述博客系统的模型定义。
要求 ：
编写Go代码，使用Gorm查询某个用户发布的所有文章及其对应的评论信息。,
编写Go代码，使用Gorm查询评论数量最多的文章信息。
*/

// QueryUserPosts 查询某个用户发布的所有文章及其对应的评论信息。,
func QueryUserPosts(db *gorm.DB, userID int) ([]Post, error) {
	var posts []Post
	err := db.Where("user_id = ?", userID).Preload("Comments").Find(&posts).Error
	if err != nil {
		return nil, err
	}
	return posts, nil
}

// QueryPostWithMostComments 查询评论数量最多的文章信息。
func QueryPostWithMostComments(db *gorm.DB) (*Post, error) {
	var post Post
	err := db.Table("comments").
		Select("post_id, COUNT(*) as comment_count").
		Group("post_id").
		Order("comment_count DESC").
		Limit(1).
		Scan(&post).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

/*
继续使用博客系统的模型。
要求 ：
为 Post 模型添加一个钩子函数，在文章创建时自动更新用户的文章数量统计字段。,
为 Comment 模型添加一个钩子函数，在评论删除时检查文章的评论数量，如果评论数量为 0，则更新文章的评论状态为 "无评论"。
*/

// 为Post模型添加钩子函数，在文章创建时自动更新用户的文章数量统计字段
func (p *Post) AfterCreate(tx *gorm.DB) error {
	// 获取关联的用户
	var user User
	if err := tx.Where("id = ?", p.UserID).First(&user).Error; err != nil {
		return err
	}

	// 更新用户的发帖数量
	user.PostCount++
	return tx.Save(&user).Error
}

// 为Comment模型添加钩子函数，在评论删除时检查文章的评论数量
func (c *Comment) AfterDelete(tx *gorm.DB) error {
	// 获取关联的文章
	var post Post
	if err := tx.Where("id = ?", c.PostID).First(&post).Error; err != nil {
		return err
	}

	// 查询当前文章的评论数量
	var count int64
	if err := tx.Model(&Comment{}).Where("post_id = ?", post.ID).Count(&count).Error; err != nil {
		return err
	}

	// 如果评论数量为0，更新文章的评论状态
	if count == 0 {
		post.CommentStatus = "无评论"
		return tx.Save(&post).Error
	}

	return nil
}
