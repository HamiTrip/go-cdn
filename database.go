package main

import (
	"database/sql"
	"gopkg.in/gorp.v1"
	"os"
	_ "github.com/ziutek/mymysql/godrv"
	"strconv"
	"strings"
)

var (
	dbmap *gorp.DbMap
)

type Group struct {
	Id   int64 `db:"id"`
	Name string `db:"name"`
}

func (group *Group) firstOrCreate() *Group{
	err := dbmap.SelectOne(group, "select * from groups where name=?", group.Name)
	if checkErr(err) {
		dbmap.Insert(group)
	}
	return group
}

type ImageDb struct {
	Id            int64 `db:"id"`
	Path          string `db:"path"`
	EncryptedPath string `db:"encrypted_path"`
	GroupId       int64 `db:"group_id"`
}

func (img *ImageDb) CreateHash() {
	img.EncryptedPath, _ = encrypt(strconv.FormatInt(img.Id, 10))
}

func (img *ImageDb) MakeDirectories() {
	folders := strings.Split(img.Path, "/")
	os.MkdirAll(strings.Join(folders[0:len(folders)-1], "/"), 0755)
}

func (img *ImageDb) MakePath() {
	img.Path = makeAddress(strconv.FormatInt(img.Id, 10))
}

func (img *ImageDb) GetUrlAddress() string{
	return BASE_URL + "images/" + img.EncryptedPath + ".jpg"
}

func (img *ImageDb) AddTags(tags []*Tag)  {
	for _, tag := range tags {
		dbmap.Insert(&TagImage{TagId:tag.Id, ImageId:img.Id})
	}
}

func (img *ImageDb)PostInsert(s gorp.SqlExecutor) error {
	img.CreateHash()
	img.MakePath()
	dbmap.Update(img)
	return nil
}

type Tag struct {
	Id   int64   `db:"id"`
	Name string  `db:"name"`
}

func (tag *Tag) firstOrCreate() *Tag {
	err := dbmap.SelectOne(tag, "select * from tags where name=?", tag.Name)
	if checkErr(err) {
		dbmap.Insert(tag)
	}
	return tag
}

type TagImage struct {
	Id      int64   `db:"id"`
	TagId   int64 `db:"tag_id"`
	ImageId int64 `db:"image_id"`
}

func connectMysql() {
	db, err := sql.Open(
		"mymysql",
		"tcp:" + os.Getenv("_MYSQL_HOST") + ":" + os.Getenv("_MYSQL_PORT") + "*" + os.Getenv("_MYSQL_DB") + "/" + os.Getenv("_MYSQL_USER") + "/" + os.Getenv("_MYSQL_PASS"))
	if err != nil {
		panic(err)
	}
	dbmap = &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}

	dbmap.AddTableWithName(Group{}, "groups").SetKeys(true, "Id")
	dbmap.AddTableWithName(ImageDb{}, "images").SetKeys(true, "Id")
	dbmap.AddTableWithName(Tag{}, "tags").SetKeys(true, "Id")
	dbmap.AddTableWithName(TagImage{}, "tag_image").SetKeys(true, "Id")

	dbmap.CreateTablesIfNotExists()
}

func checkErr(err error) bool {
	if err != nil {
		//log.Fatalln(msg, err)
		return true
	}
	return false
}