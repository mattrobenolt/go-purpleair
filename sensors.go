package purpleair

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func (c *Client) SensorData(ctx context.Context, sensorIndex int, opts *SensorDataRequestOptions) (*SensorDataResponse, error) {
	b, resp, err := c.get(ctx, fmt.Sprintf("/sensors/%d?%s", sensorIndex, opts.Encode()), 200)
	if err != nil {
		return nil, err
	}

	r := &SensorDataResponse{
		body:         b,
		httpResponse: resp,
	}
	err = decodeBody(b, r)
	return r, err
}

func (c *Client) SensorHistory(ctx context.Context, sensorIndex int, opts *SensorHistoryRequestOptions) (*SensorDataResponse, error) {
	b, resp, err := c.get(ctx, fmt.Sprintf("/sensors/%d/history?%s", sensorIndex, opts.Encode()), 200)
	if err != nil {
		return nil, err
	}

	r := &SensorDataResponse{
		body:         b,
		httpResponse: resp,
	}
	err = decodeBody(b, r)
	return r, err
}

type SensorDataRequestOptions struct {
	Fields   []string
	ReadKeys []string
}

func (o *SensorDataRequestOptions) Encode() string {
	if o == nil {
		return ""
	}
	if len(o.Fields) == 0 && len(o.ReadKeys) == 0 {
		return ""
	}
	v := url.Values{}
	if len(o.Fields) > 0 {
		v.Set("fields", strings.Join(o.Fields, ","))
	}
	if len(o.ReadKeys) > 0 {
		v.Set("read_key", strings.Join(o.ReadKeys, ","))
	}
	return v.Encode()
}

type SensorHistoryRequestOptions struct {
	Fields   []string
	ReadKeys []string

	StartTimestamp *uint64
	EndTimestamp   *uint64
	Average        *uint16
}

func (o *SensorHistoryRequestOptions) Encode() string {
	if o == nil {
		return ""
	}
	v := url.Values{}
	v.Set("fields", strings.Join(o.Fields, ","))

	if len(o.ReadKeys) > 0 {
		v.Set("read_key", strings.Join(o.ReadKeys, ","))
	}
	if o.StartTimestamp != nil {
		v.Set("start_timestamp", itoa(*o.StartTimestamp))
	}
	if o.EndTimestamp != nil {
		v.Set("end_timestamp", itoa(*o.EndTimestamp))
	}
	if o.Average != nil {
		v.Set("average", itoa(*o.Average))
	}
	return v.Encode()
}

type SensorDataResponse struct {
	ApiVersion    string  `json:"api_version"`
	TimeStamp     int64   `json:"time_stamp"`
	DataTimeStamp int64   `json:"data_time_stamp"`
	Sensor        *Sensor `json:"sensor"`

	body         []byte
	httpResponse *http.Response
}

type Sensor struct {
	SensorIndex     int64   `json:"sensor_index"`
	LastModified    *int64  `json:"last_modified"`
	DateCreated     *int64  `json:"date_created"`
	LastSeen        *int64  `json:"last_seen"`
	Private         *uint8  `json:"private"`
	IsOwner         *uint8  `json:"is_owner"`
	Name            *string `json:"name"`
	LocationType    *uint8  `json:"location_type"`
	Model           *string `json:"model"`
	Hardware        *string `json:"hardware"`
	LedBrightness   *uint8  `json:"led_brightness"`
	FirmwareVersion *string `json:"firmware_version"`
	Rssi            *int64  `json:"rssi"`
	Uptime          *uint64 `json:"uptime"`
	/* pa_latency, memory, position_rating, latitude, longitude, altitude, channel_state, channel_flags, channel_flags_manual, channel_flags_auto, confidence, confidence_auto, confidence_manual */

	Humidity     *uint8   `json:"humidity"`
	HumidityA    *uint8   `json:"humidity_a"`
	HumidityB    *uint8   `json:"humidity_b"`
	Temperature  *uint8   `json:"temperature"`
	TemperatureA *uint8   `json:"temperature_a"`
	TemperatureB *uint8   `json:"temperature_b"`
	Pressure     *float64 `json:"pressure"`
	PressureA    *float64 `json:"pressure_a"`
	PressureB    *float64 `json:"pressure_b"`
	Voc          *float64 `json:"voc"`
	VocA         *float64 `json:"voc_a"`
	VocB         *float64 `json:"voc_b"`
	AnalogInput  *float64 `json:"analog_input"`
	Pm10         *float64 `json:"pm1.0"`
	Pm10A        *float64 `json:"pm1.0_a"`
	Pm10B        *float64 `json:"pm1.0_b"`
	Pm25         *float64 `json:"pm2.5"`
	Pm25A        *float64 `json:"pm2.5_a"`
	Pm25B        *float64 `json:"pm2.5_b"`
	Pm25Alt      *float64 `json:"pm2.5_alt"`
	Pm25AltA     *float64 `json:"pm2.5_alt_a"`
	Pm25AltB     *float64 `json:"pm2.5_alt_b"`
	Pm100        *float64 `json:"pm10.0"`
	Pm100A       *float64 `json:"pm10.0_a"`
	Pm100B       *float64 `json:"pm10.0_b"`

	PrimaryIdA    *int64  `json:"primary_id_a"`
	PrimaryKeyA   *string `json:"primary_key_a"`
	PrimaryIdB    *int64  `json:"primary_id_b"`
	PrimaryKeyB   *string `json:"primary_key_b"`
	SecondaryIdA  *int64  `json:"secondary_id_a"`
	SecondaryKeyA *string `json:"secondary_key_a"`
	SecondaryIdB  *int64  `json:"secondary_id_b"`
	SecondaryKeyB *string `json:"secondary_key_b"`

	Stats  *SensorStats `json:"stats"`
	StatsA *SensorStats `json:"stats_a"`
	StatsB *SensorStats `json:"stats_b"`
}

type SensorStats struct {
	TimeStamp     int64   `json:"time_stamp"`
	Pm25          float64 `json:"pm2.5"`
	Pm25_10Minute float64 `json:"pm2.5_10minute"`
	Pm25_30Minute float64 `json:"pm2.5_30minute"`
	Pm25_60Minute float64 `json:"pm2.5_60minute"`
	Pm25_6Hour    float64 `json:"pm2.5_6hour"`
	Pm25_24Hour   float64 `json:"pm2.5_24hour"`
	Pm25_1Week    float64 `json:"pm2.5_1week"`
}

func (r *SensorDataResponse) Body() []byte                 { return r.body }
func (r *SensorDataResponse) HttpResponse() *http.Response { return r.httpResponse }
