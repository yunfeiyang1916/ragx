// Code generated; DO NOT EDIT

package repo

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/gogf/gf/v2/errors/gerror"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
	"ragx/app/internal/biz"
	"ragx/app/internal/biz/entity"
	"ragx/app/internal/biz/query"
	"ragx/app/pkg/cache/redis"
)

type KnowledgeBaseRepo struct {
	Data      biz.Data
	DB        *gorm.DB
	Rdb       *redis.Client
	Log       *log.Helper
	GormQuery *query.Query
}

func (d *KnowledgeBaseRepo) Query() *query.Query { return d.GormQuery }

// 批量创建，支持事务
func (d *KnowledgeBaseRepo) BatchCreate(ctx context.Context, list []*entity.KnowledgeBase, tx ...*query.Query) ([]*entity.KnowledgeBase, error) {
	q := d.GormQuery
	if len(tx) > 0 && tx[0] != nil {
		q = tx[0]
	}
	err := q.KnowledgeBase.WithContext(ctx).Create(list...)
	if err != nil {
		return nil, gerror.Wrap(err, "")
	}
	return list, err
}

// 创建，支持事务
func (d *KnowledgeBaseRepo) Create(ctx context.Context, obj *entity.KnowledgeBase, tx ...*query.Query) (*entity.KnowledgeBase, error) {
	q := d.GormQuery
	if len(tx) > 0 && tx[0] != nil {
		q = tx[0]
	}
	err := q.KnowledgeBase.WithContext(ctx).Create(obj)
	if err != nil {
		return nil, gerror.Wrap(err, "")
	}
	return obj, err
}

// 保存全部字段，支持事务
func (d *KnowledgeBaseRepo) Save(ctx context.Context, obj *entity.KnowledgeBase, tx ...*query.Query) (int64, error) {
	qu := d.GormQuery
	if len(tx) > 0 && tx[0] != nil {
		qu = tx[0]
	}
	q := qu.KnowledgeBase
	columns := []field.Expr{q.Name, q.Description, q.Category, q.Status, q.CreateTime, q.UpdateTime}
	res, err := q.WithContext(ctx).Where(q.ID.Eq(obj.ID)).Select(columns...).UpdateColumns(obj)
	if err != nil {
		return 0, gerror.Wrap(err, "")
	}
	return res.RowsAffected, nil
}

// 仅更新指定字段，支持表达式，表达式不能为空
func (d *KnowledgeBaseRepo) Update(ctx context.Context, obj *entity.KnowledgeBase, columns ...field.Expr) (int64, error) {
	if len(columns) == 0 {
		return 0, gerror.New("no columns to update")
	}
	q := d.GormQuery.KnowledgeBase
	res, err := q.WithContext(ctx).Where(q.ID.Eq(obj.ID)).Select(columns...).UpdateColumns(obj)
	if err != nil {
		return 0, gerror.Wrap(err, "")
	}
	return res.RowsAffected, nil
}

// 支持事务，仅更新指定字段，支持表达式，表达式不能为空
func (d *KnowledgeBaseRepo) UpdateWithTx(ctx context.Context, tx *query.Query, obj *entity.KnowledgeBase, columns ...field.Expr) (int64, error) {
	if len(columns) == 0 {
		return 0, gerror.New("no columns to update")
	}
	q := tx.KnowledgeBase
	res, err := q.WithContext(ctx).Where(q.ID.Eq(obj.ID)).Select(columns...).UpdateColumns(obj)
	if err != nil {
		return 0, gerror.Wrap(err, "")
	}
	return res.RowsAffected, nil
}

// 删除，支持事务
func (d *KnowledgeBaseRepo) Delete(ctx context.Context, id int64, tx ...*query.Query) (int64, error) {
	qu := d.GormQuery
	if len(tx) > 0 && tx[0] != nil {
		qu = tx[0]
	}
	q := qu.KnowledgeBase
	res, err := q.WithContext(ctx).Where(q.ID.Eq(id)).Delete()
	if err != nil {
		return 0, gerror.Wrap(err, "")
	}
	return res.RowsAffected, nil
}

func (d *KnowledgeBaseRepo) DeleteByConditions(ctx context.Context, conditions ...gen.Condition) (int64, error) {
	if len(conditions) == 0 {
		return 0, gerror.New("no conditions to delete")
	}
	q := d.GormQuery.KnowledgeBase
	res, err := q.WithContext(ctx).Where(conditions...).Delete()
	if err != nil {
		return 0, gerror.Wrap(err, "")
	}
	return res.RowsAffected, nil
}

func (d *KnowledgeBaseRepo) DeleteByConditionsWithTx(ctx context.Context, tx *query.Query, conditions ...gen.Condition) (int64, error) {
	if len(conditions) == 0 {
		return 0, gerror.New("no conditions to delete")
	}
	q := tx.KnowledgeBase
	res, err := q.WithContext(ctx).Where(conditions...).Delete()
	if err != nil {
		return 0, gerror.Wrap(err, "")
	}
	return res.RowsAffected, nil
}

func (d *KnowledgeBaseRepo) Get(ctx context.Context, id int64, preload ...field.RelationField) (*entity.KnowledgeBase, error) {
	q := d.GormQuery.KnowledgeBase
	obj, err := q.WithContext(ctx).Where(q.ID.Eq(id)).Preload(preload...).First()
	if err != nil {
		if gerror.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, gerror.Wrap(err, "")
	}
	return obj, nil
}

func (d *KnowledgeBaseRepo) GetByConditions(ctx context.Context, conditions ...gen.Condition) (*entity.KnowledgeBase, error) {
	if len(conditions) == 0 {
		return nil, gerror.New("no conditions to delete")
	}
	q := d.GormQuery.KnowledgeBase
	obj, err := q.WithContext(ctx).Where(conditions...).First()
	if err != nil {
		if gerror.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, gerror.Wrap(err, "")
	}
	return obj, nil
}

// 支持预加载
func (d *KnowledgeBaseRepo) GetByConditionsWithPreload(ctx context.Context, preload []field.RelationField, conditions ...gen.Condition) (*entity.KnowledgeBase, error) {
	if len(conditions) == 0 {
		return nil, gerror.New("no conditions to delete")
	}
	q := d.GormQuery.KnowledgeBase
	obj, err := q.WithContext(ctx).Where(conditions...).Preload(preload...).First()
	if err != nil {
		if gerror.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, gerror.Wrap(err, "")
	}
	return obj, nil
}

func (d *KnowledgeBaseRepo) List(ctx context.Context, page *entity.PageAndOrder, conditions ...gen.Condition) ([]*entity.KnowledgeBase, int64, error) {
	q := d.GormQuery.KnowledgeBase
	where := q.WithContext(ctx).Where(conditions...)
	count, err := where.Count()
	if err != nil {
		return nil, 0, gerror.Wrap(err, "")
	}
	if count == 0 {
		return nil, 0, nil
	}
	if page != nil {
		if page.Page <= 0 {
			page.Page = 1
		}
		if page.PageSize <= 0 {
			page.PageSize = 10
		}
		if page.Order != nil {
			where = where.Order(page.Order)
		} else {
			where = where.Order(q.ID.Desc())
		}
		where = where.Preload(page.Preload...).Offset((page.Page - 1) * page.PageSize).Limit(page.PageSize)
	}
	list, err := where.Find()
	if err != nil {
		return nil, 0, gerror.Wrap(err, "")
	}
	return list, count, nil
}

// 只需要列表，不需要总数
func (d *KnowledgeBaseRepo) ListWithoutCount(ctx context.Context, page *entity.PageAndOrder, conditions ...gen.Condition) ([]*entity.KnowledgeBase, error) {
	q := d.GormQuery.KnowledgeBase
	where := q.WithContext(ctx).Where(conditions...)
	if page != nil {
		if page.Page <= 0 {
			page.Page = 1
		}
		if page.PageSize <= 0 {
			page.PageSize = 10
		}
		if page.Order != nil {
			where = where.Order(page.Order)
		} else {
			where = where.Order(q.ID.Desc())
		}
		where = where.Preload(page.Preload...).Offset((page.Page - 1) * page.PageSize).Limit(page.PageSize)
	}
	list, err := where.Find()
	if err != nil {
		return nil, gerror.Wrap(err, "")
	}
	return list, nil
}

func (d *KnowledgeBaseRepo) ListAll(ctx context.Context, conditions ...gen.Condition) ([]*entity.KnowledgeBase, error) {
	q := d.GormQuery.KnowledgeBase
	list, err := q.WithContext(ctx).Where(conditions...).Find()
	if err != nil {
		return nil, gerror.Wrap(err, "")
	}
	return list, nil
}

// 支持预加载
func (d *KnowledgeBaseRepo) ListAllWithPreload(ctx context.Context, preload []field.RelationField, conditions ...gen.Condition) ([]*entity.KnowledgeBase, error) {
	q := d.GormQuery.KnowledgeBase
	list, err := q.WithContext(ctx).Where(conditions...).Preload(preload...).Find()
	if err != nil {
		return nil, gerror.Wrap(err, "")
	}
	return list, nil
}

func (d *KnowledgeBaseRepo) Count(ctx context.Context, conditions ...gen.Condition) (int64, error) {
	q := d.GormQuery.KnowledgeBase
	count, err := q.WithContext(ctx).Where(conditions...).Count()
	if err != nil {
		return 0, gerror.Wrap(err, "")
	}
	return count, nil
}
