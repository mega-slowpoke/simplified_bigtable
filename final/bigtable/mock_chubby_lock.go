package bigtable

//
//import (
//	"github.com/go-redis/redis/v8"
//	"golang.org/x/net/context"
//	"log"
//	"time"
//)
//
//var ctx = context.Background()
//
//type ChubbyLock struct {
//	client *redis.Client
//	key    string
//}
//
//// NewChubbyLock
//func NewChubbyLock(redisAddr string, key string) *ChubbyLock {
//	client := redis.NewClient(&redis.Options{
//		Addr: redisAddr,
//	})
//
//	_, err := client.Ping(ctx).Result()
//	if err != nil {
//		log.Fatalf("Could not connect to Redis: %v", err)
//	}
//
//	return &ChubbyLock{
//		client: client,
//		key:    key,
//	}
//}
//
//func (l *ChubbyLock) Lock(timeout time.Duration) bool {
//	success, err := l.client.SetNX(ctx, l.key, "locked", timeout).Result()
//	if err != nil {
//		log.Printf("Error trying to acquire lock: %v", err)
//		return false
//	}
//	return success
//}
//
//func (l *ChubbyLock) Unlock() error {
//	_, err := l.client.Del(ctx, l.key).Result()
//	if err != nil {
//		log.Printf("Error trying to release lock: %v", err)
//	}
//	return err
//}
//
//func (l *ChubbyLock) CheckLockExist() bool {
//	exists, err := l.client.Exists(ctx, l.key).Result()
//	if err != nil {
//		log.Printf("Error checking lock status: %v", err)
//		return false
//	}
//	return exists > 0
//}
