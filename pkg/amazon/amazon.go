package amazon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Review struct {
	ProductId string `json:"productId"`
	Rating    string `json:"rating"`
	Comment   string `json:"comment"`
}

func GenerateReviews(reviews []string, productId string, rating string) []*Review {
	var reviewsList []*Review
	for _, review := range reviews {
		reviewsList = append(reviewsList, &Review{
			ProductId: productId,
			Rating:    rating,
			Comment:   review,
		})
	}
	return reviewsList
}

func CreateReviews(host string, jwt string, req2 Review) error {
	url := host + "/api/products/" + req2.ProductId + "/reviews"

	jsonData, err := json.Marshal(req2)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Cookie", fmt.Sprintf("jwt=%s", jwt))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	fmt.Println("response Status:", resp.Status)
	return nil
}
