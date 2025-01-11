package control

import (
	"encoding/json"
	"net/http"

	"studentgit.kata.academy/Zhodaran/go-kata/internal/service"
)

type Controller struct {
	geoService service.GeoServicer
}

func NewController(geoService service.GeoServicer) *Controller {
	return &Controller{geoService: geoService}
}

type ErrorResponse struct {
	Message string `json:"message"` // Сообщение об ошибке
	Code    int    `json:"code"`    // Код ошибки
}

// @Summary Get Geo Coordinates by Address
// @Description This endpoint allows you to get geo coordinates by address.
// @Tags geo
// @Accept json
// @Produce json
// @Param address body service.RequestAddressSearch true "Address search query"
// @Success 200 {object} service.ResponseAddress "Успешное выполнение"
// @Failure 400 {object} ErrorResponse "Ошибка запроса"
// @Failure 500 {object} ErrorResponse "Ошибка подключения к серверу"
// @Security BearerAuth
// @Router /api/address/search [post]
func (c *Controller) GetGeoCoordinatesAddress(w http.ResponseWriter, r *http.Request) {
	var req service.RequestAddressSearch
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	geo, err := c.geoService.GetGeoCoordinatesAddress(req.Query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(geo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// @Summary Get Geo Coordinates by Latitude and Longitude
// @Description This endpoint allows you to get geo coordinates by latitude and longitude.
// @Tags geo
// @Accept json
// @Produce json
// @Param body body service.GeocodeRequest true "Geographic coordinates"
// @Success 200 {object} service.ResponseAddress "Успешное выполнение"
// @Failure 400 {object} ErrorResponse "Ошибка запроса"
// @Failure 500 {object} ErrorResponse "Ошибка подключения к серверу"
// @Security BearerAuth
// @Router /api/address/geocode [post]
func (c *Controller) GetGeoCoordinatesGeocode(w http.ResponseWriter, r *http.Request) {
	var req service.GeocodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	geo, err := c.geoService.GetGeoCoordinatesGeocode(req.Lat, req.Lng)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(geo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
