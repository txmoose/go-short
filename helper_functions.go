package main

import (
	"crypto/rand"
	"errors"
	"math/big"
	"strings"

	"github.com/anaskhan96/soup"
)

// GenerateSlug Found and modified from here
// https://gist.github.com/dopey/c69559607800d2f2f90b1b1ed4e550fb#gistcomment-3527095
func GenerateSlug(n int) (string, error) {
	//goland:noinspection ALL
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}

	var slug Slug
	slug, _ = GetSlugFromDB(string(ret))
	candidateSlug := string(ret)

	// Generate a slug, but if that slug matches an existing slug, recurse to generate another
	if candidateSlug == slug.Slug {
		candidateSlug, err = GenerateSlug(n)
	}
	return candidateSlug, nil
}

//GetSiteTitle Should go out and get the title of the web page at a given domain
func GetSiteTitle(url string) (string, error) {
	//Telling soup to go get the site
	resp, err := soup.Get(url)
	if err != nil {
		return "", err
	}

	//doc to hold the actual HTML document
	doc := soup.HTMLParse(resp)

	//Getting the title out of doc
	title := doc.Find("title").Text()

	if title == "" {
		return title, errors.New("no title found")
	}

	return strings.TrimSpace(title), nil
}
