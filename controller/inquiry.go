package controller

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

type RequestBodyInquiry struct {
	TimeStamp     string `json:"timeStamp"`
	MerchantToken string `json:"merchantToken"`
	ReferenceNo   string `json:"referenceNo"`
	TXid          string `json:"tXid"`
	Amt           string `json:"amt"`
	IMid          string `json:"iMid"`
}

type ResponseInquiry struct {
	TXID                    string `json:"tXid"`
	IMID                    string `json:"iMid"`
	Currency                string `json:"currency"`
	Amount                  string `json:"amt"`
	InstallmentMon          string `json:"instmntMon"`
	InstallmentType         string `json:"instmntType"`
	ReferenceNo             string `json:"referenceNo"`
	GoodsName               string `json:"goodsNm"`
	PayMethod               string `json:"payMethod"`
	BillingName             string `json:"billingNm"`
	RequestDate             string `json:"reqDt"`
	RequestTime             string `json:"reqTm"`
	Status                  string `json:"status"`
	ResultCode              string `json:"resultCd"`
	ResultMessage           string `json:"resultMsg"`
	CardNo                  string `json:"cardNo"`
	PreauthToken            string `json:"preauthToken"`
	AcquiringBankCode       string `json:"acquBankCd"`
	IssuingBankCode         string `json:"issuBankCd"`
	VirtualAccountValidDate string `json:"vacctValidDt"`
	VirtualAccountValidTime string `json:"vacctValidTm"`
	VirtualAccountNo        string `json:"vacctNo"`
	BankCode                string `json:"bankCd"`
	PaymentNo               string `json:"payNo"`
	MitraCode               string `json:"mitraCd"`
	ReceiptCode             string `json:"receiptCode"`
	CancelAmount            string `json:"cancelAmt"`
	TransactionDate         string `json:"transDt"`
	TransactionTime         string `json:"transTm"`
	RecurringToken          string `json:"recurringToken"`
	CreditCardTransType     string `json:"ccTransType"`
	PaymentValidDate        string `json:"payValidDt"`
	PaymentValidTime        string `json:"payValidTm"`
	MerchantReferenceNo     string `json:"mRefNo"`
	AcquiringStatus         string `json:"acquStatus"`
	CardExpiryYYMM          string `json:"cardExpYymm"`
	AcquiringBankName       string `json:"acquBankNm"`
	IssuingBankName         string `json:"issuBankNm"`
	DepositDate             string `json:"depositDt"`
	DepositTime             string `json:"depositTm"`
	PaymentExpireDate       string `json:"paymentExpDt"`
	PaymentExpireTime       string `json:"paymentExpTm"`
	PaymentTransactionSn    string `json:"paymentTrxSn"`
	CancelTransactionSn     string `json:"cancelTrxSn"`
	UserID                  string `json:"userId"`
	ShopID                  string `json:"shopId"`
}

func Inquiry(c echo.Context) error {

	const url = "https://dev.nicepay.co.id/nicepay/direct/v2/inquiry"
	imid := "IONPAYTEST"
	merchantKey := "33F49GnCMS1mFYlGXisbUDzVf2ATWCl9k3R++d5hDd3Frmuos/XLx8XhXpe+LDYAbpGKZYSwtlyyLOtS/8aD7A=="

	requestBody, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Failed to read request body")
		return err
	}

	var requestBody1 RequestBodyInquiry
	err = json.Unmarshal(requestBody, &requestBody1)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Failed to decode request body")
		return err
	}

	currentTime := time.Now()
	timeStamp := currentTime.Format("20060102150405")
	refno := "ord" + timeStamp

	requestBody1.IMid = imid
	requestBody1.TimeStamp = timeStamp
	requestBody1.ReferenceNo = refno
	requestBody1.MerchantToken = generateMerchantToken(timeStamp, imid, refno, requestBody1.Amt, merchantKey)

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

	var responseRegist ResponseInquiry
	if err := json.Unmarshal(responseBytes, &responseRegist); err != nil { // Parse []byte to go struct pointer
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Can not unmarshal JSON")
		return err
	}

	return c.JSON(http.StatusOK, responseRegist)
}
