/*
Copyright (c) 2026 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package fields_test

import (
	"testing"

	"github.com/golab/board/internal/assert"
	"github.com/golab/board/internal/require"
	"github.com/golab/board/pkg/core/fields"
)

func TestAddField1(t *testing.T) {
	f := &fields.Fields{}
	f.AddField("foo", "bar")

	g := f.GetField("foo")

	require.Equal(t, len(g), 1)
	assert.Equal(t, g[0], "bar")
}

func TestAddField2(t *testing.T) {
	f := &fields.Fields{}
	f.AddField("foo", "bar")
	f.AddField("foo", "baz")

	g := f.GetField("foo")

	require.Equal(t, len(g), 2)
	assert.Equal(t, g[0], "bar")
	assert.Equal(t, g[1], "baz")
}

func TestGetFieldEmpty(t *testing.T) {
	f := &fields.Fields{}
	n := f.GetField("foo")
	assert.Equal(t, len(n), 0)
}

func TestDeleteField1(t *testing.T) {
	f := &fields.Fields{}
	f.AddField("foo", "bar")
	f.DeleteField("bar")

	g := f.GetField("foo")
	require.Equal(t, len(g), 1)
	assert.Equal(t, g[0], "bar")

	f.DeleteField("foo")

	g = f.GetField("foo")
	assert.Equal(t, len(g), 0)
}

func TestAllFields(t *testing.T) {
	f := &fields.Fields{}
	f.AddField("foo", "bar")
	f.AddField("baz", "bot")

	assert.Equal(t, len(f.AllFields()), 2)
}

func TestSortFields(t *testing.T) {
	f := &fields.Fields{}
	f.AddField("foo", "bar")
	f.AddField("baz", "bot")
	f.AddField("abc", "def")
	f.SortFields()
	af := f.AllFields()
	require.Equal(t, len(af), 3)
	assert.Equal(t, af[0].Key, "abc")
	assert.Equal(t, af[1].Key, "baz")
	assert.Equal(t, af[2].Key, "foo")
}

func TestOverwriteField1(t *testing.T) {
	f := &fields.Fields{}
	f.AddField("foo", "bar")
	f.OverwriteField("foo", "baz")
	g := f.GetField("foo")
	require.Equal(t, len(g), 1)
	assert.Equal(t, g[0], "baz")
}

func TestOverwriteField2(t *testing.T) {
	f := &fields.Fields{}
	f.OverwriteField("foo", "baz")
	g := f.GetField("foo")
	require.Equal(t, len(g), 1)
	assert.Equal(t, g[0], "baz")
}

func TestRemoveField1(t *testing.T) {
	f := &fields.Fields{}
	f.AddField("foo", "bar")
	f.AddField("baz", "bot")
	f.AddField("abc", "def")
	f.RemoveField("bar", "bot")
	assert.Equal(t, len(f.AllFields()), 3)
}

func TestRemoveField2(t *testing.T) {
	f := &fields.Fields{}
	f.AddField("foo", "bar")
	f.AddField("baz", "bot")
	f.AddField("abc", "def")
	f.RemoveField("baz", "boo")
	assert.Equal(t, len(f.AllFields()), 3)
}

func TestRemoveField3(t *testing.T) {
	f := &fields.Fields{}
	f.AddField("foo", "bar")
	f.AddField("baz", "bot")
	f.AddField("abc", "def")
	f.RemoveField("baz", "bot")
	assert.Equal(t, len(f.AllFields()), 2)
}

func TestRemoveField4(t *testing.T) {
	f := &fields.Fields{}
	f.AddField("foo", "bar")
	f.AddField("baz", "bot")
	f.AddField("baz", "boo")
	f.AddField("abc", "def")
	g := f.GetField("baz")
	assert.Equal(t, len(g), 2)

	f.RemoveField("baz", "bot")
	require.Equal(t, len(f.AllFields()), 3)
	g = f.GetField("baz")
	assert.Equal(t, len(g), 1)
}

func TestSetField1(t *testing.T) {
	f := &fields.Fields{}
	f.SetField("foo", []string{"bar", "baz", "bot"})
	g := f.GetField("foo")
	require.Equal(t, len(g), 3)
	assert.Equal(t, g[0], "bar")
	assert.Equal(t, g[1], "baz")
	assert.Equal(t, g[2], "bot")
}

func TestSetField2(t *testing.T) {
	f := &fields.Fields{}
	f.AddField("foo", "boo")
	g := f.GetField("foo")
	require.Equal(t, len(g), 1)
	assert.Equal(t, g[0], "boo")

	f.SetField("foo", []string{"bar", "baz", "bot"})
	g = f.GetField("foo")
	require.Equal(t, len(g), 3)
	assert.Equal(t, g[0], "bar")
	assert.Equal(t, g[1], "baz")
	assert.Equal(t, g[2], "bot")
}

func TestHasField(t *testing.T) {
	f := &fields.Fields{}
	assert.False(t, f.HasField("foo"))
	f.AddField("foo", "boo")
	assert.True(t, f.HasField("foo"))
	assert.False(t, f.HasField("boo"))
}

func TestAppendField1(t *testing.T) {
	f := &fields.Fields{}
	f.SetField("foo", []string{"bar"})
	g := f.GetField("foo")
	require.Equal(t, len(g), 1)
	assert.Equal(t, g[0], "bar")
	f.AppendField("foo", []string{"baz", "bot"})
	g = f.GetField("foo")
	require.Equal(t, len(g), 3)
	assert.Equal(t, g[0], "bar")
	assert.Equal(t, g[1], "baz")
	assert.Equal(t, g[2], "bot")
}

func TestAppendField2(t *testing.T) {
	f := &fields.Fields{}
	f.AppendField("foo", []string{"bar", "baz", "bot"})
	g := f.GetField("foo")
	require.Equal(t, len(g), 3)
	assert.Equal(t, g[0], "bar")
	assert.Equal(t, g[1], "baz")
	assert.Equal(t, g[2], "bot")
}
