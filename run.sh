#!/bin/zsh

npx tailwindcss -i view/css/app.css -o public/styles.css --watch &
PID1=$!

wgo -file=.go -file=.templ -file=.js -file=.css -xfile=_templ.go go run -tags dev . &
PID2=$!

templ generate --watch --proxy=http://localhost:42069 &
PID3=$!

cleanup() {
  echo "Stopping all resources..."
  kill $PID1 $PID2 $PID3
}

trap cleanup SIGINT

wait $PID1 $PID2 $PID3