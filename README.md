# About

hasura-ad-webhook is used in conjunction with [Hasura GraphQL Engine](https://github.com/hasura/graphql-engine) to allow for authentication via Active Directory or a backend service using an API key. Hasura GraphQL Engine communicates with this service via its [webhook mode](https://hasura.io/docs/1.0/graphql/manual/auth/authentication/webhook.html).

# How it Works

You configure hasura-ad-webhook to map Active Directory groups to Hasura GraphQL Engine roles. Get a session ID and available roles from the `/api/1.0/auth` endpoint and send that in your GraphQL request. 

You can also map service API keys to roles. Send the API key in your GraphQL request.

## Example

```bash
# Authenticated client example
# Get a session ID and allowed roles by POSTing username and password
curl -s -X POST -H "Content-Type: application/json" \
    -d '{"username":"john.smith","password":"password"}' \
    http://<server>:<port>/api/1.0/auth | jq .

{
  "username": "john.smith",
  "display_name": "John Smith",
  "session_id": "e376011c-cfce-4119-ab59-c04c793fea3d",
  "attrs": {
    "roles": [
      "manager",
      "viewer"
    ]
  }
}

# You would normally pass these headers in your GraphQL request
# and the response would be available to Hasura GraphQL Engine
curl -s -X GET -H "Authorization: Bearer e376011c-cfce-4119-ab59-c04c793fea3d" \
    -H "X-Hasura-Role: manager" \
    http://localhost:8082/api/1.0/webhook | jq .

{
  "X-Hasura-User-Id": "john.smith",
  "X-Hasura-Role": "manager"
}


# Service example
# Your service normally pass these headers in its GraphQL request
# and the response would be available to Hasura GraphQL Engine
curl -s -X GET -H "Authorization: Bearer <api_key>" \
    -H "X-Authorization-Type: API-Key" \
    http://localhost:8082/api/1.0/webhook | jq .

{
  "X-Hasura-Role": "service"
}

```


# Install

```bash
go get github.com/korylprince/hasura-ad-webhook
```

# Configuration

The server is configured with environment variables:

```bash
LDAPSERVER="ldap.example.com"
LDAPPORT="389"
LDAPBASEDN="OU=Container,DC=example,DC=net"
LDAPSECURITY="starttls" # none, tls, or starttls
GROUPROLEMAP="Domain Admins:admin,viewer;Domain User:viewer"
    # format: <group 1 cn/DN>:<role 1>,<role 2>,...;<group 2 cn/DN>:<role 3>,<role 4>,...;...
APIKEYROLEMAP="reallylongkey:service,anotherkey:anotherrole"
    # format <key 1>:<role 1>,<key 1>:<role 2>,...
LISTENADDR=":8080"
PREFIX="/prefix" # Used to prefix all URLs
```

For more information see [config.go](https://github.com/korylprince/hasura-ad-webhook/blob/master/httpapi/config.go).

## Hasura GraphQL Engine

Configure Hasura GraphQL Engine to communicate with this service by setting the [`HASURA_GRAPHQL_AUTH_HOOK` environment variable or `--auth-hook` flag](https://hasura.io/docs/1.0/graphql/manual/auth/authentication/webhook.html#configuring-webhook-mode) to the webhook endpoint: `http[s]://<server>:<port>/api/1.0/webhook`


# Docker

You can use the pre-built Docker container, ghcr.io/korylprince/hasura-ad-webhook.

## Docker Configuration


The Docker container supports [Docker Secrets](https://docs.docker.com/engine/swarm/secrets/) by appending `_FILE` to any variable, e.g. `APIKEYROLEMAP_FILE=/run/secrets/<secret_name>`.

Additionally, you can specify individual API Key Roles using `APIKEYROLEMAP_[1-9]_KEY` and `APIKEYROLEMAP_[1-9]_ROLE`. These will be appended to whatever is in `APIKEYROLEMAP`.


## Example

```bash
docker run -d --name="hasura-ad-webhook" \
    -p 80:80 \
    -e LDAPSERVER="ldap.example.com" \
    -e LDAPPORT="389" \
    -e LDAPBASEDN="OU=Container,DC=example,DC=net" \
    -e LDAPSECURITY="starttls" \
    -e GROUPROLEMAP="Domain Admins:admin,viewer;Domain User:viewer" \
    -e APIKEYROLEMAP="reallylongkey:service,anotherkey:anotherrole" \
    -e APIKEYROLEMAP_1_ROLE="role1" \
    -e APIKEYROLEMAP_1_KEY_FILE="/run/secrets/key_1" \
    -e APIKEYROLEMAP_2_ROLE="role2" \
    -e APIKEYROLEMAP_2_KEY_FILE="/run/secrets/key_2" \
    -e LISTENADDR=":80" \
    --restart="always" \
    ghcr.io/korylprince/hasura-ad-webhook:latest
```
