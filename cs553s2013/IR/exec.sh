go test > hello.spim
line=`wc -l hello.spim | cut -f 1 -d " "`
line1=` expr $line - 2 `
head -n $line1 hello.spim > .hello.spim
mv .hello.spim hello.spim
spim -file hello.spim
