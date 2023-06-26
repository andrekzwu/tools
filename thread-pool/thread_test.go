package threadpool

import (
	"context"
	"fmt"
	"testing"
	"time"
)

type Data struct {
	Value string `json:"value"`
}

// TestTreadPool
func TestTreadPool(t *testing.T) {
	pool, err := NewThreadPool(10)
	if err != nil {
		fmt.Println("-------------err:", err)
		return
	}
	//
	fu := func(ctx context.Context, eventData *EventData) {
		data, ok := eventData.Data.(Data)
		if !ok {
			return
		}
		fmt.Println("---------------data.value", data.Value)
	}
	for i := 0; i < 100; i++ {
		pool.DispatchTask2Thread(&EventData{Func: fu, Data: Data{Value: fmt.Sprintf("test thread %d", i)}})
	}
	time.Sleep(2 * time.Second)
}
