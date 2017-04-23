package main

import (
	"github.com/speps/go-hashids"
	"strconv"
	"strings"
)

func encrypt(text string) (string, error) {
	hd := hashids.NewData()
	hd.Salt = ENCRYPTION_KEY
	hd.MinLength = 32
	h := hashids.NewWithData(hd)
	integer, _ := strconv.Atoi(text)
	numbers := []int{integer}
	return h.Encode(numbers)
}

func decrypt(cryptoText string) (string, error) {
	hd := hashids.NewData()
	hd.Salt = ENCRYPTION_KEY
	hd.MinLength = 32
	h := hashids.NewWithData(hd)
	d, err := h.DecodeWithError(cryptoText)
	return strconv.Itoa(d[0]), err
}

func makeAddress(image_id string) string {
	split_string := splitImageId(image_id)
	return BASE_IMAGE_FOLDER + "/" + strings.Join(split_string, "/") + IMAGE_SUFFIX
}

func splitImageId(image_id string) []string {
	split_string := []string{}
	for i := 0; i < len(image_id); i++ {
		if i % 2 == 1 {
			split_string = append(split_string, image_id[i - 1:i + 1])
		} else if i + 1 == len(image_id) {
			split_string = append(split_string, image_id[i:i + 1])
		}
	}
	return split_string
}