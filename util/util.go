package util

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/pimmytrousers/malpedia_cli/types"
)

// Hash is a enum to support different types of hashes
type Hash int

const (
	// MD5 is the enum to compare against for MD5 hashes
	MD5 Hash = iota
	// SHA1 is the enum to compare against for SHA1 hashes
	SHA1
	// SHA256 is the enum to compare against for SHA256 hashes
	SHA256
)

func HttpGetQuery(p types.Endpoint, apiKey string) ([]byte, error) {
	var resp *http.Response
	var err error

	urlParsed, err := url.Parse(types.APIBase + string(p))
	if err != nil {
		return nil, err
	}

	urlStr := urlParsed.String()

	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "apitoken "+apiKey)

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusNotFound {
		return nil, types.ErrResourceNotFound
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("got status code %d", resp.StatusCode)
	}
	return buf, nil
}

func GetActorName(inputActor string, apikey string) (string, error) {
	endpoint := fmt.Sprintf(types.EndpointFindActor, inputActor)
	res, err := HttpGetQuery(types.Endpoint(endpoint), apikey)
	if err != nil {
		return "", err
	}

	findActor := &types.FindActor{}
	err = json.Unmarshal(res, findActor)
	if err != nil {
		return "", err
	}

	var actorName string
	for _, actor := range *findActor {
		actorName = actor.Name
	}

	return actorName, nil
}

func GetFamilyName(inputFamily string, apikey string) (string, error) {
	endpoint := fmt.Sprintf(types.EndpointFindFamily, inputFamily)

	resp, err := HttpGetQuery(types.Endpoint(endpoint), apikey)
	if err != nil {
		return "", err
	}

	foundFamilies := &types.FindFamily{}
	err = json.Unmarshal(resp, foundFamilies)
	if err != nil {
		return "", err
	}

	// TODO: What if this returns multiple entries?
	var name string
	for _, family := range *foundFamilies {
		name = family.Name
	}

	return name, nil
}

func HttpMultipartFileUpload(p types.Endpoint, apiKey string, body io.Reader, filename string) ([]byte, error) {
	urlParsed, err := url.Parse(types.APIBase + string(p))
	if err != nil {
		return nil, err
	}

	urlStr := urlParsed.String()

	reqBody := &bytes.Buffer{}
	writer := multipart.NewWriter(reqBody)
	defer writer.Close()

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(part, body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", urlStr, reqBody)
	if err != nil {
		return nil, err
	}
	// buf, _ := ioutil.ReadAll(req.Body)

	req.Header.Add("content-type", writer.FormDataContentType())
	req.Header.Add("Authorization", "apitoken "+apiKey)
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	req.Header.Add("Connection", "keep-alive")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func HttpRawFileUpload(p types.Endpoint, apiKey string, body io.Reader, filename string) ([]byte, error) {
	urlParsed, err := url.Parse(types.APIBase + string(p))
	if err != nil {
		return nil, err
	}

	urlStr := urlParsed.String()
	req, err := http.NewRequest("POST", urlStr, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "apitoken "+apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func PrettyPrintJson(buf []byte) error {
	var prettyJSON bytes.Buffer
	err := json.Indent(&prettyJSON, buf, "", "\t")
	if err != nil {
		return err
	}

	fmt.Printf("%s", prettyJSON.Bytes())
	return nil
}

func GetHashType(hash string) (Hash, error) {
	digest, err := hex.DecodeString(hash)
	if err != nil {
		return 0, err
	}

	switch len(digest) {
	case 16:
		return MD5, nil
	case 20:
		return SHA1, nil
	case 32:
		return SHA256, nil
	}

	return 0, errors.New("invalid hash type")
}

// Base64DecodeContent takes in
func Base64DecodeContent(buf string) (*bytes.Buffer, error) {
	// TODO: what is a good balance here between size and efficiency
	b := make([]byte, 1)

	oldBuf := bytes.NewBufferString(buf)
	newBuf := bytes.NewBuffer(nil)

	decoder := base64.NewDecoder(base64.StdEncoding, oldBuf)

	for {
		_, err := decoder.Read(b)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		newBuf.Write(b)
		if err != nil {
			return nil, err
		}
	}

	return newBuf, nil
}

func Unzip(src string, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {

		fpath := filepath.Join(dest, f.Name)

		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("%s: illegal file path", fpath)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(outFile, rc)

		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}
	return nil
}

func DownloadSample(hash string, apiKey string) (types.SampleState, error) {
	hashType, err := GetHashType(hash)
	if err != nil {
		return nil, err
	}

	if hashType != MD5 && hashType != SHA256 {
		return nil, errors.New("only md5 and sha256 hashes are allowed")
	}

	formattedEndpoint := fmt.Sprintf(types.EndpointGetSampleRaw, hash)

	res, err := HttpGetQuery(types.Endpoint(formattedEndpoint), apiKey)
	if err != nil {
		return nil, err
	}

	jsonBody := make(map[string]string)

	err = json.Unmarshal(res, &jsonBody)
	if err != nil {
		return nil, err
	}

	sampleMap := make(map[string]bytes.Buffer, len(jsonBody))

	for k, v := range jsonBody {
		buf, err := Base64DecodeContent(v)
		if err != nil {
			continue
		}

		sampleMap[k] = *buf

	}

	return &sampleMap, nil
}

// DumpRaw will take the malware states object and dump all the files to
// the current directory
func DumpRaw(states types.SampleState, hash string) error {
	for k, v := range *states {
		err := ioutil.WriteFile(k+"_"+hash, v.Bytes(), 0644)
		if err != nil {
			return err
		}
	}

	return nil
}

// DumpZip will take the malware states object and dump a single
// zip file to disk that contains all of the samples
func DumpZip(states types.SampleState, hash string, outZip string) error {
	f, err := os.Create(outZip)
	if err != nil {
		return err
	}
	defer f.Close()

	w := zip.NewWriter(f)
	defer w.Close()

	for k, v := range *states {
		f, err := w.Create(k + "_" + hash)
		if err != nil {
			return err
		}
		_, err = f.Write(v.Bytes())
		if err != nil {
			return err
		}
	}

	return nil
}

// IsAPIKeyValid is a convenience function to check if there is an
// api key value
func IsAPIKeyValid(key string) bool {
	if key == "" {
		return false
	} else if res, _ := GetHashType(key); res != SHA1 {
		return false
	}

	return true
}
