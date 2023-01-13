package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Cache struct {
	cacheFile string
	cache     map[string]time.Time
	disabled  bool
}

func NewCache(path string) *Cache {
	c := &Cache{
		cacheFile: path,
		cache:     make(map[string]time.Time),
		disabled:  false,
	}
	return c
}

func (c *Cache) Disable() {
	c.disabled = true
}

func (c *Cache) Add(path string, timestamp time.Time) {
	if c.disabled {
		return
	}
	c.cache[path] = timestamp
}

func (c *Cache) Expired(path string, duration time.Duration) bool {
	if c.disabled {
		return true
	}

	timestamp, ok := c.cache[path]
	if !ok {
		return true
	}

	if timestamp.Before(time.Now().Add(-1 * duration)) {
		return true
	}

	return false
}

func (c *Cache) Dump() (err error) {
	if c.disabled {
		return nil
	}
	dmp := bytes.Buffer{}
	dmp.Reset()

	for k, v := range c.cache {
		str := fmt.Sprintf("%s::%d\n", k, v.Unix())
		dmp.WriteString(str)
	}

	err = os.Remove(c.cacheFile)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	f, err := os.OpenFile(c.cacheFile, os.O_CREATE|os.O_WRONLY, 0644)
	defer func() {
		if f != nil {
			err = f.Close()
		}
	}()

	if err != nil {
		return
	}

	_, err = f.Write(dmp.Bytes())
	if err != nil {
		return
	}

	return nil
}

func (c *Cache) Load() (err error) {
	if c.disabled {
		return nil
	}

	f, err := os.Open(c.cacheFile)
	defer func() {
		if f != nil {
			err = f.Close()
		}
	}()

	if err != nil {
		return err
	}

	b := bufio.NewScanner(f)
	b.Split(bufio.ScanLines)

	for b.Scan() {
		splitLine := strings.Split(b.Text(), "::")
		mpath := splitLine[0]
		timestamp := splitLine[1]
		timeInt, err := strconv.ParseInt(timestamp, 10, 64)
		if err != nil {
			return err
		}
		timeParsed := time.Unix(timeInt, 0)
		c.Add(mpath, timeParsed)
	}

	return nil
}
