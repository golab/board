if [[ "$1" == "" ]]
then
    echo "no input"
    exit
fi

echo curl -X POST http://localhost:8080/apps/twitch/callback -d "{\"event\": {\"message\": {\"text\": \"$1\"}}}"

curl -X POST http://localhost:8080/apps/twitch/callback -d "{\"event\": {\"message\": {\"text\": \"$1\"}}}"

