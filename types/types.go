package types

import "errors"

type Endpoint string

const (
	EndpointGetSampleRaw          = "/get/sample/%s/raw"
	EndpointGetSampleZip          = "/get/sample/%s/zip"
	EndpointVersion               = "/get/version"
	EndpointGetActors             = "/list/actors"
	EndpointGetActor              = "/get/actor/%s"
	EndpointGetYaraRaw            = "/get/yara/%s/raw"
	EndpointGetYaraZip            = "/get/yara/%s/zip"
	EndpointFindActor             = "/find/actor/%s"
	EndpointGetFamilies           = "/get/families"
	EndpointGetFamily             = "/get/family/%s"
	EndpointFindFamily            = "/find/family/%s"
	EndpointGetYaraRulesForFamily = "/get/yara/%s/zip"
	EndpointListFamilySamples     = "/list/samples/%s"
	EndpointScanBinary            = "/scan/binary"
)

const (
	APIBase         = "https://malpedia.caad.fkie.fraunhofer.de/api"
	Dumped   string = "dumped"
	Packed   string = "packed"
	Unpacked string = "unpacked"
)

var ErrResourceNotFound = errors.New("Resource Not Found")

type actorSearchElement struct {
	CommonName string   `json:"common_name"`
	Synonyms   []string `json:"synonyms"`
	Name       string   `json:"name"`
}

type FamilySamples map[string][]FamilySample

type FamilySample struct {
	Status  string `json:"status"`
	Sha256  string `json:"sha256"`
	Version string `json:"version"`
}

type ZippedResults struct {
	Zipped string `json:"zipped"`
}

// GetVersion -----------------------------------------------------------
type Version struct {
	Date    string `json:"date"`
	Version int    `json:"version"`
}

type Actor struct {
	Value       string            `json:"value"`
	Meta        Meta              `json:"meta"`
	Families    map[string]Family `json:"families"`
	Description string            `json:"description"`
	Related     []Related         `json:"related"`
	UUID        string            `json:"uuid"`
}

type Meta struct {
	CfrSuspectedVictims      []string `json:"cfr-suspected-victims"`
	Country                  string   `json:"country"`
	Refs                     []string `json:"refs"`
	CfrTargetCategory        []string `json:"cfr-target-category"`
	CfrTypeOfIncident        string   `json:"cfr-type-of-incident"`
	Synonyms                 []string `json:"synonyms"`
	CfrSuspectedStateSponsor string   `json:"cfr-suspected-state-sponsor"`
	AttributionConfidence    string   `json:"attribution-confidence"`
}

type Related struct {
	DestUUID string   `json:"dest-uuid"`
	Type     string   `json:"type"`
	Tags     []string `json:"tags"`
}

type Actors []string

// FindActor -----------------------------------------------------------
type FindActorElement struct {
	CommonName string   `json:"common_name"`
	Synonyms   []string `json:"synonyms"`
	Name       string   `json:"name"`
}

type FindActor []FindActorElement

// GetFamilies ----------------------------------------------------------
type Families map[string]Family

type Family struct {
	Updated     string        `json:"updated"`
	Attribution []string      `json:"attribution"`
	Description string        `json:"description"`
	Notes       []string      `json:"notes"`
	AltNames    []string      `json:"alt_names"`
	Sources     []interface{} `json:"sources"`
	Urls        []string      `json:"urls"`
	CommonName  string        `json:"common_name"`
	UUID        *string       `json:"uuid,omitempty"`
	Properties  *Properties   `json:"properties,omitempty"`
}

type Properties struct {
	ProgrammingLanguage string        `json:"programming_language"`
	Iat                 string        `json:"iat"`
	Obfuscation         []interface{} `json:"obfuscation"`
	PEHeader            string        `json:"pe_header"`
	ValidTimestamp      string        `json:"valid_timestamp"`
}

// FindFamily ----------------------------------------------------------
type FindFamily []FindFamilyElement

type FindFamilyElement struct {
	Name     string   `json:"name"`
	AltNames []string `json:"alt_names"`
}

// ListFamilySamples ----------------------------------------------------------
type ListFamilySamples struct {
	Status  string `json:"status"`
	Sha256  string `json:"sha256"`
	Version string `json:"version"`
}

// ScanBinary ----------------------------------------------------------
type YaraMatchesValue struct {
	MatchedStrings int64 `json:"matched_strings"`
	MatchedHits    int64 `json:"matched_hits"`
	Match          bool  `json:"match"`
}
