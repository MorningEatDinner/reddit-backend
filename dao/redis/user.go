package redis

import (
	"context"
	"fmt"
)

func SetProfileStatus(uid string, status bool) error {
	if err := RDB.Client.HSet(context.Background(), getRedisKey(KeyProfileStatus), uid, status).Err(); err != nil {
		fmt.Printf("Error setting profile status: %v\n", err)
		return err
	}
	return nil
}
