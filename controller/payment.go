package controller

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func Payment(c echo.Context) error {
	currentTime := time.Now()
	timeStamp := currentTime.Format("20060102150405")

	imid := "IONPAYTEST"
	merchantKey := "33F49GnCMS1mFYlGXisbUDzVf2ATWCl9k3R++d5hDd3Frmuos/XLx8XhXpe+LDYAbpGKZYSwtlyyLOtS/8aD7A=="

	refno := "ord" + timeStamp

	baseURL := "https://dev.nicepay.co.id/nicepay/direct/v2/payment"
	tXid := c.QueryParam("tXid")
	merchantToken := generateMerchantToken(timeStamp, imid, refno, c.QueryParam("amt"), merchantKey)
	cardNo := c.QueryParam("cardNo")
	cardExpYymm := c.QueryParam("cardExpYymm")
	cardCvv := c.QueryParam("cardCvv")
	cardHolderNm := c.QueryParam("cardHolderNm")
	callBackUrl := c.QueryParam("callBackUrl")
	amt := c.QueryParam("amt")

	url := fmt.Sprintf("%s?timeStamp=%s&tXid=%s&merchantToken=%s&&amt=%s&cardNo=%s&cardExpYymm=%s&cardCvv=%s&cardHolderNm=%s&callBackUrl=%s",
		baseURL, timeStamp, tXid, merchantToken, amt, cardNo, cardExpYymm, cardCvv, cardHolderNm, callBackUrl)

	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	return c.String(http.StatusOK, string(body))
}
