package entity

import (
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gen/field"
	"gorm.io/gorm"
)

type PageAndOrder struct {
	PageData
	// 排序字段,这是个接口，不能直接赋值，需要使用field.NewString("", order).Desc()
	Order field.Expr `json:"-" form:"-" copier:"-"`
	// 提供预加载关联数据支持
	Preload []field.RelationField `json:"-" form:"-" copier:"-"`
}

// 设置排序字段
// order: 排序字段
// desc: 是否降序
// 例如:	SetOrderBy("created_at", true)
func (p *PageAndOrder) SetOrderBy(order string, desc bool) {
	if desc {
		p.Order = field.NewString("", order).Desc()
	} else {
		p.Order = field.NewString("", order).Asc()
	}
}

type PageData struct {
	Page     int `json:"page" form:"page"`
	PageSize int `json:"page_size" form:"page_size"`
	// 游标值，格式如 "EventTime|id"
	Cursor string `json:"cursor" form:"cursor"`
}

// 判断是否重复键冲突错误
// 不能使用gorm.ErrDuplicatedKey，不同的数据库可能会返回不同的错误信息
func IsDuplicateKey(err error) bool {
	var pgErr *pgconn.PgError
	if gerror.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}

// IsForeignKeyViolation 判断错误是否为外键约束冲突错误。
// 不同数据库可能返回不同的错误信息，当前实现仅处理 PostgreSQL 数据库。
// 参数 err 为需要判断的错误对象。
// 返回值为布尔类型，若为外键约束冲突错误则返回 true，否则返回 false。
func IsForeignKeyViolation(err error) bool {
	var pgErr *pgconn.PgError
	if gerror.As(err, &pgErr) {
		return pgErr.Code == "23503"
	}
	return false
}

// IsSiteIDFKeyViolation 判断错误是否为 site_id_fkey 外键约束冲突错误。
// 不同数据库可能返回不同的错误信息，当前实现仅处理 PostgreSQL 数据库。
// 参数 err 为需要判断的错误对象。
// 返回值为布尔类型，若为 site_id_fkey 外键约束冲突错误则返回 true，否则返回 false。
func IsSiteIDFKeyViolation(err error) bool {
	var pgErr *pgconn.PgError
	if gerror.As(err, &pgErr) {
		// 检查错误代码是否为外键冲突，并且错误信息包含 site_id_fkey
		return pgErr.Code == "23503" && strings.Contains(pgErr.Message, "site_id_fkey")
	}
	return false
}

func IsNotFound(err error) bool {
	return gerror.Is(err, gorm.ErrRecordNotFound)
}

// 生成忽略大小写的表达式
func ILike(f field.String, s string) field.Expr {
	return field.NewUnsafeFieldRaw(f.ColumnName().String()+" ilike ?", "%"+s+"%")
}
