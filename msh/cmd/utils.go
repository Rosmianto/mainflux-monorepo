package cmd

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"fmt"

	"github.com/hokaccha/go-prettyjson"
)

func PrettyJson(body string) string {
	pj, err := prettyjson.Format([]byte(body))
	if err != nil {
		return err.Error()
	} else {
		return string(pj)
	}
}

func PrettyHttpResp(resp *http.Response, err error) string {
	if err != nil {
		return `{"error": "` + err.Error() + `"}`
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return `{"error": "` + err.Error() + `"}`
	}

	str := "Status code: " + strconv.Itoa(resp.StatusCode) + " - " +
		   http.StatusText(resp.StatusCode) + "\n\n"

	if resp.StatusCode == 201 {
		return str + fmt.Sprintf("Resource location: %s",
			                     resp.Header.Get("Location"))
	} else {
		return str + PrettyJson(string(body))
	}
}
