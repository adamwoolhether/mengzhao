#!/bin/zsh

ngrok http 127.0.0.1:42069 --host-header=mengzhao.test > /dev/null &
PID0=$!
sleep 1
if ! kill -0 $PID0 > /dev/null 2>&1; then
    echo "ngrok failed to start. Exiting..."
    exit 1
else
    printf "ngrok\t\tPID: %s\n" $PID0
fi

get_public_url() {
    start_time=$(date +%s)
    timeout=5

    while true; do
        current_time=$(date +%s)
        elapsed_time=$((current_time - start_time))

        if [ "$elapsed_time" -ge "$timeout" ]; then
            echo "Timeout reached without obtaining a public URL."
            return 1
        fi

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
if [ $? -ne 0 ]; then
    echo "Exiting due to failure in get_public_url" >&2
    exit 1
fi
printf "ngrok URL:\t%s\n" "$PUBLIC_URL"

export WEBHOOK_URL=$PUBLIC_URL/replicate/callback

wgo -file=.go -file=.templ -file=.js -xfile=_templ.go npx tailwindcss -i view/css/app.css -o public/styles.css :: go run -tags dev . &
PID1=$!
printf "wgo\t\tPID: %s\n" $PID1

templ generate --watch --proxy=http://localhost:42069 &
PID3=$!
printf "templ\t\tPID: %s\n" $PID3

cleanup() {
  for pid in $PID3 $PID1 $PID0; do
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

wait $PID3 $PID1 $PID0

osascript -e 'tell application "Google Chrome" to close (tabs of window 1 whose URL contains "http://127.0.0.1:7331/")'