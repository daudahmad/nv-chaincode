# Nostro Vostro blockchain demo

* A nostro is our account of our money, held by the other bank
* A vostro is our account of other bank money, held by us

For more information on nostro vostro check [this wikipedia link!](https://en.wikipedia.org/wiki/Nostro_and_vostro_accounts)

The `Init` function creates 3 financial institutions in the ledger/world state:
- Financial institution 1
  - Owner: BANKA (American bank: USD), holding the money of the following banks:
	  1. BANKB - 250000
		2. BANKC - 360000
- Financial institution 2
	- Owner: BANKB (Australian bank: AUD), holding the money of the following banks:
	  1. BANKA - 250000 * 1.34 (USD to AUD exchange rate)
		2. BANKC - 120000
- Financial institution 3 
  - Owner: BANKC (European bank: EUR), holding the money of the following banks:
    1. BANKA - 360000 * 0.9 (USR to EUR exchange rate)
	  2. BANKB - 120000 * 0.67 (AUD to EUR exchange rate)

An `Account` in chaincode is represented by the following type, sample values are shown in comments:

```go
type Account struct {
	Holder      string  `json:"holder"`         //BANKA
	Currency    string  `json:"currency"`       //USD
	CashBalance float64 `json:"cashBalance"`    //250000
}
```

An `FinancialInstitution` in chaincode is represented by the following type, `Accounts []Account` is the array of accounts of other banks that this financial institution is holding, also called **vostro** in financial lingo:

```go
type FinancialInst struct {
	Owner    string    `json:"owner"`           //BANKA
	Accounts []Account `json:"accounts"`
}
```

### getFinancialInstitutionDetails
To view accounts information (bank name, currency and cash balance) held by a particular financial institution i.e. the **vostro** accounts we will use the `getFinancialInstitutionDetails` function of the chaincode. To view vostro details for BANKA we can use the following json request. Replace chaincode name with actual chaindcodeID, secureContext with actual user.

**Request**
```json
{
  "jsonrpc": "2.0",
  "method": "query",
  "params": {
    "type": 1,
    "chaincodeID": {
      "name": "mycc"
    },
    "ctorMsg": {
      "function": "getFinancialInstitutionDetails",
      "args": [
        "BANKA"
      ]
    },
    "secureContext": "jim"
  },
  "id": 1
}
```

**Response**
```json
{
  "jsonrpc": "2.0",
  "result": {
    "status": "OK",
    "message": 
      "{
        \"owner\":\"BANKA\",
        \"accounts\":[
          {
            \"holder\":\"BANKB\",
            \"currency\":\"USD\",
            \"cashBalance\":250000
          },
          {
            \"holder\":\"BANKC\",
            \"currency\":\"USD\",
            \"cashBalance\":360000
          }
        ]
       }"
  },
  "id": 1
}
```

### getNostroVostroAccounts
To view accounts information (bank name, currency and cash balance) held by a particular financial institution i.e. the **vostro** accounts as well as the accounts of that bank that are being held by other financial institutions (banks) i.e. the **nostro** accounts, we will use the `getNostroVostroAccounts` function of the chaincode. To view *nostro* and *vostro* details for BANKA we can use the following json request. Replace chaincode name with actual chaindcodeID, secureContext with actual user.

**Request**
```json
{
  "jsonrpc": "2.0",
  "method": "query",
  "params": {
    "type": 1,
    "chaincodeID": {
      "name": "mycc"
    },
    "ctorMsg": {
      "function": "getNostroVostroAccounts",
      "args": [
        "BANKA"
      ]
    },
    "secureContext": "jim"
  },
  "id": 1
}
```

**Response**
```json
{
  "jsonrpc": "2.0",
  "result": {
    "status": "OK",
    "message": 
      "{
          \"owner\":\"BANKA\",
          \"nostro\":[
            {
              \"owner\":\"BANKB\",
              \"accounts\":[
                {
                  \"holder\":\"BANKA\",
                  \"currency\":\"AUD\",
                  \"cashBalance\":335000
                 }
               ]
            },
            {
              \"owner\":\"BANKC\",
              \"accounts\":[
                {
                  \"holder\":\"BANKA\",
                  \"currency\":\"EUR\",
                  \"cashBalance\":324000}
                  ]
                }
          ],
          \"vostro\":[
            {
              \"owner\":\"BANKA\",
              \"accounts\":[
                {
                  \"holder\":\"BANKB\",
                  \"currency\":\"USD\",
                  \"cashBalance\":250000
                },
                {
                  \"holder\":\"BANKC\",
                  \"currency\":\"USD\",
                  \"cashBalance\":360000
                }
              ]
            }
           ]
        }"
  },
  "id": 1
}
```

### submitTransaction
Send money from one country to another by calling this function. Let's say Alice wants to send Bob USD 1000, Alice has a bank account in BANKA (American bank) and Bob has a bank account in BANKB (Australian bank). We can do this by calling `submitTransaction` with the arguments given below. This will result in the following:
* Credit the vostro account (BANKB's money held by BANKA)
* Debit the nostro account (BANKA's money held by BANKB)
* Add a new transaction in the `AllTransactions` struct 

```go
type AllTransactions struct {
	Transactions []Transaction `json:"transactions"`
}
```

A Transaction has the following structure:

```go
type Transaction struct {
	RefNumber  string  `json:"refNumber"`
	OpCode     string  `json:"opCode"`
	VDate      string  `json:"vDate"`
	Currency   string  `json:"currency"`
	Amount     float64 `json:"amount"`
	Sender     string  `json:"sender"`
	Receiver   string  `json:"receiver"`
	OrdCust    string  `json:"ordcust"`
	BenefCust  string  `json:"benefcust"`
	DetCharges string  `json:"detcharges"`
	StatusCode int     `json:"statusCode"`
	StatusMsg  string  `json:"statusMsg"`
}
```

**Request**
```json
{
  "jsonrpc": "2.0",
  "method": "invoke",
  "params": {
    "type": 1,
    "chaincodeID": {
      "name": "mycc"
    },
    "ctorMsg": {
      "function": "submitTransaction",
      "args": [
        "NWVFZ3NY2HB4YJ",
        "CRED",
        "2017-01-01",
        "USD",
        "1000",
        "BANKA",
        "BANKB",
        "Alice",
        "Bob",            
        "OUR"
      ]
    },
    "secureContext": "jim"
  },
  "id": 1
}
```

### getTransactions
We can get information on  transactions by calling this function. If we want to fetch all transactions in which *BANKA* was either a sender or a receiver we will pass *BANKA* as an arugment. If we want to view all transactions we will pass *AUDITOR* as the argument.
Even if the transaction fails to credit/debit any accounts due to insufficient balance or invaild currency, the transaction will still be recorded in the Transactions type.

**Request**
```json
{
  "jsonrpc": "2.0",
  "method": "query",
  "params": {
    "type": 1,
    "chaincodeID": {
      "name": "mycc"
    },
    "ctorMsg": {
      "function": "getTransactions",
      "args": [
        "BANKA" //BANKB, BANKC OR AUDITOR
      ]
    },
    "secureContext": "jim"
  },
  "id": 1
}
```

**Response**
```json
{
  "jsonrpc": "2.0",
  "result": {
    "status": "OK",
    "message": "{\"transactions\":[{\"refNumber\":\"REFNO-001\",\"opCode\":\"OPCODE-001\",\"vDate\":\"VDATE-001\",\"currency\":\"USD\",\"amount\":10000,\"sender\":\"BANKA\",\"receiver\":\"BANKB\",\"ordcust\":\"OrdCust\",\"benefcust\":\"BenefCust\",\"detcharges\":\"DetCharges\",\"statusCode\":1,\"statusMsg\":\"Transaction Completed\"},{\"refNumber\":\"NWVFZ3NY2HB4YJ\",\"opCode\":\"CRED\",\"vDate\":\"2017-01-01\",\"currency\":\"USD\",\"amount\":260000,\"sender\":\"BANKA\",\"receiver\":\"BANKB\",\"ordcust\":\"Alice\",\"benefcust\":\"Bob\",\"detcharges\":\"OUR\",\"statusCode\":0,\"statusMsg\":\"Insufficient funds on Nostro Account\"},{\"refNumber\":\"NWVFZ3NY2HB4YJ\",\"opCode\":\"CRED\",\"vDate\":\"2017-01-01\",\"currency\":\"USD\",\"amount\":200000,\"sender\":\"BANKA\",\"receiver\":\"BANKB\",\"ordcust\":\"Alice\",\"benefcust\":\"Bob\",\"detcharges\":\"OUR\",\"statusCode\":1,\"statusMsg\":\"Transaction Completed\"}]}"
  },
  "id": 1
}
```