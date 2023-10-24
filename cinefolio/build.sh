npm install .
npm run build-github &
pid=$!
sleep 20

if ps -p $pid > /dev/null; then
    kill -9 $pid
    exit 0  # Return error code 0
else
    wait $pid
    exit $?
fi
