package output

import (
	"encoding/base64"
	"encoding/csv"
	"os"
	"strconv"

	"github.com/ffuf/ffuf/pkg/ffuf"
)

var header = []string{"input", "position", "status_code", "content_length", "content_words"}

func writeCSV(config *ffuf.Config, res []Result, encode bool) error {
	f, err := os.Create(config.OutputFile)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	if err := w.Write(header); err != nil {
		return err
	}
	for _, r := range res {
		if encode {
			inputs := make(map[string][]byte, 0)
			for k, v := range r.Input {
				inputs[k] = []byte(base64encode(v))
			}
			r.Input = inputs
		}

		err := w.Write(toCSV(r))
		if err != nil {
			return err
		}
	}
	return nil
}

func base64encode(in []byte) string {
	return base64.StdEncoding.EncodeToString(in)
}

func toCSV(r Result) []string {
	res := make([]string, 0)
	for k, v := range r.Input {
		res = append(res, k+":"+string(v))
	}
	res = append(res, strconv.Itoa(r.Position))
	res = append(res, strconv.FormatInt(r.StatusCode, 10))
	res = append(res, strconv.FormatInt(r.ContentLength, 10))
	res = append(res, strconv.FormatInt(r.ContentWords, 10))
	return res
}
