# Plan: Add uses_weight / uses_duration flags to exercise master

## Context
ExercisePicker แสดง Weight (kg) และ Duration (s) ทุก exercise ทั้งที่บางท่าไม่ใช้ (Plank ไม่ต้องใส่ weight, Bench Press ไม่ต้องใส่ duration) ต้องการให้ master data บอกว่าแต่ละ exercise ใช้ weight/duration หรือไม่ และให้ Settings page สามารถ Add/Edit exercise พร้อม flags เหล่านี้ได้ด้วย

## Files to modify

| File | Change |
|---|---|
| `backend/cmd/api/main.go` | เพิ่ม uses_weight/uses_duration ใน seed + upsert logic |
| `frontend/src/constants/master.ts` | เพิ่ม boolean fields ใน defaultFields config |
| `frontend/src/components/MasterDataCard.tsx` | support type: "boolean" (checkbox) ใน InlineForm + display |
| `frontend/src/components/ExercisePicker.tsx` | lookup flags จาก master และซ่อน field ตาม flags |

---

## Step 1 — Backend: seed upsert + flags

**File:** `backend/cmd/api/main.go`

เพิ่ม `usesWeight bool` และ `usesDuration bool` ใน `ex` struct:

```go
type ex struct {
  name        string
  muscle      string
  sets        int
  reps        int
  rest        int
  usesWeight  bool
  usesDuration bool
}
```

กำหนดค่าให้แต่ละ exercise:
- Plank → `usesWeight: false, usesDuration: true`
- ที่เหลือทั้งหมด → `usesWeight: true, usesDuration: false`

เพิ่ม `"uses_weight"` และ `"uses_duration"` ใน Defaults map

เปลี่ยน seed logic เป็น **upsert** (ไม่ใช้ clear เพราะ user อาจเพิ่ม exercise เองผ่าน Settings):
```go
existing, _ := repo.FindByCategory(ctx, "exercise")
byName := map[string]*domain.MasterRecord{}
for _, r := range existing { byName[r.Name] = r }

for _, record := range records {
  if old, found := byName[record.Name]; found {
    repo.Update(ctx, old.ID, record)  // อัพ defaults ใหม่
  } else {
    repo.Insert(ctx, record)
  }
}
```

---

## Step 2 — Frontend constants: เพิ่ม boolean type

**File:** `frontend/src/constants/master.ts`

1. เพิ่ม `"boolean"` ใน `MasterDefaultField.type`:
   ```ts
   type: "number" | "string" | "select" | "boolean";
   ```

2. เพิ่ม fields ใน `defaultFields` ของ workout config:
   ```ts
   { key: "uses_weight",   label: "Uses Weight",   type: "boolean", defaultValue: true },
   { key: "uses_duration", label: "Uses Duration",  type: "boolean", defaultValue: false },
   ```

---

## Step 3 — MasterDataCard: support boolean field

**File:** `frontend/src/components/MasterDataCard.tsx`

1. `setField` function: เพิ่ม branch สำหรับ boolean (string `"true"`/`"false"` → boolean)

2. `InlineForm`: เพิ่ม branch สำหรับ `f.type === "boolean"` → render `<input type="checkbox">` แทน text input

3. Record display (list view): boolean fields แสดงเป็น Yes/No badge แทนตัวเลข

---

## Step 4 — ExercisePicker: conditional fields

**File:** `frontend/src/components/ExercisePicker.tsx`

ใน render loop `selected.map((entry) => ...)`:

```ts
const masterRec = records.find((r) => r.name === entry.name);
const usesWeight   = masterRec ? Boolean(masterRec.defaults.uses_weight   ?? true)  : true;
const usesDuration = masterRec ? Boolean(masterRec.defaults.uses_duration ?? false) : false;
```

**Planned fields** (บรรทัด ~364): filter ก่อน render:
```ts
{ label: "Weight (kg)", field: "weight", min: 0, show: usesWeight },
{ label: "Duration (s)", field: "duration_seconds", min: 0, show: usesDuration },
```
เพิ่ม `.filter(f => f.show !== false)` ก่อน `.map()`

**Per-set log columns** (บรรทัด ~485): เช่นกัน — filter weight/duration ตาม flags

**Set summary line** (บรรทัด ~455): ซ่อน weight/duration ใน Plan: text ถ้าไม่ใช้

---

## Verification

1. Restart backend → exercises ใน DB มี `uses_weight`/`uses_duration` fields
2. Frontend → Workouts → สร้าง workout → เลือก **Plank** → เห็น Duration แต่ไม่เห็น Weight
3. เลือก **Bench Press** → เห็น Weight แต่ไม่เห็น Duration
4. Edit workout เก่าที่มี Plank → set log ก็ต้องซ่อน weight column
5. Settings → Exercises → Add exercise → form มี checkbox "Uses Weight" / "Uses Duration"
