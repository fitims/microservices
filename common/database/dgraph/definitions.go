package dgraph

import (
	"errors"
	"fmt"
	"strconv"
)

var (
	ErrInvalidUID = errors.New("The data provided is not valid UID")
)

const DATE_FORMAT = "2006-01-02"
const DATE_TIME_FORMAT = "2006-01-02T15:04:05"

type UID uint64
type UIDHex string

func (u *UID) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`%d`, *u)), nil
}

func (u *UID) UnmarshalJSON(data []byte) error {

	strData := string(data[3 : len(data)-1])
	id, err := strconv.ParseUint(strData, 16, 64)
	if err != nil {
		fmt.Println("Error parsing data : ", strData)
		return err
	}

	*u = UID(id)
	return nil
}

func (u *UIDHex) MarshalJSON() ([]byte, error) {
	return []byte(string(*u)), nil
}

func (u *UIDHex) UnmarshalJSON(data []byte) error {

	strData := string(data)

	_, err := strconv.ParseUint(strData[3 : len(data)-1], 16, 64)
	if err != nil {
		fmt.Println("Error parsing data : ", strData, " parsing : ", strData[3: len(data)-1])
		return ErrInvalidUID
	}

	*u = UIDHex(strData)
	return nil
}

func (u *UID) ToHex() UIDHex {
	return UIDHex(fmt.Sprintf("0x%x", *u))
}

func (u *UID) ToTag() string {
	return fmt.Sprintf("<0x%x>", *u)
}

func ToUID(id string) (UID, error) {
	if len(id) < 3 {
		return 0, ErrInvalidUID
	}
	uintID, err := strconv.ParseUint(id[2:], 16, 64)
	if err != nil {
		return 0, err
	}

	return UID(uintID), nil
}
