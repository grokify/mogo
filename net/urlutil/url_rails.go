package urlutil

import (
	"encoding/json"
	"fmt"

	"github.com/derekstavis/go-qs"
	"github.com/grokify/mogo/errors/errorsutil"
)

// UnmarshalRailsQS unmarshals a Rails query string to a Go struct.
func UnmarshalRailsQS(railsQuery string, i interface{}) error {
	query, err := qs.Unmarshal(railsQuery)
	if err != nil {
		return errorsutil.Wrap(err,
			fmt.Sprintf("urlutil.UnmarshalRailsQS__err__qs.Unmarshal[%s]", railsQuery))
	}
	bytes, err := json.Marshal(query)
	if err != nil {
		return errorsutil.Wrap(err,
			fmt.Sprintf("urlutil.UnmarshalRailsQS__err__json.Marshal[%s]", railsQuery))
	}
	err = json.Unmarshal(bytes, i)
	if err != nil {
		return errorsutil.Wrap(err,
			fmt.Sprintf("urlutil.UnmarshalRailsQS__err__json.Unmarshal[%s]", railsQuery))
	}
	return nil
}
