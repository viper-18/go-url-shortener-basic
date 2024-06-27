package routes

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/viper-18/go-url-shortener-basic/models"
)

func ShortenURL(c *fiber.Ctx) error {
	body := models.URL{}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
	}
	// modified original_url and trimmed_url
	body.SimplifyURL()
	//now we have to hash the trimmed url
	body.HashURL()
	//flush it in the database
	models.FlushIntoDB(&body)

	godotenv.Load()
	shortUrl := body.ShortenedURL[:7]
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"newUrl": os.Getenv("HOST_SCHEME") + "://" + os.Getenv("HOST_NAME") + "/" + shortUrl,
	})
}

func NormaliseURL(c *fiber.Ctx) error {
	shortUrl := c.Params("url")
	targetURL, _ := models.GetOriginalURL(shortUrl)
	return c.Redirect(targetURL)
}
