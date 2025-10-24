# UI Style Guide

This UI layer is built from a small set of composable building blocks. Reuse these classes whenever possible before introducing new selectors.

## Design Tokens

- Colors, spacing, radii, and typography live in CSS custom properties at the top of `static/app.css`.
- Prefer using the variables (e.g. `var(--color-primary)`) instead of hard-coding new values to keep the palette consistent.

## Layout

- `ui-shell`, `ui-shell__sidebar`, `ui-shell__main`, `ui-shell__content` define the page frame.
- Use `panel-grid` for responsive column layouts (auto-fit grid with consistent gaps).

## Panels and Lists

- Wrap feature sections in `panel`. Add `panel--muted` when you need a lower contrast surface.
- Use `panel__header` with `panel__title` and `panel__actions` for section headers.
- Lists that belong to a panel should use `panel__list` (a gap-based, bulletless list).

## Cards and Collapsibles

- `card` is the base surface. Variants:
  - `card--collapsible` pairs with `<details class="collapsible">`.
  - `card--compact` trims padding for dense content (waiting tasks, etc.).
  - `card--form` is used for workflow step editors.
- Inside `details`, use:
  - `collapsible__summary` for the clickable row.
  - `collapsible__title` for the main label (flexes to fill).
  - `collapsible__content` for the expanded body.
- For inline lists inside a card use `data-list` so status badges align to the right automatically.

## Buttons

- Always start with the base `button` class and layer modifiers:
  - `button--accent` (primary green) and `button--primary` (blue action).
  - `button--ghost` (neutral outline), `button--outline` (accent outline).
  - `button--danger` for destructive emphasis, `button--text` to render as an inline link.
  - `button--icon` creates a pill-shaped square (used for the sidebar “Create” button).
  - `button--sm` tightens padding for compact actions.
- Combine modifiers as needed (e.g. `button button--ghost button--sm`).

## Status Badges

- Use `status-badge` to display state chips. Apply one of:
  - `status-badge--ready`, `--succeeded`, `--failed`, `--pending`, `--unknown`.
- The Go view model helpers already emit these modifier classes.

## Forms

- Group inputs with `form__field`; helper text goes in `form__hint`.
- `form__actions` right-aligns submit/cancel controls and keeps spacing consistent.
- Read-only inputs rely on the `[readonly]` selector—no extra class is needed.

## Modals

- Structure: `modal` (backdrop) → `modal__dialog` → `modal__header`/`modal__body`.
- Close buttons use `button button--ghost button--icon modal__close` and `data-close-modal`.
- Error messaging inside modals reuses `alert alert--error` with the `is-visible` toggle.

## Waiting Queue

- Waiting sections stay inside `panel panel--muted`.
- Group headers use `waiting-group__header` with `waiting-group__divider` for the rule.
- Tasks render as `card card--compact waiting-task`; optional prompt content uses `waiting-task__prompt`.

## Workflow Step Builder

- The step list is a grid (`#workflow-steps`) of `card card--form workflow-step`.
- `workflow-steps__header` controls the “Steps” title row and add button alignment.
- `workflow-step__header` contains the per-step title and remove action (`button--text button--danger`).

## General Guidelines

- Prefer combining existing modifiers over creating new bespoke classes.
- If a new style truly cannot reuse existing pieces, document it here when you add the CSS.
- Keep HTML markup semantic (`<section>`, `<header>`, `<article>`), the utility classes supply layout and visuals.
