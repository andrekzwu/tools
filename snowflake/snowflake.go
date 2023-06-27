package snowflake

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

// GetUUID 获取uuid
func GetUUID() string {
	uid, _ := uuid.NewUUID()
	return uid.String()
}

// snowflake
const (
	epoch         = int64(1672502400000)        // 设置起始时间(时间戳/毫秒)：2023-01-01 00:00:00
	workerIdBits  = 10                          //机器ID位数
	sequenceBits  = 12                          //序列号位数
	workerIdMax   = -1 ^ (-1 << workerIdBits)   // 最大机器ID
	sequenceMask  = -1 ^ (-1 << sequenceBits)   // 序列号掩码
	timeShiftBits = workerIdBits + sequenceBits // 时间戳左移位数
	workerIdShift = sequenceBits                // 机器人ID左移位数
)

type Snowflake struct {
	sync.Mutex
	lastTimestamp int64
	workerId      uint16
	sequence      uint16
}

var snow *Snowflake // 默认雪花算法机器人用于单机使用

func init() {
	snow = New(1)
}

// New 生成雪花算法生成器
func New(workId int) *Snowflake {
	if workId < 0 || workId > workerIdMax {
		panic("invalid work id")
	}
	return &Snowflake{
		Mutex:         sync.Mutex{},
		lastTimestamp: 0,
		workerId:      uint16(workId),
		sequence:      0,
	}
}

// GetSnowflakeID 获取单机雪花算法ID
func GetSnowflakeID() int64 {
	return snow.NextID()
}

func (s *Snowflake) NextID() int64 {
	s.Lock()
	defer s.Unlock()
	//
	var currTimestamp = int64(time.Now().UnixNano() / int64(time.Millisecond))
	if currTimestamp < s.lastTimestamp {
		panic("Invalid timestamp")
	}
	if currTimestamp == s.lastTimestamp {
		s.sequence = (s.sequence + 1) & sequenceMask
		if s.sequence == 0 {
			currTimestamp = s.waitNextMillis(currTimestamp)
		}
	} else {
		s.sequence = 0
	}
	s.lastTimestamp = currTimestamp
	return ((currTimestamp - epoch) << timeShiftBits) | (int64(s.workerId) << workerIdShift) |
		int64(s.sequence)
}

func (s *Snowflake) waitNextMillis(currTimestamp int64) int64 {
	for currTimestamp < s.lastTimestamp {
		currTimestamp = int64(time.Now().UnixNano() / int64(time.Millisecond))
	}
	return currTimestamp
}
