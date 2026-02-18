# fullstack app

Explain your service in here. This is fulltsack project related Payment using golang as backend and nuxt as frontend....

list of tools version of your machine:

```bash
go version go1.25.5 darwin/arm64
node v24.13.1
make
sqlite3 (optional)
```

Install all related requirements:

```bash
# Backend
cd backend
cp env.sample .env
make tool-openapi
make openapi-gen
make dep
make gen-secret

# Frontend
cd ../frontend
cp .env.example .env
npm install
```

How to run backend server on local:

```bash
cd backend
make run
# Runs on localhost:8080 (default)
# Database seeds automatically on first run
```

How to run backend server on production build:

```bash
cd backend
make build
./bin/mygolangapp
```

How to run frontend on local:

```bash
cd frontend
npm run dev
# Runs on localhost:3000
```

How to run frontend on production build:

```bash
cd frontend
npm run build
npm run preview
```

To checking openapi documentations, you can visit this url after backend running.

```bash
# The OpenAPI spec is located at:
openapi.yaml

# You can view it using any OpenAPI/Swagger editor or VS Code extension.
# You can also visit the built-in Swagger UI at:
http://localhost:8080/docs
```

Login to frontend by visiting:

```bash
visit: http://localhost:3000/auth/login
email: cs@test.com
password: password

email: operation@test.com
password: password
```

## Evidences

[Video Evidence](https://drive.google.com/file/d/1oPMaROSwlLBl-Ur5bwITp3ES_zUOWyML/view?usp=sharing)

see backend [README.md](backend/README.md)
see frontend [README.md](frontend/README.md)
