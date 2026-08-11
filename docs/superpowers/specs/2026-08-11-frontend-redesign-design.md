# Frontend Redesign Design Document

**Date**: 2026-08-11
**Author**: ZCode
**Status**: Approved

## Overview

Complete frontend redesign of Git Sync Service to achieve a modern SaaS-style UI using Ant Design Vue.

## Goals

1. **Modern SaaS Style**: Clean, professional UI similar to Notion, Linear, or Vercel
2. **Switch to Ant Design Vue**: Replace Element Plus with Ant Design Vue for better components and styling
3. **Visual Refresh Only**: Keep existing functionality, improve only the UI/UX
4. **Light with Blue Accents**: Clean white background with blue (#1677FF) accent color
5. **Collapsible Sidebar**: Sidebar navigation that can collapse to icon-only mode

## Design System

### Color Palette

| Color | Hex | Usage |
|-------|-----|-------|
| Primary Blue | #1677FF | Buttons, links, active states |
| Background | #F5F5F5 | Page background |
| Card Background | #FFFFFF | Cards, modals, content areas |
| Text Primary | #141414 | Main text |
| Text Secondary | #8C8C8C | Labels, captions |
| Border | #F0F0F0 | Borders, dividers |
| Success | #52C41A | Success states, active badges |
| Warning | #FAAD14 | Warning states |
| Error | #FF4D4F | Error states, delete actions |

### Typography

| Element | Size | Weight |
|---------|------|--------|
| Heading 1 | 24px | 600 |
| Heading 2 | 20px | 600 |
| Heading 3 | 16px | 600 |
| Body | 14px | 400 |
| Caption | 12px | 400 |

**Font Family**: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif

### Spacing System

| Name | Size |
|------|------|
| Base unit | 4px |
| Small | 8px |
| Medium | 16px |
| Large | 24px |
| X-Large | 32px |

### Shadows

**Card Shadow**:
```css
box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.03), 0 1px 6px -1px rgba(0, 0, 0, 0.02), 0 2px 4px 0 rgba(0, 0, 0, 0.02);
```

**Card Hover Shadow**:
```css
box-shadow: 0 6px 16px 0 rgba(0, 0, 0, 0.08);
```

**Modal Shadow**:
```css
box-shadow: 0 6px 16px 0 rgba(0, 0, 0, 0.08), 0 3px 6px -4px rgba(0, 0, 0, 0.12), 0 9px 28px 8px rgba(0, 0, 0, 0.05);
```

## Layout Structure

### Overall Layout

```
┌─────────────────────────────────────────────────────────────┐
│ Sidebar (Collapsible)    │         Header                   │
│ ┌─────────────────────┐  │  ┌─────────────────────────────┐ │
│ │ Logo                │  │  │ Breadcrumb    User Profile  │ │
│ │                     │  │  └─────────────────────────────┘ │
│ │ Dashboard           │  │  ┌─────────────────────────────┐ │
│ │ Sync Tasks          │  │  │                             │ │
│ │ Sync History        │  │  │                             │ │
│ │ Webhook Rules       │  │  │        Content Area         │ │
│ │ Webhook Events      │  │  │                             │ │
│ │ Repos               │  │  │                             │ │
│ │ Settings            │  │  │                             │ │
│ │                     │  │  │                             │ │
│ └─────────────────────┘  │  └─────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

### Sidebar Design

- **Width**: 240px (expanded) / 80px (collapsed)
- **Background**: #FFFFFF with subtle border-right (#F0F0F0)
- **Logo**: Git Sync icon + text (hidden when collapsed)
- **Navigation**: Icon + text, active state with blue highlight
- **Collapse button**: Bottom of sidebar

### Header Design

- **Height**: 56px
- **Background**: #FFFFFF with border-bottom (#F0F0F0)
- **Left**: Breadcrumb navigation
- **Right**: User avatar, notifications, settings

### Content Area

- **Padding**: 24px
- **Background**: #F5F5F5
- **Max width**: 1400px (centered)

## Component Design

### Cards

**Stat Cards**:
```css
.stat-card {
  background: #FFFFFF;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.03), 0 1px 6px -1px rgba(0, 0, 0, 0.02), 0 2px 4px 0 rgba(0, 0, 0, 0.02);
  transition: box-shadow 0.2s;
}

.stat-card:hover {
  box-shadow: 0 6px 16px 0 rgba(0, 0, 0, 0.08);
}
```

**Content Cards**:
```css
.content-card {
  background: #FFFFFF;
  border-radius: 8px;
  box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.03);
  overflow: hidden;
}

.card-header {
  padding: 16px 24px;
  border-bottom: 1px solid #F0F0F0;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-body {
  padding: 24px;
}
```

### Tables

- Use Ant Design Table component
- Striped rows for readability
- Hover effect on rows
- Action buttons in last column
- Pagination at bottom

### Forms

- Use Ant Design Form component
- Labels on top (not side)
- Validation messages below inputs
- Submit button aligned right

### Modals

- Use Ant Design Modal component
- 600px width for forms
- Footer with Cancel + Submit buttons
- Close button in top-right corner

### Buttons

- **Primary**: #1677FF (Ant Design blue)
- **Default**: White with border
- **Danger**: #FF4D4F
- **Ghost**: Transparent with border

### Status Badges

```css
.success { background: #F6FFED; color: #52C41A; }
.running { background: #E6F7FF; color: #1677FF; }
.failed  { background: #FFF2F0; color: #FF4D4F; }
.idle    { background: #F5F5F5; color: #8C8C8C; }
```

## Page-by-Page Redesign

### 1. Dashboard

**Current**: Basic stat cards + recent tasks/repos lists

**New Design**:
- 4 stat cards in a row (same as current, but with Ant Design Card)
- Recent tasks table (Ant Design Table)
- Recent repos table (Ant Design Table)
- Quick actions section

### 2. Sync Task List

**Current**: Stats row + basic table

**New Design**:
- Stats row (4 cards)
- Search + filter bar
- Ant Design Table with:
  - Task name (bold)
  - Source → Target branch (with arrow icon)
  - Sync mode badge
  - Status badge
  - Last run time (relative)
  - Actions (Run, Edit, Delete)

### 3. Sync History

**Current**: Basic list

**New Design**:
- Filter by task
- Ant Design Table with:
  - Task name
  - Trigger source
  - Status badge
  - Start/End time
  - Commit range
  - Details expandable row

### 4. Webhook Rules

**Current**: Basic table

**New Design**:
- Filter by repo
- Ant Design Table with:
  - Rule name
  - Repo key
  - Event type badge
  - Branch pattern
  - Status toggle
  - Actions

### 5. Webhook Events

**Current**: Basic list

**New Design**:
- Filter by repo, event type, status
- Ant Design Table with:
  - Event ID (truncated)
  - Event type badge
  - Source
  - Actor
  - Branch
  - Status badge
  - Processed time
  - Retry action

### 6. Repo List

**Current**: Card-based layout

**New Design**:
- Search bar
- Ant Design Table with:
  - Repo name + URL
  - Platform badge
  - Owner/Repo
  - Default branch
  - Status badge
  - Actions (Test, Edit, Delete)

### 7. Repo Config

**Current**: Form-based

**New Design**:
- Ant Design Form with sections:
  - Basic Info (name, URL, token)
  - Platform Config (platform, owner, repo)
  - Advanced (default branch, SSH URL)
- Save + Test Connection buttons

### 8. Settings

**Current**: Basic form

**New Design**:
- Ant Design Form with sections:
  - General Settings
  - Notification Settings
  - Advanced Settings
- Save button

## Technical Stack

| Component | Current | New |
|-----------|---------|-----|
| UI Library | Element Plus 2.5 | Ant Design Vue 4.x |
| State Management | Pinia 2.1 | Pinia 2.1 (no change) |
| Router | Vue Router 4.3 | Vue Router 4.3 (no change) |
| HTTP Client | Axios 1.6 | Axios 1.6 (no change) |
| Charts | ECharts 5.5 | ECharts 5.5 (no change) |
| Build Tool | Vite 5.1 | Vite 5.1 (no change) |
| CSS Preprocessor | Sass 1.72 | Sass 1.72 (no change) |

## Migration Steps

1. **Phase 1: Setup**
   - Install Ant Design Vue
   - Remove Element Plus
   - Update main.ts to register Ant Design Vue

2. **Phase 2: Layout**
   - Redesign AppLayout.vue with new sidebar
   - Add header with breadcrumb
   - Update responsive behavior

3. **Phase 3: Components**
   - Create shared components (StatusBadge, PageHeader, etc.)
   - Update all pages to use Ant Design components
   - Apply new design system

4. **Phase 4: Pages**
   - Dashboard
   - Sync Task List
   - Sync History
   - Webhook Rules
   - Webhook Events
   - Repo List
   - Repo Config
   - Settings

5. **Phase 5: Polish**
   - Add animations/transitions
   - Test responsive behavior
   - Fix any styling issues

## Success Criteria

1. All pages use Ant Design Vue components
2. Consistent design system across all pages
3. Collapsible sidebar working
4. All existing functionality preserved
5. No console errors
6. Responsive on desktop (1280px+)

## Risks and Mitigations

| Risk | Mitigation |
|------|------------|
| Ant Design Vue API differences | Use Ant Design Vue documentation |
| Styling conflicts | Remove all Element Plus styles |
| Missing components | Use Ant Design Vue equivalents |
| Performance issues | Optimize component imports |

## Timeline

- **Phase 1**: 2 hours
- **Phase 2**: 4 hours
- **Phase 3**: 4 hours
- **Phase 4**: 8 hours
- **Phase 5**: 2 hours
- **Total**: ~20 hours (2.5 days)
