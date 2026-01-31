# Natural Language Scheduling Guide

> **Guide for**: Pryx Users  
> **Topic**: Creating scheduled tasks using natural language  
> **Version**: 1.0  

---

## Overview

Pryx understands natural language for scheduling tasks. Instead of memorizing cron syntax like `0 9 * * 1-5`, you can simply say:

```
"every weekday at 9am"
```

This guide explains what patterns work, what doesn't, and how to troubleshoot issues.

---

## Quick Reference

### Patterns That Work ✅

| What You Want | What You Say |
|---------------|--------------|
| Run every 5 minutes | `"every 5 minutes"` |
| Run every hour | `"every hour"` or `"every 60 minutes"` |
| Daily at 9 AM | `"every day at 9am"` |
| Daily at 2:30 PM | `"every day at 2:30pm"` or `"14:30"` |
| Weekdays only | `"every weekday at 8pm"` |
| Weekends only | `"every weekend at 10am"` |
| Specific day | `"every Monday at noon"` |
| Once in future | `"in 30 minutes"` or `"in 2 hours"` |
| Tomorrow | `"tomorrow at 9am"` |
| Next week | `"next Monday at 10am"` |
| Monthly | `"every month on the 1st at 9am"` |

### Patterns That Don't Work ❌

| What You Want | Why It Doesn't Work | Alternative |
|---------------|---------------------|-------------|
| `"every payday"` | Business logic, not time | `"every 15th at 9am"` |
| `"every 3-5 minutes"` | Variable interval | `"every 4 minutes"` |
| `"9am Tokyo time"` | Timezone in text | Use timezone selector |
| `"every day except weekends"` | Complex exceptions | Create two separate jobs |
| `"soon"` or `"later"` | Too vague | Specify exact time |
| `"if it's raining"` | Needs weather API | Check weather API manually |

---

## Detailed Examples

### 1. Fixed Intervals

**Every X minutes/hours/days:**

```
"every 5 minutes"           → Runs every 5 minutes
"every 30 minutes"          → Runs every 30 minutes  
"every 2 hours"             → Runs every 2 hours
"every day"                 → Runs every 24 hours
"every week"                → Runs every 7 days
```

**Use for**: Monitoring, health checks, periodic updates

---

### 2. Daily Schedules

**Specific time each day:**

```
"every day at 9am"          → 9:00 AM every day
"daily at 14:30"            → 2:30 PM every day
"at midnight"               → 12:00 AM every day
"at noon"                   → 12:00 PM every day
```

**Time formats supported:**
- **12-hour**: `9am`, `3:30pm`, `12pm`
- **24-hour**: `09:00`, `15:30`, `14:00`
- **Keywords**: `noon`, `midnight`

**Use for**: Daily reports, backups, morning briefings

---

### 3. Weekly Schedules

**Specific day(s) of the week:**

```
"every Monday"              → Every Monday at midnight
"every Monday at 9am"       → Every Monday at 9:00 AM
"every Tuesday at 3pm"      → Every Tuesday at 3:00 PM
"weekends at 10am"          → Saturday & Sunday at 10:00 AM
"weekdays at 8pm"           → Monday-Friday at 8:00 PM
```

**Day names:** Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday (case insensitive)

**Shortcuts:**
- `weekdays` = Monday-Friday
- `weekends` = Saturday-Sunday

**Use for**: Weekly reports, Monday morning standups, weekend maintenance

---

### 4. One-Time Schedules

**Future events:**

```
"in 30 minutes"             → 30 minutes from now
"in 2 hours"                → 2 hours from now
"in 1 day"                  → 24 hours from now
"tomorrow at 9am"           → Tomorrow at 9:00 AM
"tomorrow at noon"          → Tomorrow at 12:00 PM
"next Monday at 10am"       → Next Monday at 10:00 AM
```

**Use for**: Reminders, one-off tasks, future notifications

---

### 5. Monthly Schedules

**Day of month:**

```
"every month on the 1st"    → 1st of every month
"every month on the 15th at 9am" → 15th of every month at 9am
"monthly on the first"      → 1st of every month
```

**Use for**: Monthly reports, billing reminders, monthly maintenance

---

## Time Zone Handling

### Default Behavior

By default, schedules use your **local system time**.

### Specifying Timezones

If you need a specific timezone, use the timezone selector in the TUI or specify in standard cron syntax:

```bash
# Using TUI
Schedule: "every day at 9am"
Timezone: "America/New_York"

# This will run at 9am New York time, regardless of where you're located
```

**Common timezones:**
- `America/New_York` (Eastern)
- `America/Chicago` (Central)
- `America/Denver` (Mountain)
- `America/Los_Angeles` (Pacific)
- `Europe/London`
- `Europe/Paris`
- `Asia/Tokyo`
- `Asia/Singapore`
- `Australia/Sydney`

**Note**: Natural language like "9am Tokyo time" is **not supported**. Use the timezone selector instead.

---

## Common Use Cases

### Daily Morning Briefing

```
"every day at 8am"
Task: Summarize calendar and emails
```

### Weekly Report

```
"every Friday at 4pm"
Task: Generate weekly work summary
```

### System Monitoring

```
"every 5 minutes"
Task: Check system health and alert if issues
```

### Backup Reminder

```
"every day at 2am"
Task: Run backup script
```

### Standup Prep

```
"every weekday at 9:30am"
Task: Prepare standup notes
```

### Weekend Digest

```
"every Saturday at 10am"
Task: Weekly personal digest
```

---

## Troubleshooting

### "I don't understand that schedule"

**Problem**: Pryx couldn't parse your natural language pattern.

**Solution**: 
1. Check if your pattern matches the [supported patterns](#patterns-that-work-) above
2. Try rephrasing: `"every 2 hours"` instead of `"every few hours"`
3. Be more specific: `"tomorrow at 9am"` instead of `"tomorrow morning"`
4. Use standard cron syntax as fallback (see below)

---

### "Schedule is ambiguous"

**Problem**: Your pattern could mean multiple things.

**Example**: `"every morning"` → Is that 8am? 9am?

**Solution**: Specify the exact time:
- Instead of `"every morning"` → `"every day at 8am"`
- Instead of `"every evening"` → `"every day at 7pm"`

---

### "Pattern not supported"

**Problem**: You're trying to use business logic or complex conditions.

**Examples that won't work:**
- ❌ `"every payday"`
- ❌ `"if the stock market is up"`
- ❌ `"when I get an important email"`

**Solution**: Use fixed times:
- Instead of `"every payday"` → `"every 15th at 9am"` (if paid monthly)
- Instead of `"when I get an important email"` → `"every 15 minutes"` (check frequently)

---

## Fallback: Standard Cron Syntax

When natural language doesn't work, you can use standard 5-field cron syntax:

```
# Format: minute hour day month weekday
#         *     *    *   *     *

Examples:
0 9 * * *       → Every day at 9am
0 9 * * 1-5     → Weekdays at 9am
0 12 * * 0      → Every Sunday at noon
*/5 * * * *     → Every 5 minutes
0 0 1 * *       → First day of every month
```

**Field meanings:**
| Field | Values | Description |
|-------|--------|-------------|
| minute | 0-59 | Minute of the hour |
| hour | 0-23 | Hour of the day |
| day | 1-31 | Day of the month |
| month | 1-12 | Month of the year |
| weekday | 0-6 | Day of week (0=Sunday, 6=Saturday) |

**Special characters:**
| Char | Meaning | Example |
|------|---------|---------|
| `*` | Any value | `*` in hour = every hour |
| `-` | Range | `1-5` = Monday to Friday |
| `,` | List | `1,3,5` = Mon, Wed, Fri |
| `/` | Step | `*/5` = every 5 units |

---

## Advanced: Combining Schedules

### Multiple Times per Day

```
# Create two separate jobs:
Job 1: "every day at 9am"
Job 2: "every day at 5pm"
```

### Different Schedules for Different Days

```
# Create separate jobs:
Job 1: "weekdays at 9am"      → Monday-Friday
Job 2: "weekends at 10am"     → Saturday-Sunday
```

### Complex Patterns with Cron

```
# Every 15 minutes during business hours (9am-5pm)
*/15 9-17 * * 1-5

# First Monday of every month
0 9 1-7 * 1
```

---

## Best Practices

### 1. Be Specific

❌ **Vague**: `"every morning"`
✅ **Specific**: `"every day at 8am"`

### 2. Use Timezones for Remote Teams

If your team is distributed, always specify timezone:
```
Schedule: "every day at 9am"
Timezone: "America/New_York"
```

### 3. Test Your Schedule

Before relying on a schedule, verify the next run time:
```
$ pryx cron test "every weekday at 9am"
Next run: Monday at 9:00 AM
Subsequent: Tuesday at 9:00 AM, Wednesday at 9:00 AM...
```

### 4. Start Simple

If natural language doesn't work, start with simple patterns:
```
"every day at 9am"          ✓ Works
"every day at 9am except holidays"  ✗ Too complex
```

### 5. Document Your Jobs

Give jobs descriptive names:
```
❌ "cron-job-1"
✅ "Daily morning email summary"
```

---

## Reference: Day Names

| Full Name | Short | Cron Number |
|-----------|-------|-------------|
| Sunday | Sun | 0 |
| Monday | Mon | 1 |
| Tuesday | Tue | 2 |
| Wednesday | Wed | 3 |
| Thursday | Thu | 4 |
| Friday | Fri | 5 |
| Saturday | Sat | 6 |

---

## Getting Help

### In Pryx

Type:
```
/help cron
```

Or ask:
```
"What schedule patterns do you support?"
```

### Documentation

- **This guide**: `docs/guides/cron-natural-language.md`
- **System reference**: `docs/templates/SYSTEM.md`
- **Implementation**: `pkg/cron/parser/`

---

## Summary

**Remember:**

1. ✅ Use simple, specific time expressions
2. ✅ Stick to documented patterns
3. ❌ Don't use business logic or complex conditions
4. ❌ Don't use vague times like "soon" or "later"
5. 🔧 Use standard cron syntax as fallback

**When in doubt:**
- Start with `"every day at [time]"`
- Test with `pryx cron test "your schedule"`
- Ask Pryx: "Can you parse this schedule: [your schedule]?"

---

*Last updated: 2026-01-29*
