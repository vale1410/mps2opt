#/bin/zsh

go build 

echo success
for x in MIPLIB/new*/*; 
do 
    ./mps2opt -lp -f $x >/dev/null; 
    r=$?
    if [ $r = 0 ]; then 
        echo \* $(basename $x .mps)
    fi 
done

echo fail
for x in MIPLIB/new*/*; 
do 
    ./mps2opt -lp -f $x >/dev/null; 
    r=$?
    if [ $r = 100 ]; then 
        echo $x 
        ./mps2opt -lp -f $x; 
    fi 
done

