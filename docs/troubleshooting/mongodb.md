# MongoDB Troubleshooting

## Connection refused

**Error:**
```
connection refused: localhost:27017
```

**Cause:**
MongoDB not running.

**Solution:**
```bash
docker-compose up mongo
# Or check if container is healthy
docker-compose ps
```

---

## Authentication failed

**Error:**
```
authentication failed
```

**Cause:**
Wrong credentials.

**Solution:**
- Check `MONGO_URI` env var
- Verify username/password in docker-compose
- Reset: `docker-compose down -v && docker-compose up`

---

## Duplicate key error

**Error:**
```
E11000 duplicate key error collection
```

**Cause:**
Inserting document with existing unique field.

**Solution:**
- Use `upsert` if updating
- Check unique index fields
- Generate new ID if needed

---

## Document too large

**Error:**
```
document exceeds maximum size (16MB)
```

**Solution:**
- Use GridFS for large files
- Store file references instead of content
- Split into multiple documents

---

## Cursor timeout

**Error:**
```
cursor not found
```

**Cause:**
Cursor expired during iteration.

**Solution:**
```go
opts := options.Find().SetNoCursorTimeout(true)
// Or batch with Skip/Limit
```

---

## Index not used

**Cause:**
Query not matching index.

**Solution:**
```bash
# Check query explain
db.collection.find({...}).explain("executionStats")
```
- Create compound index matching query fields
- Check field order in compound index
