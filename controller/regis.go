package controller

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

type RequestBody struct {
	BankCd          string `json:"bankCd"`
	TimeStamp       string `json:"timeStamp"`
	IMid            string `json:"iMid"`
	PayMethod       string `json:"payMethod"`
	Currency        string `json:"currency"`
	Amt             string `json:"amt"`
	ReferenceNo     string `json:"referenceNo"`
	GoodsNm         string `json:"goodsNm"`
	DbProcessUrl    string `json:"dbProcessUrl"`
	Description     string `json:"description"`
	MerchantToken   string `json:"merchantToken"`
	ReqDt           string `json:"reqDt"`
	ReqTm           string `json:"reqTm"`
	ReqDomain       string `json:"reqDomain"`
	ReqServerIP     string `json:"reqServerIP"`
	ReqClientVer    string `json:"reqClientVer"`
	UserIP          string `json:"userIP"`
	UserSessionID   string `json:"userSessionID"`
	UserAgent       string `json:"userAgent"`
	UserLanguage    string `json:"userLanguage"`
	CartData        string `json:"cartData"`
	InstmntType     string `json:"instmntType"`
	InstmntMon      string `json:"instmntMon"`
	RecurrOpt       string `json:"recurrOpt"`
	VacctValidDt    string `json:"vacctValidDt"`
	VacctValidTm    string `json:"vacctValidTm"`
	MerFixAcctId    string `json:"merFixAcctId"`
	BillingNm       string `json:"billingNm"`
	BillingPhone    string `json:"billingPhone"`
	BillingEmail    string `json:"billingEmail"`
	BillingAddr     string `json:"billingAddr"`
	BillingCity     string `json:"billingCity"`
	BillingState    string `json:"billingState"`
	BillingPostCd   string `json:"billingPostCd"`
	BillingCountry  string `json:"billingCountry"`
	DeliveryNm      string `json:"deliveryNm"`
	DeliveryPhone   string `json:"deliveryPhone"`
	DeliveryAddr    string `json:"deliveryAddr"`
	DeliveryCity    string `json:"deliveryCity"`
	DeliveryState   string `json:"deliveryState"`
	DeliveryPostCd  string `json:"deliveryPostCd"`
	DeliveryCountry string `json:"deliveryCountry"`
}

type RegistResponse struct {
	ResultCd     string `json:"resultCd"`
	ResultMsg    string `json:"resultMsg"`
	TXid         string `json:"tXid"`
	ReferenceNo  string `json:"referenceNo"`
	PayMethod    string `json:"payMethod"`
	Amt          string `json:"amt"`
	TransDt      string `json:"transDt"`
	TransTm      string `json:"transTm"`
	Description  string `json:"description"`
	BankCd       string `json:"bankCd"`
	VacctNo      string `json:"vacctNo"`
	MitraCd      string `json:"mitraCd"`
	PayNo        string `json:"payNo"`
	Currency     string `json:"currency"`
	GoodsNm      string `json:"goodsNm"`
	BillingNm    string `json:"billingNm"`
	VacctValidDt string `json:"vacctValidDt"`
	VacctValidTm string `json:"vacctValidTm"`
	PayValidDt   string `json:"payValidDt"`
	PayValidTm   string `json:"payValidTm"`
	RequestURL   string `json:"requestURL"`
	PaymentExpDt string `json:"paymentExpDt"`
	PaymentExpTm string `json:"paymentExpTm"`
	QrContent    string `json:"qrContent"`
	QrUrl        string `json:"qrUrl"`
}

func generateMerchantToken(timeStamp, iMid, referenceNo, amt, merchantKey string) string {
	data := timeStamp + iMid + referenceNo + amt + merchantKey
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

func Regis(c echo.Context) error {

	const url = "https://www.nicepay.co.id/nicepay/direct/v2/registration"
	imid := "TNICEEW051"
	merchantKey := "33F49GnCMS1mFYlGXisbUDzVf2ATWCl9k3R++d5hDd3Frmuos/XLx8XhXpe+LDYAbpGKZYSwtlyyLOtS/8aD7A=="

	requestBody, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Failed to read request body")
		return err
	}

	var requestBody1 RequestBody
	err = json.Unmarshal(requestBody, &requestBody1)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Failed to decode request body")
		return err
	}

	currentTime := time.Now()
	timeStamp := currentTime.Format("20060102150405")
	requestBody1.IMid = imid
	requestBody1.TimeStamp = timeStamp
	requestBody1.ReferenceNo = "ord" + timeStamp
	requestBody1.MerchantToken = generateMerchantToken(timeStamp, imid, requestBody1.ReferenceNo, requestBody1.Amt, merchantKey)

	requestBytes, err := json.Marshal(requestBody1)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Failed to encode request body")
		return err
	}

	log.WithFields(log.Fields{
		"request_body": string(requestBytes),
	}).Info("API request")

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestBytes))
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Failed to create request")
		return err
	}

	// Set header raw
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Failed to send request")
		return err
	}
	defer response.Body.Close()

	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Failed to read response body")
		return err
	}
	log.WithFields(log.Fields{
		"response_body": string(responseBytes),
	}).Info("API response")

	var responseRegist RegistResponse
	if err := json.Unmarshal(responseBytes, &responseRegist); err != nil { // Parse []byte to go struct pointer
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Can not unmarshal JSON")
		return err
	}

	return c.JSON(http.StatusOK, responseRegist)
}
