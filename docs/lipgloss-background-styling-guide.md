# Lipgloss Background Styling: Why Order Matters

## Problem Statement

When trying to apply a background color to a full row of text in a Bubble Tea list, you might encounter an issue where the background only appears on the padding/edges but not on the text itself. This happens when you try to **wrap already-styled text** with a background style.

---

## The Two Approaches

### ❌ Approach 1: Wrapping (Doesn't Work)

```go
// Step 1: Render text with foreground colors
typ := lg.NewStyle().Foreground(styles.MutedForeground).Render(item.Type)
task := lg.NewStyle().Foreground(styles.PrimaryForeground).Render(item.Task)
priority := lg.NewStyle().Foreground(priorityColor).Render(item.Priority)

// Step 2: Join them together
content := typ + " " + task + padding + priority

// Step 3: Try to wrap with background
if selected {
    baseStyle = baseStyle.Background(styles.SelectedBackground)
}
fmt.Fprint(w, baseStyle.Render(content))
```

**Result:** Background only appears on padding, not on text.

### ✅ Approach 2: Apply Background First (Works)

```go
// Step 1: Create background style
bg := lg.NewStyle()
if selected {
    bg = bg.Background(styles.SelectedBackground)
}

// Step 2: Apply background to EACH segment
typ := bg.Foreground(styles.MutedForeground).Render(item.Type)
space := bg.Render(" ")  // Background on space too!
task := bg.Foreground(styles.PrimaryForeground).Render(item.Task)
padding := bg.Render(strings.Repeat(" ", paddingWidth))
priority := bg.Foreground(priorityColor).Render(item.Priority)

// Step 3: Join pre-styled segments
content := typ + space + task + padding + priority

// Step 4: Wrap with layout style (border, outer padding)
fmt.Fprint(w, baseStyle.Render(content))
```

**Result:** Background appears on entire row!

---

## Technical Explanation: ANSI Escape Codes

### What Happens When You Render

When lipgloss renders styled text, it converts it to **ANSI escape codes**:

```go
lg.NewStyle().Foreground(lg.Color("240")).Render("bug")
```

Becomes:
```
"\x1b[38;5;240mbug\x1b[0m"
 ^^^^^^^^^^^    ^^^^^^
 Foreground     RESET code
 color code
```

The `\x1b[0m` is the **ANSI reset code** that clears ALL styling.

### Why Wrapping Fails

**Approach 1 - Wrapping:**

```go
typ := lg.NewStyle().Foreground(color1).Render("bug")
// "\x1b[38;5;240mbug\x1b[0m"

task := lg.NewStyle().Foreground(color2).Render("Fix login")
// "\x1b[38;5;205mFix login\x1b[0m"

content := typ + " " + task
// "\x1b[38;5;240mbug\x1b[0m \x1b[38;5;205mFix login\x1b[0m"
//                      ^^^^^^                       ^^^^^^
//                      These resets prevent background from penetrating!

baseStyle.Background(bg).Render(content)
// Background can't penetrate the reset codes
// Only applies to unstyled parts (padding)
```

When you try to wrap this with a background, lipgloss sees:
- "bug" with its own styling sealed by reset → can't add background
- " " plain space → can add background (but it's just one space)
- "Fix login" with its own styling sealed by reset → can't add background

### Why Applying Background First Works

**Approach 2 - Background First:**

```go
bg := lg.NewStyle().Background(selectedBg)

typ := bg.Foreground(color1).Render("bug")
// "\x1b[48;5;8m\x1b[38;5;240mbug\x1b[0m"
//  ^^^^^^^^^^  Background code is INSIDE the render!

space := bg.Render(" ")
// "\x1b[48;5;8m \x1b[0m"
//  Background on space!

task := bg.Foreground(color2).Render("Fix login")
// "\x1b[48;5;8m\x1b[38;5;205mFix login\x1b[0m"
//  Background code embedded in the render
```

Each segment has **both background and foreground codes** applied together during rendering. The background is **embedded** in each segment, not added afterwards.

---

## The Key Principle

> **Lipgloss applies styles during rendering, not after.**

Once text is rendered with ANSI codes and reset sequences, those codes are **sealed**. You cannot retroactively add styling to already-rendered text.

### Order of Operations

**❌ Wrong Order:**
1. Apply foreground color → Render → RESET
2. Join segments
3. Try to apply background → **Can't penetrate resets**

**✅ Correct Order:**
1. Create background style
2. For each segment: Combine background + foreground → Render together
3. Join segments (all have background embedded)
4. Wrap with layout style (border, padding)

---

## Visual Analogy

### Approach 1: Painting Then Sealing

```
1. Paint text red     [bug]
2. Apply clear coat   [bug]  ← RESET code seals it
3. Try to paint blue background underneath
   → Blue can't get under the seal!
```

### Approach 2: Canvas First

```
1. Start with blue canvas   [    ]
2. Paint red text on it    [bug]  ← Both colors present
3. Text has blue background because canvas was blue from start
```

---

## Code Pattern to Follow

### Single Segment

```go
bg := lg.NewStyle()
if selected {
    bg = bg.Background(styles.SelectedBackground)
}

// Apply background + foreground together
styledText := bg.Foreground(textColor).Render(content)
```

### Multiple Segments (Full Row)

```go
// 1. Create background base
bg := lg.NewStyle()
if selected {
    bg = bg.Background(styles.SelectedBackground)
}

// 2. Apply to each segment
segment1 := bg.Foreground(color1).Render(text1)
space := bg.Render(" ")
segment2 := bg.Foreground(color2).Render(text2)
padding := bg.Render(strings.Repeat(" ", width))
segment3 := bg.Foreground(color3).Render(text3)

// 3. Join pre-styled segments
content := segment1 + space + segment2 + padding + segment3

// 4. Wrap with layout (border, outer padding)
fmt.Fprint(w, layoutStyle.Render(content))
```

---

## Common Mistakes

### Mistake 1: Using SetString Before Background

```go
// ❌ Wrong
typ := lg.NewStyle().
    Foreground(color).
    SetString(item.Type)  // Already rendered!

// Later trying to add background won't work on typ
```

```go
// ✅ Correct
bg := lg.NewStyle().Background(selectedBg)
typ := bg.Foreground(color).SetString(item.Type)  // Background applied first
```

### Mistake 2: Joining First, Styling Later

```go
// ❌ Wrong
left := styledText1 + " " + styledText2
right := styledText3
content := left + padding + right
styledContent := bg.Render(content)  // Too late!
```

```go
// ✅ Correct
left := bg.Foreground(c1).Render(text1) +
        bg.Render(" ") +
        bg.Foreground(c2).Render(text2)
right := bg.Foreground(c3).Render(text3)
padding := bg.Render(strings.Repeat(" ", width))
content := left + padding + right  // All segments already have background
```

### Mistake 3: Forgetting Spaces/Padding

```go
// ❌ Wrong - spaces between segments won't have background
content := styledSeg1 + " " + styledSeg2  // Plain space!
```

```go
// ✅ Correct - explicitly style the space
space := bg.Render(" ")
content := styledSeg1 + space + styledSeg2
```

---

## Why the Outer Wrapper Still Works

You might wonder: "If wrapping doesn't work for background, why does the outer `baseStyle.Render()` work?"

```go
baseStyle := lg.NewStyle().
    Border(lg.NormalBorder(), false, false, true, false).
    PaddingLeft(4).
    PaddingRight(4)

fmt.Fprint(w, baseStyle.Render(content))
```

**Answer:** The outer wrapper adds **layout styling** (borders, padding) **around** the content, not **modifying** the ANSI codes within it.

Think of it as:
- Inner rendering: Applies colors/backgrounds to the text itself
- Outer rendering: Adds a frame/container around the colored text

The outer wrapper doesn't try to change the colors inside - it just adds structural elements around them.

---

## Real-World Example: gh-dash

The gh-dash project uses this exact pattern:

```go
// From gh-dash prrow.go
func (pr *PullRequest) RenderLines(isSelected bool) string {
    baseStyle := lipgloss.NewStyle()
    if isSelected {
        baseStyle = baseStyle.Background(pr.Ctx.Theme.SelectedBackground)
    }

    // Apply to each segment
    additionsText := baseStyle.Foreground(additionsFg).Render(additions)
    deletionsText := baseStyle.Foreground(deletionsFg).Render(deletions)
    space := baseStyle.Render(" ")  // Don't forget the space!

    // Join segments that already have background
    return lipgloss.JoinHorizontal(lipgloss.Left,
        additionsText,
        space,
        deletionsText,
    )
}
```

**Result:** Full-width background when selected, because each segment (including spaces) has the background applied before joining.

---

## Quick Reference

### When to Apply Background

| Scenario | When to Apply Background |
|----------|-------------------------|
| Single color text | Apply together with foreground |
| Multi-color text in one row | Apply to each color segment individually |
| Spaces between segments | Apply to spaces explicitly |
| Padding/whitespace | Apply to padding strings |
| Layout (borders) | Apply in outer wrapper (separate from background) |

### Style Layering Order

```
1. Background (if selected)
   ↓
2. Foreground color
   ↓
3. Render segment
   ↓
4. Join all segments
   ↓
5. Wrap with layout style (border, padding)
```

---

## Summary

**Key Takeaway:** In lipgloss, you must apply the background **during rendering**, not after. Once text is rendered with ANSI codes, it's sealed and cannot be retroactively styled.

**The Fix:**
1. Create a base style with background (if selected)
2. Apply that base style + foreground to each segment
3. Join the pre-styled segments
4. Wrap with layout style for borders/padding

**Remember:** Every piece of the row needs the background applied individually - text segments, spaces, and padding. Only then will you get a full-width highlighted row.
