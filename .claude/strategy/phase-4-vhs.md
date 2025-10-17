# Phase 4: Documentation & Demo - VHS Integration

**Duration**: 1 day | **Complexity**: Low | **Risk**: Minimal | **Status**: Not Started

## Overview
Add VHS (Video HeadShot) scripts for capturing TUI demos and create animated GIFs for the README to showcase Bento's capabilities.

## Goals
- Create VHS tape scripts for key workflows
- Generate animated GIFs for README
- Document key use cases visually
- Make it easy to regenerate demos

## Related Original TODOs
- TODO #5: VHS Charm functionality

## Demo Scripts to Create

### 1. demo-overview.tape
- Show main interface
- Navigate between tabs
- Browse existing bentos
- Quick tour of all features

### 2. demo-create-bento.tape
- Create new bento from scratch
- Add multiple nodes (HTTP, shell, etc.)
- Configure node parameters
- Save and execute bento

### 3. demo-editor.tape
- Show new editor interface (if Phase 3 complete)
- Navigate table view
- Edit node with form
- Add new node

### 4. demo-recipes.tape
- Browse example bentos (Recipes tab)
- View example details
- Copy example to own collection
- Customize and run

## Implementation Details

### Directory Structure
```
.vhs/
├── README.md              - Instructions
├── demo-overview.tape     - Main overview
├── demo-create-bento.tape - Creating bentos
├── demo-editor.tape       - Editor features
└── demo-recipes.tape      - Example bentos
```

### VHS Script Template
```tape
Output demo-overview.gif

# Settings
Set FontSize 14
Set Width 1200
Set Height 800
Set Theme "Catppuccin Mocha"
Set Shell "bash"
Set PlaybackSpeed 1.0

# Demo
Type "bento taste"
Sleep 2s
Enter
Sleep 1s

# Navigate tabs
Type "1"  # Bentos tab
Sleep 1s
Type "2"  # Recipes tab
Sleep 1s
Type "3"  # Mise tab
Sleep 1s
Type "4"  # Sensei tab
Sleep 1s

# Screenshot
Screenshot
```

### VHS Configuration
- **Terminal Size**: 1200x800 (readable in README)
- **Font Size**: 14-16 (clear text)
- **Theme**: Match Bento theme (Catppuccin Mocha or similar)
- **Timing**: Realistic pauses for readability
- **Length**: 5-15 seconds per demo

## Files to Create

### .vhs/README.md
```markdown
# Bento TUI Demos

This directory contains VHS tape scripts for generating animated GIF demos.

## Prerequisites
- Install VHS: `brew install vhs` (macOS) or see https://github.com/charmbracelet/vhs

## Generating GIFs
```bash
# Generate all demos
vhs demo-overview.tape
vhs demo-create-bento.tape
vhs demo-editor.tape
vhs demo-recipes.tape

# Or use make target
make demos
```

## Demos
- `demo-overview.tape` - Main interface overview
- `demo-create-bento.tape` - Creating a new bento
- `demo-editor.tape` - Editor features
- `demo-recipes.tape` - Using example bentos
```

### Files to Modify
- `README.md` - Add animated GIF demos
- `Makefile` (optional) - Add demo generation target

## README Integration

### Demo Section
```markdown
## Demos

### Creating a Bento
![Create Demo](.vhs/demo-create-bento.gif)

### Browsing Examples
![Recipes Demo](.vhs/demo-recipes.tape.gif)

### Editing Nodes
![Editor Demo](.vhs/demo-editor.gif)
```

## Dependencies
- VHS CLI tool (development dependency only)
- Working Bento installation

## Testing Requirements
- [ ] All tape scripts run without errors
- [ ] GIFs are clear and readable
- [ ] Demos show key features
- [ ] Timing is appropriate
- [ ] Theme matches Bento style

## Success Criteria
- [ ] .vhs/ directory created
- [ ] All 4 tape scripts created
- [ ] VHS README documentation written
- [ ] All GIFs generated successfully
- [ ] Main README updated with demos
- [ ] GIFs committed to repo
- [ ] Regeneration instructions clear

## References
- https://github.com/charmbracelet/vhs
- https://github.com/charmbracelet/vhs/tree/main/examples

## Notes
- This is a documentation-only phase
- No code changes to Bento itself
- Can be done in parallel with Phase 5
- VHS is only needed for development, not users

---

## Claude Code Prompt for Guilliman

```
Implement Phase 4 of the Bento TUI enhancement plan: Documentation & Demo - VHS Integration.

IMPORTANT: Read the phase document at .claude/strategy/phase-4-vhs.md for complete context, requirements, and success criteria before starting implementation.

Requirements:
1. Create .vhs/ directory with VHS tape scripts:

   a. demo-overview.tape:
      - Show main tab navigation
      - Browse existing bentos
      - Quick tour of all tabs

   b. demo-create-bento.tape:
      - Create new bento from scratch
      - Add multiple nodes
      - Save and run

   c. demo-editor.tape:
      - Show new editor interface (if Phase 3 complete)
      - Table navigation
      - Form editing

   d. demo-recipes.tape:
      - Browse example bentos
      - Copy example to own collection
      - Customize and run

2. VHS script features:
   - Set appropriate terminal size (1200x800)
   - Use readable font size (14-16)
   - Include pauses for readability
   - Add typed commands with realistic timing
   - Capture key interactions

3. Documentation:
   - Create .vhs/README.md with instructions
   - Document how to generate GIFs
   - List prerequisites (VHS installation)
   - Provide regeneration commands

4. Update main README.md:
   - Add animated GIF demos
   - Show key features visually
   - Link to detailed documentation

5. Reference:
   - https://github.com/charmbracelet/vhs

Example tape structure:
```tape
Output demo-overview.gif
Set FontSize 14
Set Width 1200
Set Height 800
Set Theme "Catppuccin Mocha"

Type "bento taste"
Sleep 2s
Enter
Sleep 1s
Type "1"  # Switch to Bentos tab
Sleep 1s
Screenshot
```

6. Generate all GIFs and verify:
   - Clear and readable
   - Show key features
   - Appropriate length (5-15s each)
   - High quality rendering

Return a summary of:
- Tape scripts created
- GIFs generated
- README updates made
- Instructions for regenerating demos
```
