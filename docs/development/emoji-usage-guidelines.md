# Emoji Usage Guidelines for 7zarch-go TUI

## 🎯 **Core Principle: Monospace Reliability**

**Problem:** Emoji can render as 1 or 2 characters depending on terminal/font, breaking monospace layouts.

**Solution:** Use text indicators for structural UI elements, emoji only where layout is safe.

---

## ✅ **SAFE EMOJI ZONES**

### **1. Command Responses (Single Line)**
```bash
:push complete ✅             # Safe - no alignment dependency
:analyze finished 🎯          # Safe - standalone message
Error: File not found ❌      # Safe - error message
```

### **2. Progress Messages**
```bash
Uploading episode-423.7z... 📤  # Safe - progress doesn't need alignment
Download complete 📥            # Safe - status message
Archive verified ✅             # Safe - confirmation message
```

### **3. Confirmation Dialogs**
```bash
┌─ Confirm Delete ─────────────┐
│ Delete episode-423.7z? 🗑️   │  # Safe - dialog box content
│                             │
│ [y] Yes  [n] No             │
└─────────────────────────────┘
```

### **4. Help Text and Documentation**
```bash
🎵 For podcast episodes, use --profile media
📝 For documents, use --profile documents  
⚖️ For mixed content, use --profile balanced
```

---

## ❌ **EMOJI DANGER ZONES**

### **1. List/Table Layouts**
```bash
# BAD - Emoji breaks monospace alignment
> episode-423.7z    89 MB  2h ago  ✓     # ✓ might be 1 or 2 chars
  episode-422.7z    92 MB  1d ago  ✓     # Misaligns columns

# GOOD - Text indicators maintain alignment  
> episode-423.7z    89 MB  2h ago  OK    # Always 2 chars
  episode-422.7z    92 MB  1d ago  OK    # Perfect alignment
```

### **2. Status Columns**
```bash
# BAD - Inconsistent width
Status
------
✓      # Variable width
?      # Breaks table
X      # Alignment lost

# GOOD - Fixed width
Status
------
OK     # Always 2-4 chars
MISS   # Predictable layout
DEL    # Reliable alignment
```

### **3. Selection Indicators in Lists**
```bash
# BAD - Selection breaks when emoji width varies
[✓] episode-423.7z             # ✓ width unpredictable
[ ] episode-422.7z             # Spacing breaks

# GOOD - ASCII characters are monospace-safe
[x] episode-423.7z             # Always 1 char
[ ] episode-422.7z             # Consistent spacing
```

---

## 📋 **ESTABLISHED TEXT INDICATORS**

### **Status Indicators (Monospace-Safe)**
- **OK** - Present archives, functioning correctly
- **MISS** - Missing files, archive not found
- **DEL** - Deleted archives, in trash

### **Location Indicators (Established)**
- **MANAGED** - In MAS local storage (staging area)
- **EXTERNAL** - External file system locations

### **Future Indicators (Proposed)**
- **REMOTE** - TrueNAS remote storage (when implemented)
- **STAGED** - Ready for upload/push (status, not location)
- **SYNC** - Synchronized between local and remote

---

## 🎨 **IMPLEMENTATION STRATEGY**

### **TUI List Layout (Monospace-Safe)**
```
7zarch-go

Archives: 247

> episode-423.7z        89 MB   2h ago   OK   MANAGED
  episode-422.7z        92 MB   1d ago   OK   MANAGED  
  vacation-photos.7z   3.8 GB   1w ago   OK   EXTERNAL
  old-backup.7z         2.1 GB   6m ago   MISS EXTERNAL

[Enter] Details  [Space] Select  [d] Delete  [q] Quit
```

### **Safe Emoji Usage**
- **Command feedback:** `:push complete ✅`
- **Progress messages:** `Uploading... 📤`
- **Help documentation:** `🎵 Use media profile for audio`
- **Confirmation dialogs:** Content that doesn't affect layout

**This maintains beautiful theming while ensuring rock-solid monospace reliability!** 

Should I update the current TUI implementation with these proper text indicators? 🎯
