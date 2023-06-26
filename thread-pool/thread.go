package threadpool

import (
	"context"
	"errors"
	"sync"
	"time"
)

const (
	maxThreadLimit = 100                  // 协程最大限制数量
	queueInterval  = 5 * time.Millisecond // 队列处理频率
)

var ErrMaxThreadLimit = errors.New("input thread nums more than Max nums or is zero")

// EventData 处理公共知识库事件
type EventData struct {
	Func HandleEvent // 处理喊出
	Data interface{} // 处理事件数据
}

// HandleEvent 任务调度处理函数模型
type HandleEvent func(ctx context.Context, eventData *EventData)

// ThreadInfo 知识库处理线程
type ThreadInfo struct {
	sem chan *EventData
}

// ThreadPool 公共知识库处理协程池
type ThreadPool struct {
	ThreadNums  int
	FreeChans   chan int
	Mutex       sync.Mutex
	QueuesMutex sync.Mutex
	Queues      []*EventData
	ThreadInfos []*ThreadInfo
	Context     context.Context
	Cancel      context.CancelFunc
	IsOver      bool
}

// NewThreadPool 创建公共知识库协程池
func NewThreadPool(threadNums int) (*ThreadPool, error) {
	if threadNums <= 0 || threadNums > maxThreadLimit {
		return nil, ErrMaxThreadLimit
	}
	ctx, cancel := context.WithCancel(context.Background())
	//  thread pool
	pool := &ThreadPool{
		ThreadNums:  threadNums,
		FreeChans:   make(chan int, threadNums),
		Mutex:       sync.Mutex{},
		ThreadInfos: make([]*ThreadInfo, 0),
		Context:     ctx,
		Cancel:      cancel,
		IsOver:      false,
	}
	// init thread
	for i := int(0); i < threadNums; i++ {
		pool.ThreadInfos = append(pool.ThreadInfos, &ThreadInfo{
			sem: make(chan *EventData),
		})
		go pool.threadFunc(pool.Context, pool.ThreadInfos[i])
	}
	// queue monitor
	go pool.queueMonitor(ctx)
	return pool, nil
}

// DistoryThreadPool 销毁公共知识库协程池
func (tp *ThreadPool) DistoryThreadPool() {
	tp.IsOver = true
	// stop func
	tp.Cancel()
	// close free chans
	close(tp.FreeChans)
	// close thread sem
	for _, v := range tp.ThreadInfos {
		close(v.sem)
	}
	tp.ThreadInfos = nil
	// queue
	tp.Queues = nil
}

// DispatchTask2Thread 调度任务到协程处理
func (tp *ThreadPool) DispatchTask2Thread(eventData *EventData) {
	if tp.IsOver {
		return
	}
	//  空闲协程满了后放入等待队列
	if len(tp.FreeChans) == tp.ThreadNums {
		tp.QueuesMutex.Lock()
		tp.Queues = append(tp.Queues, eventData)
		tp.QueuesMutex.Unlock()
		return
	}
	// free sem
	tp.FreeChans <- 0
	// thread infos
	if len(tp.ThreadInfos) == 0 {
		return
	}
	// add task to thread
	tp.Mutex.Lock()
	threadInfo := tp.ThreadInfos[len(tp.ThreadInfos)-1]
	tp.ThreadInfos = tp.ThreadInfos[:len(tp.ThreadInfos)-1]
	tp.Mutex.Unlock()
	// add event data
	threadInfo.sem <- eventData
}

// threadFunc 处理协程
func (tp *ThreadPool) threadFunc(ctx context.Context, threadInfo *ThreadInfo) {
	for {
		// 检测是否需要退出
		select {
		case <-ctx.Done():
			return
		default:
		}
		// 处理协程,若协程关闭，则直接返回
		eventData, ok := <-threadInfo.sem
		if !ok {
			return
		}
		if eventData == nil || eventData.Func == nil {
			continue
		}
		// 执行协程处理
		eventData.Func(ctx, eventData)
		// 释放线程到线程池
		tp.Mutex.Lock()
		tp.ThreadInfos = append(tp.ThreadInfos, threadInfo)
		tp.Mutex.Unlock()
		//
		<-tp.FreeChans
	}
}

// queueMonitor
func (tp *ThreadPool) queueMonitor(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.NewTicker(queueInterval).C:
		}
		if len(tp.Queues) == 0 {
			continue
		}
		tp.QueuesMutex.Lock()
		eventData := tp.Queues[len(tp.Queues)-1]
		tp.Queues = tp.Queues[:len(tp.Queues)-1]
		tp.QueuesMutex.Unlock()
		//
		tp.DispatchTask2Thread(eventData)
	}
}
