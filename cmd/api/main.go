package main

import (
	"context"

	"github.com/go-redis/redis/v8"
)

var (
	ctx    = context.Background()
	client = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
)

type Product struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string `gorm:"size:100"`
	Price int
}

// TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>
func main() {

}
