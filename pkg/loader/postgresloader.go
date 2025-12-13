/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package loader

import (
	"database/sql"
	"time"

	// blank import necessary for postgres driver
	_ "github.com/lib/pq"
)

func NewPostgresLoader(dsn string) (*BaseLoader, error) {
	var err error
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Verify connection is good
	healthy := false
	for i := 0; i < 10; i++ {
		if err = db.Ping(); err == nil {
			healthy = true
			break
		}
		time.Sleep(time.Second)
	}

	if !healthy {
		return nil, err
	}

	ldr := &BaseLoader{
		db:     db,
		dbType: DBTypePostgres,
	}
	err = ldr.setup()

	if err != nil {
		return nil, err
	}

	return ldr, nil
}
