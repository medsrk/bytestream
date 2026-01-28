# Bytestream Video Service

## Running

```bash
go run cmd/main.go
curl -H "Authorization: Bearer any-token" http://localhost:8080/video/46325
```

## Production Considerations

**Resilience** - retries with backoff, circuit breakers, per-service timeouts.

**Observability** - metrics instrumentation and distributed tracing.

**Availability timezones** - Might want to consider that UTC comparison may not match user expectations.