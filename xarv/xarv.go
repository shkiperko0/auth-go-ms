package xarv

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type ProxyHttpHandler struct {
	entry     *echo.Echo
	auth_host string
}

type REQ_SesionCheck struct {
	Url string `json:"url"`
}

func CheckSession(url string, refreshToken *string, auth_host *string) error {

	jsonBody, _ := json.Marshal(REQ_SesionCheck{Url: url})
	bodyReader := bytes.NewReader(jsonBody)
	auth_url := (*auth_host) + "/api/v1/session/check"
	check_req, _ := http.NewRequest(http.MethodPost, auth_url, bodyReader)

	if refreshToken != nil && *refreshToken != "undefined" {
		check_req.Header.Set("Authorization", "Bearer "+*refreshToken)
		//fmt.Println(auth_url, " << ", *refreshToken)
	} else {
		//fmt.Println(auth_url, " << ", "<no token>")
	}

	check_req.Header.Set("App", "eamapp")
	check_req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(check_req)

	if err != nil {
		return err
	} else {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		type Message struct {
			Message string
		}

		var msg Message
		json.Unmarshal(body, &msg)

		if res.StatusCode != http.StatusOK {
			if msg.Message == "auth-ms.session.not-found" {
				return CheckSession(url, nil, auth_host)
			}

			fmt.Println("SessionCheck:", res.Status, " ", msg.Message)
			return errors.New(msg.Message)
		}

		var json string = string(body)
		fmt.Println("SessionCheck:", res.Status, " ", json)
	}

	return nil
}

func GetBearerFromHeaders(headers *http.Header) *string {
	auths := (*headers)["Authorization"]
	if len(auths) > 0 {
		res := strings.Split(auths[0], " ")
		if len(res) == 2 {
			return &res[1]
		}
	}
	return nil
}

func ServeProxyCall(req *http.Request, res *echo.Response, entry *MSRedirectEntryRecord, proxy *httputil.ReverseProxy) error {

	var Relative string
	fmt.Sscanf(req.URL.Path, entry.GateURL+"%s", &Relative)
	req.URL.Path = entry.ApiURL + Relative

	proxy.ServeHTTP(res, req)
	//res.Header().Set("Access-Control-Allow-Origin", origin)
	//res.Header().Del("Access-Control-Allow-Origin")
	//res.Header().Del("Access-Control-Allow-Origin")
	//allowOrigin := res.Header().Values("Access-Control-Allow-Origin")
	//origin := req.Header.Get("Origin")
	//if len(allowOrigin) > 1 {
	//res.Writer.Header().Set("Access-Control-Allow-Origin", origin+"_2")
	//res.Writer.Header().Add("Access-Control-Allow-Origin", origin+"_3")
	//res.Header().Del("Access-Control-Allow-Origin")
	//}
	return nil
}

func (h *ProxyHttpHandler) DeligateHTTPRequest(ms *MSRecord, proxy *httputil.ReverseProxy, entry *MSRedirectEntryRecord) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		res := c.Response()

		err := CheckSession(req.URL.Path, GetBearerFromHeaders(&req.Header), &h.auth_host)
		if err != nil {
			return err
		}
		ServeProxyCall(req, res, entry, proxy)
		return nil
	}
}

func (h *ProxyHttpHandler) DeligateHTTPRequestOptions(ms *MSRecord, proxy *httputil.ReverseProxy, entry *MSRedirectEntryRecord) echo.HandlerFunc {
	return func(c echo.Context) error {
		ServeProxyCall(c.Request(), c.Response(), entry, proxy)
		return nil
	}
	//}
}

func (h *ProxyHttpHandler) Service(ms *MSRecord) error {

	host_url, err := url.Parse(ms.Host)
	if err != nil {
		log.Print(err)
		return err
	}

	proxy := httputil.NewSingleHostReverseProxy(host_url)

	fmt.Println(ms.Name, "Entries: ", len(ms.Entries))

	for i := range ms.Entries {
		entry := &ms.Entries[i]
		switch entry.Method {
		case "ANY":
			{
				h.entry.Any(entry.GateURL, (h.DeligateHTTPRequest(ms, proxy, entry)))
				h.entry.Add("OPTIONS", entry.GateURL, h.DeligateHTTPRequestOptions(ms, proxy, entry))
			}
		case "GET", "POST", "DELETE", "PUT", "PATCH":
			{
				h.entry.Add(entry.Method, entry.GateURL, (h.DeligateHTTPRequest(ms, proxy, entry)))
				h.entry.Add("OPTIONS", entry.GateURL, h.DeligateHTTPRequestOptions(ms, proxy, entry))
			}

		default:
			{
				fmt.Println("Warning, unknown method ", entry.Method, "\t", entry.GateURL, " => ", entry.ApiURL)
			}
		}

		if strings.HasSuffix(entry.GateURL, "/*") == true {
			entry.GateURL = entry.GateURL[0 : len(entry.GateURL)-2]
		}
	}

	return nil
}

type MSRecord struct {
	ID          uint           `gorm:"primaryKey, autoIncrement"`
	Name        string         `gorm:"not null"`
	Host        string         `gorm:"index:idx_host,unique"`
	AutoEntries bool           `gorm:"not null, default:false"`
	Description string         `gorm:"default:null"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`

	Entries []MSRedirectEntryRecord `gorm:"foreignKey:MSID"`
}

func (MSRecord) TableName() string {
	return "MS_Records"
}
func (MSRedirectEntryRecord) TableName() string {
	return "MS_EntryRoutes"
}

type MSRedirectEntryRecord struct {
	ID          uint           `gorm:"primaryKey, autoIncrement"`
	MSID        uint           `gorm:"not null"`
	Name        string         `gorm:"not null"`
	Description string         `gorm:"default:null"`
	GateURL     string         `gorm:"index:idx_gate_url,unique"`
	ApiURL      string         `gorm:"not null"`
	Method      string         `gorm:"not null"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

const API_V1 = "/api/v1"

func NewProxyHttpHandler(db *gorm.DB, e *echo.Echo, auth_host string) {
	handler := ProxyHttpHandler{e, auth_host}

	db.AutoMigrate(&MSRecord{})
	db.AutoMigrate(&MSRedirectEntryRecord{})

	mss := []MSRecord{}
	db.Find(&mss)

	if len(mss) == 0 {
		mss = []MSRecord{
			// {
			// 	Name: "EAM Store MS",
			// 	Host: "http://127.0.0.1:4001",
			// 	Entries: []MSRedirectEntryRecord{
			// 		{GateURL: "/api/items", ApiURL: "/items"},
			// 	},
			// },
			{
				Name: "EAM Content MS",
				Host: "http://cd-api-s:4000",
				Entries: []MSRedirectEntryRecord{
					{GateURL: API_V1 + "/users/get-avatare-by-id/*", ApiURL: API_V1 + "/users/avatare-by-id", Method: "GET"},
					{GateURL: API_V1 + "/users/avatare", ApiURL: API_V1 + "/users/avatare", Method: "ANY"},
					{GateURL: API_V1 + "/users/settings", ApiURL: API_V1 + "/users/settings", Method: "ANY"},

					{GateURL: API_V1 + "/content/categories", ApiURL: API_V1 + "/posts/categories", Method: "GET"},
					{GateURL: API_V1 + "/content/create-category", ApiURL: API_V1 + "/posts/categories/create", Method: "POST"},

					{GateURL: API_V1 + "/content/get-comments", ApiURL: API_V1 + "/comments/get", Method: "POST"},
					{GateURL: API_V1 + "/content/create-comment", ApiURL: API_V1 + "/comments/create", Method: "POST"},
					{GateURL: API_V1 + "/content/edit-comment/*", ApiURL: API_V1 + "/comments/edit", Method: "POST"},
					{GateURL: API_V1 + "/content/delete-comment/*", ApiURL: API_V1 + "/comments/delete", Method: "DELETE"},

					{GateURL: API_V1 + "/content/posts/*", ApiURL: API_V1 + "/posts", Method: "GET"},
					{GateURL: API_V1 + "/content/create-post", ApiURL: API_V1 + "/posts/create", Method: "POST"},
					{GateURL: API_V1 + "/content/delete-post/*", ApiURL: API_V1 + "/posts", Method: "DELETE"},
					{GateURL: API_V1 + "/content/edit-post/*", ApiURL: API_V1 + "/posts/edit", Method: "DELETE"},
					{GateURL: API_V1 + "/content/posts-list", ApiURL: API_V1 + "/posts/list", Method: "POST"},

					{GateURL: API_V1 + "/vfiles/upload/*", ApiURL: API_V1 + "/vfiles", Method: "POST"},
					{GateURL: API_V1 + "/vfiles/ls/*", ApiURL: API_V1 + "/vfiles/ls", Method: "GET"},
					{GateURL: API_V1 + "/vfiles/*", ApiURL: API_V1 + "/vfiles", Method: "GET"},
					{GateURL: API_V1 + "/vfiles/rn/*", ApiURL: API_V1 + "/vfiles/rn", Method: "POST"},
					{GateURL: API_V1 + "/vfiles/rm/*", ApiURL: API_V1 + "/vfiles/rm", Method: "POST"},
					{GateURL: API_V1 + "/vfiles/mv/*", ApiURL: API_V1 + "/vfiles/mv", Method: "POST"},

					{GateURL: API_V1 + "/content/health", ApiURL: "/api/health", Method: "GET"},
				},
			},
		}
		db.Create(&mss)
	}

	//handler.Service(&store)
	//handler.Service(&content)

	for i := range mss {
		ms := &mss[i]
		db.Find(&ms.Entries, MSRedirectEntryRecord{MSID: ms.ID})
		handler.Service(ms)
	}
}
