/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package loader

type MemoryLoader struct {
}

func NewMemoryLoader() *MemoryLoader {
	return &MemoryLoader{}
}

func (ml *MemoryLoader) Setup() {
}

func (ml *MemoryLoader) TwitchGetRoom(_ string) string {
	return ""
}

func (ml *MemoryLoader) TwitchSetRoom(_, _ string) error {
	return nil
}

func (ml *MemoryLoader) SaveRoom(_ string, _ *LoadJSON) error {
	return nil
}

func (ml *MemoryLoader) LoadRoom(_ string) (*LoadJSON, error) {
	return nil, nil
}

func (ml *MemoryLoader) LoadAllRooms() []*LoadJSON {
	return nil
}

func (ml *MemoryLoader) DeleteRoom(_ string) error {
	return nil
}

func (ml *MemoryLoader) LoadAllMessages() []*MessageJSON {
	return nil
}

func (ml *MemoryLoader) DeleteAllMessages() {
}
