#!/bin/zsh

npx tailwindcss -i view/css/app.css -o public/styles.css --watch &
PID1=$!
echo "tailwindcss started with PID: $PID1"

(nohup wgo -file=.go -file=.templ -file=.js -file=.css -xfile=_templ.go go run -tags dev .) &
PID2=$!
echo "go run started with PID: $PID2"

templ generate --watch --proxy=http://localhost:42069 &
PID3=$!
echo "templ generate started with PID: $PID3"



cleanup() {
  for pid in $PID3 $PID2 $PID1; do
    if kill -0 $pid 2>/dev/null; then
      echo "Stopping $pid..."
      kill $pid || echo "Failed to stop process $pid"
    else
      echo "Process $pid already stopped."
    fi
  done
}

trap cleanup SIGINT

wait $PID3 $PID2 $PID1

osascript -e 'tell application "Google Chrome" to close (tabs of window 1 whose URL contains "http://127.0.0.1:7331")'