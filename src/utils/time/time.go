package main

/*
sort.Slice 函数的实现原理：

首先它会接受一个切片和一个比较函数
它会创建一个 sort.Interface 类型的值，该值保存了切片的长度、比较函数和一个交换元素的函数
sort.Slice 会调用 sort.Sort 函数对该 sort.Interface 类型的值进行排序
sort.Sort 函数会使用快速排序算法对切片进行排序
*/

import (
	"fmt"
	"sort"
	"time"
)

type Cid struct {
	ID         string
	UpdateTime time.Time
}

func main() {
	// 示例数据
	cids := []*Cid{
		{
			ID:         "cid1",
			UpdateTime: time.Unix(1569247774, 0),
		},
		{
			ID:         "cid2",
			UpdateTime: time.Unix(1569247775, 0),
		},
		{
			ID:         "cid3",
			UpdateTime: time.Unix(1569247772, 0),
		},
	}

	// 根据 UpdateTime 字段对 Cid 切片进行排序
	sort.Slice(cids, func(i, j int) bool {
		return cids[i].UpdateTime.After(cids[j].UpdateTime)
	})

	// 输出排序后的 Cid 切片
	for _, cid := range cids {
		fmt.Printf("ID: %s, UpdateTime: %v\n", cid.ID, cid.UpdateTime)
	}
}

/*
ID: cid2, UpdateTime: 2019-09-23 22:09:35 +0800 CST
ID: cid1, UpdateTime: 2019-09-23 22:09:34 +0800 CST
ID: cid3, UpdateTime: 2019-09-23 22:09:32 +0800 CST
*/
