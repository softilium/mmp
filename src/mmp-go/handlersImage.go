package main

import (
	"bytes"
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/hashicorp/golang-lru/v2/expirable"
	"github.com/nfnt/resize"
)

var imagesCache *expirable.LRU[string, *[]byte]
var thumbsCache *expirable.LRU[string, *[]byte]

func initRouterImages(router *http.ServeMux) {
	router.HandleFunc("/api/goods/images", func(w http.ResponseWriter, r *http.Request) {
		handleImages(w, r, false, "good")
	})
	router.HandleFunc("/api/goods/thumbs", func(w http.ResponseWriter, r *http.Request) {
		handleImages(w, r, true, "good")
	})
	router.HandleFunc("/api/shops/images", func(w http.ResponseWriter, r *http.Request) {
		handleImages(w, r, false, "shop")
	})
	router.HandleFunc("/api/shops/thumbs", func(w http.ResponseWriter, r *http.Request) {
		handleImages(w, r, true, "shop")
	})
}

func handleImages(w http.ResponseWriter, r *http.Request, downScale bool, prefix string) {

	ref := r.URL.Query().Get("ref")
	if ref == "" {
		HandleErr(w, http.StatusBadRequest, fmt.Errorf("ref is required"))
		return
	}
	n := r.URL.Query().Get("n")
	if n == "" {
		HandleErr(w, http.StatusBadRequest, fmt.Errorf("n (image number) is required"))
		return
	}
	fn := fmt.Sprintf("%s/%sImage-%s-%s", Cfg.ImagesFolder, prefix, ref, n)

	cache := imagesCache
	if downScale {
		cache = thumbsCache
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
				img, _, err := image.Decode(bytes.NewReader(*imageData))
				if err != nil {
					HandleErr(w, http.StatusInternalServerError, fmt.Errorf("failed to decode image: %v", err))
					return
				}

				var newImage image.Image

				img_bnd := img.Bounds()
				img_w := img_bnd.Max.X - img_bnd.Min.X
				img_h := img_bnd.Max.Y - img_bnd.Min.Y
				if img_w > img_h {
					newImage = resize.Resize(60, 0, img, resize.Lanczos3)
				} else {
					newImage = resize.Resize(0, 60, img, resize.Lanczos3)
				}

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

	if r.Method == http.MethodDelete {
		if _, err = os.Stat(fn); os.IsNotExist(err) {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		if err := os.Remove(fn); err != nil {
			HandleErr(w, http.StatusInternalServerError, fmt.Errorf("failed to delete image: %v", err))
			return
		}
		imagesCache.Remove(fn)
		thumbsCache.Remove(fn)
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if r.Method == http.MethodPost {

		err = r.ParseMultipartForm(10 << 20) // Limit to 10MB
		if err != nil {
			HandleErr(w, http.StatusBadRequest, fmt.Errorf("failed to parse form: %v", err))
			return
		}
		file, _, err := r.FormFile("image")
		if err != nil {
			HandleErr(w, http.StatusBadRequest, fmt.Errorf("failed to get image from form: %v", err))
			return
		}
		defer func() {
			_ = file.Close()
		}()
		imageData, err := io.ReadAll(file)
		if err != nil {
			HandleErr(w, http.StatusInternalServerError, fmt.Errorf("failed to read image data: %v", err))
			return
		}
		if len(imageData) == 0 {
			HandleErr(w, http.StatusBadRequest, fmt.Errorf("image data cannot be empty"))
			return
		}
		if err := os.WriteFile(fn, imageData, 0644); err != nil {
			HandleErr(w, http.StatusInternalServerError, fmt.Errorf("failed to write image data: %v", err))
			return
		}
		imagesCache.Remove(fn)
		thumbsCache.Remove(fn)
		w.WriteHeader(http.StatusCreated)
		return
	}
}

func init() {
	imagesCache = expirable.NewLRU[string, *[]byte](1000, nil, time.Minute*60*24)
	thumbsCache = expirable.NewLRU[string, *[]byte](1000, nil, time.Minute*60*24)
}
