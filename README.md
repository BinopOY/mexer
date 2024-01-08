# Mexer

Mexer (or Mekseri for the finnish mouth) or mex-tester is a tool for checking exam packages for Abitti\*-exams. It validates exam integrity and authenticity and checks that provided unpacking codes are valid.

> Readme in finnish at `doc/LUEMINUT.md`

## Usage

Get the zip file including all the exams you want to check and input their unpacking codes to a single text file, seperated by new lines.

```bash
# Windows
mexer_amd64.exe <zip_file_path> <code_file_path>
# Linux
./mexer_amd64 <zip_file_path> <code_file_path>
```

## Examples

Example unpacking code file `codes.txt`:

```txt
annostaa syvyytys toukokuu panettaa
kittaus labiili kaatua asettelu
```

Example command:

```bash
# zip file at exams.zip and codes at codes.txt

mexer_amd64 ./exams.zip ./codes.txt
```
