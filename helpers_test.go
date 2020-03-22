/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package orderedmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveNestedField(t *testing.T) {
	x := New()
	x.Set("y", 1)
	x.Set("a", "foo")
	obj := New()
	obj.Set("x", x)

	RemoveNestedField(obj, "x", "a")
	assert.Equal(t, 1, obj.Entry("x").(*OrderedMap).Len())
	RemoveNestedField(obj, "x", "y")
	assert.True(t, obj.Entry("x").(*OrderedMap).IsZero())
	RemoveNestedField(obj, "x")
	assert.True(t, obj.IsZero())
	RemoveNestedField(obj, "x") // Remove of a non-existent field
	assert.True(t, obj.IsZero())
}

func TestNestedFieldNoCopy(t *testing.T) {
	target := New()
	target.Set("foo", "bar")

	f := New()
	f.Set("f", "bar")
	a := New()
	a.Set("b", target)
	a.Set("c", nil)
	a.Set("d", []interface{}{"foo"})
	a.Set("e", []interface{}{f})
	obj := New()
	obj.Set("a", a)

	// case 1: field exists and is non-nil
	res, exists, err := NestedFieldNoCopy(obj, "a", "b")
	assert.True(t, exists)
	assert.Nil(t, err)
	assert.Equal(t, target, res)
	target.Set("foo", "baz")
	assert.Equal(t, target.Entry("foo"), res.(*OrderedMap).Entry("foo"), "result should be a reference to the expected item")

	// case 2: field exists and is nil
	res, exists, err = NestedFieldNoCopy(obj, "a", "c")
	assert.True(t, exists)
	assert.Nil(t, err)
	assert.Nil(t, res)

	// case 3: error traversing obj
	res, exists, err = NestedFieldNoCopy(obj, "a", "d", "foo")
	assert.False(t, exists)
	assert.NotNil(t, err)
	assert.Nil(t, res)

	// case 4: field does not exist
	res, exists, err = NestedFieldNoCopy(obj, "a", "g")
	assert.False(t, exists)
	assert.Nil(t, err)
	assert.Nil(t, res)

	// case 5: intermediate field does not exist
	res, exists, err = NestedFieldNoCopy(obj, "a", "g", "f")
	assert.False(t, exists)
	assert.Nil(t, err)
	assert.Nil(t, res)

	// case 6: intermediate field is null
	//         (background: happens easily in YAML)
	res, exists, err = NestedFieldNoCopy(obj, "a", "c", "f")
	assert.False(t, exists)
	assert.Nil(t, err)
	assert.Nil(t, res)

	// case 7: array/slice syntax is not supported
	//         (background: users may expect this to be supported)
	res, exists, err = NestedFieldNoCopy(obj, "a", "e[0]")
	assert.False(t, exists)
	assert.Nil(t, err)
	assert.Nil(t, res)
}

func TestNestedFieldCopy(t *testing.T) {
	target := New()
	target.Set("foo", "bar")

	a := New()
	a.Set("b", target)
	a.Set("c", nil)
	a.Set("d", []interface{}{"foo"})
	obj := New()
	obj.Set("a", a)

	// case 1: field exists and is non-nil
	res, exists, err := NestedFieldCopy(obj, "a", "b")
	assert.True(t, exists)
	assert.Nil(t, err)
	assert.Equal(t, target, res)
	target.Set("foo", "baz")
	assert.NotEqual(t, target.Entry("foo"), res.(*OrderedMap).Entry("foo"), "result should be a copy of the expected item")

	// case 2: field exists and is nil
	res, exists, err = NestedFieldCopy(obj, "a", "c")
	assert.True(t, exists)
	assert.Nil(t, err)
	assert.Nil(t, res)

	// case 3: error traversing obj
	res, exists, err = NestedFieldCopy(obj, "a", "d", "foo")
	assert.False(t, exists)
	assert.NotNil(t, err)
	assert.Nil(t, res)

	// case 4: field does not exist
	res, exists, err = NestedFieldCopy(obj, "a", "e")
	assert.False(t, exists)
	assert.Nil(t, err)
	assert.Nil(t, res)
}
