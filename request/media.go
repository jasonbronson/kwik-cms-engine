package request

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jasonbronson/kwik-cms-engine/config"
	"github.com/jasonbronson/kwik-cms-engine/model"
	"github.com/jasonbronson/kwik-cms-engine/repositories"

	S3AWS "github.com/jasonbronson/kwik-cms-engine/library/helpers/s3"
	"github.com/jasonbronson/kwik-cms-engine/request/response"
	_ "github.com/jasonbronson/kwik-cms-engine/request/response"
)

func GetMedia(g *gin.Context) {
	db := config.Cfg.GormDB
	db = db.WithContext(g)
	params := defaultParameters(g)
	r := repositories.GetMedia(db, params)
	response.Standard(g.Writer, http.StatusOK, r)
}

func SetMedia(g *gin.Context) {
	db := config.Cfg.GormDB
	db = db.WithContext(g)

	g.Request.ParseMultipartForm(300 << 20) // Max 300 MB memory used

	s3URL, name, e := mediaUpload(g)
	if e != nil {
		response.ErrorResponse(g.Writer, http.StatusInternalServerError, e)
		return

	}
	log.Println("DEBUG URL:", s3URL)

	media := model.Media{
		URL:  s3URL,
		Name: name,
	}
	mediaID, e := repositories.CreateMedia(db, media)
	if e != nil {
		response.ErrorResponse(g.Writer, http.StatusInternalServerError, e)
		return
	}

	response.Action(g.Writer, http.StatusOK, fmt.Sprint(mediaID))
}

func mediaUpload(g *gin.Context) (url, name string, err error) {
	r := g.Request

	// Maximum upload of 10 MB files
	//r.ParseMultipartForm(10 << 20)

	// Get handler for filename, size and headers
	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Println("Error Retrieving the File")
		log.Println(err)
		return "", "", err
	}

	defer file.Close()
	log.Printf("Uploaded Filename: %+v\n", handler.Filename)
	log.Printf("File Size: %+v\n", handler.Size)
	log.Printf("MIME Header: %+v\n", handler.Header)

	// Create temp file locally
	dst, err := os.Create(handler.Filename)
	if _, err := io.Copy(dst, file); err != nil {
		return "", "", err
	}
	log.Println("File Created ", dst, handler)
	defer dst.Close()
	if err != nil {
		return "", "", err
	}

	//Upload to S3
	url, err = S3AWS.S3FileUpload(handler.Filename, "")
	if err != nil {
		log.Println("can not upload file to s3", err)
		e := os.Remove(handler.Filename)
		if e != nil {
			log.Println(e)
			return "", "", e
		}
		return "", "", err
	}

	return url, handler.Filename, nil
}

func DeleteMedia(g *gin.Context) {
	db := config.Cfg.GormDB
	db = db.WithContext(g)
	id := g.Param("mediaid")
	e := repositories.DeleteMedia(db, id)
	if e != nil {
		response.ErrorResponse(g.Writer, http.StatusInternalServerError, e)
		return
	}
	response.Action(g.Writer, http.StatusOK, "success")
}
