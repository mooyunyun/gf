package gmap_test

import (
	"github.com/gogf/gf/g/container/gmap"
	"github.com/gogf/gf/g/test/gtest"
	"testing"
)

func getInt() int {
	return 123
}
func intIntCallBack(int, int) bool {
	return true
}
func Test_IntIntMap_Basic(t *testing.T) {
	gtest.Case(t, func() {
		m := gmap.NewIntIntMap()
		m.Set(1, 1)

		gtest.Assert(m.Get(1), 1)
		gtest.Assert(m.Size(), 1)
		gtest.Assert(m.IsEmpty(), false)

		gtest.Assert(m.GetOrSet(2, 2), 2)
		gtest.Assert(m.SetIfNotExist(2, 2), false)

		gtest.Assert(m.SetIfNotExist(3, 3), true)

		gtest.Assert(m.Remove(2), 2)
		gtest.Assert(m.Contains(2), false)

		gtest.AssertIN(3, m.Keys())
		gtest.AssertIN(1, m.Keys())
		gtest.AssertIN(3, m.Values())
		gtest.AssertIN(1, m.Values())
		m.Flip()
		gtest.Assert(m.Map(), map[int]int{1: 1, 3: 3})

		m.Clear()
		gtest.Assert(m.Size(), 0)
		gtest.Assert(m.IsEmpty(), true)

		m2 := gmap.NewIntIntMapFrom(map[int]int{1: 1, 2: 2})
		gtest.Assert(m2.Map(), map[int]int{1: 1, 2: 2})
		m3 := gmap.NewIntIntMapFromArray([]int{1, 2}, []int{1, 2})
		gtest.Assert(m3.Map(), map[int]int{1: 1, 2: 2})

	})
}
func Test_IntIntMap_Set_Fun(t *testing.T) {
	m := gmap.NewIntIntMap()

	m.GetOrSetFunc(1, getInt)
	m.GetOrSetFuncLock(2, getInt)
	gtest.Assert(m.Get(1), 123)
	gtest.Assert(m.Get(2), 123)
	gtest.Assert(m.SetIfNotExistFunc(1, getInt), false)
	gtest.Assert(m.SetIfNotExistFuncLock(2, getInt), false)
}

func Test_IntIntMap_Batch(t *testing.T) {
	m := gmap.NewIntIntMap()

	m.BatchSet(map[int]int{1: 1, 2: 2, 3: 3})
	m.Iterator(intIntCallBack)
	gtest.Assert(m.Map(), map[int]int{1: 1, 2: 2, 3: 3})
	m.BatchRemove([]int{1, 2})
	gtest.Assert(m.Map(), map[int]int{3: 3})
}

func Test_IntIntMap_Clone(t *testing.T) {
	//clone 方法是深克隆
	m := gmap.NewIntIntMapFrom(map[int]int{1: 1, 2: 2})

	m_clone := m.Clone()
	m.Remove(1)
	//修改原 map,clone 后的 map 不影响
	gtest.AssertIN(1, m_clone.Keys())

	m_clone.Remove(2)
	//修改clone map,原 map 不影响
	gtest.AssertIN(2, m.Keys())
}
func Test_IntIntMap_Merge(t *testing.T) {
	m1 := gmap.NewIntIntMap()
	m2 := gmap.NewIntIntMap()
	m1.Set(1, 1)
	m2.Set(2, 2)
	m1.Merge(m2)
	gtest.Assert(m1.Map(), map[int]int{1: 1, 2: 2})
}
