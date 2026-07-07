package geo

import (
	"encoding/json"
	"errors"
	"os"
)

type Cache struct {
	path string
	data map[string]Coordinates
}

func NewCache(path string) (*Cache, error) {
	c := &Cache{
		path: path,
		data: make(map[string]Coordinates),
	}

	file, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return c, nil
		}
		return nil, err
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&c.data); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Cache) Get(location string) (Coordinates, bool) {
	coords, ok := c.data[location]
	return coords, ok
}

func (c *Cache) Set(location string, coords Coordinates) {
	c.data[location] = coords
}

func (c *Cache) Save() error {
	file, err := os.Create(c.path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	return encoder.Encode(c.data)
}
