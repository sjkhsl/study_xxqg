git pull
go mod tidy
go build ./
nohup ./study_xxqg > study_xxqg.log 2>&1 & echo $!>pid.pid