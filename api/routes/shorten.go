package routes

import (
	"os"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/karthikbhandary2/url-shortener/database"
	"github.com/karthikbhandary2/url-shortener/helpers"
)

type request struct {
	URL         string        `json:"url"`
	CustomShort string        `json:"short"`
	Expiry      time.Duration `json:"expiry"`
}

type response struct {
	URL             string        `json:"url"`
	CustomShort     string        `json:"short"`
	Expiry          time.Duration `json:"expiry"`
	XRateRemaining  int64         `json:"rate_limit"`
	XRateLimitReset time.Duration `json:"rate_limit_reset"`
}

func ShortenURL(c *fiber.Ctx) error {
	body := new(request)
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
	}

	//rate limiting
	redisClient := database.CreateClient(1)
	defer redisClient.Close()

	value, err := redisClient.Get(database.Ctx, c.IP()).Result()
	if err == redis.Nil {
		_ = redisClient.Set(database.Ctx, c.IP(), os.Getenv("API_QUOTA"), 30*time.Minute).Err()
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot connect to the DB"})
	} else {
		value, _ = redisClient.Get(database.Ctx, c.IP()).Result()
		val, _ := strconv.Atoi(value)

		if val <= 0 {
			limit, _ := redisClient.TTL(database.Ctx, c.IP()).Result()
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "rate limit exceeded", "rate_limit_reset": limit / time.Nanosecond / time.Minute})
		}
	}
	// check if the input is an actual url
	if !govalidator.IsURL(body.URL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid URL"})
	}

	//check for domain error
	if !helpers.RemoveDomainError(body.URL) {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "you cant hack the system"})
	}

	// enforce https, SSL
	body.URL = helpers.EnforceHTTP(body.URL)

	// check if the custom short url is already in use
	var id string
	if body.CustomShort == "" {
		id = uuid.New().String()[:6]
	} else {
		id = body.CustomShort
	}

	r := database.CreateClient(0)
	defer r.Close()

	value, _ = r.Get(database.Ctx, id).Result()
	if value != "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "URL custom short is already in use"})
	}

	//set the default expiry time to 24 hours if user does not provide one
	if body.Expiry == 0 {
		body.Expiry = 24
	}

	err = r.Set(database.Ctx, id, body.URL, body.Expiry*time.Hour).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot connect to the DB"})
	}

	// response
	resp := response{
		URL:             body.URL,
		CustomShort:     "",
		Expiry:          body.Expiry,
		XRateRemaining:  10,
		XRateLimitReset: 30,
	}

	//decrease the quota after func call
	redisClient.Decr(database.Ctx, c.IP())

	val, _ := redisClient.Get(database.Ctx, c.IP()).Result()
	if err == redis.Nil {
    resp.XRateRemaining = 0  // or some default value
	} else if err != nil {
		// handle other errors
	} else {
		intVal, _ := strconv.Atoi(val)
		resp.XRateRemaining = int64(intVal)
	}
	ttl, _ := redisClient.TTL(database.Ctx, c.IP()).Result()
	resp.XRateLimitReset = ttl/ time.Nanosecond / time.Minute

	resp.CustomShort = os.Getenv("DOMAIN") + "/" + id
	return c.Status(fiber.StatusOK).JSON(resp)
}
