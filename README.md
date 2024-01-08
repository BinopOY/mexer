# Mexer

Mexer (or Mekseri for the finnish mouth) or mex-tester is a tool for checking exam packages for Abitti\*-exams. It validates exam integrity and authenticity and checks that provided unpacking codes are valid.

> Readme in finnish at `doc/LUEMINUT.md`

> \* Abitti is a registered EU trademark (015833742, 015838915) of the Finnish Matriculation Examination Board.

## Installation

Download the latest executable for windows or linux from https://github.com/BinopOY/mexer/releases/latest.

## Usage

What counts as an exam package:

-   zip file with multiple exams
-   zip file with a single exams
-   mex file as a single exam

What counts as unpacking code input:

-   Single code as a command line argument
-   `*.txt`-file containing a list of codes separated with newline

There should be at least as many unpacking codes as there are exams to check.

```bash
# Windows
mexer_amd64.exe <exam_package_path> <code_input>
# Linux
./mexer_amd64 <exam_package_path> <code_input>
```

## Examples

Example unpacking code file `codes.txt`:

```txt
annostaa syvyytys toukokuu panettaa
kittaus labiili kaatua asettelu
```

Example commands:

```bash
# zip file at exams.zip and codes at codes.txt
mexer_amd64 ./exams.zip ./codes.txt

# Single exam in zip and codes as input
mexer_amd64 ./exams.zip "kittaus labiili kaatua asettelu"

# Single exam as mex and codes as input
mexer_amd64 ./exams.mex "kittaus labiili kaatua asettelu"

# Single exam as mex and codes at codes.txt
mexer_amd64 ./exams.mex ./codes.txt
```
