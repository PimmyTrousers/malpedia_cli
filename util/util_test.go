package util

import (
	"crypto/sha256"
	"fmt"
	"os"
	"testing"
)

var PrettyPrintJSON = [][]byte{
	[]byte("[\"[unnamed_group]\",\"zoopark\"]"),
}

var Base64Decode = []struct {
	in  string
	out string
}{
	{"YXNkYXNmYWFnYWRzZmFzZGZhZmFzZmQ=", "asdasfaagadsfasdfafasfd"},
	{"dGVzdDEyMw==", "test123"},
	{"dGVzdDEyMzQ1", "test12345"},
}

var FamilyNameTranslation = []struct {
	in  string
	out string
}{
	{"revil", "win.sodinokibi"},
	{"stuxnet", "win.stuxnet"},
	{"ursnif", "win.snifula"},
}

var ActorNameTranslation = []struct {
	in  string
	out string
}{
	{"apt28", "sofacy"},
	{"cozybear", "apt_29"},
	{"lazarus", "lazarus_group"},
}

var DownloadSamples = []string{
	"e0eea847f58efe604287a0fa9abe84576235dbbfa5f3e9636dcda10092c015b1",
	"7fb5df4519284e44b621265aca19e61e8e5f00351087ddf2cf838935261774e9",
	"c4a7f8b8046a6623cd7909bacb1cbef13471a4efd8adb4aedbf7fc1377ab502d",
}

var HashType = []struct {
	in  string
	out Hash
}{
	{"e0eea847f58efe604287a0fa9abe84576235dbbfa5f3e9636dcda10092c015b1", SHA256},
	{"0800fc577294c34e0b28ad2839435945", MD5},
	{"2346ad27d7568ba9896f1b7da6b5991251debdf2", SHA1},
}

var APIKeyValidity = []struct {
	in  string
	out bool
}{
	{"2346ad27d7568ba9896f1b7da6b5991251debdf2", true},
	{"2346ad27d7568ba9896f1b7da6b5991251debdf", false},
	{"", false},
}

func TestPrettyPrintJSON(t *testing.T) {
	for _, k := range PrettyPrintJSON {
		err := PrettyPrintJson(k)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestGetFamilyName(t *testing.T) {
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		t.Error("apikey not set")
		return
	}

	for _, familyName := range FamilyNameTranslation {
		serviceName, err := GetFamilyName(familyName.in, apiKey)
		if err != nil {
			t.Error(err)
			return
		}

		if serviceName != familyName.out {
			t.Errorf("expected %s but got %s", familyName.out, serviceName)
		}
	}
}

func TestGetActorName(t *testing.T) {
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		t.Error("apikey not set")
		return
	}

	for _, actorName := range ActorNameTranslation {
		serviceName, err := GetActorName(actorName.in, apiKey)
		if err != nil {
			t.Error(err)
			return
		}

		if serviceName != actorName.out {
			t.Errorf("expected %s but got %s", actorName.out, serviceName)
		}
	}
}

func TestIsValidAPIKey(t *testing.T) {
	for _, apiKey := range APIKeyValidity {
		if IsAPIKeyValid(apiKey.in) != apiKey.out {
			t.Errorf("expected result %t but got %t", apiKey.out, IsAPIKeyValid(apiKey.in))
		}
	}
}

func TestBase64DecodeContent(t *testing.T) {
	for _, test := range Base64Decode {
		buf, err := Base64DecodeContent(test.in)
		if err != nil {
			t.Error(err)
			return
		}

		if buf.String() != test.out {
			t.Errorf("expected %s but got %s", test.out, buf.String())
			return
		}
	}
}

func TestGetHashType(t *testing.T) {
	for _, hash := range HashType {
		hashtype, err := GetHashType(hash.in)
		if err != nil {
			t.Error(err)
			return
		}

		if hashtype != hash.out {
			t.Error("incorrect hash type")
		}
	}
}

func TestDownloadSample(t *testing.T) {
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		t.Error("apikey not set")
	}
	for _, revilHash := range DownloadSamples {
		samples, err := DownloadSample(revilHash, apiKey)
		if err != nil {
			t.Error(err)
		}

		match := false

		for _, sample := range *samples {
			h := sha256.New()
			h.Write(sample.Bytes())
			hashStr := fmt.Sprintf("%x", h.Sum(nil))
			if hashStr == revilHash {
				match = true
			}
		}

		if !match {
			t.Error("failed to find a sample with a matching hash")
		}
	}
}

func TestDumpRaw(t *testing.T) {
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		t.Error("apikey not set")
	}
	for _, revilHash := range DownloadSamples {
		samples, err := DownloadSample(revilHash, apiKey)
		if err != nil {
			t.Error(err)
		}
		err = DumpRaw(samples, revilHash)
		if err != nil {
			t.Error(err)
			return
		}

		for k := range *samples {
			err := os.Remove(k + "_" + revilHash)
			if err != nil {
				t.Error(err)
				return
			}
		}
	}
}

func TestDumpZip(t *testing.T) {
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		t.Error("apikey not set")
	}
	for _, revilHash := range DownloadSamples {
		samples, err := DownloadSample(revilHash, apiKey)
		if err != nil {
			t.Error(err)
		}
		err = DumpZip(samples, revilHash, "test.zip")
		if err != nil {
			t.Error(err)
			return
		}

		err = os.Remove("test.zip")
		if err != nil {
			t.Error(err)
			return
		}
	}
}
