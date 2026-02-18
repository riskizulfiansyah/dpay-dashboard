Scaffold demonstrating entity / repository / usecase / service separation.

generate openapi:

```bash
make openapi-gen
```

generate JWT_SECRET:

```bash
make gen-secret
```

Run server:

```bash
cp env.sample .env
make tool-openapi
make openapi-gen
make dep
make gen-secret
make run
```

API:

- POST /dashboard/v1/auth/login {email,password}
- GET /dashboard/v1/payments?sort=sort,status=status,id=id -> not implemented like this
- PUT /dashboard/v1/payment/{id}/review -> not implement this
- GET /dashboard/v1/payments?page=1&limit=20&sort=-created_at&status=completed&from_date=2026-01-01&to_date=2026-02-01
