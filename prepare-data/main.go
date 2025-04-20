package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"log/slog"
	"os"
	"slices"

	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
	"github.com/gocarina/gocsv"
)

const (
	pathJSON  = "../dataset/restaurants_list.json"
	pathCSV   = "../dataset/restaurants_info.csv"
	indexName = "restaurants"
)

type Restaurants []*Restaurant

func (rs Restaurants) toRecord() []map[string]any {
	records := make([]map[string]any, len(rs))
	for i, r := range rs {
		records[i] = r.toRecord()
	}

	return records
}

type Restaurant struct {
	ObjectID int `json:"objectID,omitempty" csv:"objectID"`

	// info given only in the JSON

	Name             string         `json:"name,omitempty" csv:"-"`
	Address          string         `json:"address,omitempty" csv:"-"`
	Area             string         `json:"area,omitempty" csv:"-"`
	City             string         `json:"city,omitempty" csv:"-"`
	Country          string         `json:"country,omitempty" csv:"-"`
	ImageURL         string         `json:"image_url,omitempty" csv:"-"`
	MobileReserveURL string         `json:"mobile_reserve_url,omitempty" csv:"-"`
	PaymentOptions   PaymentOptions `json:"payment_options,omitempty" csv:"-"`
	Phone            string         `json:"phone,omitempty" csv:"-"`
	PostalCode       string         `json:"postal_code,omitempty" csv:"-"`
	Price            int            `json:"price,omitempty" csv:"-"`
	ReserveURL       string         `json:"reserve_url,omitempty" csv:"-"`
	State            string         `json:"state,omitempty" csv:"-"`
	Geolocation      Geolocation    `json:"_geoloc" csv:"-"`

	// info only given in the CSV

	FoodType     string     `json:"-" csv:"food_type"`
	StarsCount   float32    `json:"-" csv:"stars_count"`
	ReviewsCount int        `json:"-" csv:"reviews_count"`
	Neighborhood string     `json:"-" csv:"neighborhood"`
	PhoneNumber  string     `json:"-" csv:"phone_number"`
	PriceRange   PriceRange `json:"-" csv:"price_range"`
	DiningStyle  string     `json:"-" csv:"dining_style"`
}

func (r Restaurant) toRecord() map[string]any {
	return map[string]any{
		"objectID":           r.ObjectID,
		"name":               r.Name,
		"address":            r.Address,
		"area":               r.Area,
		"city":               r.City,
		"country":            r.Country,
		"image_url":          r.ImageURL,
		"mobile_reserve_url": r.MobileReserveURL,
		"payment_options":    r.PaymentOptions,
		"phone":              r.Phone,
		"postal_code":        r.PostalCode,
		"price":              r.Price,
		"reserve_url":        r.ReserveURL,
		"state":              r.State,
		"_geoloc": map[string]float64{
			"lat": r.Geolocation.Lat,
			"lng": r.Geolocation.Lng,
		},
		"food_type":     r.FoodType,
		"stars_count":   r.StarsCount,
		"reviews_count": r.ReviewsCount,
		"neighborhood":  r.Neighborhood,
		"phone_number":  r.PhoneNumber,
		"price_range":   r.PriceRange,
		"dining_style":  r.DiningStyle,
	}
}

type Geolocation struct {
	Lat float64 `json:"lat,omitempty"`
	Lng float64 `json:"lng,omitempty"`
}

type PaymentOptions []PaymentOption

// PaymentOption represents a standardized payment method.
type PaymentOption string

const (
	PaymentOptionAMEX       PaymentOption = "AMEX" // i.e. American Express
	PaymentOptionVisa       PaymentOption = "Visa"
	PaymentOptionDiscover   PaymentOption = "Discover"
	PaymentOptionMasterCard PaymentOption = "MasterCard"
)

func (opt *PaymentOption) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	// unmarshal the string
	s := ""
	if err := json.UnmarshalDecode(dec, &s); err != nil {
		return err
	}

	// check if the payment option is valid
	switch PaymentOption(s) {
	case PaymentOptionAMEX, PaymentOptionVisa, PaymentOptionDiscover, PaymentOptionMasterCard:
		// requirement: we should only have: AMEX/American Express, Visa, Discover, and MasterCard
		*opt = PaymentOption(s)
	case "Diners Club", "Carte Blanche":
		// requirement: for our purpose, Diners Club and Carte Blanche are Discover cards
		*opt = PaymentOptionDiscover
	case "Pay with OpenTable", "JCB", "Cash Only":
		// unknown options, keep payment option blank
		slog.Debug("skipping unknown payment option", slog.String("payment_option", s))
	default:
		// NOTE: when injesting data, add unknown options to the above and run again
		// this way, we have a list of these options and can communicate with stakeholders
		return fmt.Errorf("unkown payment option %q", s)
	}

	return nil
}

type PriceRange string

const (
	PriceRangeUnder30 PriceRange = "$30 and under"
	PriceRange30to50  PriceRange = "$31 to $50"
	PriceRangeOver50  PriceRange = "$50 and over"
)

func (pr *PriceRange) UnmarshalCSV(s string) error {
	// check if the price range is valid
	switch PriceRange(s) {
	case PriceRangeUnder30, PriceRange30to50, PriceRangeOver50:
		*pr = PriceRange(s)
	default:
		// we rule out any misspells or other errors
		// if need be, we change this code to add another enum or correct a false value
		return fmt.Errorf("unkown price range %q", s)
	}

	return nil
}

func main() {
	if err := prepareData(); err != nil {
		log.Fatal(err)
	}
}

func prepareData() error {
	restaurants, err := getRestaurantsList()
	if err != nil {
		return err
	}

	restsByID := map[int]*Restaurant{}
	for _, r := range restaurants {
		restsByID[r.ObjectID] = r
	}

	infos, err := getRestaurantsInfo()
	if err != nil {
		return err
	}

	// merge the data
	for _, info := range infos {
		r, ok := restsByID[info.ObjectID]
		if !ok {
			// NOTE: this return is just to confirm that all infos are matched to restaurants
			// alternatively, we can just continue
			return fmt.Errorf("unknown restaurant with ID %d", info.ObjectID)
		}

		r.FoodType = info.FoodType
		r.StarsCount = info.StarsCount
		r.ReviewsCount = info.ReviewsCount
		r.Neighborhood = info.Neighborhood
		r.PhoneNumber = info.PhoneNumber
		r.PriceRange = info.PriceRange
		r.DiningStyle = info.DiningStyle
	}

	// initiate client
	client, err := search.NewClient(os.Getenv("ALGOLIA_APP_ID"), os.Getenv("ALGOLIA_API_KEY"))
	if err != nil {
		return fmt.Errorf("initialize Algolia client: %w", err)
	}

	// push data to algolia
	result, err := client.SaveObjects(indexName, restaurants.toRecord())
	if err != nil {
		return fmt.Errorf("saving objects: %w", err)
	}

	slog.Info("done uploading records", slog.Int("num_batches", len(result)))

	if err := setIndexSettings(client); err != nil {
		return fmt.Errorf("setting index: %w", err)
	}

	return nil
}

func getRestaurantsList() (Restaurants, error) {
	f, err := os.Open(pathJSON)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rests := Restaurants{}
	if err := json.UnmarshalRead(f, &rests,
		// data validation
		json.RejectUnknownMembers(true)); err != nil {
		return nil, fmt.Errorf("unmarshal JSON: %w", err)
	}

	return rests, nil
}

func getRestaurantsInfo() (Restaurants, error) {
	// Read CSV
	f, err := os.Open(pathCSV)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.Comma = ';'

	// NOTE: if memory becomes an issue, there are strategies to parse just one line at a time
	// this way, we can merge with the JSON right away
	rests := Restaurants{}
	if err := gocsv.UnmarshalCSV(r, &rests); err != nil {
		return nil, fmt.Errorf("unmarshal CSV: %w", err)
	}

	return rests, err
}

// setIndexSettings sets the facets for food type and payment options
func setIndexSettings(c *search.APIClient) error {
	cfg, err := c.GetSettings(c.NewApiGetSettingsRequest(indexName))
	if err != nil {
		return fmt.Errorf("getting settings: %w", err)
	}

	attr := cfg.AttributesForFaceting

	if slices.Contains(attr, "food_type") &&
		slices.Contains(attr, "payment_options") {
		return nil // settings up to date
	}

	attr = append(attr, "food_type", "payment_options")
	slices.Sort(attr)
	attr = slices.Compact(attr)

	rsp, err := c.SetSettings(c.NewApiSetSettingsRequest("restaurants", &search.IndexSettings{
		AttributesForFaceting: attr,
	}))
	if err != nil {
		return err
	}

	slog.Info("index settings updated", slog.Int64("taskID", rsp.TaskID))

	return nil
}
