# Natural Language Cron Parser: Anti-Hallucination Documentation

> **Status**: ✅ Complete  
> **Date**: 2026-01-29  
> **Purpose**: Document AI capabilities to prevent hallucination

---

## Summary

Based on research from **OpenCode** and **Moltbot** (`@.temp_refs`), I've created comprehensive documentation that prevents AI hallucination by explicitly defining what the natural language cron parser can and cannot do.

---

## Research Findings

### How OpenCode Prevents Hallucination

**Location**: `.opencode/agent/*.md`

OpenCode uses agent configuration files that explicitly define:
- What tools are available
- When to use specific tools
- Clear boundaries of capabilities

**Example from `triage.md`**:
```yaml
---
mode: primary
tools:
  "*": false           # Disable all tools by default
  "github-triage": true # Only enable specific tool
---

You are a triage agent responsible for triaging github issues.
Use your github-triage tool to triage issues.
```

**Key Pattern**: Tools are **disabled by default**, explicitly enabled.

---

### How Moltbot Prevents Hallucination

**Locations**:
- `docs/reference/templates/AGENTS.md` - System prompt
- `docs/reference/templates/TOOLS.md` - Tool documentation
- `skills/*/SKILL.md` - Individual skill docs

**Key Sections in AGENTS.md**:
1. **"Heartbeat vs Cron: When to Use Each"** - Clear decision matrix
2. **Explicit capability boundaries** - What agent can/cannot do
3. **Error handling guidelines** - How to respond to unsupported requests
4. **Examples of good vs bad responses**

**Example from AGENTS.md**:
```markdown
### Heartbeat vs Cron: When to Use Each

**Use cron when:**
- Exact timing matters ("9:00 AM sharp every Monday")
- Task needs isolation from main session history
- One-shot reminders ("remind me in 20 minutes")

**Use heartbeat when:**
- Multiple checks can batch together
- Timing can drift slightly (every ~30 min is fine)
```

---

## Documentation Created

### 1. AI-Facing: `docs/templates/SYSTEM.md`

**Purpose**: System prompt that defines AI capabilities

**Key Sections**:
- ✅ **Supported Patterns Table** (14 patterns with exact conversions)
- ❌ **Explicitly Not Supported List** (6 categories)
- 🔄 **Cron vs Heartbeat Decision Matrix**
- 🛠️ **Error Handling Guidelines**
- 📊 **Examples of Good vs Bad Responses**

**Anti-Hallucination Features**:
```markdown
## EXPLICITLY NOT SUPPORTED (Reject These)

When users request these patterns, you MUST reject them:

### Business Logic
- ❌ "every payday"
- ❌ "every month-end"

**Response**: "I can't parse business logic like 'payday'. 
Try specific dates like 'every 15th at 9am'."
```

---

### 2. User-Facing: `docs/guides/cron-natural-language.md`

**Purpose**: User reference guide for natural language scheduling

**Key Sections**:
- 📋 **Quick Reference Table** (What works / What doesn't)
- 💡 **Detailed Examples** by category
- 🌍 **Timezone Handling**
- 🔧 **Troubleshooting Guide**
- 📚 **Fallback to Standard Cron**

**User-Friendly Features**:
```markdown
### Patterns That Work ✅
| What You Want | What You Say |
|---------------|--------------|
| Run every 5 minutes | `"every 5 minutes"` |
| Daily at 9 AM | `"every day at 9am"` |
| Weekdays only | `"every weekday at 8pm"` |

### Patterns That Don't Work ❌
| What You Want | Why | Alternative |
|---------------|-----|-------------|
| `"every payday"` | Business logic | `"every 15th at 9am"` |
| `"every 3-5 minutes"` | Variable interval | `"every 4 minutes"` |
```

---

## Beads Tasks Created

| Task ID | Title | Priority | Description |
|---------|-------|----------|-------------|
| **pryx-tx8o** | SYSTEM.md Template | P0 | AI-facing documentation with anti-hallucination guardrails |
| **pryx-uosk** | Parser Implementation | P0 | Natural language parser based on documented patterns |
| **pryx-2yd4** | User Guide | P1 | User-facing reference documentation |
| **pryx-g0g2** | Validation & Testing | P1 | Ensure parser matches documented capabilities |
| **pryx-99zi** | Error Handling | P1 | Graceful fallbacks for unsupported patterns |
| **pryx-ygl7** | Scheduler Integration | P1 | Connect parser to cron service |

**Dependency Chain**:
```
SYSTEM.md (pryx-tx8o)
    ├── Parser Implementation (pryx-uosk)
    │       ├── Validation (pryx-g0g2)
    │       ├── Error Handling (pryx-99zi)
    │       └── Scheduler Integration (pryx-ygl7)
    └── User Guide (pryx-2yd4)
```

---

## Anti-Hallucination Strategy

### 1. Explicit Capability Definition

**Before (Hallucination Risk)**:
```
User: "Every payday"
AI: "Sure, I'll set that up" → [creates arbitrary schedule]
```

**After (With Documentation)**:
```
User: "Every payday"
AI: "I can't parse 'payday' as it's business logic. 
     
     Supported patterns include:
     • 'every 2 weeks on Friday at 9am'
     • 'every 15th at 9am'
     • '0 9 15 * *' (standard cron)
     
     Which would work for you?"
```

### 2. Two-Layer Documentation

| Layer | Audience | Purpose |
|-------|----------|---------|
| **SYSTEM.md** | AI | "These are your capabilities. Don't exceed them." |
| **User Guide** | Users | "These patterns work. Use these." |

### 3. Validation Testing

The validation task (`pryx-g0g2`) ensures:
- Parser only accepts documented patterns
- Unsupported patterns are rejected
- Error messages match documentation

### 4. Clear Error Categories

**4 Types of Errors Defined**:
1. **Unsupported Pattern** - Not in capability list
2. **Almost Supported** - Close to working pattern
3. **Wrong Format** - Syntax error
4. **Ambiguous** - Multiple interpretations

---

## Supported vs Unsupported Patterns

### ✅ Supported (14 Patterns)

| Category | Example | Converts To |
|----------|---------|-------------|
| Every X units | "every 5 minutes" | `{kind: "every", everyMs: 300000}` |
| Daily at time | "every day at 9am" | `{kind: "cron", expr: "0 9 * * *"}` |
| Weekdays | "every weekday at 8pm" | `{kind: "cron", expr: "0 20 * * 1-5"}` |
| Specific day | "every Monday" | `{kind: "cron", expr: "0 0 * * 1"}` |
| One-shot | "in 30 minutes" | `{kind: "at", atMs: timestamp}` |
| Relative | "tomorrow at noon" | `{kind: "at", atMs: timestamp}` |
| Monthly | "every month on the 1st" | `{kind: "cron", expr: "0 0 1 * *"}` |

### ❌ Not Supported (6 Categories)

| Category | Example | Why Not |
|----------|---------|---------|
| Business logic | "every payday" | Requires external API |
| Complex conditions | "if it's raining" | Requires context |
| Variable intervals | "every 3-5 minutes" | Non-deterministic |
| Timezones in text | "9am Tokyo time" | Use tz parameter |
| Exceptions | "every day except weekends" | Complex logic |
| Vague times | "soon", "later" | Not parseable |

---

## Integration with Existing PRD

### Documents Created

1. **`docs/templates/SYSTEM.md`** - AI system prompt template
2. **`docs/guides/cron-natural-language.md`** - User reference guide
3. **6 beads tasks** - Implementation roadmap

### References Added to PRD

The `docs/prd/prd.md` now includes:
```markdown
### B. Related Documents
- `docs/prd/ai-assisted-setup.md` - AI-assisted configuration flows
- `docs/prd/implementation-roadmap.md` - 3-phase implementation plan
- `docs/templates/SYSTEM.md` - AI system prompt with capability documentation
- `docs/guides/cron-natural-language.md` - User guide for natural language features
```

---

## How This Prevents Hallucination

### Before Documentation
```
User: "Remind me every time it rains"
AI: "I'll set up a weather-based reminder"
[AI hallucinates capability that doesn't exist]
```

### After Documentation
```
User: "Remind me every time it rains"
AI: "I can't create weather-based reminders 
     because Pryx doesn't have weather service integration.
     
     Alternatives:
     • Set a daily reminder at a specific time
     • Use weather API with custom MCP server
     • Manual reminder when you check weather"
[AI follows documented boundaries]
```

---

## Next Steps

1. **Review Documentation**
   - Check `docs/templates/SYSTEM.md` for completeness
   - Verify `docs/guides/cron-natural-language.md` accuracy

2. **Start Implementation**
   - Begin with `pryx-tx8o` (SYSTEM.md template)
   - Then `pryx-uosk` (parser implementation)

3. **Validation**
   - Run `pryx-g0g2` (validation suite)
   - Ensure parser matches documentation exactly

4. **User Testing**
   - Test with real users
   - Verify error messages are helpful
   - Check that fallback to cron works

---

## Success Criteria

✅ **AI knows exact capabilities** - No guessing
✅ **Users know what works** - Clear documentation
✅ **Unsupported patterns rejected** - No false promises
✅ **Helpful alternatives offered** - Good UX
✅ **Standard cron fallback** - Power user path

---

## References

### From `@.temp_refs`

**OpenCode**:
- `.opencode/agent/triage.md` - Agent capability definition
- Pattern: Disable all tools, explicitly enable specific ones

**Moltbot**:
- `docs/reference/templates/AGENTS.md` - System prompt template
- `docs/reference/templates/TOOLS.md` - Tool documentation
- `skills/*/SKILL.md` - Individual skill docs
- Pattern: Clear decision matrices, explicit boundaries

### Created Documents

1. `docs/templates/SYSTEM.md` (AI-facing)
2. `docs/guides/cron-natural-language.md` (User-facing)
3. 6 beads tasks for implementation

---

*This documentation ensures the AI knows exactly what it can and cannot do, preventing hallucination while providing a great user experience.*
