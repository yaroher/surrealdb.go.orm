# Access annotations

Use a standalone access type with `//orm:access`:

```go
//orm:access name=app_access on=database type=record signup="create user" signin="select user" refresh=true
//orm:access name=jwt_access on=namespace type=jwt alg=HS256 key="mysecret"
//orm:access name=jwt_url on=namespace type=jwt url="https://example.com/jwks"
//orm:access name=bearer_user on=database type=bearer alg=user
```

Supported keys:
- name (or access)
- on: namespace|database|root
- type: jwt|record|bearer
- overwrite: true
- if_not_exists: true
- alg: jwt algorithm, or for bearer use alg=user|record
- key: jwt key
- url: jwt url
- signup / signin (record)
- issuer (record issuer key)
- refresh: true
- authenticate: raw expr
- duration_grant / duration_token / duration_session
