package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
)

/////////////////////////////////////////////////////////////
//
// Structures for Posts
//
/////////////////////////////////////////////////////////////

type PostCountry struct {
	Id   string `json:"id"`
	Name string `json:"name"`
} // PostCountry

type PostRegion struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	PostCountry PostCountry
} // PostRegion

type PostType struct {
	TypeId    int    `json:"typeId"`
	TypeConst string `json:"const"`
	Name      string `json:"name"`
} // PostType

type PostGroup struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	UniqueName string `json:"uniqueName"`
	Latitude   int    `json:"latitude"`
	Longitude  int    `json:"longitude"`
	Timezone   string `json:"timezone"`
} // PostGroup

type Post struct {
	Id           int       `json:"id"`
	UserId       int       `json:"userId"`
	Subject      string    `json:"subject"`
	Location     string    `json:"location"`
	Description  string    `json:"description"`
	IsApproved   int       `json:"isApproved"`
	RejectReason string    `json:"rejectReason"`
	Date         string    `json:"date"`
	Time         string    `json:"time"`
	PostType     PostType  `json:"type"`
	PostGroup    PostGroup `json:"group"`
	Static       string    `json:"static"`
	Image        string    `json:"image"`
	Thumb        string    `json:"thumb"`
	Images       []string  `json:"images"`
	Thumbs       []string  `json:"thumbs"`
	Tags         string    `json:"tags"`
} // Post

/////////////////////////////////////////////////////////////
//
// Structures for Groups
//
/////////////////////////////////////////////////////////////

type GroupCountry struct {
	Id   int    `json:"country_id"`
	Name string `json:"country_name"`
} // GroupCountry

type GroupRegion struct {
	Id      int          `json:"region_id"`
	Name    string       `json:"region_name"`
	Country GroupCountry `json:"country"`
} // GroupRegion

type Group struct {
	Id                            int         `json:"group_id"`
	Name                          string      `json:"group_name"`
	StatusId                      int         `json:"region_id"`
	RegionId                      int         `json:"group_status_id"`
	NgaApproved                   int         `json:"nga_approved"`
	NgaApprovedBy                 string      `json:"nga_approved_by"`
	Deleted                       int         `json:"deleted"`
	DeletedBy                     string      `json:"deleted_by"`
	YahooGroup                    int         `json:"yahoo_group"`
	InvitationOnly                int         `json:"invitation_only"`
	NumMembers                    int         `json:"num_members"`
	NumPosts7Day                  int         `json:"num_posts_7day"`
	NGAApprovedDate               string      `json:"nga_approved_date"`
	Description                   string      `json:"description"`
	MaxOtherGroups                int         `json:"max_other_groups"`
	MembersRequireApproval        int         `json:"members_require_approval"`
	MaxWantedsPerMemeber          int         `json:"max_wanteds_per_member"`
	YahooGroupName                string      `json:"yahoo_group_name"`
	DefaultTimezone               string      `json:"default_tz"`
	MyFreeCycleOptOut             int         `json:"myfreecycle_optout"`
	AllowImagePosting             int         `json:"allow_image_posting"`
	EmergencyModeration           string      `json:"emergency_moderation"`
	AutoModNewMemebers            int         `json:"auto_mod_new_members"`
	DefaultEmailDelivery          int         `json:"default_email_delivery"`
	DefaultLang                   string      `json:"default_lang"`
	Latitude                      int         `json:"latitude"`
	Longitude                     int         `json:"longitude"`
	TreatDigestResponsesAsBounces int         `json:"treat_digest_responses_as_bounces"`
	ExpireWantedDays              int         `json:"expire_wanteds_days"`
	ExpireOfferDays               int         `json:"expire_offers_days"`
	Region                        GroupRegion `json:"region"`
} // Group

/////////////////////////////////////////////////////////////
//
// Mainstructure for JSON response from FreeCycle
//
/////////////////////////////////////////////////////////////

type FreeCyleJSON struct {
	Group Group  `json:"group"`
	Posts []Post `json:"posts"`
} // FreeCycleJSON

/////////////////////////////////////////////////////////////
//
// Filtering Functions
//
/////////////////////////////////////////////////////////////
func filterOffers(posts []Post) []Post {
	ans := make([]Post, 0)
	for _, post := range posts {
		if post.PostType.Name == "OFFER" {
			ans = append(ans, post)
		} // if
	} // for
	return ans
} // filterOffers

func freeItems(w http.ResponseWriter, r *http.Request) {
	// Can someone inject malicious query params from this?
	var town string = strings.ToTitle(r.URL.Query().Get("town"))
	var state_symbol string = strings.ToUpper(r.URL.Query().Get("state_symbol"))

	// Instantiate default collector
	freeCycleCollector := colly.NewCollector()

	freeCycleCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("Requested: ", r.URL.String())
	})

	freeCycleCollector.OnResponse(func(r *colly.Response) {
		fmt.Println(r.Request.URL, "Responded!")
	})

	// Parse throught the HTML to get all posts
	freeCycleCollector.OnHTML("div.item-list-view", func(e *colly.HTMLElement) {
		// Retrieve the JSON stored in the fc-data html element's ":data" attribute
		attrVal, _ := e.DOM.ChildrenFiltered("fc-data").Attr(":data")

		// Stored the JSON data into a struct
		var freecycle FreeCyleJSON
		json.Unmarshal([]byte(attrVal), &freecycle)

		// Fileter Posts for only "Offers"
		var filt_posts = filterOffers(freecycle.Posts)

		// Respond with JSON
		log.Info("API's Crawler has a response!")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(filt_posts)
	})

	freeCycleCollector.Visit("https://www.freecycle.org/town/" + town + state_symbol)
} // flips
