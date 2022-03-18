## Geo-Service
You can use this package to parse CSV files and persist them in a database.

## CSV File Format
```csv
ip_address,country_code,country,city,latitude,longitude,mystery_value
200.106.141.15,SI,Nepal,DuBuquemouth,-84.87503094689836,7.206435933364332,7823011346
160.103.7.140,CZ,Nicaragua,New Neva,-68.31023296602508,-37.62435199624531,7301823115
```
- The first row should always be the header since the package will ignore it.
- The package will also discard the row on:
  - Empty `ip_address`
  - Empty `latitude` or `longitude`
  - `latitude` and `longitude` not being of type Float
  - `mystery_value` not being of type int
  - More than 7 columns in a row
    - There's an exception to this, for `country` & `city` you can use more than one comma escaped by double/single quotes ("/').
    - For example: `"Virgin Islands, U.S."`, `'Virgin Islands, U.S.'`.

## Repository Methods
The package has an interface called `Repository` with 3 functions to store or retrieve data:
- Store
- StoreMany
- Retrieve