package models

type Snake struct {
	SnakeId string `json:"sid"`
	Species string `json:"species"`
	Sex     string `json:"sex"`
	Age     int    `json:"age"`
	Genes   string `json:"genes"`
	Notes   string `json:"notes"`
}

type SnakesListItem struct {
	Sid     string `json:"sid"`
	Species string `json:"species"`
}
type UpdateSnake struct {
	Sid     string  `json:"sid"` // required
	Species *string `json:"species,omitempty"`
	Sex     *string `json:"sex,omitempty"`
	Age     *int    `json:"age,omitempty"`
	Genes   *string `json:"genes,omitempty"`
	Notes   *string `json:"notes,omitempty"`
}

type SnakeFeed struct {
	Sid      string `json:"sid"`       // snake id (user-given)
	FeedDate string `json:"feed_date"` // YYYY-MM-DD
	PreyType string `json:"prey_type"`
	PreySize string `json:"prey_size"`
	Notes    string `json:"notes"`
}

type SnakeFeedDeleteRequest struct {
	Sid  string `json:"sid"`       // snake ID or name
	Date string `json:"feed_date"` // feeding record date (match the format stored in DB)
}

type UpdateSnakeFeed struct {
	Sid      string  `json:"sid"`
	FeedDate string  `json:"feed_date"` // YYYY-MM-DD
	PreyType *string `json:"prey_type,omitempty"`
	PreySize *string `json:"prey_size,omitempty"`
	Notes    *string `json:"notes,omitempty"`
}

type SnakeSid struct {
	Sid string `json:"sid"`
}

type HealthRecord struct {
	Sid       string `json:"sid"`
	CheckDate string `json:"check_date"`
	Weight    string `json:"weight"`
	Length    string `json:"length"`
	Topic     string `json:"topic"`
	Notes     string `json:"notes"`
}

type UpdateHealthRecord struct {
	Sid       string `json:"sid"`        // snake id (user-given)
	CheckDate string `json:"check_date"` // YYYY-MM-DD
	Weight    string `json:"weight,omitempty"`
	Length    string `json:"length,omitempty"`
	Topic     string `json:"topic,omitempty"`
	Notes     string `json:"notes,omitempty"`
}

type DeleteSnakeHealth struct {
	Sid       string `json:"sid"`
	CheckDate string `json:"check_date"`
}

type FindBreedingEvent struct {
	Sid      string `json:"sid"`
	Event_id string `json:"event_id"`
}
type FindBreedingEvent2 struct {
	Sid      string `json:"sid"`
	Event_id string `json:"event_id"`
}

type BreedingEventItem struct {
	Sid           string   `json:"sid"`
	BreedingYears []string `json:"breeding_years"`
}

type BreedingEvent struct {
	Event_id       string  `json:"event_id"`
	FemaleSid      string  `json:"female_sid"`          // snake's ID
	Male1Sid       string  `json:"ma1e1_sid"`           // mate's ID
	Male2Sid       string  `json:"ma1e2_sid,omitempty"` // mate's ID
	Male3Sid       string  `json:"ma1e3_sid,omitempty"` // mate's ID
	Male4Sid       string  `json:"ma1e4_sid,omitempty"` // mate's ID
	BreedingYear   string  `json:"breeding_year"`
	BreedingSeason string  `json:"breeding_season"`
	FemaleWeight   *string `json:"female_weight,omitempty"`
	Male1Weight    *string `json:"male1_weight,omitempty"`
	Male2Weight    *string `json:"male2_weight,omitempty"`
	Male3Weight    *string `json:"male3_weight,omitempty"`
	Male4Weight    *string `json:"male4_weight,omitempty"`
	CoolingStart   *string `json:"cooling_start,omitempty"`
	CoolingEnd     *string `json:"cooling_end,omitempty"`
	WarmingStart   *string `json:"warming_start,omitempty"`
	WarmingEnd     *string `json:"warming_end,omitempty"`
	PairingDate1   *string `json:"pairing1_date,omitempty"`
	PairingDate2   *string `json:"pairing2_date,omitempty"`
	PairingDate3   *string `json:"pairing3_date,omitempty"`
	PairingDate4   *string `json:"pairing4_date,omitempty"`
	GravidDate     *string `json:"gravid_date,omitempty"`
	LayDate        *string `json:"lay_date,omitempty"`
	ClutchSize     *int    `json:"clutch_size,omitempty"`
	ClutchSurvive  *string `json:"clutch_survive,omitempty"`
	Outcome        *string `json:"outcome,omitempty"`
	Notes          *string `json:"notes,omitempty"`
}
type BreedingSummary struct {
	BreedingUUID string `json:"breeding_uuid"`
	EventID      string `json:"event_id"`
	Year         string `json:"breeding_year"`
	Season       string `json:"breeding_season"`
}

type Offspring struct {
	Uuid          string  `json:"uuid"`
	Baby_Suid     string  `json:"baby_suid"`
	Baby_Sid      string  `json:"baby_sid"`
	Mother_Suid   string  `json:"mother_suid"`
	Mother_Sid    string  `json:"mother_sid"`
	Father_suid   string  `json:"father_suid"`
	Father_sid    string  `json:"father_sid"`
	Breeding_uuid string  `json:"breeding_uuid"`
	Breed_year    string  `json:"breeding_year"`
	Baby_genes    string  `json:"baby_genes"`
	Mother_genes  string  `json:"mother_genes"`
	Father_genes  string  `json:"father_genes"`
	Notes         *string `json:"notes,omitempty"`
}
