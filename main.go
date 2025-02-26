package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/urfave/cli/v2"
	"gopkg.in/auth0.v5/management"
)

func main() {
	app := &cli.App{
		Name:  "auth0-user-export",
		Usage: "Export Auth0 users to CSV format",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "domain",
				Aliases:  []string{"d"},
				Usage:    "Auth0 domain",
				Required: true,
				EnvVars:  []string{"AUTH0_DOMAIN"},
			},
			&cli.StringFlag{
				Name:     "client-id",
				Aliases:  []string{"i"},
				Usage:    "Auth0 client ID",
				Required: true,
				EnvVars:  []string{"AUTH0_CLIENT_ID"},
			},
			&cli.StringFlag{
				Name:     "client-secret",
				Aliases:  []string{"s"},
				Usage:    "Auth0 client secret",
				Required: true,
				EnvVars:  []string{"AUTH0_CLIENT_SECRET"},
			},
			&cli.StringFlag{
				Name:    "fields",
				Aliases: []string{"f"},
				Usage:   "Comma-separated list of user fields to include in CSV (default: user_id,email,name,app_metadata)",
				Value:   "user_id,email,name,app_metadata",
			},
		},
		Action: exportUsers,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func exportUsers(c *cli.Context) error {
	domain := c.String("domain")
	clientID := c.String("client-id")
	clientSecret := c.String("client-secret")
	fields := strings.Split(c.String("fields"), ",")

	// Initialize Auth0 management API client
	m, err := management.New(
		domain,
		management.WithClientCredentials(clientID, clientSecret),
	)
	if err != nil {
		return fmt.Errorf("failed to initialize Auth0 management client: %w", err)
	}

	// Create CSV writer
	w := csv.NewWriter(os.Stdout)
	defer w.Flush()

	// Write CSV header
	if err := w.Write(fields); err != nil {
		return fmt.Errorf("error writing CSV header: %w", err)
	}

	// Fetch users from Auth0
	var page int
	for {
		users, err := m.User.List(
			management.Page(page),
			management.PerPage(100),
			management.IncludeFields("user_id", "email", "name", "app_metadata"),
		)
		if err != nil {
			return fmt.Errorf("failed to fetch users: %w", err)
		}

		if len(users.Users) == 0 {
			break
		}

		// Write each user to CSV
		for _, user := range users.Users {
			record := make([]string, len(fields))
			for i, field := range fields {
				value := ""

				switch field {
				case "user_id":
					value = user.GetID()
				case "email":
					value = user.GetEmail()
				case "name":
					value = user.GetName()
				case "app_metadata":
					if user.AppMetadata != nil {
						metadata, err := json.Marshal(user.AppMetadata)
						if err == nil {
							value = string(metadata)
						}
					}
				default:
					// For other fields, try to extract as string
					val, _ := getFieldAsString(user, field)
					value = val
				}

				record[i] = value
			}

			if err := w.Write(record); err != nil {
				return fmt.Errorf("error writing user record: %w", err)
			}
		}

		page++
	}

	return nil
}

// Helper function to extract field values from the user struct
func getFieldAsString(user *management.User, field string) (string, error) {
	val := reflect.ValueOf(user).Elem()
	fieldVal := val.FieldByName(strings.Title(field))

	if !fieldVal.IsValid() {
		return "", fmt.Errorf("invalid field: %s", field)
	}

	switch fieldVal.Kind() {
	case reflect.String:
		return fieldVal.String(), nil
	case reflect.Ptr:
		if fieldVal.IsNil() {
			return "", nil
		}
		return fmt.Sprintf("%v", fieldVal.Elem().Interface()), nil
	case reflect.Map, reflect.Struct, reflect.Slice:
		b, err := json.Marshal(fieldVal.Interface())
		if err != nil {
			return "", err
		}
		return string(b), nil
	default:
		return fmt.Sprintf("%v", fieldVal.Interface()), nil
	}
}
