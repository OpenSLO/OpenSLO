#!/usr/bin/awk -f

BEGIN {
    FS = ":"
    printf "Targets:\n"
}

!/^##/ && f > 0 {
    printf "    \033[36m%-30s\033[0m  %s\n", $1, r
    f = 0
    r = ""
}
/^##/ {
    f = f + 1
    sub(/^## /,"",$0)
}
f == 1  { r = $0 }
f > 1   { r = r " " $0 }