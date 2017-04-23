package main

import (
	"strings"
	"os"
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
	"io"
)

func serveImage(c *gin.Context) {
	name := c.Param("image")
	image_id_encrypted := strings.Split(name, ".")[0]

	image_id, err := decrypt(image_id_encrypted);
	if err != nil {
		c.String(http.StatusOK, "The requestd file not exist! (%s) - ERROR 1", name)
		return
	}

	file_address := makeAddress(image_id)
	if _, err := os.Stat(file_address); os.IsNotExist(err) {
		c.String(http.StatusOK, "The requestd file not exist! (%s) - ERROR 2", name)
		return
	}

	showImage(c, file_address)
}

func encryptId(c *gin.Context) {
	id := c.Param("id")
	data, _ := encrypt(id)
	c.String(http.StatusOK, data)
}

func upload(c *gin.Context) {
	group := c.DefaultPostForm("group", "No Group")
	tags := c.PostFormArray("tags")

	insert_tags := []*Tag{}
	for _, tag := range tags {
		temp := &Tag{Name:tag}
		insert_tags = append(insert_tags, temp.firstOrCreate())
	}

	insert_group := &Group{Name:group}
	insert_group.firstOrCreate()

	insert_image := &ImageDb{GroupId:insert_group.Id}
	dbmap.Insert(insert_image)

	insert_image.AddTags(insert_tags)

	file, _, _ := c.Request.FormFile("file")

	insert_image.MakeDirectories()

	save_file, _ := os.Create(insert_image.Path)
	defer save_file.Close()

	io.Copy(save_file, file)

	gin_tags := []string{}
	for _, tag := range insert_tags {
		gin_tags = append(gin_tags, tag.Name)
	}

	c.JSON(http.StatusOK, gin.H{
		"encrypted_id" : insert_image.EncryptedPath,
		"file_address": insert_image.GetUrlAddress(),
		"tags": gin_tags,
		"group": gin.H{"name":group},
	})
}
