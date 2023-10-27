package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

    "rest-service/models"
)

func enrich(person *models.Person) error {
    var err error
    person.Age, err = getAge(person.Name)
    if err != nil {
        return err
    }
    person.Gender, err = getGender(person.Name)
    if err != nil {
        return err
    }
    person.Country, err = getCountry(person.Name)
    if err != nil {
        return err
    }

    return nil
}

func getAge(name string) (int, error) {
    resp, err := http.Get("https://api.agify.io/?name=" + name)
    if err != nil {
        return 0, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return 0, errors.New("Age API failed")
    }

    // Parse the JSON response
    var ageInfo struct {
        Age *int `json:"age"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&ageInfo); err != nil {
        return 0, err
    }
    
    if ageInfo.Age == nil {
        return 0, errors.New("Age is nil")
    }

    return *ageInfo.Age, nil
}

func getGender(name string) (string, error) {
    resp, err := http.Get("https://api.genderize.io/?name=" + name)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return "", errors.New("Gender API failed")
    }

    // Parse the JSON response
    var genderInfo struct {
        Gender *string `json:"gender"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&genderInfo); err != nil {
        return "", err
    }

    if genderInfo.Gender == nil {
        return "", errors.New("Gender is nil")
    }

    return *genderInfo.Gender, nil
}

func getCountry(name string) (string, error) {
    resp, err := http.Get("https://api.nationalize.io/?name=" + name)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return "", errors.New("Country API failed")
    }

    // Parse the JSON response
    var countryInfo struct {
        Country []struct {
            CountryID string `json:"country_id"`
        } `json:"country"`
    }

    if err := json.NewDecoder(resp.Body).Decode(&countryInfo); err != nil {
        return "", err
    }

    if len(countryInfo.Country) == 0 {
        return "", errors.New("No Country")
    }

    return countryInfo.Country[0].CountryID, nil
}
