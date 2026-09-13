# Rules for agents

1. NEVER use abbreviations in variable, function, or type names (except for conventional ones like `i` in loops). Names must be descriptive and unabbreviated for clarity and easier refactoring.

2. Always split large markup or complex UI pages into small, descriptive, self-contained subcomponents under the `components/` directory, extracting logical blocks into reusable components or composables as needed.

3. NEVER use raw pixel ("px") units in inline styles or style attributes. Always utilize standard Tailwind CSS spacing and sizing utility classes (e.g., `p-4`, `h-4`, `w-4`, `gap-2`, etc.) to build layouts and maintain visual consistency.
