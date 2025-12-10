/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package fields

import (
	"sort"
)

type Field struct {
	Key    string   `json:"key"`
	Values []string `json:"values"`
}

type Fields struct {
	fields []Field
}

func (f *Fields) AddField(key, value string) {
	for i := range f.fields {
		if f.fields[i].Key == key {
			f.fields[i].Values = append(f.fields[i].Values, value)
			return
		}
	}
	f.fields = append(f.fields, Field{Key: key, Values: []string{value}})
}

func (f *Fields) GetField(key string) []string {
	for i := range f.fields {
		if f.fields[i].Key == key {
			return f.fields[i].Values
		}
	}
	return nil
}

func (f *Fields) DeleteField(key string) {
	i := -1
	for j := range f.fields {
		if f.fields[j].Key == key {
			i = j
		}
	}
	if i == -1 {
		return
	}
	f.fields = append(f.fields[:i], f.fields[i+1:]...)
}

func (f *Fields) AllFields() []Field {
	return f.fields
}

func (f *Fields) SortFields() {
	sort.Slice(f.fields, func(i, j int) bool {
		return f.fields[i].Key < f.fields[j].Key
	})
}

func (f *Fields) OverwriteField(key, value string) {
	for i := range f.fields {
		if f.fields[i].Key == key {
			f.fields[i].Values = []string{value}
			return
		}
	}
	f.fields = append(f.fields, Field{Key: key, Values: []string{value}})
}

func (f *Fields) RemoveField(key, value string) {
	// find the index of the key
	i := -1
	for z := range f.fields {
		if f.fields[z].Key == key {
			i = z
		}
	}

	// if the key is not present, done
	if i == -1 {
		return
	}

	// now find if the value is present
	j := -1
	for z := range f.fields[i].Values {
		if f.fields[i].Values[z] == value {
			j = z
		}
	}

	// if the value is not present, done
	if j == -1 {
		return
	}

	// take the value out
	f.fields[i].Values = append(f.fields[i].Values[:j], f.fields[i].Values[j+1:]...)

	// if there are no values left, take the key out
	if len(f.fields[i].Values) == 0 {
		f.fields = append(f.fields[:i], f.fields[i+1:]...)
	}
}

func (f *Fields) SetField(key string, values []string) {
	for i := range f.fields {
		if f.fields[i].Key == key {
			f.fields[i].Values = values
			return
		}
	}
	f.fields = append(f.fields, Field{Key: key, Values: values})
}
