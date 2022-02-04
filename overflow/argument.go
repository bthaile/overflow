package overflow

import (
	"fmt"
	"time"

	"github.com/araddon/dateparse"
	"github.com/onflow/cadence"
	"github.com/onflow/flow-go-sdk"
)

// FlowArgumentsBuilder used to create a builder pattern for a transaction
type FlowArgumentsBuilder struct {
	Overflow  *Overflow
	Arguments []cadence.Value
}

func (f *Overflow) Arguments() *FlowArgumentsBuilder {
	return &FlowArgumentsBuilder{
		Overflow:  f,
		Arguments: []cadence.Value{},
	}
}

func (a *FlowArgumentsBuilder) Build() []cadence.Value {
	return a.Arguments
}

// RawAccount add an account from a string as an argument
func (a *FlowArgumentsBuilder) RawAccount(key string) *FlowArgumentsBuilder {
	account := flow.HexToAddress(key)
	accountArg := cadence.BytesToAddress(account.Bytes())
	return a.Argument(accountArg)
}

// Account add an account as an argument
func (a *FlowArgumentsBuilder) Account(key string) *FlowArgumentsBuilder {
	f := a.Overflow

	account := f.Account(key)
	return a.Argument(cadence.BytesToAddress(account.Address().Bytes()))
}

// String add a String Argument to the transaction
func (a *FlowArgumentsBuilder) String(value string) *FlowArgumentsBuilder {
	return a.Argument(cadence.String(value))
}

// Boolean add a Boolean Argument to the transaction
func (a *FlowArgumentsBuilder) Boolean(value bool) *FlowArgumentsBuilder {
	return a.Argument(cadence.NewBool(value))
}

// Bytes add a Bytes Argument to the transaction
func (a *FlowArgumentsBuilder) Bytes(value []byte) *FlowArgumentsBuilder {
	return a.Argument(cadence.NewBytes(value))
}

// Int add an Int Argument to the transaction
func (a *FlowArgumentsBuilder) Int(value int) *FlowArgumentsBuilder {
	return a.Argument(cadence.NewInt(value))
}

// Int8 add an Int8 Argument to the transaction
func (a *FlowArgumentsBuilder) Int8(value int8) *FlowArgumentsBuilder {
	return a.Argument(cadence.NewInt8(value))
}

// Int16 add an Int16 Argument to the transaction
func (a *FlowArgumentsBuilder) Int16(value int16) *FlowArgumentsBuilder {
	return a.Argument(cadence.NewInt16(value))
}

// Int32 add an Int32 Argument to the transaction
func (a *FlowArgumentsBuilder) Int32(value int32) *FlowArgumentsBuilder {
	return a.Argument(cadence.NewInt32(value))
}

// Int64 add an Int64 Argument to the transaction
func (a *FlowArgumentsBuilder) Int64(value int64) *FlowArgumentsBuilder {
	return a.Argument(cadence.NewInt64(value))
}

// Int128 add an Int128 Argument to the transaction
func (a *FlowArgumentsBuilder) Int128(value int) *FlowArgumentsBuilder {
	return a.Argument(cadence.NewInt128(value))
}

// Int256 add an Int256 Argument to the transaction
func (a *FlowArgumentsBuilder) Int256(value int) *FlowArgumentsBuilder {
	return a.Argument(cadence.NewInt256(value))
}

// UInt add an UInt Argument to the transaction
func (a *FlowArgumentsBuilder) UInt(value uint) *FlowArgumentsBuilder {
	return a.Argument(cadence.NewUInt(value))
}

// UInt8 add an UInt8 Argument to the transaction
func (a *FlowArgumentsBuilder) UInt8(value uint8) *FlowArgumentsBuilder {
	return a.Argument(cadence.NewUInt8(value))
}

// UInt16 add an UInt16 Argument to the transaction
func (a *FlowArgumentsBuilder) UInt16(value uint16) *FlowArgumentsBuilder {
	return a.Argument(cadence.NewUInt16(value))
}

// UInt32 add an UInt32 Argument to the transaction
func (a *FlowArgumentsBuilder) UInt32(value uint32) *FlowArgumentsBuilder {
	return a.Argument(cadence.NewUInt32(value))
}

// UInt64 add an UInt64 Argument to the transaction
func (a *FlowArgumentsBuilder) UInt64(value uint64) *FlowArgumentsBuilder {
	return a.Argument(cadence.NewUInt64(value))
}

// UInt128 add an UInt128 Argument to the transaction
func (a *FlowArgumentsBuilder) UInt128(value uint) *FlowArgumentsBuilder {
	return a.Argument(cadence.NewUInt128(value))
}

// UInt256 add an UInt256 Argument to the transaction
func (a *FlowArgumentsBuilder) UInt256(value uint) *FlowArgumentsBuilder {
	return a.Argument(cadence.NewUInt256(value))
}

// Word8 add a Word8 Argument to the transaction
func (a *FlowArgumentsBuilder) Word8(value uint8) *FlowArgumentsBuilder {
	return a.Argument(cadence.NewWord8(value))
}

// Word16 add a Word16 Argument to the transaction
func (a *FlowArgumentsBuilder) Word16(value uint16) *FlowArgumentsBuilder {
	return a.Argument(cadence.NewWord16(value))
}

// Word32 add a Word32 Argument to the transaction
func (a *FlowArgumentsBuilder) Word32(value uint32) *FlowArgumentsBuilder {
	return a.Argument(cadence.NewWord32(value))
}

// Word64 add a Word64 Argument to the transaction
func (a *FlowArgumentsBuilder) Word64(value uint64) *FlowArgumentsBuilder {
	return a.Argument(cadence.NewWord64(value))
}

// Fix64 add a Fix64 Argument to the transaction
func (a *FlowArgumentsBuilder) Fix64(value string) *FlowArgumentsBuilder {
	amount, err := cadence.NewFix64(value)
	if err != nil {
		panic(err)
	}
	return a.Argument(amount)
}

// DateStringAsUnixTimestamp sends a dateString parsed in the timezone as a unix timeszone ufix
func (a *FlowArgumentsBuilder) DateStringAsUnixTimestamp(dateString string, timezone string) *FlowArgumentsBuilder {
	value := parseTime(dateString, timezone)
	amount, err := cadence.NewUFix64(value)
	if err != nil {
		panic(err)
	}
	return a.Argument(amount)
}

func parseTime(timeString string, location string) string {
	loc, err := time.LoadLocation(location)
	if err != nil {
		panic(err)
	}

	time.Local = loc
	t, err := dateparse.ParseLocal(timeString)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%d.0", t.Unix())
}

// UFix64 add a UFix64 Argument to the transaction
func (a *FlowArgumentsBuilder) UFix64(input float64) *FlowArgumentsBuilder {
	value := fmt.Sprintf("%f", input)
	amount, err := cadence.NewUFix64(value)
	if err != nil {
		panic(err)
	}
	return a.Argument(amount)
}

// Argument add an argument to the transaction
func (a *FlowArgumentsBuilder) Argument(value cadence.Value) *FlowArgumentsBuilder {
	a.Arguments = append(a.Arguments, value)
	return a
}

// Argument add an argument to the transaction
func (a *FlowArgumentsBuilder) StringArray(value ...string) *FlowArgumentsBuilder {
	array := []cadence.Value{}
	for _, val := range value {
		stringVal, err := cadence.NewString(val)
		if err != nil {
			//TODO: what to do with errors here? Accumulate in builder?
			panic(err)
		}
		array = append(array, stringVal)
	}
	a.Arguments = append(a.Arguments, cadence.NewArray(array))
	return a
}

// Argument add an argument to the transaction
func (a *FlowArgumentsBuilder) UInt64Array(value ...uint64) *FlowArgumentsBuilder {
	array := []cadence.Value{}
	for _, val := range value {
		cadenceVal := cadence.NewUInt64(val)
		array = append(array, cadenceVal)
	}
	a.Arguments = append(a.Arguments, cadence.NewArray(array))
	return a
}

// Argument add an argument to the transaction
func (a *FlowArgumentsBuilder) UInt8Array(value ...uint8) *FlowArgumentsBuilder {
	array := []cadence.Value{}
	for _, val := range value {
		cadenceVal := cadence.NewUInt8(val)
		array = append(array, cadenceVal)
	}
	a.Arguments = append(a.Arguments, cadence.NewArray(array))
	return a
}