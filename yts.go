package yts

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// APIEndpoint represents the API's endpoint
const APIEndpoint = "https://yts.ag/api/v2"

// Sort options
const (
	SortByTitle     = "title"
	SortByYear      = "year"
	SortByRating    = "rating"
	SortByPeers     = "peers"
	SortBySeeds     = "seeds"
	SortByDownload  = "download_count"
	SortByLike      = "like_count"
	SortByDateAdded = "date_added"
)

// Order options
const (
	OrderAsc  = "asc"
	OrderDesc = "desc"
)

// Movie represents the movies
type Movie struct {
	DateUploaded     string    `json:"date_uploaded"`
	DateUploadedUnix int64     `json:"date_uploaded_unix"`
	Genres           []string  `json:"genres"`
	ID               int       `json:"id"`
	ImdbID           string    `json:"imdb_code"`
	Language         string    `json:"language"`
	MediumCover      string    `json:"medium_cover_image"`
	Rating           float64   `json:"rating"`
	Runtime          int       `json:"runtime"`
	SmallCover       string    `json:"small_cover_image"`
	State            string    `json:"state"`
	Title            string    `json:"title"`
	TitleLong        string    `json:"title_long"`
	Torrents         []Torrent `json:"torrents"`
	Year             int       `json:"year"`
}

// Torrent represents the quality for a torrent
type Torrent struct {
	DateUploaded     string `json:"date_uploaded"`
	DateUploadedUnix int64  `json:"date_uploaded_unix"`
	Hash             string `json:"hash"`
	Peers            int    `json:"peers"`
	Quality          string `json:"quality"`
	Seeds            int    `json:"seeds"`
	Size             string `json:"size"`
	SizeBytes        int    `json:"size_bytes"`
	URL              string `json:"url"`
}

// Data represents the data inside the response body
type Data struct {
	PageNumber int     `json:"int"`
	Movies     []Movie `json:"movies"`
}

// Result represents the response from the API
type Result struct {
	Status        string `json:"status"`
	StatusMessage string `json:"status_message"`
	Data          Data   `json:"data"`
}

func getMovieList(v url.Values) ([]Movie, error) {
	endpoint := APIEndpoint + "/list_movies.json?" + v.Encode()

	var httpClient = &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := httpClient.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("got error %d while getting movie list", resp.StatusCode)
	}

	var result Result
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result.Data.Movies, nil
}

// GetList gets a list of movies from a page number
func GetList(pageNumber, minRating int, sort, order string) ([]Movie, error) {
	v := url.Values{}
	v.Set("limit", "50")
	v.Set("sort_by", sort)
	v.Set("order_by", order)
	v.Set("minimum_rating", strconv.Itoa(minRating))
	v.Set("page", strconv.Itoa(pageNumber))
	return getMovieList(v)
}

// Search searches movies
func Search(movieTitle string) ([]Movie, error) {
	v := url.Values{}
	v.Set("query_term", movieTitle)
	return getMovieList(v)
}

// Status returns an error if the YTS API is not responding correctly
func Status() error {
	v := url.Values{}
	v.Set("limit", "1")
	v.Set("sort_by", SortByPeers)
	v.Set("order_by", OrderDesc)
	v.Set("minimum_rating", "6")

	movies, err := getMovieList(v)
	if err != nil {
		return err
	}
	if len(movies) == 0 {
		return fmt.Errorf("no movies returned")
	}

	return nil
}
