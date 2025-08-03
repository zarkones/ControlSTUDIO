package profiles

import (
	"common/slices"
	"encoding/base64"
	"encoding/hex"
	"errors"

	"github.com/zarkones/netescape"
)

var (
	ErrUnknownOp = errors.New("unknown operation")
)

func OperateData(operations *[]Operation, body *string) (data string, err error) {
	data = *body

	for _, operation := range *operations {
		switch operation.Action {

		default:
			return data, ErrUnknownOp

		// SET PREFIX
		case "prefix":
			data = slices.Rand(&operation.Value) + data

		// SET SUFFIX
		case "suffix":
			data = data + slices.Rand(&operation.Value)

		// HEX
		case "hex_d":
			decoded, err := hex.DecodeString(data)
			if err != nil {
				return "", err
			}
			data = string(decoded)
		case "hex":
			data = hex.EncodeToString([]byte(data))

		// BASE64 STD
		case "base64_std_d":
			decoded, err := base64.StdEncoding.DecodeString(data)
			if err != nil {
				return "", err
			}
			data = string(decoded)
		case "base64_std":
			data = base64.StdEncoding.EncodeToString([]byte(data))

		// BASE64 URL
		case "base64_url_d":
			decoded, err := base64.URLEncoding.DecodeString(data)
			if err != nil {
				return "", err
			}
			data = string(decoded)
		case "base64_url":
			data = base64.URLEncoding.EncodeToString([]byte(data))

		// CSV
		case "csv_d":
			decoded, err := netescape.FromCSV(&data)
			if err != nil {
				return "", err
			}
			data = decoded
		case "csv":
			decoded, err := netescape.ToCsv(&data)
			if err != nil {
				return "", err
			}
			data = decoded

		}
	}

	return data, nil
}
