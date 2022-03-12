package nandlers

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	models "groupie-tracker/models"
)

// Error для вывода ошибок,через html
type Error struct {
	Message string
	Code    int
}

// Parse функции для парсинга json api файлов
func Parse(url string, j interface{}) error {
	spaceClient := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", "spacecount-tutorial")

	res, getErr := spaceClient.Do(req)
	if getErr != nil {
		return getErr
	}

	if res.Body != nil {
		defer res.Body.Close()
	}
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return readErr
	}
	strct := j
	jsonErr := json.Unmarshal(body, &strct)
	if jsonErr != nil {
		return jsonErr
	}
	return nil
}

// MethodGet проверяет используется ли метод GET
func MethodGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		CustomError(http.StatusMethodNotAllowed, w)
		return
	}
}

// func Parsert(url string, strct interface{}) {
// 	res, err := http.Get(url)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		os.Exit(0)
// 	}
// 	body, err := ioutil.ReadAll(res.Body)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		os.Exit(0)
// 	}
// 	json.Unmarshal(body, &strct)
// }

// Home handler
func Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		CustomError(http.StatusNotFound, w)
		return
	}
	url := "https://groupietrackers.herokuapp.com/api/artists"
	err := Parse(url, &models.General.Artists)
	if err != nil {
		CustomError(http.StatusInternalServerError, w)
		return
	}
	err = Parse("https://groupietrackers.herokuapp.com/api/locations", &models.General.Locations)
	if err != nil {
		log.Println(err.Error())
		CustomError(http.StatusInternalServerError, w)
		return
	}
	MethodGet(w, r)
	for i := 0; i < len(models.General.Artists); i++ {
		models.General.Artists[i].Location = models.General.Locations.Index[i].Locations
	}

	ts, err := template.ParseFiles("./ui/html/index.html")
	if err != nil {
		log.Println(err.Error())
		CustomError(http.StatusInternalServerError, w)
		return
	}
	err = ts.Execute(w, models.General.Artists)
	if err != nil {
		log.Println(err.Error())
		CustomError(http.StatusInternalServerError, w)
		return
	}
}

// Artist handler
func Artist(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/artist/" {
		CustomError(http.StatusBadRequest, w)
		return
	}
	fr := strings.TrimPrefix(r.URL.Path, "/artist/")
	if Atoi(fr) <= 0 || Atoi(fr) > 52 {
		CustomError(http.StatusBadRequest, w)
		return
	}
	err := Parse("https://groupietrackers.herokuapp.com/api/artists", &models.General.Artists)
	if err != nil {
		log.Println(err.Error())
		CustomError(http.StatusInternalServerError, w)
		return
	}
	err = Parse("https://groupietrackers.herokuapp.com/api/relation", &models.General.Relation)
	if err != nil {
		log.Println(err.Error())
		CustomError(http.StatusInternalServerError, w)
		return
	}
	MethodGet(w, r)
	// Меняю значение в Artists.DatesLocation на Relation.DatesLocation
	models.General.Artists[Atoi(fr)-1].DatesLocation = models.General.Relation.Index[Atoi(fr)-1].DatesLocation
	ts, err := template.ParseFiles("./ui/html/artist.html")
	if err != nil {
		log.Println(err.Error())
		CustomError(http.StatusInternalServerError, w)
		return
	}
	err = ts.Execute(w, models.General.Artists[Atoi(fr)-1])
	if err != nil {
		log.Println(err.Error())
		CustomError(http.StatusInternalServerError, w)
		return
	}
}

// Atoi atoi
func Atoi(s string) int {
	res := 0
	if s[0] == '-' {
		return 0
	}
	for _, i := range s {
		if i >= 48 && i <= 57 {
			i = i - 48
			res = res*10 + int(i)
		} else {
			return 0
		}
	}
	return res
}

// CustomError обработка ошибок
func CustomError(code int, w http.ResponseWriter) {
	message := Error{Code: code, Message: http.StatusText(code)}
	gh, _ := template.ParseFiles("./ui/html/error.html")
	w.WriteHeader(code)
	gh.Execute(w, message)
}

// Search handler
func Search(w http.ResponseWriter, r *http.Request) {
	err := Parse("https://groupietrackers.herokuapp.com/api/artists", &models.General.Artists)
	if err != nil {
		log.Println(err.Error())
		CustomError(http.StatusInternalServerError, w)
		return
	}
	id := strings.ToLower(r.URL.Query().Get("id"))
	tag := r.FormValue("tags")
	res := []int{}
	if id == "" {
		CustomError(http.StatusBadRequest, w)
		return
	}
	check := false
	// обработка запросов без тегов
	if check == false {
		for _, i := range models.General.Artists {
			if tag == "all" {
				if strings.Contains(strings.ToLower(i.Name), id) {
					if Proverka(res, i.ID) {
						res = append(res, i.ID)
					}
				} else if strings.Contains(strings.ToLower(i.FirstAlbum), id) {
					if Proverka(res, i.ID) {
						res = append(res, i.ID)
					}
				} else if strings.Contains(strings.ToLower(strconv.Itoa(i.CreationDate)), id) {
					if Proverka(res, i.ID) {
						res = append(res, i.ID)
					}
				}
				for _, j := range i.Members {
					if strings.Contains(strings.ToLower(j), id) {
						if Proverka(res, i.ID) {
							res = append(res, i.ID)
						}
					}
				}
				for _, k := range i.Location {
					if strings.Contains(strings.ToLower(k), id) {
						if Proverka(res, i.ID) {
							res = append(res, i.ID)
						}
					}
				}
			} else if tag == "artist/band" {
				if strings.Contains(strings.ToLower(i.Name), id) {
					if Proverka(res, i.ID) {
						res = append(res, i.ID)
					}
				}
			} else if tag == "first album date" {
				if strings.Contains(strings.ToLower(i.FirstAlbum), id) {
					if Proverka(res, i.ID) {
						res = append(res, i.ID)
					}
				}
			} else if tag == "creation date" {
				if strings.Contains(strings.ToLower(strconv.Itoa(i.CreationDate)), id) {
					if Proverka(res, i.ID) {
						res = append(res, i.ID)
					}
				}
			} else if tag == "members" {
				for _, j := range i.Members {
					if strings.Contains(strings.ToLower(j), id) {
						if Proverka(res, i.ID) {
							res = append(res, i.ID)
						}
					}
				}
			} else if tag == "locations" {
				for _, k := range i.Location {
					if strings.Contains(strings.ToLower(k), id) {
						if Proverka(res, i.ID) {
							res = append(res, i.ID)
						}
					}
				}
			} else {
				CustomError(http.StatusNotFound, w)
				return
			}
		}
	}
	if len(res) == 0 {
		tt, err := template.ParseFiles("./ui/html/not_found.html")
		if err != nil {
			CustomError(http.StatusInternalServerError, w)
			return
		}
		tt.Execute(w, nil)
		return
	} else if len(res) != 0 {
		for i := 0; i < len(res); i++ {
			gh, err := template.ParseFiles("./ui/html/search.html")
			if err != nil {
				CustomError(http.StatusInternalServerError, w)
				return
			}
			gh.Execute(w, models.General.Artists[res[i]-1])
		}
	}

}

// Proverka функция для того чтобы , ответы не дублировались
func Proverka(s []int, n int) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == n {
			return false
		}
	}
	return true
}
