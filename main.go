package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"github.com/joho/godotenv"
)

const (
	sessionExpiredStatus = 403
	mainapiUrl           = "https://api.writedom.com/writer"
)

type Job struct {
	Jobid string `json:"jobid"`
}

func init(){
	//load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

}

func InitDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&Job{})
	return db
}

func apiurl(path string) string {
	return mainapiUrl + path
}

func Login() (*LoginResponse, []*http.Cookie, error) {
	endpoint := apiurl("/login")
	params := url.Values{}
	// // params.Add("email", "academiawriting27@gmail.com")
	// params.Add("email", godotenv.Get("EMAIL"))

	// // params.Add("password", "HotGrb@5328")
	// params.Add("password", godotenv.Get("PASSWORD"))
	// // params.Add("_token", "4mPVoJhMnw3cpVY8KC43HglbDoD0nCfVPMZwLvtZ")
	// params.Add("_token", godotenv.Get("_token"))
	email := os.Getenv("EMAIL")
	password := os.Getenv("PASSWORD")
	token := os.Getenv("_token")
	params.Add("email", email)
	params.Add("password", password)
	params.Add("_token", token)

	fmt.Println(email, password, token)
	resp, err := http.PostForm(endpoint, params)
	if err != nil {
		time.Sleep(5 * time.Second)
		Login()
	}
	defer resp.Body.Close()

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	var loginResponse LoginResponse
	err = json.Unmarshal(response, &loginResponse)
	if err != nil {
		// return nil, nil, err
		fmt.Println("Relogging in...")
		time.Sleep(3 * time.Second)
		Login()
	}
	return &loginResponse, resp.Cookies(), nil
}

func CheckAvailableJobs(cookiess []*http.Cookie) (*AvailableJobsResponse, error) {
	available_url := apiurl("/assignments")
	endpoint, err := url.Parse(available_url)
	if err != nil {
		return nil, err
	}
	//
	params := endpoint.Query()
	params.Set("page", "1")
	params.Set("perPage", "10")

	endpoint.RawQuery = params.Encode()

	client := http.DefaultClient

	req, err := http.NewRequest(http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Referer", "https://writedom.com/")
	for _, cookie := range cookiess {
		req.AddCookie(cookie)
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	response, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var availableJobsResponse AvailableJobsResponse
	err = json.Unmarshal(response, &availableJobsResponse)
	if err != nil {
		var alt AlternativeResponse
		err = json.Unmarshal(response, &alt)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(alt.Errors[0])
	}

	return &availableJobsResponse, nil
}

func SendBid(cookiess []*http.Cookie, id string) (*AppliedResponse, error) {
	applyUrl := apiurl("/assignments/" + id + "/apply")
	endpoint, err := url.Parse(applyUrl)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	//https://api.writedom.com/writer/assignments/2816467/apply?user_id=1490371&assignment_id=2816467&user_id=1490371&access_token=none&app_id=3&_token=4mPVoJhMnw3cpVY8KC43HglbDoD0nCfVPMZwLvtZ&is_new_wd=true&local_time=2022-07-10+19%3A41%3A47
	params := endpoint.Query()
	params.Set("user_id", "1490371")
	params.Set("assignment_id", id)
	params.Set("user_id", "1490371")
	params.Set("access_token", "none")
	params.Set("app_id", "3")
	params.Set("_token", "4mPVoJhMnw3cpVY8KC43HglbDoD0nCfVPMZwLvtZ")
	params.Set("is_new_wd", "true")
	params.Set("local_time", "2022-07-10+19%3A41%3A47")

	endpoint.RawQuery = params.Encode()

	client := http.DefaultClient

	req, err := http.NewRequest(http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Referer", "https://writedom.com/")
	for _, cookie := range cookiess {
		req.AddCookie(cookie)
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	response, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(response))
	var bidResponse AppliedResponse
	err = json.Unmarshal(response, &bidResponse)
	if err != nil {
		return nil, err
	}
	log.Println("Successfully applied for job")
	return &bidResponse, nil
}

func ParseJobs(jobs *AvailableJobsResponse) ([]string, error) {
	var jobsId []string
	for _, job := range jobs.Data.AvailableAssignments {
		jobsId = append(jobsId, strconv.Itoa(job.ID))
	}
	return jobsId, nil
}

func main() {
	log.Println("Starting...")
	fmt.Println("Ctrl + C to exit")
	_, cookies, err := Login()
	if err != nil {

		log.Fatal(err)
		log.Fatal("Contact support (0111694419)")
		return
	}

	for {
		jobs, err := CheckAvailableJobs(cookies)
		if err != nil {
			_, cookies, err = Login()
			if err != nil {
				log.Fatal(err)
				log.Fatal("Contact support (0111694419)")
				return
			}
		}
		var jobsId []string
		if jobs == nil {
			log.Println("No jobs available")
			log.Println("Checking for jobs...")
			time.Sleep(time.Second * 2)
			continue
		}
		jobsId, err = ParseJobs(jobs)

		db := InitDB()

		jobdb := new(Job)

		for _, job := range jobsId {
			if err = db.Where("jobid = ?", job).First(jobdb).Error; err != nil {
				// if err != nil {
				_, _ = SendBid(cookies, job)
				jobdb.Jobid = job
				err = db.Create(jobdb).Error
				if err != nil {
					panic(err)
				}
			}
			// }

		}
		log.Println("All caught up...")
		log.Println("Checking for new jobs in 5 seconds..")
		time.Sleep(time.Second * 5)
	}
}

// log.Println(AsPrettyJson(jobs))

type LoginResponse struct {
	Data struct {
		User struct {
			ID                 int         `json:"id"`
			Username           string      `json:"username"`
			Email              string      `json:"email"`
			FirstName          string      `json:"first_name"`
			LastName           string      `json:"last_name"`
			Phone              string      `json:"phone"`
			Confirmed          int         `json:"confirmed"`
			LinkID             interface{} `json:"link_id"`
			PriceGroupID       int         `json:"price_group_id"`
			Avail247           int         `json:"avail24_7"`
			Activated          int         `json:"activated"`
			PhoneUnformatted   string      `json:"phone_unformatted"`
			ClientDeletedAt    interface{} `json:"client_deleted_at"`
			IsTest             int         `json:"is_test"`
			RiskLevelID        int         `json:"risk_level_id"`
			OtherSocialLink    string      `json:"other_social_link"`
			PhoneConfirmed     int         `json:"phone_confirmed"`
			WriterRating       float64     `json:"writer_rating"`
			FacebookSocialLink string      `json:"facebook_social_link"`
			VkSocialLink       string      `json:"vk_social_link"`
			Unsubscribed       int         `json:"unsubscribed"`
			WmDisplayAt        string      `json:"wm_display_at"`
			OrderFormToken     interface{} `json:"order_form_token"`
			OrderFormAuth      int         `json:"order_form_auth"`
			WrongContact       int         `json:"wrong_contact"`
			Inactive           interface{} `json:"inactive"`
			LastLogsID         string      `json:"last_logs_id"`
			Balance            int         `json:"balance"`
			NotificationsToken interface{} `json:"notifications_token"`
			DeviceIDToken      interface{} `json:"device_id_token"`
			St                 string      `json:"st"`
			Autoassign         bool        `json:"autoassign"`
			NewWriterRating    string      `json:"new_writer_rating"`
			Rating15Days       json.Number `json:"rating_15_days"`
			RejectStat         struct {
				Reject int `json:"reject"`
				Apply  int `json:"apply"`
			} `json:"rejectStat"`
			Alerts            interface{}   `json:"alerts"`
			OrderID           interface{}   `json:"order_id"`
			Cashback          string        `json:"cashback"`
			CashbackBalance   int           `json:"cashbackBalance"`
			IsReferral        bool          `json:"is_referral"`
			LoggedAsSuperUser bool          `json:"logged_as_super_user"`
			FullName          string        `json:"full_name"`
			Webmaster         []interface{} `json:"webmaster"`
			Friendsmaster     interface{}   `json:"friendsmaster"`
			PayoutRegularity  interface{}   `json:"payout_regularity"`
			OwnDocuments      []struct {
				ID                 int    `json:"id"`
				UserDocumentTypeID string `json:"user_document_type_id"`
				UserID             string `json:"user_id"`
				Path               string `json:"path"`
				Name               string `json:"name"`
				CreatedAt          string `json:"created_at"`
				UpdatedAt          string `json:"updated_at"`
			} `json:"own_documents"`
			PriceGroup struct {
				ID    int     `json:"id"`
				Title string  `json:"title"`
				Rate  float64 `json:"rate"`
			} `json:"price_group"`
		} `json:"user"`
	} `json:"data"`
	Status int           `json:"status"`
	Errors []interface{} `json:"errors"`
	Alerts []interface{} `json:"alerts"`
}

type AvailableJobsResponse struct {
	Data struct {
		Total                int `json:"total"`
		PerPage              int `json:"perPage"`
		Page                 int `json:"page"`
		AvailableAssignments []struct {
			ID                int     `json:"id"`
			Topic             string  `json:"topic"`
			Price             float64 `json:"price"`
			Pages             int     `json:"pages"`
			Slides            int     `json:"slides"`
			Problems          int     `json:"problems"`
			Sources           int     `json:"sources"`
			OrderID           int     `json:"order_id"`
			Deadline          string  `json:"deadline"`
			FinalDeadline     string  `json:"final_deadline"`
			ConfirmedByWriter int     `json:"confirmed_by_writer"`
			PriceGroupID      int     `json:"price_group_id"`
			AcademicLevels    string  `json:"academicLevels"`
			Subjects          string  `json:"subjects"`
			Spacings          string  `json:"spacings"`
			PaperFormats      string  `json:"paperFormats"`
			WordCount         int     `json:"word_count"`
			Questions         int     `json:"questions"`
			Description       string  `json:"description"`
		} `json:"available_assignments"`
	} `json:"data"`
	Status int           `json:"status"`
	Errors []interface{} `json:"errors"`
	Alerts []interface{} `json:"alerts"`
}

type AlternativeResponse struct {
	Data   []interface{} `json:"data"`
	Status int           `json:"status"`
	Errors []string      `json:"errors"`
	Alerts []interface{} `json:"alerts"`
}

func AsPrettyJson(input interface{}) string {
	jsonB, _ := json.MarshalIndent(input, "", "  ")
	return string(jsonB)
}

func AsJson(input interface{}) string {
	jsonB, _ := json.Marshal(input)
	return string(jsonB)
}

type AppliedResponse struct {
	Data   []interface{} `json:"data"`
	Status int           `json:"status"`
	Errors []interface{} `json:"errors"`
	Alerts []string      `json:"alerts"`
}
