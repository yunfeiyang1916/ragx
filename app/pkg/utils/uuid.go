package utils

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/google/uuid"
	"strings"
)

func NewUUID() string {
	// v7 基于时间的 UUID，使用当前时间戳生成，保证了生成的 UUID 是唯一的，并且可以根据时间戳进行排序。
	// 理论上不会出错，如果真的出错就使用 v4 生成
	id, err := uuid.NewV7()
	if err != nil {
		log.Errorf("uuid.NewV7() error: %v", gerror.Wrap(err, ""))
		return uuid.New().String()
	}
	return id.String()
}

// UniqueID 生成唯一ID，去掉uuid中的短横线
func UniqueID() string {
	return strings.Replace(NewUUID(), "-", "", -1)
}
