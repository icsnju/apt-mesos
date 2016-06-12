for i in $(seq $1);do
  curl -X POST  -H 'content-type: application/json'  -d @job2.json http://114.212.189.126:3030/api/jobs
done
