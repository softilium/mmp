package main

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"net/http"
	"os"
	"time"

	"github.com/hashicorp/golang-lru/v2/expirable"
	"github.com/nfnt/resize"
)

var images *expirable.LRU[string, *[]byte]
var thumbs *expirable.LRU[string, *[]byte]

func initRouterImages(router *http.ServeMux) {
	router.HandleFunc("/api/goods/images", func(w http.ResponseWriter, r *http.Request) {
		handleImages(w, r, false)
	})
	router.HandleFunc("/api/goods/thumbs", func(w http.ResponseWriter, r *http.Request) {
		handleImages(w, r, true)
	})
}

func handleImages(w http.ResponseWriter, r *http.Request, downScale bool) {

	ref := r.URL.Query().Get("ref")
	if ref == "" {
		HandleErr(w, http.StatusBadRequest, fmt.Errorf("goodref is required"))
		return
	}
	n := r.URL.Query().Get("n")
	if n == "" {
		HandleErr(w, http.StatusBadRequest, fmt.Errorf("image number is required"))
		return
	}
	fn := fmt.Sprintf("%s/goodImage-%s-%s", Cfg.ImagesFolder, ref, n)

	cache := images
	if downScale {
		cache = thumbs
	}

	if r.Method == http.MethodGet {

		imageData, ok := cache.Get(fn)
		if !ok {
			i, err := os.ReadFile(fn)
			if err != nil {
				cache.Add(fn, nil)
				w.WriteHeader(http.StatusNoContent) // 204 means no image found
				return
			}
			imageData = &i

			if downScale {
				// Decoding gives you an Image.
				// If you have an io.Reader already, you can give that to Decode
				// without reading it into a []byte.
				image, _, err := image.Decode(bytes.NewReader(*imageData))
				if err != nil {
					HandleErr(w, http.StatusInternalServerError, fmt.Errorf("failed to decode image: %v", err))
					return
				}

				newImage := resize.Resize(60, 60, image, resize.Lanczos3)

				var buf bytes.Buffer
				err = jpeg.Encode(&buf, newImage, nil)
				if err != nil {
					HandleErr(w, http.StatusInternalServerError, fmt.Errorf("failed to encode resized image: %v", err))
					return
				}
				i := buf.Bytes()
				imageData = &i
			}
			cache.Add(fn, imageData)
		}

		if imageData == nil {
			w.WriteHeader(http.StatusNoContent) // 204 means no image found
			return
		}

		w.Header().Set("Content-Type", "image/jpeg") // Adjust based on your image type
		_, err = w.Write(*imageData)
		if err != nil {
			HandleErr(w, http.StatusInternalServerError, fmt.Errorf("failed to write image data: %v", err))
			return
		}
	}

	//user, err := models.UserFromHttpRequest(r)

}

func init() {
	images = expirable.NewLRU[string, *[]byte](1000, nil, time.Minute*60*24)
	thumbs = expirable.NewLRU[string, *[]byte](1000, nil, time.Minute*60+24)
}
