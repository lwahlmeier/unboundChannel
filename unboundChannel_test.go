package unboundchannel

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

var letterRunes = []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func makeString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func TestUnboundChannel(t *testing.T) {
	count := int(math.Max(400, float64(rand.Intn(1000))))
	fmt.Println(count)
	duc := NewUnboundChannel[int]()
	for i := 0; i < count; i++ {
		duc.Add() <- i
	}
	assert.Equal(t, count, <-duc.Count())
	waiter := sync.WaitGroup{}
	waiter.Add(count)
	for i := 0; i < 50; i++ {
		go func() {
			for {
				_, ok := <-duc.Get()
				if !ok {
					break
				}
				waiter.Done()
			}
		}()
	}
	waiter.Wait()
	assert.Equal(t, 0, <-duc.Count())
	count = int(math.Max(400, float64(rand.Intn(1000))))
	waiter.Add(count)
	for i := 0; i < count; i++ {
		duc.Add() <- i
	}
	waiter.Wait()
	assert.Equal(t, 0, <-duc.Count())
	duc.Stop()
	assert.False(t, duc.IsRunning())
}

// func TestDynamicMsgChannel(t *testing.T) {
// 	count := int(math.Max(400, float64(rand.Intn(1000))))
// 	dmc := NewDynamicMsgChannel()
// 	for i := 0; i < count; i++ {
// 		dmc.Add() <- &pubsub.ReceivedMessage{}
// 	}
// 	et := time.Now().Add(time.Millisecond * 5)
// 	for time.Until(et) > 0 && count != <-dmc.Count() {
// 		runtime.Gosched()
// 	}
// 	assert.Equal(t, count, <-dmc.Count())
// 	waiter := sync.WaitGroup{}
// 	waiter.Add(count)
// 	for i := 0; i < 50; i++ {
// 		go func() {
// 			for {
// 				<-dmc.Get()
// 				<-dmc.Count()
// 				waiter.Done()
// 			}
// 		}()
// 	}
// 	waiter.Wait()
// 	assert.Equal(t, 0, <-dmc.Count())
// 	count = int(math.Max(400, float64(rand.Intn(1000))))
// 	waiter.Add(count)
// 	for i := 0; i < count; i++ {
// 		dmc.Add() <- &pubsub.ReceivedMessage{}
// 	}
// 	waiter.Wait()
// 	assert.Equal(t, 0, <-dmc.Count())
// }
