package main

import (
	"encoding/csv"
	"flag"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/garyburd/redigo/redis"
)

var (
	err       error
	logStderr *log.Logger

	redisHost = flag.String("redis_host", "127.0.0.1", "-redis_host=127.0.0.1")
	redisPort = flag.String("redis_port", "6379", "-redis_port=6379")
	redisAuth = flag.String("redis_auth", "", "-redis_auth=MyPasswd")
	flushall  = flag.Bool("flushall", false, "-flushall")
	debug     = flag.Bool("debug", false, "-debug")
)

func main() {
	logStderr = log.New(os.Stderr, "", log.Ldate|log.Lmicroseconds|log.Lshortfile)

	flag.Parse()

	pool := newRedisPool(*redisHost+":"+*redisPort, *redisAuth)
	// Check redis connect
	if _, err = pool.Dial(); err != nil {
		logStderr.Fatalln("Redis Driver Error", err)
	}

	redispool := pool.Get()
	defer redispool.Close()

	rcsv := csv.NewReader(os.Stdin)
	rcsv.Comma = ';'
	rcsv.Comment = '#'
	rcsv.LazyQuotes = true

	if *flushall {
		_, err = redispool.Do("FLUSHALL")
		if err != nil {
			logStderr.Fatal(err)
		}
	}

	for {
		record, err := rcsv.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			logStderr.Fatal(err)
		}

		parsePrefix(redispool, "7"+strings.TrimSpace(record[0])+strings.TrimSpace(record[1]), "7"+strings.TrimSpace(record[0])+strings.TrimSpace(record[2]), strings.TrimSpace(record[4])+";"+strings.TrimSpace(record[5]))
	}

	err = redispool.Flush()
	if err != nil {
		logStderr.Fatal(err)
	}
}

// загружает данные в Redis
func load_csv(w http.ResponseWriter, r *http.Request) {

}

func parsePrefix(redispool redis.Conn, min, max, operator string) {
	if *debug {
		logStderr.Printf("%s\t%s\t%s\n", min, max, operator)
	}

	min_len := len(min)
	max_len := len(max)
	if min_len != max_len {
		logStderr.Fatalf("Invalid len min %d != len max %d\n", min_len, max_len)
	}

	minuint64, err := strconv.ParseUint(min, 10, 64)
	if err != nil {
		logStderr.Fatalln(err)
	}

	maxuint64, err := strconv.ParseUint(max, 10, 64)
	if err != nil {
		logStderr.Fatalln(err)
	}

	for minuint64 < maxuint64 {
		var mask uint64 = 0
		var pow uint64 = 0
		for i := 0; ; i++ {
			pow = uint64(math.Pow10(min_len - i))
			mask = minuint64 / pow
			if minuint64+pow-1 <= maxuint64 {
				mask = minuint64 / pow

				if *debug {
					logStderr.Printf("%d %s\n", mask, operator)
				}

				err = redispool.Send("SET", mask, operator)
				if err != nil {
					logStderr.Fatal(err)
				}
				break
			}
		}
		minuint64 = (mask + 1) * pow
	}
}

func newRedisPool(server, password string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		MaxActive:   1024,
		IdleTimeout: 180 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
