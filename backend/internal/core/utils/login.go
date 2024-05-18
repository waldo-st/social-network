package util

import (
	"errors"
	"fmt"
	"regexp"
	"social/internal/core/domain"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/html"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(html.EscapeString(password)), bcrypt.DefaultCost)

	return string(hashedPassword), err
}

func ComparePassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func CheckEmail(email string) error {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]{2,}\.[a-zA-Z]{2,}$`)

	if re.MatchString(email) {
		return nil
	}
	return errors.New("invalid email")
}

func CheckName(values []string) error {
	err := map[int]string{
		0: "Firstname",
		1: "Lastname",
		2: "Nickname",
	}

	for i, v := range values {
		if len(values) == 3 && i != 2 {
			continue
		}
		fmt.Println("len:", len(values), v)
		if len(v) > 20 || len(v) < 2 {
			return errors.New(err[i] + " must be between 2 and 20 characters")
		}

		if v == "" || strings.Contains(v, " ") {
			return errors.New(err[i] + " must not be empty and must not contain spaces")
		}

		re := regexp.MustCompile(`^[a-zA-Z]+$`)
		if !re.MatchString(v) && i != 2 {
			return errors.New(err[i] + " must only contain alphabet letters")
		}

	}

	return nil
}

func CheckOptional(m map[string]string) error {
	exts := map[string]bool{
		"jpeg": true,
		"png ": true,
		"gif":  true,
		"svg":  true,
		"jpg ": true,
	}

	if m["avatar"] != "" {
		if strings.TrimSpace(m["avatar"]) == "" {
			return errors.New("invalid file name")
		} else {
			// ext := filepath.Ext(m["avatar"])
			ext := "jpeg"
			if !exts[ext] {
				return errors.New("invalid extension")
			}
		}
	}

	if m["nickname"] != "" {
		err := CheckName([]string{"", "", m["nickname"]})
		if err != nil {
			return err
		}
	}

	if m["about"] != "" {
		if strings.TrimSpace(m["avatar"]) == "" {
			return errors.New("should have at least one character")
		}
	}

	return nil
}

func CheckDateOfBirth(date time.Time) error {
	now := time.Now()
	age := now.Year() - date.Year()
	// Check if the birthday has not occurred in the current year and, if so, subtract one year from the age
	if now.YearDay() < date.YearDay() {
		age--
	}
	if age < 18 || age > 120 {
		return errors.New("sorry, you must be between 18 and 120 years old")
	}

	return nil
}

func CheckFields(user *domain.User) error {

	if err := CheckName([]string{user.FirstName, user.LastName}); err != nil {
		return err
	}

	if err := CheckDateOfBirth(user.DateOfBirth); err != nil {
		return err
	}

	if err := CheckEmail(user.Email); err != nil {
		return err
	}

	m := map[string]string{
		"about":    user.About,
		"avatar":   user.Avatar,
		"nickname": user.Nickname,
	}

	if len(user.Nickname) > 0 && m["nickname"] == "" {
		return errors.New("nickname must not be empty")
	}
	if err := CheckOptional(m); err != nil {
		return err
	}

	return nil
}
