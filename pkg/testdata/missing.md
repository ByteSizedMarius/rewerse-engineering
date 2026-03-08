# Missing Unmarshal Test Coverage

Responses where I have no source data from my api dump. Need to capture actual API responses to add proper tests at some point.

## No fixture at all

- `bulkyGoodsResponse`: need a successful response from a delivery market
- `productRecommendationsResponse`: need a response with actual products

## Fixture exists but doesn't exercise key fields

- `basketResponse`: fixture is an empty newly-created basket. Missing coverage for:
  - `lineItems` with nested `Product`, `Listing`, `Attributes`
  - `violations` with non-null `detailMessage` (*string)
  - `orderId` non-null (*string)
  - `fees` with non-null `*int` fields (beverageSurcharge, serviceFee, etc.)
  - `timeSlotInformation` with populated startTime/endTime/timeSlotPrice
- `productDetailResponse`: fixture has `featureBenefit: null` and `additionalImageURLs: null`. Partly missing coverage for non-null arrays
- `servicePortfolioResponse`: fixture has a `deliveryMarket`. Missing coverage for `deliveryMarket: null` case
