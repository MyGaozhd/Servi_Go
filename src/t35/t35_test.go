package t35

import (
	"errors"
	"testing"
	"time"
)

func Test35_0(t *testing.T) {
	//初始化十个对象的对象池
	pool := NewResumeAbleObject(10)

	for i := 0; i < 10; i++ {
		if obj, err := pool.GetObj(time.Millisecond * 10); err != nil {
			t.Error(err)
		} else {
			t.Log(obj)
			if err := pool.ReleaseObj(obj); err != nil {
				t.Error(err)
			}
		}
	}
}

type ResumeAbleObj struct {
}

type ObjectPool struct {
	bufChan chan *ResumeAbleObj
}

func NewResumeAbleObject(num int) *ObjectPool {
	objectPool := ObjectPool{}
	objectPool.bufChan = make(chan *ResumeAbleObj, num)
	for i := 0; i < num; i++ {
		objectPool.bufChan <- &ResumeAbleObj{}
	}
	return &objectPool
}

func (p *ObjectPool) GetObj(timeout time.Duration) (*ResumeAbleObj, error) {
	select {
	case ret := <-p.bufChan:
		return ret, nil
	case <-time.After(timeout):
		return nil, errors.New("time out")
	}
}

func (p *ObjectPool) ReleaseObj(obj *ResumeAbleObj) error {
	select {
	case p.bufChan <- obj:
		return nil
	default:
		return errors.New("over fllow")
	}
}
