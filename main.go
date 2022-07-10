package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

const (
	sessionExpiredStatus = 403
	mainapiUrl           = "https://api.writedom.com/writer"
)

func apiurl(path string) string {
	return mainapiUrl + path
}

func Login() (*LoginResponse, []*http.Cookie, error) {
	endpoint := apiurl("/login")
	params := url.Values{}
	params.Add("email", "charlesproficientwriternyak1@gmail.com")
	params.Add("password", "Doshab@21")

	resp, err := http.PostForm(endpoint, params)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	var loginResponse LoginResponse
	err = json.Unmarshal(response, &loginResponse)
	if err != nil {
		return nil, nil, err
	}
	return &loginResponse, resp.Cookies(), nil
}

func CheckAvailableJobs(cookiess []*http.Cookie) (*AvailableJobsResponse, error) {
	available_url := apiurl("/assignments")
	endpoint, err := url.Parse(available_url)
	if err != nil {
		return nil, err
	}

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

func SendBid(cookie *[]http.Cookie){
	applyUrl := apiurl("/apply")
	endpoint,err := url.Parse(applyUrl)
	if err != nil {
		log.Fatal(err)
	}
	jobId := "234375"
	//https://api.writedom.com/writer/assignments/2816449/apply?user_id=1490371&assignment_id=2816449&user_id=1490371&access_token=none&app_id=3&_token=4mPVoJhMnw3cpVY8KC43HglbDoD0nCfVPMZwLvtZ&is_new_wd=true&local_time=2022-07-10+19:11:24
	params := endpoint.Query()
	params.Set("user_id", "1490371")
	params.Set("assignment_id", jobId)

	endpoint.RawQuery = params.Encode()
	client := http.DefaultClient
	req, err := http.NewRequest(http.MethodGet, endpoint.String(), nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Referer", "https://writedom.com/")
	for _, cookie := range *cookie {
		req.AddCookie(&cookie)
	}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	response, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(response)


}


func main() {
	_, cookies, err := Login()
	if err != nil {
		log.Fatal(err)
	}
	// log.Println(AsPrettyJson(loginResponse))

	jobs, err := CheckAvailableJobs(cookies)
	if err != nil {
		panic(err)
	}
	log.Println(AsPrettyJson(jobs))
}






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