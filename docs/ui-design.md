# UI Design Guidelines

## Button Color System

All buttons follow a semantic color hierarchy based on their purpose. This ensures visual consistency and helps users quickly understand the intent of each action.

### Color Tiers

| Tier | Color | Tailwind Classes | Usage |
|------|-------|------------------|-------|
| **Primary** | Blue | `bg-primary-600 hover:bg-primary-700` | Main actions: Create, Submit, Confirm (default) |
| **Secondary** | Stone (neutral) | `bg-stone-100 hover:bg-stone-200` or `bg-stone-800 hover:bg-stone-700` | Supporting actions: Cancel, Apply Filter, Back |
| **Danger** | Red | `bg-red-600 hover:bg-red-700` | Destructive actions: Delete |
| **Warning** | Amber | `bg-amber-600 hover:bg-amber-700` | Caution actions: Stop |

### Rules

- **Primary actions** (the most important action on a page) always use `primary-600`.
- **Only one primary button** should appear per logical section to maintain clear visual hierarchy.
- **Secondary buttons** use the `stone` palette and should not compete visually with the primary action.
- **Danger/Warning buttons** use semantic colors (`red`/`amber`) and should be reserved for destructive or cautionary operations.

### Common Patterns

```tsx
// Primary action
className="... text-white bg-primary-600 rounded-lg hover:bg-primary-700 transition-colors"

// Secondary action (light)
className="... text-stone-600 bg-stone-100 rounded-lg hover:bg-stone-200 transition-colors"

// Secondary action (dark, e.g., filter bar)
className="... text-white bg-stone-800 rounded-lg hover:bg-stone-700 transition-colors"

// Danger action
className="... text-white bg-red-600 rounded-lg hover:bg-red-700 transition-colors"

// Warning action
className="... text-white bg-amber-600 rounded-lg hover:bg-amber-700 transition-colors"
```

### ConfirmDialog Variants

The `ConfirmDialog` component supports three variants that map to this color system:

- `default` → `primary-600` (blue)
- `danger` → `red-600`
- `warning` → `amber-600`
