if [[ "$1" == "build" ]]
then
    rm build/*
    go build -C frontend -o ../build/
    go build -C backend -o ../build/
elif [[ "$1" == "frontend" ]]
then
    rm build/frontend
    go build -C frontend -o ../build/
    cd frontend && ../build/frontend
elif [[ "$1" == "backend" ]]
then
    rm build/backend
    go build -C backend -o ../build/
    ./build/backend
elif [[ "$1" == "all" ]]
then
    rm build/*
    go build -C frontend -o ../build/
    go build -C backend -o ../build/
    ./build/backend & cd frontend && ../build/frontend
fi
