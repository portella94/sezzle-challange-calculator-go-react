# Prompt: Engineering-Quality Review

Review this repository as a senior engineer doing a design/code-quality
review before merge. Judge the codebase the way an experienced reviewer would
against general industry standards for the languages and paradigms in use
(idiomatic Go per *Effective Go* and the Google Go Style Guide; idiomatic
React/TypeScript; REST API design.
Sound architecture, sound object/module design, and sound API design should
fall out of that lens naturally, rather than being scored as separate
checkboxes.

## Scope

Review:
- Backend: package layout, layering, dependency direction, error handling,
  concurrency/composition root.
- Frontend: component/module structure, state and data-fetching patterns,
  type usage.
- The HTTP API surface: resource/endpoint design, status codes, error
  shape, versioning.
- Cross-cutting: naming, duplication, configuration, testability.

## What to produce

1. **Overall impression** — 2-3 sentences on the codebase's engineering
   maturity and whether the architecture fits the problem size.
2. **Strengths** — what is genuinely well-designed, and why it works
   (cite the file/pattern).
3. **Issues**, grouped by theme (not by principle name), each with:
   - Where it is (file/function).
   - What is actually wrong in concrete terms (not "violates DRY" but,
     e.g., "this validation logic is duplicated in X and Y and will drift").
   - Why it matters in practice (maintenance cost, coupling, surprise
     behavior, risk).
   - A concrete fix.
   - A severity: blocking / worth fixing / nice-to-have.
4. **API design notes** — anything a consumer of this API would find
   inconsistent or non-standard versus common REST conventions.
5. **One next step** — the single highest-leverage change to make first.

## Constraints

- Ground every finding in a specific file, function, or endpoint — no
  generic advice.
- Where the code already follows a sound pattern, say so explicitly rather
  than staying silent; don't manufacture issues to fill sections.
- Prefer showing a short before/after snippet over describing a fix in the
  abstract, when it clarifies the point.
- Keep the review scannable: short paragraphs, no filler.
