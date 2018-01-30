package main

import (
	"log"
	"os"
	"testing"
)

func cleanRedis() error {
	client := newRedisClient()
	return client.Del(StatusKey).Err()
}

func TestRedisSetAndGet(t *testing.T) {
	// get when no key is set
	val, err := redisGetStatus()
	if err != nil {
		t.Errorf("get when no record is set, %s", err)
	} else if val != StatusNormal {
		t.Errorf("get when no record is set returned %s, expected %s", val, StatusNormal)
	}

	if err = redisSetStatus("hoge"); err != nil {
		t.Errorf("set status failed %s", err)
	}
	val, err = redisGetStatus()
	if err != nil {
		t.Errorf("get status failed %s", err)
	}
	if val != "hoge" {
		t.Errorf("get status returned %s, expected %s", val, "hoge")
	}
}

func TestMain(m *testing.M) {
	if err := cleanRedis(); err != nil {
		log.Fatalf("cleaning redis before TestRedisSetAndGet, %s", err)
	}

	code := m.Run()

	if err := cleanRedis(); err != nil {
		log.Fatalf("cleaning redis after TestRedisSetAndGet, %s", err)
	}

	os.Exit(code)
}
