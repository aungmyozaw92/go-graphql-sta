package utils

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"regexp"

	"github.com/aungmyozaw92/go-graphql/config"
	"github.com/ttacon/libphonenumber"
)


func IsValidEmail(email string) bool {
	// Basic email validation regex pattern
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(email)
}


func ValidatePhoneNumber(phoneNumber, countryCode string) error {
	p, err := libphonenumber.Parse(phoneNumber, countryCode)
	if err != nil {
		return err // Phone number is invalid
	}

	if !libphonenumber.IsValidNumber(p) {
		return fmt.Errorf("phone number is not valid")
	}

	return nil // Phone number is valid for the specified country code
}

// check if id exists, using ctx's business_id in WHERE, return RecordNOtFound Error
func ValidateResourceId[T any](ctx context.Context, id interface{}) error {

	count, err := ResourceCountWhere[T](ctx, "id = ?", id)
	if err != nil {
		return err
	}
	if count <= 0 {
		typeName := GetTypeName[T]()
		return fmt.Errorf("%s record not found", typeName)
	}

	return nil
}

// check if ALL id exists, using ctx's business_id in WHERE, return RecordNOtFound Error
func ValidateResourcesId[M any, ID comparable](ctx context.Context, ids []ID) error {
	unqIds := UniqueSlice(ids)

	count, err := ResourceCountWhere[M](ctx, "id IN ?", unqIds)
	if err != nil {
		return err
	}
	if count != int64(len(unqIds)) {
		typeName := GetTypeName[M]()
		return fmt.Errorf("%s record not found", typeName)
	}

	return nil
}

func ValidateUnique[T any](ctx context.Context, column string, value interface{}, exceptId interface{}) error {
	var count int64
	var err error
	if reflect.ValueOf(exceptId).IsZero() {
		count, err = ResourceCountWhere[T](ctx, column+" = ?", value)
	} else {
		count, err = ResourceCountWhere[T](ctx, column+" = ? AND NOT id = ?", value, exceptId)
	}

	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("duplicate " + column)
	}
	return nil
}

// count records, using WHERE business_id = ? AND $condition
// business_id can be blank for admin user
func ResourceCountWhere[T any](ctx context.Context, condition string, value ...interface{}) (int64, error) {
	var model T

	db := config.GetDB()
	dbCtx := db.WithContext(ctx).Model(&model)
	var count int64
	dbCtx.Where(condition, value...)
	if err := dbCtx.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
