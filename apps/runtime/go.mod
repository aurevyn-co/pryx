module pryx-core

go 1.24.0

require (
	github.com/go-chi/chi/v5 v5.0.12
	github.com/google/uuid v1.6.0
	github.com/mattn/go-sqlite3 v1.14.22
	github.com/playwright-community/playwright-go v0.5200.1
	github.com/zalando/go-keyring v0.2.4
	github.com/stretchr/testify v1.10.0
)

require (
	github.com/alessio/shellescape v1.4.1
	github.com/danieljoos/wincred v1.2.0
	github.com/deckarep/golang-set/v2 v2.7.0
	github.com/go-jose/go-jose/v3 v0.6.4
	github.com/go-stack/stack v1.8.1
	github.com/godbus/dbus/v5 v5.1.0
	github.com/kr/text v0.2.0
)

# Test Coverage Configuration

[profile.test.coverage]
go test -coverprofile=test.coverage ./...

[cover]
# Coverage reporting configuration
go test -cover ./... -coverpkg=./...

[report]
# Coverage reporting configuration
go test -coverprofile=test.coverage ./... -covermode=atomic -json > coverage/coverage.json

# Coverage threshold
