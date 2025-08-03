package profiles

import (
	"common/slices"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"strings"

	"github.com/zarkones/netescape"
)

var (
	ErrUnknownOp = errors.New("unknown operation")
)

func OperateDataReverseOperations(operations *[]Operation, body *string, obfuscate bool) (data string, err error) {
	reversedOperations := make([]Operation, len(*operations))
	reversedOpsIndex := 0
	originalOpsIndex := len(*operations) - 1
	for reversedOpsIndex < len(*operations) {
		reversedOperations[reversedOpsIndex] = Operation{
			Action: (*operations)[originalOpsIndex].Action,
			Value:  (*operations)[originalOpsIndex].Value,
		}
		reversedOpsIndex++
		originalOpsIndex--
	}

	return OperateData(&reversedOperations, body, obfuscate)
}

func OperateData(operations *[]Operation, body *string, obfuscate bool) (data string, err error) {
	data = *body

	if obfuscate {
		for _, operation := range *operations {
			switch operation.Action {

			default:
				return data, ErrUnknownOp

			case "prefix":
				data = slices.Rand(&operation.Value) + data

			case "suffix":
				data = data + slices.Rand(&operation.Value)

			case "hex":
				data = hex.EncodeToString([]byte(data))

			case "base64_std":
				data = base64.StdEncoding.EncodeToString([]byte(data))

			case "base64_url":
				data = base64.URLEncoding.EncodeToString([]byte(data))

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

	for _, operation := range *operations {
		switch operation.Action {

		default:
			return data, ErrUnknownOp

		case "prefix":
			for _, value := range operation.Value {
				data = strings.TrimPrefix(data, value)
			}

		case "suffix":
			for _, value := range operation.Value {
				data = strings.TrimSuffix(data, value)
			}

		case "hex":
			decoded, err := hex.DecodeString(data)
			if err != nil {
				return "", err
			}
			data = string(decoded)

		case "base64_std":
			decoded, err := base64.StdEncoding.DecodeString(data)
			if err != nil {
				return "", err
			}
			data = string(decoded)

		case "base64_url":
			decoded, err := base64.URLEncoding.DecodeString(data)
			if err != nil {
				return "", err
			}
			data = string(decoded)

		case "csv":
			decoded, err := netescape.FromCSV(&data)
			if err != nil {
				return "", err
			}
			data = decoded

		}
	}

	return data, nil
}
