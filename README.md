# Go CDN
this is a mini cdn server that can save and manage your image files.
This utility can be used as stand alone Content Delivery Network, using Mysql as database .

# features
  - on the fly resize image/jpeg
  - cache strategy
  - additional metadata for each file `category` and `tags`
  
# Environments
   - _BASE_URL=http://localhost:8080/ `#Base url`
   - _HOST=0.0.0.0:8080 `#server ip and port`
   - _MYSQL_HOST=127.0.0.1  `#mysql config`
   - _MYSQL_PORT=3306 `#mysql config`
   - _MYSQL_DB=go_cdn `#mysql config`
   - _MYSQL_USER=root `#mysql config`
   - _MYSQL_PASS=root `#mysql config`
   - _ENCRYPTION_KEY=TESTKEY `#encrypt image id with this key`
   - _BASE_IMAGE_FOLDER=/your/image/store/path `#save images to this folder`
   

# dependencies

   - gopkg.in/gin-gonic/gin.v1
   - github.com/gin-contrib/cache
   - github.com/gin-contrib/cache/persistence
   - gopkg.in/gorp.v1
   - github.com/ziutek/mymysql/godrv
   - github.com/speps/go-hashids