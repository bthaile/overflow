package overflow

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/onflow/flow-go-sdk"
	"github.com/sanity-io/litter"
)

// Event parsing
//
// a type alias to an OverflowEventFilter to filter out all events with a given suffix and the fields with given suffixes
type OverflowEventFilter map[string][]string

type OverflowEventList []OverflowEvent

// a type holding all events that are emitted from a Transaction
type OverflowEvents map[string]OverflowEventList

type OverflowEvent struct {
	Fields        map[string]interface{} `json:"fields"`
	TransactionId string                 `json:"transactionID"`
	Name          string                 `json:"name"`
}

// Check if an event exist in the other events
func (o OverflowEvent) ExistIn(events []OverflowEvent) bool {
	for _, ev := range events {
		if litter.Sdump(o.Fields) == litter.Sdump(ev.Fields) {
			return true
		}
	}
	return false
}

func (e OverflowEventList) MarshalAs(marshalTo interface{}) error {
	bytes, err := json.Marshal(e)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, marshalTo)
	if err != nil {
		return err
	}
	return nil
}

func (e OverflowEvent) MarshalAs(marshalTo interface{}) error {
	bytes, err := json.Marshal(e)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, marshalTo)
	if err != nil {
		return err
	}
	return nil
}

// Parse raw flow events into a list of events and a fee event
func parseEvents(events []flow.Event) (OverflowEvents, OverflowEvent) {
	overflowEvents := OverflowEvents{}
	fee := OverflowEvent{}
	for _, event := range events {

		var fieldNames []string

		for _, eventTypeFields := range event.Value.EventType.Fields {
			fieldNames = append(fieldNames, eventTypeFields.Identifier)
		}

		finalFields := map[string]interface{}{}

		for id, field := range event.Value.Fields {
			name := fieldNames[id]
			value := CadenceValueToInterface(field)
			if value != nil {
				finalFields[name] = value
			}
		}

		events, ok := overflowEvents[event.Type]
		if !ok {
			events = []OverflowEvent{}
		}
		events = append(events, OverflowEvent{
			Fields:        finalFields,
			Name:          event.Type,
			TransactionId: event.TransactionID.String(),
		})
		overflowEvents[event.Type] = events

		if strings.HasSuffix(event.Type, "FlowFees.FeesDeducted") {
			fee = OverflowEvent{
				Fields:        finalFields,
				Name:          event.Type,
				TransactionId: event.TransactionID.String(),
			}
		}

	}
	return overflowEvents, fee
}

// Filter out temp withdraw deposit events
func (overflowEvents OverflowEvents) FilterTempWithdrawDeposit() OverflowEvents {
	filteredEvents := overflowEvents
	for name, events := range overflowEvents {
		if strings.HasSuffix(name, "TokensWithdrawn") {

			withDrawnEvents := []OverflowEvent{}
			for _, value := range events {
				if value.Fields["from"] != nil {
					withDrawnEvents = append(withDrawnEvents, value)
				}
			}
			if len(withDrawnEvents) != 0 {
				filteredEvents[name] = withDrawnEvents
			} else {
				delete(filteredEvents, name)
			}
		}

		if strings.HasSuffix(name, "TokensDeposited") {
			depositEvents := []OverflowEvent{}
			for _, value := range events {
				if value.Fields["to"] != nil {
					depositEvents = append(depositEvents, value)
				}
			}
			if len(depositEvents) != 0 {
				filteredEvents[name] = depositEvents
			} else {
				delete(filteredEvents, name)
			}
		}
	}
	return filteredEvents
}

// Filtter out fee events
func (overflowEvents OverflowEvents) FilterFees(fee float64) OverflowEvents {

	filteredEvents := overflowEvents
	for name, events := range overflowEvents {
		if strings.HasSuffix(name, "FlowFees.FeesDeducted") {
			delete(filteredEvents, name)
		}

		if strings.HasSuffix(name, "FlowToken.TokensWithdrawn") {

			withDrawnEvents := []OverflowEvent{}
			for _, value := range events {
				if value.Fields["amount"].(float64) != fee {
					withDrawnEvents = append(withDrawnEvents, value)
				}
			}
			if len(withDrawnEvents) != 0 {
				filteredEvents[name] = withDrawnEvents
			} else {
				delete(filteredEvents, name)
			}
		}

		if strings.HasSuffix(name, "FlowToken.TokensDeposited") {
			depositEvents := []OverflowEvent{}
			for _, value := range events {
				if value.Fields["amount"].(float64) != fee {
					depositEvents = append(depositEvents, value)
				}
			}
			if len(depositEvents) != 0 {
				filteredEvents[name] = depositEvents
			} else {
				delete(filteredEvents, name)
			}

		}
	}
	return filteredEvents
}

func printOrLog(t *testing.T, s string) {
	if t == nil {
		fmt.Println(s)
	} else {
		t.Log(s)
		t.Helper()
	}
}
func (overflowEvents OverflowEvents) Print(t *testing.T) {
	if t != nil {
		t.Helper()
	}

	printOrLog(t, "=== Events ===")
	for name, eventList := range overflowEvents {
		for _, event := range eventList {
			printOrLog(t, name)
			length := 0
			for key := range event.Fields {
				keyLength := len(key)
				if keyLength > length {
					length = keyLength
				}
			}

			format := fmt.Sprintf("%%%ds -> %%v", length+2)
			for key, value := range event.Fields {
				printOrLog(t, fmt.Sprintf(format, key, value))
			}
		}
	}
}

// Filter out events given the sent in filter
func (overflowEvents OverflowEvents) FilterEvents(ignoreFields OverflowEventFilter) OverflowEvents {
	filteredEvents := OverflowEvents{}
	for name, events := range overflowEvents {

		//find if we should ignore fields
		ignoreFieldNames := []string{}
		for ignoreEvent, fields := range ignoreFields {
			if strings.HasSuffix(name, ignoreEvent) {
				ignoreFieldNames = fields
			}
		}

		eventList := []OverflowEvent{}
		for _, ev := range events {
			event := OverflowEvent{Fields: map[string]interface{}{}}
			for key, value := range ev.Fields {
				valid := true
				for _, ig := range ignoreFieldNames {
					if strings.HasSuffix(key, ig) {
						valid = false
					}
				}
				if valid {
					event.Fields[key] = value
				}
			}
			if len(event.Fields) != 0 {
				eventList = append(eventList, event)
			}
		}
		if len(eventList) != 0 {
			filteredEvents[name] = eventList
		}
	}
	return filteredEvents
}
