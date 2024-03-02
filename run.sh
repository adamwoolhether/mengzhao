#!/bin/zsh

ngrok http 127.0.0.1:42069 --host-header=mengzhao.test > /dev/null 2>&1 &
PID0=$!
printf "ngrok\t\tPID: %s\n" $PID0

get_public_url() {
    start_time=$(date +%s)
    timeout=5

    while true; do
        current_time=$(date +%s)
        elapsed_time=$((current_time - start_time))

        curl_response=$(curl -s http://127.0.0.1:4040/api/tunnels)
        if [ "$curl_response" != "" ]; then
            public_url=$(echo "$curl_response" | jq -r '.tunnels[0].public_url')
            if [ "$public_url" != "null" ]; then
                echo "$public_url"
                break
            fi
        fi
        sleep 1
    done
}

PUBLIC_URL=$(get_public_url)
printf "ngrok URL:\t%s\n" "$PUBLIC_URL"

export WEBHOOK_URL=$PUBLIC_URL/replicate/callback
sed -i '' -e "s|WEBHOOK_URL=.*|WEBHOOK_URL=${MY_ENV}|" .env

npx tailwindcss -i view/css/app.css -o public/styles.css --watch &
PID1=$!
printf "tailwindcss\tPID: %s\n" $PID1

wgo -file=.go -file=.templ -file=.js -file=.css -xfile=_templ.go go run -tags dev . &
PID2=$!
printf "wgo\t\tPID: %s\n" $PID2

templ generate --watch --proxy=http://localhost:42069 &
PID3=$!
printf "templ\t\tPID: %s\n" $PID3

cleanup() {
  for pid in $PID3 $PID2 $PID1 $PID0; do
    if kill -0 $pid 2>/dev/null; then
      printf "PID\t\t%s stopping...\n" $pid
      kill $pid || printf "PID\t\t%s not stopped\n" $pid
    else
      printf "PID\t\t%s stopped\n" $pid
    fi

    # ensure wgo is stopped
    pkill wgo
  done
}

trap cleanup SIGINT

wait $PID3 $PID2 $PID1 $PID0

osascript -e 'tell application "Google Chrome" to close (tabs of window 1 whose URL contains "http://127.0.0.1:7331/")'