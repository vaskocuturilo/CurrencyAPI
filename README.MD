Currency API is a free, open-source exchange rate API for 217 currencies from 3 central banks and official sources. Rates are updated daily.

### Providers (Banks):
  - Banca d'Italia BDI (1999 – present)
  - European Central Bank ECB (1999 – present)
  - Banca Națională a României (2005 – present)

## Coverage 

* The existing `/v1/currencies` returns all currencies as of the current date. 
* The existing `/v1/currencies?scope=currency` returns the exchange rate for the current date (It is necessary to use the currency encoding according to [ISO 4217](https://www.iban.com/currency-codes)
* The existing `/v1/currencies?scope=all` catalogue covers eligible daily sources plus pegged currencies. 
* While `/v1/currencies?providers=BIS` describes that provider's recognized currencies. 
* Those currency dates, and the dates in `/v1/providers`, describe publication coverage. 

## Deployment

The public API runs at https://api.******.dev. If you prefer to host your own instance, you can run the Currency API with Docker.

## Using Docker

The quickest way to get started:

```bash
docker run -d --init -p 8080:8080 vaskocuturilo/currency-service
```

## How it works:

        [  Bank Feeds ]
              │ (Daily Cron Job)
              ▼
    [ Database / Memory Cache ] 
              │
              ▼
     [ REST API (Go) ]
              │
              ▼
        [ Response ]
