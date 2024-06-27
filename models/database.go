package models

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/url"
	"strings"
	"time"

	"gorm.io/gorm"
)

type URL struct {
	ID           string    `json:"id"`            //done
	OriginalURL  string    `json:"original_url"`  //done
	SchemeURL    string    `json:"scheme_url"`    //done
	TrimmedURL   string    `json:"trimmed_url"`   //done
	ShortenedURL string    `json:"shortened_url"` //done
	CreationTime time.Time `json:"creation_time"` //will be passed from FE
}

var DbCon *gorm.DB

func (longURL *URL) SimplifyURL() {
	u, err := url.Parse(longURL.OriginalURL)
	if err != nil {
		panic(err)
	}
	// Hostname trimmed
	u.Host = strings.TrimPrefix(u.Host, "www.")
	// Scheme explicitly enforced as https if none
	if u.Scheme == "" {
		u.Scheme = "https"
	}
	// Removed trailing '/' end of Path
	u.Path = strings.TrimSuffix(u.Path, "/")

	longURL.OriginalURL = u.Scheme + "://" + u.Host + u.Path
	longURL.TrimmedURL = u.Host + u.Path
}

func (url *URL) HashURL() {
	hasher := md5.New()
	// Write the simplified URL to the hasher
	hasher.Write([]byte(url.TrimmedURL))
	// Finalize the hash and retrieve the hash bytes
	hashBytes := hasher.Sum(nil)
	// Convert hash bytes to hexadecimal string
	url.ShortenedURL = hex.EncodeToString(hashBytes)
}

func ConnectDB(db *gorm.DB) {
	DbCon = db
}

func MigrateSchema() error {
	err := DbCon.AutoMigrate(URL{})
	if err != nil {
		return fmt.Errorf("failed to migrate database schema: %v", err)
	}
	return nil
}

func FlushIntoDB(url *URL) error {
	result := DbCon.Create(url)
	if result.Error != nil {
		return fmt.Errorf("failed to insert URL into database: %v", result.Error)
	}
	return nil
}

func GetOriginalURL(short7 string) (string, error) {
	var url URL
	// Assuming DbCon is your *gorm.DB connection instance
	result := DbCon.Where("LEFT(shortened_url, 7) = ?", short7).First(&url)
	if result.Error != nil {
		return "", fmt.Errorf("failed to retrieve original URL: %v", result.Error)
	}
	return url.OriginalURL, nil
}
