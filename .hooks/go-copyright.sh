#!/bin/bash

set -ex

copyright_header='/*
Copyright © 2023 Jayson Grace <jayson.e.grace@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/'

while IFS= read -r -d '' file; do
    if ! grep -qF "$copyright_header" "$file"; then
        echo "Adding copyright header to ${file}"
        temp_file=$(mktemp)
        echo "${copyright_header}" > "${temp_file}"
        # Add an empty line after the copyright header
        echo "" >> "${temp_file}"
        cat "${file}" >> "${temp_file}"
        mv "${temp_file}" "${file}"
    fi
done < <(find . -type f -name "*.go" -print0)
