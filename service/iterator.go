package main

import (
	"baseservices/kvservice/service/rocksdb"
	"sync"

	"github.com/tecbot/gorocksdb"
	"baseservices/kvservice"
)

type Iterator interface {
	Next() bool
	Value() *kvservice.KVInfo
	Close()
}

type _Iterator struct {
	channel   chan *rocksdb.DataInfo
	waitGroup *sync.WaitGroup
	allClose  bool
	value     *kvservice.KVInfo
	isClose   bool
	prefixKey []byte
}

func newIterator(size int, prefixKey []byte) *_Iterator {
	return &_Iterator{
		channel:   make(chan *rocksdb.DataInfo, size),
		waitGroup: new(sync.WaitGroup),
		allClose:  false,
		isClose:   false,
		prefixKey: prefixKey,
	}
}

func (this *_Iterator) add(iter *gorocksdb.Iterator) {
	this.waitGroup.Add(1)
	go func() {
		defer func() {
			recover()
			iter.Close()
			this.waitGroup.Done()
		}()
		for iter.SeekToFirst(); iter.ValidForPrefix(this.prefixKey); iter.Next() {
			this.channel <- &rocksdb.DataInfo{
				Key:   copyData(iter.Key()),
				Value: copyData(iter.Value()),
			}
		}
	}()
}

func (this *_Iterator) wait() {
	go func() {
		this.waitGroup.Wait()
		this.allClose = true
		if len(this.channel) == 0 {
			this.Close()
		}
	}()
}

func (this *_Iterator) Next() bool {
	if !this.isClose && (!this.allClose || len(this.channel) != 0) {
		value := <-this.channel
		this.value = &kvservice.KVInfo{
			Key:string(value.Key),
			Value:string(value.Value),
		}
		return true
	}
	this.value = nil
	return false
}

func (this *_Iterator) Value() *kvservice.KVInfo{
	return this.value

}
func (this *_Iterator) Close() {
	close(this.channel)
}

func copyData(slice *gorocksdb.Slice) []byte {
	data := make([]byte, slice.Size())
	copy(data, slice.Data())
	slice.Free()
	return data
}
