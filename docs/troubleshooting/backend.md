# Go Backend Troubleshooting

## Module not found

**Error:**
```
cannot find module providing package github.com/xxx
```

**Solution:**
```bash
go mod tidy
go mod download
```

---

## Nil pointer dereference

**Error:**
```
panic: runtime error: invalid memory address or nil pointer dereference
```

**Cause:**
Accessing uninitialized pointer/interface.

**Solution:**
- Check if dependency injection wired correctly
- Verify repository/service initialization in `cmd/`
- Add nil checks before access

---

## Context deadline exceeded

**Error:**
```
context deadline exceeded
```

**Cause:**
Operation took longer than context timeout.

**Solution:**
- Increase timeout for slow operations
- Check MongoDB connection
- Check external API latency

---

## Port binding failed

**Error:**
```
listen tcp :8080: bind: address already in use
```

**Solution:**
```bash
# Find and kill process
netstat -ano | findstr :8080
taskkill /PID <pid> /F
```

---

## JSON marshal/unmarshal error

**Error:**
```
json: cannot unmarshal string into Go struct field
```

**Cause:**
Type mismatch between JSON and struct.

**Solution:**
- Check struct tags match JSON keys
- Verify field types (string vs int, etc.)
- Use `json.RawMessage` for dynamic fields

---

## Import cycle

**Error:**
```
import cycle not allowed
```

**Cause:**
Package A imports B, B imports A.

**Solution:**
- Move shared types to separate package
- Use interfaces to break dependency
- Review architecture layers
