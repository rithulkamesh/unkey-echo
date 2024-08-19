# Unkey-Echo Middleware Example

This project demonstrates the usage of Unkey-Echo middleware for handling API requests.

## Environment Setup

Create a `.env` file in the root directory with the following variables:

```
UNKEY_ROOT_KEY=<Your root key from app.unkey.com/settings/root-keys>
UNKEY_API_ID=<Your Project API ID>
```

## Description

This middleware handles any request thrown to the server.

## Usage

To use the API, make a GET request to the server with an Authorization header containing your API key:

```bash
curl -X GET \
  http://localhost:8080/ \
  -H 'Authorization: key_your_api_key'
```

Replace `key_your_api_key` with your actual API key.

## Note

Ensure that your server is running on `localhost:8080` or update the URL in the curl command accordingly.
