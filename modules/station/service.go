package station

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/arjunstein/mrt-schedules/common/client"
)

type Service interface {
	GetAllStations() (response []StationResponse, err error)
	CheckSchedulesByStation(id string) (response []ScheduleResponse, err error)
}

type service struct {
	client *http.Client
}

func NewService() Service {
	return &service{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *service) GetAllStations() (response []StationResponse, err error) {
	url := "https://www.jakartamrt.co.id/id/val/stasiuns"

	byteResponse, err := client.DoRequest(s.client, url)
	if err != nil {
		return
	}

	var stations []Station
	err = json.Unmarshal(byteResponse, &stations)

	for _, item := range stations {
		response = append(response, StationResponse{
			Id:   item.Id,
			Name: item.Name,
		})
	}

	return
}

func (s *service) CheckSchedulesByStation(id string) (response []ScheduleResponse, err error) {
	url := "https://www.jakartamrt.co.id/id/val/stasiuns"

	byteResponse, err := client.DoRequest(s.client, url)
	if err != nil {
		return
	}

	var schedule []Schedule
	err = json.Unmarshal(byteResponse, &schedule)

	if err != nil {
		return
	}

	// schedule selected by station id
	var scheduleSelected Schedule
	for _, item := range schedule {
		if item.StationId == id {
			scheduleSelected = item
			break
		}
	}
	if scheduleSelected.StationId == "" {
		err = errors.New("Station not found")
		return
	}

	response, err = ConvertDataToResponse(scheduleSelected)
	if err != nil {
		return
	}

	return
}

func ConvertDataToResponse(schedule Schedule) (response []ScheduleResponse, err error) {
	var (
		LebakBulusTripName = "Stasiun Lebak Bulus Grab"
		BundaranHiTripName = "Stasiun Bundaran HI Bank DKI"
	)

	scheduleLebakBulus := schedule.ScheduleLebakBulus
	scheduleBundaranHi := schedule.ScheduleBundaranHi

	scheduleLebakBulusParsed, err := ConvertScheduleToTimeFormat(scheduleLebakBulus)
	if err != nil {
		return
	}
	scheduleBundaranHiParsed, err := ConvertScheduleToTimeFormat(scheduleBundaranHi)
	if err != nil {
		return
	}

	// convert to response
	for _, item := range scheduleLebakBulusParsed {
		if item.Format("15:04") > time.Now().Format("15:04") {
			response = append(response, ScheduleResponse{
				StationName: LebakBulusTripName,
				Time:        item.Format("15:04"),
			})
		}
	}

	for _, item := range scheduleBundaranHiParsed {
		if item.Format("15:04") > time.Now().Format("15:04") {
			response = append(response, ScheduleResponse{
				StationName: BundaranHiTripName,
				Time:        item.Format("15:04"),
			})
		}
	}
	return
}

func ConvertScheduleToTimeFormat(schedule string) (response []time.Time, err error) {
	var (
		parsedTime time.Time
		schedules = strings.Split(schedule, ",") // split schedule by comma
	)

	for _, item := range schedules {
		trimmedTime := strings.TrimSpace(item) // remove space
		if trimmedTime == "" {
			continue
		}

		parsedTime, err = time.Parse("15:04", trimmedTime)
		if err != nil {
			err = errors.New("Invalid time format " + trimmedTime)
			return
		}
		response = append(response, parsedTime)
	}
	return
}