# To lint, install Earthly and run `earth +lint`
# This ensures the usage of the same version of golangci-lint

FROM golangci/golangci-lint:v1.32

WORKDIR /tmpl

lint:
    COPY . .
    RUN golangci-lint --timeout 5m run