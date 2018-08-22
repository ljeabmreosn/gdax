package gdax

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	gock "gopkg.in/h2non/gock.v1"
)

const (
	accountsJson = `
		[
		    {
		        "id": "71452118-efc7-4cc4-8780-a5e22d4baa53",
		        "currency": "BTC",
		        "balance": "0.0000000000000000",
		        "available": "0.0000000000000000",
		        "hold": "0.0000000000000000",
		        "profile_id": "75da88c5-05bf-4f54-bc85-5c775bd68254"
		    },
		    {
		        "id": "e316cb9a-0808-4fd7-8914-97829c1925de",
		        "currency": "USD",
		        "balance": "80.2301373066930000",
		        "available": "79.2266348066930000",
		        "hold": "1.0035025000000000",
		        "profile_id": "75da88c5-05bf-4f54-bc85-5c775bd68254"
		    }
		]
	`
	accountJson = `
		{
		    "id": "6cf2b1ba-3705-40e6-a41e-69be033514f7",
		    "balance": "1.100",
		    "holds": "0.100",
		    "available": "1.00",
		    "currency": "USD"
		}
	`
	accountHistoryJson1 = `
		[
		    {
		        "id": 100,
		        "created_at": "2014-11-07T08:19:27.028459Z",
		        "amount": "0.001",
		        "balance": "239.669",
		        "type": "fee",
		        "details": {
		            "order_id": "d50ec984-77a8-460a-b958-66f114b0de9b",
		            "trade_id": "74",
		            "product_id": "BTC-USD"
		        }
		    }
		]
	`
	accountHistoryJson2 = `
		[
		    {
		        "id": 100,
		        "created_at": "2014-11-07T08:19:29.028459Z",
		        "amount": "0.001",
		        "balance": "170.322",
		        "type": "fee",
		        "details": {
		            "order_id": "62087add-1eea-47fc-b79f-8cde52b458d6",
		            "trade_id": "75",
		            "product_id": "BTC-USD"
		        }
		    }
		]
	`
	accountHoldJson1 = `
		[
		    {
		        "id": "82dcd140-c3c7-4507-8de4-2c529cd1a28f",
		        "account_id": "e0b3f39a-183d-453e-b754-0c13e5bab0b3",
		        "created_at": "2014-11-06T10:34:47.123456Z",
		        "updated_at": "2014-11-06T10:40:47.123456Z",
		        "amount": "4.23",
		        "type": "order",
		        "ref": "0a205de4-dd35-4370-a285-fe8fc375a273"
		    },
		    {
		        "id": "1fa18826-8f96-4640-b73a-752d85c69326",
		        "account_id": "e0b3f39a-183d-453e-b754-0c13e5bab0b3",
		        "created_at": "2014-11-06T10:34:47.123456Z",
		        "updated_at": "2014-11-06T10:40:47.123456Z",
		        "amount": "5.25",
		        "type": "order",
		        "ref": "ba2a968c-17f9-4fcb-90d7-eb6f2ac49538"
		    }
		]
	`
	accountHoldJson2 = `
		[
		    {
		        "id": "e6b60c60-42ed-4329-a311-694d6c897d9b",
		        "account_id": "e0b3f39a-183d-453e-b754-0c13e5bab0b3",
		        "created_at": "2014-11-06T10:34:47.123456Z",
		        "updated_at": "2014-11-06T10:40:47.123456Z",
		        "amount": "6.34",
		        "type": "order",
		        "ref": "ba2a968c-17f9-4fcb-90d7-eb6f2ac49538"
		    }
		]
	`
)

func TestGetAccounts(t *testing.T) {
	defer gock.Off()
	assert := assert.New(t)

	accessInfo, err := RetrieveAccessInfoFromEnvironmentVariables()
	assert.NoError(err)

	gock.New(EndPoint).
		Get("/accounts").
		Reply(http.StatusOK).
		BodyString(accountsJson)

	var ids = [...]string{"71452118-efc7-4cc4-8780-a5e22d4baa53", "e316cb9a-0808-4fd7-8914-97829c1925de"}

	for idx, accounts := 0, accessInfo.GetAccounts(); accounts.HasNext(); idx++ {
		account, err := accounts.Next()
		assert.NoError(err)

		parsedId, err := uuid.Parse(ids[idx])
		assert.NoError(err)

		assert.Equal(*account.Id, parsedId)
	}
}

func TestGetAccount(t *testing.T) {
	defer gock.Off()
	assert := assert.New(t)

	accessInfo, err := RetrieveAccessInfoFromEnvironmentVariables()
	assert.NoError(err)

	const id = "6cf2b1ba-3705-40e6-a41e-69be033514f7"
	gock.New(EndPoint).
		Get(fmt.Sprintf("/accounts/%s", id)).
		Reply(http.StatusOK).
		BodyString(accountJson)

	parsedId, err := uuid.Parse(id)
	assert.NoError(err)

	account, err := accessInfo.GetAccount(&parsedId)
	assert.NoError(err)

	assert.Equal(*account.Id, parsedId)
}

func TestGetAccountHistory(t *testing.T) {
	defer gock.Off()
	assert := assert.New(t)

	accessInfo, err := RetrieveAccessInfoFromEnvironmentVariables()
	assert.NoError(err)

	var cursors = [...]int{10, 20}
	const accountId = "6cf2b1ba-3705-40e6-a41e-69be033514f7"
	gock.New(EndPoint).
		Get(fmt.Sprintf("/accounts/%s/ledger", accountId)).
		Reply(http.StatusOK).
		BodyString(accountHistoryJson1).
		SetHeader("CB-AFTER", strconv.Itoa(cursors[0]))
	gock.New(EndPoint).
		Get(fmt.Sprintf("/accounts/%s/ledger", accountId)).
		MatchParam("after", strconv.Itoa(cursors[0])).
		Reply(http.StatusOK).
		BodyString(accountHistoryJson2).
		SetHeader("CB-AFTER", strconv.Itoa(cursors[1]))
	gock.New(EndPoint).
		Get(fmt.Sprintf("/accounts/%s/ledger", accountId)).
		MatchParam("after", strconv.Itoa(cursors[1])).
		Reply(http.StatusOK).
		BodyString("[]")

	var orderIds = [...]string{"d50ec984-77a8-460a-b958-66f114b0de9b", "62087add-1eea-47fc-b79f-8cde52b458d6"}

	parsedAccountId, err := uuid.Parse(accountId)
	assert.NoError(err)

	for idx, accountHistories := 0, accessInfo.GetAccountHistory(&parsedAccountId); accountHistories.HasNext(); idx++ {
		accountHistory, err := accountHistories.Next()
		assert.NoError(err)
		t.Log(accountHistory)

		parsedId, err := uuid.Parse(orderIds[idx])
		assert.NoError(err)

		assert.Equal(*accountHistory.Details.OrderId, parsedId)
	}
}

func TestGetAccountHolds(t *testing.T) {
	defer gock.Off()
	assert := assert.New(t)

	accessInfo, err := RetrieveAccessInfoFromEnvironmentVariables()
	assert.NoError(err)

	var cursors = [...]int{10, 20}
	const accountId = "e0b3f39a-183d-453e-b754-0c13e5bab0b3"
	gock.New(EndPoint).
		Get(fmt.Sprintf("/accounts/%s/holds", accountId)).
		Reply(http.StatusOK).
		BodyString(accountHoldJson1).
		SetHeader("CB-AFTER", strconv.Itoa(cursors[0]))
	gock.New(EndPoint).
		Get(fmt.Sprintf("/accounts/%s/holds", accountId)).
		MatchParam("after", strconv.Itoa(cursors[0])).
		Reply(http.StatusOK).
		BodyString(accountHoldJson2).
		SetHeader("CB-AFTER", strconv.Itoa(cursors[1]))
	gock.New(EndPoint).
		Get(fmt.Sprintf("/accounts/%s/holds", accountId)).
		MatchParam("after", strconv.Itoa(cursors[1])).
		Reply(http.StatusOK).
		BodyString("[]")

	var ids = [...]string{"82dcd140-c3c7-4507-8de4-2c529cd1a28f", "1fa18826-8f96-4640-b73a-752d85c69326", "e6b60c60-42ed-4329-a311-694d6c897d9b"}
	var amounts = [...]float64{4.23, 5.25, 6.34}

	parsedAccountId, err := uuid.Parse(accountId)
	assert.NoError(err)

	for idx, accountHolds := 0, accessInfo.GetAccountHolds(&parsedAccountId); accountHolds.HasNext(); idx++ {
		accountHold, err := accountHolds.Next()
		assert.NoError(err)
		t.Log(accountHold)

		parsedId, err := uuid.Parse(ids[idx])
		assert.NoError(err)

		assert.Equal(*accountHold.Id, parsedId)
		assert.Equal(accountHold.Amount, amounts[idx])
	}
}