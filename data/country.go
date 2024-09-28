package data

import (
	"context"
	"time"
)

type Country struct {
	Id               int       `json:"id"`
	CountryName      string    `json:"country_name"`
	CountryAlphaName string    `json:"country_alpha_name"`
	RegionName       string    `json:"region_name"`
	RegionAlphaName  string    `json:"region_alpha_name"`
	UpdatedAt        time.Time `db:"updated_at"`
	UpdatedBy        string    `db:"updated_by"`
	RegisteredAt     time.Time `db:"registered_at"`
	RegisteredBy     string    `db:"registered_by"`
}

func GetCountries(ctx context.Context, connection DbActions) ([]Country, error) {
	countryCount := connection.QueryRow(ctx, "SELECT count(*) FROM authmantledb.country")
	ctrlCount := new(int)
	err := countryCount.Scan(ctrlCount)
	if err != nil {
		return nil, err
	}

	ctrlIndex := 0
	countries := make([]Country, *ctrlCount)
	rows, err := connection.Query(ctx, "SELECT * FROM authmantledb.country")
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var country Country
		err := rows.Scan(
			&country.Id,
			&country.CountryName,
			&country.CountryAlphaName,
			&country.RegionName,
			&country.RegionAlphaName,
			&country.UpdatedAt,
			&country.UpdatedBy,
			&country.RegisteredAt,
			&country.RegisteredBy,
		)
		if err != nil {
			return nil, err
		}
		countries[ctrlIndex] = country
		ctrlIndex++
	}

	return countries, nil

}
