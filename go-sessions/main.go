package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/coherence"
	"log"
	"strings"
	"time"
)

var store *session.Store

func main() {
	// create new coherence session store using defaults of localhost:1408
	storage, err := coherence.New()

	if err != nil {
		log.Fatal("unable to connect to Coherence", err)
	}
	defer storage.Close()

	// initialize the gofiber session store using the Coherence storage driver
	store = session.New(session.Config{
		Storage:    storage,
		Expiration: time.Duration(60) * time.Second,
	})

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		// retrieve the session
		sess, err1 := store.Get(c)
		if err1 != nil {
			log.Println(err1)
			return c.Status(500).SendString(err1.Error())
		}

		if sess.Fresh() {
			// new session
			sess.Set("accessCount", 1)
			sess.Set("firstAccess", time.Now().Format(time.ANSIC))
		} else {
			// increment the number of times we have hit this endpoint
			count := sess.Get("accessCount").(int)
			count++
			sess.Set("accessCount", count)
			sess.Set("lastAccess", time.Now().Format(time.ANSIC))
		}

		var sb strings.Builder

		sb.WriteString(fmt.Sprintf("Session: %s, new=%v\nSession values:\n", sess.ID(), sess.Fresh()))
		for _, k := range sess.Keys() {
			sb.WriteString(fmt.Sprintf("   %s=%v\n", k, sess.Get(k)))
		}

		if err1 = sess.Save(); err1 != nil {
			log.Println(err1)
			return c.Status(500).SendString(err1.Error())
		}

		return c.SendString(sb.String())
	})

	app.Get("/destroy", func(c *fiber.Ctx) error {
		// retrieve the session
		sess, err1 := store.Get(c)
		if err1 != nil {
			log.Println(err1)
			return c.Status(500).SendString(err1.Error())
		}
		id := sess.ID()

		// remove the session
		err1 = sess.Destroy()
		if err1 != nil {
			log.Println(err1)
			return c.Status(500).SendString(err1.Error())
		}
		return c.Status(200).SendString(fmt.Sprintf("session %v destroyed", id))
	})

	panic(app.Listen("127.0.0.1:2000"))
}
