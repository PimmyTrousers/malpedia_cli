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

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "apitoken "+apiKey)

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

func Base64DecodeContent(buf *string) (*bytes.Buffer, error) {
	b := make([]byte, 8)

	oldBuf := bytes.NewBufferString(*buf)
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

		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
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
