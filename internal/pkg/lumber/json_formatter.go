package lumber

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
)

// HACK: this class remove the key "msg" if message is empty.
//
//       https://github.com/sirupsen/logrus/blob/v1.4.2/json_formatter.go#L89

type JSONFormatter struct {
	*logrus.JSONFormatter
}

func (f *JSONFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	original, err := f.JSONFormatter.Format(entry)
	if err != nil {
		return nil, err
	}
	values, err := f.decode(original)
	if err != nil {
		return nil, err
	}
	return f.encode(values)
}

func (f *JSONFormatter) decode(data []byte) (map[string]interface{}, error) {
	var values map[string]interface{}
	if err := json.Unmarshal(data, &values); err != nil {
		return nil, err
	}

	switch v := values[logrus.FieldKeyMsg].(type) {
	case string:
		if len(v) == 0 {
			delete(values, logrus.FieldKeyMsg)
		}
	default:
		break
	}

	return values, nil
}

func (f *JSONFormatter) encode(values map[string]interface{}) ([]byte, error) {
	if f.PrettyPrint {
		compact, err := json.MarshalIndent(values, "", "  ")
		if err != nil {
			return nil, err
		}
		return append(compact, '\n'), nil
	} else {
		compact, err := json.Marshal(values)
		if err != nil {
			return nil, err
		}
		return append(compact, '\n'), nil
	}
}
