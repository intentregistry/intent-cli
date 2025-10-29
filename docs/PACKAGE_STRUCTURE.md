# Understanding the Package Structure

## What Gets Created Where

### Project Structure (what you should have):
```
example3/
├── itpkg.json          # ✅ Project manifest (stays in project)
├── intents/           # ✅ Your intent files go here
│   └── hello.itml
├── policies/          # ✅ Your policy files go here
│   └── ...
└── dist/              # ✅ Output directory (you create this)
    └── scope-example3-0.1.0.itpkg  # ✅ Generated package file
```

### Generated Files:
- **`itpkg.json`**: Created in project root (stays there, part of your project)
- **`scope-example3-0.1.0.itpkg`**: Created in `dist/` (or current dir if no --out)

### Required Directories:
- **`intents/`**: Your `.itml` intent files go here
- **`policies/`**: Your policy files go here

## Correct Workflow

```bash
# Step 1: Create project structure and manifest
cd example3
intent package . --scaffold --unsigned --out dist/

# This creates:
# - itpkg.json in current directory (project root)
# - intents/ directory (if missing)
# - policies/ directory (if missing)  
# - dist/scope-example3-0.1.0.itpkg (the package file)

# Step 2: Move your files into the right places
mv hello.itml intents/  # if it's an intent file
mv *.itml policies/      # if they're policy files

# Step 3: Rebuild package with your files included
intent package . --out dist/
```

## What's Inside the .itpkg File?

When you extract the `.itpkg` file, it contains:
```
scope-example3-0.1.0.itpkg (extracted):
├── itpkg.json          # Manifest
├── MANIFEST.sha256     # Checksums
├── SIGNATURE           # Signature
├── intents/            # Your intent files
│   └── hello.itml
├── policies/           # Your policy files
│   └── ...
└── [other project files]
```

## Summary

- **`itpkg.json`**: Stays in project root (version control this!)
- **`intents/` and `policies/`**: Project directories (version control these!)
- **`.itpkg` file**: Goes to `dist/` (don't version control, it's a build artifact)
- **`dist/`**: Create this directory to organize your built packages

You should commit `itpkg.json`, `intents/`, and `policies/` to git, but NOT `dist/` or `.itpkg` files (add them to `.gitignore`).

