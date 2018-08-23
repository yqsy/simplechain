package main

import (
	"github.com/boltdb/bolt"
	"fmt"
	"os"
	"math/rand"
)

func getRandomBit(len int) []byte {
	token := make([]byte, len)
	rand.Read(token)
	return token
}

func main() {
	os.Remove("my.db")
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		panic(err)
	}

	defer func() {
		db.Close()
		os.Remove("my.db")
	}()

	// 增
	update_(db)
	select_(db)

	// 改
	modify_(db)
	select_(db)

	// 删
	delete_(db)
	select_(db)

	// 遍历所有键
	batch_(db)
}

func batch_(db *bolt.DB) {
	if err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("MyBucket"))

		if b == nil {
			panic("not exist")
		}

		for i := 0; i < 10; i++ {
			b.Put([]byte{byte(i)}, []byte{byte(i)})
		}

		return nil
	}); err != nil {
		panic(err)
	}

	if err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("MyBucket"))

		if b == nil {
			panic("not exist")
		}

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			fmt.Printf("key=%x, value=%x\n", k, v)
		}

		return nil
	}); err != nil {
		panic(err)
	}
}

func update_(db *bolt.DB) {
	if err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte("MyBucket"))
		if err != nil {
			panic(err)
		}
		err = b.Put([]byte("2018"), []byte("hello"))
		return err
	}); err != nil {
		panic(err)
	}
}

func delete_(db *bolt.DB) {
	if err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("MyBucket"))

		if b == nil {
			panic("not exist")
		}

		err := b.Delete([]byte("2018"))
		return err
	}); err != nil {
		panic(err)
	}
}

func modify_(db *bolt.DB) {
	if err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("MyBucket"))

		if b == nil {
			panic("not exist")
		}

		err := b.Put([]byte("2018"), []byte("world"))
		return err
	}); err != nil {
		panic(err)
	}
}

func select_(db *bolt.DB) {
	if err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("MyBucket"))

		if b == nil {
			panic("not exist")
		}

		v := b.Get([]byte("2018"))

		fmt.Printf("%s\n", v)

		return nil
	}); err != nil {
		panic(err)
	}
}
