package sorting

import (
	"encoding/json"
	"errors"
	"net/http"
)

type SortingApi struct {
	service SortingService
}

type SortingRequest struct {
	Numbers []int `json:"numbers"`
	Sorted  []int `json:"sorted"`
}

func NewSortingApi(s SortingService) *SortingApi {
	return &SortingApi{
		service: s,
	}
}

func (api *SortingApi) BubbleSortHandler(w http.ResponseWriter, r *http.Request) (int, error) {
	if r.Method != http.MethodPost {
		return http.StatusMethodNotAllowed, errors.New("method not allowed")
	}
	var req SortingRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return http.StatusBadRequest, errors.New("An error occured while decoding the json \n")
	}

	req.Sorted = make([]int, len(req.Numbers))
	copy(req.Sorted, req.Numbers)

	api.service.BubbleSort(req.Sorted)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(req)
	if err != nil {
		return http.StatusInternalServerError, errors.New("failed to write http response")
	}

	return http.StatusOK, nil
}
