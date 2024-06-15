package controllers

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"net/http"
)

func GetProfile(c *gin.Context) {
	url := "https://vip-reseller.co.id/api/profile"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("key", "eAJioW6S9VG2VIFl7Oj9vYcJvlO89qN10UYvKOuQJcVC2Ksgigaank8Z1Sk3w3BQ")
	_ = writer.WriteField("sign", "00162806611de6dfe265c255b8c2c3e2")
	err := writer.Close()
	if err != nil {
		c.String(http.StatusInternalServerError, "Error creating form data: %s", err)
		return
	}

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error creating request: %s", err)
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error sending request: %s", err)
		return
	}
	defer resp.Body.Close()

	var responseData map[string]interface{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&responseData)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error decoding response: %s", err)
		return
	}

	c.JSON(resp.StatusCode, responseData)
}

func GetGameOrder(c *gin.Context) {
	var data struct {
		Trxid string
	}

	c.Bind(&data)

	url := "https://vip-reseller.co.id/api/game-feature"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("key", "eAJioW6S9VG2VIFl7Oj9vYcJvlO89qN10UYvKOuQJcVC2Ksgigaank8Z1Sk3w3BQ")
	_ = writer.WriteField("sign", "00162806611de6dfe265c255b8c2c3e2")
	_ = writer.WriteField("type", "status")
	_ = writer.WriteField("trxid", data.Trxid)
	_ = writer.WriteField("limit", "")
	err := writer.Close()
	if err != nil {
		c.String(http.StatusInternalServerError, "Error creating form data: %s", err)
		return
	}

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error creating request: %s", err)
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error sending request: %s", err)
		return
	}
	defer resp.Body.Close()

	var responseData map[string]interface{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&responseData)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error decoding response: %s", err)
		return
	}

	c.JSON(resp.StatusCode, responseData)
}

func ListGameHarga(c *gin.Context) {
	var data struct {
		FilterType   string
		FilterValue  string
		FilterStatus string
	}

	c.Bind(&data)

	url := "https://vip-reseller.co.id/api/game-feature"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("key", "eAJioW6S9VG2VIFl7Oj9vYcJvlO89qN10UYvKOuQJcVC2Ksgigaank8Z1Sk3w3BQ")
	_ = writer.WriteField("sign", "00162806611de6dfe265c255b8c2c3e2")
	_ = writer.WriteField("type", "services")
	_ = writer.WriteField("filter_type", data.FilterType)
	_ = writer.WriteField("filter_value", data.FilterValue)
	_ = writer.WriteField("filter_status", data.FilterStatus)
	err := writer.Close()
	if err != nil {
		c.String(http.StatusInternalServerError, "Error creating form data: %s", err)
		return
	}

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error creating request: %s", err)
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error sending request: %s", err)
		return
	}
	defer resp.Body.Close()

	var responseData map[string]interface{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&responseData)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error decoding response: %s", err)
		return
	}

	c.JSON(resp.StatusCode, responseData)
}

func GetNickGame(c *gin.Context) {
	var data struct {
		Code             string
		Target           string
		AdditionalTarget string
	}

	c.Bind(&data)

	url := "https://vip-reseller.co.id/api/game-feature"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("key", "eAJioW6S9VG2VIFl7Oj9vYcJvlO89qN10UYvKOuQJcVC2Ksgigaank8Z1Sk3w3BQ")
	_ = writer.WriteField("sign", "00162806611de6dfe265c255b8c2c3e2")
	_ = writer.WriteField("type", "get-nickname")
	_ = writer.WriteField("code", data.Code)
	_ = writer.WriteField("target", data.Target)
	_ = writer.WriteField("additional_target", data.AdditionalTarget)
	err := writer.Close()
	if err != nil {
		c.String(http.StatusInternalServerError, "Error creating form data: %s", err)
		return
	}

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error creating request: %s", err)
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error sending request: %s", err)
		return
	}
	defer resp.Body.Close()

	var responseData map[string]interface{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&responseData)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error decoding response: %s", err)
		return
	}

	c.JSON(resp.StatusCode, responseData)
}

func TopUpGame(c *gin.Context) {
	var data struct {
		Service  string
		DataNo   string
		DataZone string
	}

	c.Bind(&data)

	url := "https://vip-reseller.co.id/api/game-feature"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("key", "eAJioW6S9VG2VIFl7Oj9vYcJvlO89qN10UYvKOuQJcVC2Ksgigaank8Z1Sk3w3BQ")
	_ = writer.WriteField("sign", "00162806611de6dfe265c255b8c2c3e2")
	_ = writer.WriteField("type", "order")
	_ = writer.WriteField("service", data.Service)
	_ = writer.WriteField("data_no", data.DataNo)
	_ = writer.WriteField("data_zone", data.DataZone)
	err := writer.Close()
	if err != nil {
		c.String(http.StatusInternalServerError, "Error creating form data: %s", err)
		return
	}

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error creating request: %s", err)
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error sending request: %s", err)
		return
	}
	defer resp.Body.Close()

	var responseData map[string]interface{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&responseData)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error decoding response: %s", err)
		return
	}

	c.JSON(resp.StatusCode, responseData)
}

func TopUpPrepaid(c *gin.Context) {
	var data struct {
		Service string
		DataNo  string
	}

	c.Bind(&data)

	url := "https://vip-reseller.co.id/api/prepaid"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("key", "eAJioW6S9VG2VIFl7Oj9vYcJvlO89qN10UYvKOuQJcVC2Ksgigaank8Z1Sk3w3BQ")
	_ = writer.WriteField("sign", "00162806611de6dfe265c255b8c2c3e2")
	_ = writer.WriteField("type", "order")
	_ = writer.WriteField("service", data.Service)
	_ = writer.WriteField("data_no", data.DataNo)
	err := writer.Close()
	if err != nil {
		c.String(http.StatusInternalServerError, "Error creating form data: %s", err)
		return
	}

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error creating request: %s", err)
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error sending request: %s", err)
		return
	}
	defer resp.Body.Close()

	var responseData map[string]interface{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&responseData)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error decoding response: %s", err)
		return
	}

	c.JSON(resp.StatusCode, responseData)
}

func ListPrepaid(c *gin.Context) {
	var data struct {
		FilterType  string
		FilterValue string
	}

	c.Bind(&data)

	url := "https://vip-reseller.co.id/api/prepaid"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("key", "eAJioW6S9VG2VIFl7Oj9vYcJvlO89qN10UYvKOuQJcVC2Ksgigaank8Z1Sk3w3BQ")
	_ = writer.WriteField("sign", "00162806611de6dfe265c255b8c2c3e2")
	_ = writer.WriteField("type", "services")
	_ = writer.WriteField("filter_type", data.FilterType)
	_ = writer.WriteField("filter_value", data.FilterValue)
	err := writer.Close()
	if err != nil {
		c.String(http.StatusInternalServerError, "Error creating form data: %s", err)
		return
	}

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error creating request: %s", err)
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error sending request: %s", err)
		return
	}
	defer resp.Body.Close()

	var responseData map[string]interface{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&responseData)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error decoding response: %s", err)
		return
	}

	c.JSON(resp.StatusCode, responseData)
}

func GetPrepaidOrder(c *gin.Context) {
	var data struct {
		Trxid string
		Limit string
	}

	c.Bind(&data)

	url := "https://vip-reseller.co.id/api/prepaid"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("key", "eAJioW6S9VG2VIFl7Oj9vYcJvlO89qN10UYvKOuQJcVC2Ksgigaank8Z1Sk3w3BQ")
	_ = writer.WriteField("sign", "00162806611de6dfe265c255b8c2c3e2")
	_ = writer.WriteField("type", "status")
	_ = writer.WriteField("trxid", data.Trxid)
	_ = writer.WriteField("limit", data.Limit)
	err := writer.Close()
	if err != nil {
		c.String(http.StatusInternalServerError, "Error creating form data: %s", err)
		return
	}

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error creating request: %s", err)
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error sending request: %s", err)
		return
	}
	defer resp.Body.Close()

	var responseData map[string]interface{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&responseData)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error decoding response: %s", err)
		return
	}

	c.JSON(resp.StatusCode, responseData)
}
