# Docker Troubleshooting

## Port already in use

**Error:**
```
Bind for 0.0.0.0:3000 failed: port is already allocated
```

**Cause:**
Another process using the port.

**Solution:**
```bash
# Find process
netstat -ano | findstr :3000  # Windows
lsof -i :3000                 # Linux/Mac

# Kill or change port in docker-compose.yml
```

---

## Container won't start

**Error:**
```
container exited with code 1
```

**Cause:**
App crash on startup.

**Solution:**
```bash
docker-compose logs <service>
# Check for missing env vars or config errors
```

---

## Image build fails

**Error:**
```
failed to solve: process "/bin/sh -c ..." did not complete successfully
```

**Cause:**
Build step failed (missing deps, syntax error).

**Solution:**
```bash
# Rebuild without cache
docker-compose build --no-cache
```

---

## Volume permission denied

**Error:**
```
permission denied: '/data/db'
```

**Cause:**
Host/container user mismatch.

**Solution:**
```bash
# Reset volumes
docker-compose down -v
docker-compose up
```

---

## Out of disk space

**Error:**
```
no space left on device
```

**Solution:**
```bash
docker system prune -a
docker volume prune
```
