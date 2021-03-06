package types

import (
	"bytes"
	"errors"
)

type Endpoint string

// SampleState represents a acquired file from malpedia.
// Malpedia will return a sample at multiple states so you can see
// it dumped at different addresses and being packed or unpacked
type SampleState *map[string]bytes.Buffer

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
	EndpointScanYara              = "/scan/yara"
	EndpointScanYaraAgainstFamily = "/scan/yara/%s"
)

const (
	// APIBase is the default endpoint for malpedia
	APIBase         = "https://malpedia.caad.fkie.fraunhofer.de/api"
	Dumped   string = "dumped"
	Packed   string = "packed"
	Unpacked string = "unpacked"
)

// ErrResourceNotFound is the generic error returned for resources that aren't found whether that be entries or malware samples
var ErrResourceNotFound = errors.New("Resource Not Found")

type actorSearchElement struct {
	CommonName string   `json:"common_name"`
	Synonyms   []string `json:"synonyms"`
	Name       string   `json:"name"`
}

type FamilySamples []FamilySample

// ListFamilySamples ----------------------------------------------------------

// FamilySample is the struct to represent the result from the API
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

// Actor is the struct that is marshalled into from the API
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

// ScanBinary ----------------------------------------------------------

// BinaryScanMatches is the datatype used for the ScanBinary command. It is a struct for the json that is returned by the API
type BinaryScanMatches struct {
	MatchedStrings int64 `json:"matched_strings"`
	MatchedHits    int64 `json:"matched_hits"`
	Match          bool  `json:"match"`
}
