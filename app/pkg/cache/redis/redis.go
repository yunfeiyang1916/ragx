package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/redis/go-redis/v9"
)

func NewClient(rdb redis.UniversalClient) *Client {
	return &Client{
		rdb: rdb,
	}
}

type Client struct {
	rdb redis.UniversalClient
}

// Exists 检查指定的键是否存在于 Redis 中。
// 参数 ctx 为上下文，可用于控制请求的超时和取消等操作。
// 参数 key 为要检查的 Redis 键。
// 返回值 bool 表示键是否存在，error 表示操作过程中可能出现的错误。
func (r *Client) Exists(ctx context.Context, key string) (bool, error) {
	intCmd := r.rdb.Exists(ctx, key)
	if intCmd == nil {
		return false, gerror.New("redisDao Exists cmd is nil")
	}
	return intCmd.Val() > 0, nil
}

// Get 从 Redis 中获取指定键的值。
// 参数 ctx 为上下文，可用于控制请求的超时和取消等操作。
// 参数 key 为要获取值的 Redis 键。
// 返回值 string 为键对应的值，error 表示操作过程中可能出现的错误。若键不存在，返回空字符串和 nil。
func (r *Client) Get(ctx context.Context, key string) (string, error) {
	strCmd := r.rdb.Get(ctx, key)
	if strCmd == nil {
		return "", gerror.New("redisDao Get cmd is nil")
	}
	val, err := strCmd.Result()
	if err != nil {
		if gerror.Is(err, redis.Nil) {
			return "", nil
		}
		return "", gerror.Wrap(err, fmt.Sprintf("get key %s from redis failed", key))
	}
	return val, nil
}

// GetDel 从 Redis 中获取指定键的值，并在获取后删除该键。
// 参数 ctx 为上下文，可用于控制请求的超时和取消等操作。
// 参数 key 为要操作的 Redis 键。
// 返回值 string 为键对应的值，error 表示操作过程中可能出现的错误。若键不存在，返回空字符串和 nil。
func (r *Client) GetDel(ctx context.Context, key string) (string, error) {
	strCmd := r.rdb.GetDel(ctx, key)
	if strCmd == nil {
		return "", gerror.New("redisDao GetDel cmd is nil")
	}
	val, err := strCmd.Result()
	if err != nil {
		if gerror.Is(err, redis.Nil) {
			return "", nil
		}
		return "", gerror.Wrap(err, fmt.Sprintf("GetDel key %s from redis failed", key))
	}
	return val, nil
}

// SetEx 设置指定键的值，并为该键设置过期时间。
// 参数 ctx 为上下文，可用于控制请求的超时和取消等操作。
// 参数 key 为要设置值的 Redis 键。
// 参数 val 为要设置的值。
// 参数 expire 为键的过期时间。
// 返回值 error 表示操作过程中可能出现的错误。
func (r *Client) SetEx(ctx context.Context, key string, val string, expire time.Duration) error {
	cmd := r.rdb.SetEx(ctx, key, val, expire)
	if cmd == nil {
		return gerror.New("redisDao SetEx cmd is nil")
	}

	_, err := cmd.Result()
	if err != nil && gerror.Is(err, redis.Nil) {
		return err
	}

	return nil
}

// SetNX 仅当指定键不存在时，设置该键的值，并为其设置过期时间。
// 参数 ctx 为上下文，可用于控制请求的超时和取消等操作。
// 参数 key 为要设置值的 Redis 键。
// 参数 val 为要设置的值。
// 参数 expire 为键的过期时间。
// 返回值 bool 表示键是否成功设置，error 表示操作过程中可能出现的错误。
func (r *Client) SetNX(ctx context.Context, key string, val string, expire time.Duration) (bool, error) {
	cmd := r.rdb.SetNX(ctx, key, val, expire)
	if cmd == nil {
		return false, gerror.New("redisDao SetEx cmd is nil")
	}
	return cmd.Result()
}

// Del 从 Redis 中删除指定的键。
// 参数 ctx 为上下文，可用于控制请求的超时和取消等操作。
// 参数 key 为要删除的 Redis 键。
// 返回值 error 表示操作过程中可能出现的错误。
func (r *Client) Del(ctx context.Context, key string) error {
	cmd := r.rdb.Del(ctx, key)
	if cmd == nil {
		return gerror.New("redisDao Del cmd is nil")
	}

	_, err := cmd.Result()
	if err != nil && gerror.Is(err, redis.Nil) {
		return err
	}

	return nil
}

// Mget 从 Redis 中批量获取指定键的值。
// 参数 ctx 为上下文，可用于控制请求的超时和取消等操作。
// 参数 keys 为要获取值的 Redis 键列表。
// 返回值 []interface{} 为获取到的值列表，error 表示操作过程中可能出现的错误。
func (r *Client) Mget(ctx context.Context, keys []string) ([]interface{}, error) {
	vals, err := r.rdb.MGet(ctx, keys...).Result()
	if err != nil && gerror.Is(err, redis.Nil) {
		return nil, err
	}
	return vals, nil
}

// HSet 设置 Redis 哈希表中指定键的一个或多个字段的值。
// 参数 ctx 为上下文，可用于控制请求的超时和取消等操作。
// 参数 key 为 Redis 哈希表的键名。
// 参数 val 为要设置的字段名和对应值的映射。
// 返回值 error 表示操作过程中可能出现的错误。
func (r *Client) HSet(ctx context.Context, key string, val map[string]interface{}) error {
	cmd := r.rdb.HSet(ctx, key, val)
	if err := cmd.Err(); err != nil && gerror.Is(err, redis.Nil) {
		return err
	}

	return nil
}

// HGetAll 获取 Redis 哈希表中指定键的所有字段和对应的值。
// 参数 ctx 为上下文，可用于控制请求的超时和取消等操作。
// 参数 key 为 Redis 哈希表的键名。
// 返回值 map[string]string 为哈希表中所有字段名和对应值的映射，error 表示操作过程中可能出现的错误。
func (r *Client) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	res, err := r.rdb.HGetAll(ctx, key).Result()
	if err != nil && gerror.Is(err, redis.Nil) {
		return nil, err
	}

	return res, nil
}

// HDel 删除 Redis 哈希表中指定键的一个或多个字段。
// 参数 ctx 为上下文，可用于控制请求的超时和取消等操作。
// 参数 key 为 Redis 哈希表的键名。
// 参数 field 为要删除的字段名列表。
// 返回值 error 表示操作过程中可能出现的错误。
func (r *Client) HDel(ctx context.Context, key string, field ...string) error {
	_, err := r.rdb.HDel(ctx, key, field...).Result()
	if err != nil && gerror.Is(err, redis.Nil) {
		return err
	}

	return nil
}

// Expire 为 Redis 中指定的键设置过期时间。
// 参数 ctx 为上下文，可用于控制请求的超时和取消等操作。
// 参数 key 为要设置过期时间的 Redis 键。
// 参数 exp 为键的过期时长。
// 返回值 error 表示操作过程中可能出现的错误。
func (r *Client) Expire(ctx context.Context, key string, exp time.Duration) error {
	_, err := r.rdb.Expire(ctx, key, exp).Result()
	if err != nil && gerror.Is(err, redis.Nil) {
		return err
	}
	return nil
}

// MSet 批量设置 Redis 中的多个键值对。
// 参数 ctx 为上下文，可用于控制请求的超时和取消等操作。
// 参数 values 为要设置的键值对映射，键为 Redis 键，值为要设置的值。
// 返回值 error 表示操作过程中可能出现的错误。
func (r *Client) MSet(ctx context.Context, values map[string]interface{}) error {
	_, err := r.rdb.MSet(ctx, values).Result()
	if err != nil {
		return err
	}
	return nil
}

// ExpireAt 为 Redis 中指定的键设置过期时间。
// 参数 ctx 为上下文，可用于控制请求的超时和取消等操作。
// 参数 key 为要设置过期时间的 Redis 键。
// 参数 t 为键的过期时间点。
// 返回值 error 表示操作过程中可能出现的错误。
func (r *Client) ExpireAt(ctx context.Context, key string, t time.Time) error {
	_, err := r.rdb.ExpireAt(ctx, key, t).Result()
	if err != nil {
		return err
	}
	return nil
}

// Increment 对 Redis 中的指定键进行自增操作。
// 参数 ctx 为上下文，可用于控制请求的超时和取消等操作。
// 参数 key 为要自增的 Redis 键。
// 返回值 int64 为自增后的键值，error 表示操作过程中可能出现的错误。
func (r *Client) Increment(ctx context.Context, key string) (int64, error) {
	cmd := r.rdb.Incr(ctx, key)
	return cmd.Result()
}

// ConfigSet 设置 Redis 服务器的配置参数。
// 参数 ctx 为上下文，用于控制请求的生命周期，可进行超时控制、取消操作等。
// 参数 param 为要设置的 Redis 配置参数名称。
// 参数 value 为该配置参数要设置的值。
// 返回值 error 表示操作过程中可能出现的错误，若操作成功则返回 nil。
func (r *Client) ConfigSet(ctx context.Context, param string, value string) error {
	cmd := r.rdb.ConfigSet(ctx, param, value)
	if cmd == nil {
		return gerror.New("redisDao ConfigSet cmd is nil")
	}

	_, err := cmd.Result()
	if err != nil && gerror.Is(err, redis.Nil) {
		return err
	}

	return nil
}

// ConfigGet 获取 Redis 服务器的配置参数值。
// 参数 ctx 为上下文，用于控制请求的生命周期，可进行超时控制、取消操作等。
// 参数 param 为要获取的 Redis 配置参数名称。支持使用通配符，如 * 匹配所有配置。
// 返回值 map[string]string 包含配置参数名及其对应的值，error 表示操作过程中可能出现的错误，若操作成功则返回 nil。
func (r *Client) ConfigGet(ctx context.Context, param string) (map[string]string, error) {
	cmd := r.rdb.ConfigGet(ctx, param)
	if cmd == nil {
		return nil, gerror.New("redisDao ConfigGet cmd is nil")
	}

	result, err := cmd.Result()
	if err != nil && gerror.Is(err, redis.Nil) {
		return nil, err
	}
	return result, nil
}

// ProduceMsg 向 Redis 列表右侧追加消息，模拟消息队列的生产者操作。
// 参数 ctx 为上下文，用于控制请求的生命周期，可进行超时控制、取消操作等。
// 参数 queue 为 Redis 列表的键名，代表消息队列的名称。
// 参数 msg 为要添加到消息队列中的消息内容，以字符串形式存储。
// 返回值 error 表示操作过程中可能出现的错误，若操作成功则返回 nil。
func (r *Client) ProduceMsg(ctx context.Context, queue, msg string) error {
	return r.rdb.RPush(ctx, queue, msg).Err()
}

// ConsumeMsg 从 Redis 列表左侧阻塞式地获取消息，模拟消息队列的消费者操作。
// 参数 ctx 为上下文，用于控制请求的生命周期，可进行超时控制、取消操作等。
// 参数 queue 为 Redis 列表的键名，代表要消费消息的队列名称。
// 返回值 string 为获取到的消息内容，error 表示操作过程中可能出现的错误，若操作成功则返回 nil。
// 当列表为空时，该方法会一直阻塞，直到有新消息加入或上下文被取消。
func (r *Client) ConsumeMsg(ctx context.Context, queue string) (string, error) {
	res, err := r.rdb.BLPop(ctx, 0, queue).Result()
	if err != nil {
		return "", err
	}
	return res[1], nil
}

// RunScript 在 Redis 服务器上执行 Lua 脚本。
// 参数 ctx 为上下文，用于控制请求的生命周期，可进行超时控制、取消操作等。
// 参数 script 为要执行的 Lua 脚本内容。
// 参数 keys 为脚本中会使用到的 Redis 键名列表，脚本中可通过 KEYS 数组访问。
// 参数 args 为传递给脚本的额外参数，脚本中可通过 ARGV 数组访问。
// 返回值 *redis.Cmd 为 Redis 命令对象，可通过该对象获取脚本执行结果或错误信息。
func (r *Client) RunScript(ctx context.Context, script string, keys []string, args ...interface{}) *redis.Cmd {
	incr := redis.NewScript(script)
	return incr.Run(ctx, r.rdb, keys, args...)
}

func (r *Client) RPush(ctx context.Context, key string, values ...interface{}) error {
	return r.rdb.RPush(ctx, key, values).Err()
}

// Publish 向指定的 Redis 频道发布消息，实现消息的广播功能。
// 参数 ctx 为上下文，用于控制请求的生命周期，可进行超时控制、取消操作等。
// 参数 channel 为 Redis 频道的名称，订阅该频道的客户端将能接收到此消息。
// 参数 message 为要发布的消息内容，支持多种数据类型。
// 返回值 error 表示操作过程中可能出现的错误，若操作成功则返回 nil。
func (r *Client) Publish(ctx context.Context, channel string, message interface{}) error {
	return r.rdb.Publish(ctx, channel, message).Err()
}

// GetSubscriberCount 获取指定 Redis 频道的订阅者数量。
// 参数 ctx 为上下文，用于控制请求的生命周期，可进行超时控制、取消操作等。
// 参数 channel 为 Redis 频道的名称。
// 返回值 int64 为该频道当前的订阅者数量，若操作出错则返回 0。
func (r *Client) GetSubscriberCount(ctx context.Context, channel string) (count int64) {
	statusCmd := r.rdb.PubSubNumSub(ctx, channel)
	if statusCmd.Err() != nil {
		// Subscribe 订阅指定的 Redis 频道，返回一个用于接收消息的订阅对象。
		// 参数 ctx 为上下文，用于控制请求的生命周期，可进行超时控制、取消操作等。
		// 参数 channel 为 Redis 频道的名称。
		// 返回值 *redis.PubSub 为 Redis 订阅对象，可通过该对象接收频道中的消息。
		return
	}
	count = statusCmd.Val()[channel]
	return
}

// Subscribe 订阅指定的 Redis 频道，用于接收该频道发布的消息。
// 调用此方法后，会返回一个 *redis.PubSub 对象，可通过该对象的方法（如 ReceiveMessage）
// 来获取频道中发布的消息，也可用于取消订阅等操作。
// 参数 ctx 为上下文，用于控制请求的生命周期，可进行超时控制、取消操作等。
// 参数 channel 为要订阅的 Redis 频道名称。
// 返回值 *redis.PubSub 为 Redis 订阅对象，通过该对象可对订阅进行管理和接收消息。
func (r *Client) Subscribe(ctx context.Context, channel string) *redis.PubSub {
	return r.rdb.Subscribe(ctx, channel)
}
