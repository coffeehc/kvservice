package main

import (
	"github.com/coffeehc/baseservices/kvservice/service/rocksdb"

	"github.com/coffeehc/baseservices/kvservice"
	"github.com/tecbot/gorocksdb"
	"sort"
	"github.com/coffeehc/baseservices/kvservice/modules"
)

func memcmp(i, j []byte) int {
	s, is, js := 0, len(i), len(j)
	if is < js {
		s = is
	} else {
		s = js
	}
	for k := 0; k < s; k++ {
		if i[k] == j[k] {
			continue
		}
		if i[k] > j[k] {
			return 1
		}
		return -1
	}
	return 0
}

type Order interface {
	Less(i, j []byte) bool
	InitIterator(iter *gorocksdb.Iterator)
	Next(iter *gorocksdb.Iterator)
}

type Order_asc struct {
}

func (this Order_asc) Less(i, j []byte) bool {
	return memcmp(j, i) > 0
}
func (this Order_asc) InitIterator(iter *gorocksdb.Iterator) {
	iter.SeekToFirst()
}
func (this Order_asc) Next(iter *gorocksdb.Iterator) {
	iter.Next()
}

type Order_desc struct {
}

func (this Order_desc) Less(i, j []byte) bool {
	return memcmp(i, j) > 0
}
func (this Order_desc) InitIterator(iter *gorocksdb.Iterator) {
	iter.SeekToLast()
}
func (this Order_desc) Next(iter *gorocksdb.Iterator) {
	iter.Prev()
}

type Iterator interface {
	Next() bool
	Value() *modules.KVInfo
	Close()
}

type value struct {
	data      *rocksdb.DataInfo
	partition int
}

type sortValue struct {
	data  []*value
	order Order
}

func (c *sortValue) Len() int {
	return len(c.data)
}
func (c *sortValue) Swap(i, j int) {
	c.data[i], c.data[j] = c.data[j], c.data[i]
}
func (c *sortValue) Less(i, j int) bool {
	return c.order.Less(c.data[i].data.Key, c.data[j].data.Key)
}

func (c *sortValue) addData(v *value) {
	c.data = append(c.data, v)
}

type _Iterator struct {
	value        *modules.KVInfo
	prefixKey    []byte
	startKey     []byte
	order        Order
	iterators    []*gorocksdb.Iterator
	valueCache   *sortValue
	prefixKeyLen int
	valueChan    chan *modules.KVInfo
	isClose      bool
}

func newIterator(partition int, prefixKey, startKey []byte, orderStr string) *_Iterator {
	var order Order = &Order_asc{}
	if orderStr == string(kvservice.Order_DESC) {
		order = &Order_desc{}
	}
	return &_Iterator{
		prefixKey:    prefixKey,
		startKey:     startKey,
		order:        order,
		iterators:    make([]*gorocksdb.Iterator, partition),
		prefixKeyLen: len(prefixKey),
		valueCache: &sortValue{
			data:  make([]*value, 0),
			order: order,
		},
		valueChan: make(chan *modules.KVInfo, 1024),
		isClose:   false,
	}
}

func (this *_Iterator) add(partition int, iter *gorocksdb.Iterator) {
	this.order.InitIterator(iter)
	var dataInfo *rocksdb.DataInfo = nil
	start := this.startKey
	for iter.ValidForPrefix(this.prefixKey) {
		if len(this.startKey) == 0 || this.order.Less(start, copyData(iter.Key())) {
			dataInfo = &rocksdb.DataInfo{
				Key:   copyData(iter.Key()),
				Value: copyData(iter.Value()),
			}
			break
		}
		this.order.Next(iter)
	}
	if dataInfo != nil {
		this.iterators[partition] = iter
		this.valueCache.addData(&value{data: dataInfo, partition: partition})
	}
}

func (this *_Iterator) addEnd() {
	sort.Sort(this.valueCache)
	go this.getValue()
}

func (this *_Iterator) getValue() {
	defer close(this.valueChan)
	for !this.isClose {
		if len(this.valueCache.data) == 0 {
			return
		}
		v := this.valueCache.data[0]
		this.valueChan <- &modules.KVInfo{
			Key:   v.data.Key,
			Value: v.data.Value,
		}
		partition := v.partition
		nextValue := this.getValueWhitIterator(partition)
		if nextValue != nil {
			quickSort(this.valueCache.data, &value{data: nextValue, partition: partition}, this.prefixKeyLen, this.order.Less)
		} else {
			this.valueCache.data = this.valueCache.data[1:]
		}
	}
}

func (this *_Iterator) getValueWhitIterator(partition int) *rocksdb.DataInfo {
	iter := this.iterators[partition]
	if iter == nil {
		return nil
	}
	this.order.Next(iter)
	if iter.ValidForPrefix(this.prefixKey) {
		return &rocksdb.DataInfo{
			Key:   copyData(iter.Key()),
			Value: copyData(iter.Value()),
		}
	}
	this.iterators[partition] = nil
	iter.Close()
	return nil
}

func (this *_Iterator) Next() bool {
	this.value = <-this.valueChan
	return this.value != nil
}

func (this *_Iterator) Value() *modules.KVInfo {
	return this.value

}
func (this *_Iterator) Close() {
	this.isClose = true
	for _, i := range this.iterators {
		if i != nil {
			i.Close()
		}
	}

}

func copyData(slice *gorocksdb.Slice) []byte {
	data := make([]byte, slice.Size())
	copy(data, slice.Data())
	slice.Free()
	return data
}

func quickSort(vs []*value, newValue *value, prefixLen int, orderFunc func(i, j []byte) bool) {
	key := newValue.data.Key[prefixLen:]
	size := len(vs)
	for i := 1; i < size; i++ {
		if orderFunc(vs[i].data.Key[prefixLen:], key) {
			vs[i-1] = vs[i]
		} else {
			vs[i-1] = newValue
			return
		}
	}
	vs[size-1] = newValue
}
