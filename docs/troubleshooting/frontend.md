# React Frontend Troubleshooting

## Module not found

**Error:**
```
Module not found: Can't resolve 'xxx'
```

**Solution:**
```bash
rm -rf node_modules package-lock.json
npm install
```

---

## CORS error

**Error:**
```
Access to fetch blocked by CORS policy
```

**Cause:**
Backend not allowing frontend origin.

**Solution:**
- Add CORS middleware to Go backend
- Check `Access-Control-Allow-Origin` header
- Verify API URL in frontend config

---

## Invalid hook call

**Error:**
```
Invalid hook call. Hooks can only be called inside of a function component.
```

**Cause:**
- Hook called outside component
- Multiple React versions
- Rules of hooks violated

**Solution:**
```bash
npm ls react  # Check for duplicate React
```
- Ensure hooks at top level of component
- Don't call hooks in conditions/loops

---

## State not updating

**Cause:**
- Mutating state directly
- Async state update

**Solution:**
```jsx
// Wrong
state.items.push(item)

// Correct
setItems([...items, item])
```

---

## useEffect infinite loop

**Cause:**
Missing or wrong dependency array.

**Solution:**
```jsx
// Wrong - runs every render
useEffect(() => { fetchData() })

// Correct - runs once
useEffect(() => { fetchData() }, [])
```

---

## Build fails in production

**Error:**
```
'X' is not defined
```

**Cause:**
ESLint stricter in production build.

**Solution:**
- Fix all eslint warnings
- Remove unused imports/variables
