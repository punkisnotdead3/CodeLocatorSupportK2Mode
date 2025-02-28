# CodeLocatorSupportK2Mode

# IDEA Plugin K2 Adapter

A simple tool to repackage IntelliJ IDEA plugins for K2 mode compatibility.

## Build

Use make command to build executables for all platforms:

```bash
make
```

This will generate executables for:
- Windows: `k2adapter.exe`
- Linux: `k2adapter`
- macOS: `k2adapter_mac`

## Usage

1. Choose the appropriate executable for your platform
2. Run the command with your plugin zip file:

For Windows:
```bash
k2adapter.exe -zip "path/to/your/plugin.zip"
```

For Linux:
```bash
./k2adapter -zip "path/to/your/plugin.zip"
```

For macOS:
```bash
./k2adapter_mac -zip "path/to/your/plugin.zip"
```

3. After execution, a new file named `plugin-k2.zip` will be generated in the same directory as your original plugin

## Example

```bash
# Windows
k2adapter.exe -zip "C:\Downloads\my-plugin.zip"

# Linux/macOS
./k2adapter -zip "./my-plugin.zip"
```

## Output

The tool will create a new zip file with "-k2" suffix. For example:
- Original file: `my-plugin.zip`
- Generated file: `my-plugin-k2.zip`

The generated file is ready to use with IDEA's K2 mode.

## Note

Make sure you have the correct permissions to read the input file and write to the output directory.