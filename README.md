# sal-backend

Mac OS 
in sal/sal-backend
mkdir src
cd src
ln -s /Users/preeti.dave/sal/sal-backend/aws/ aws
GOOS=LINUX go build -o github.com/build/user github.com/sal-group/sal-backend/user/user.go
zip build/user and upload to lambda functions code

To Test: 
curl --request GET https://sf7z9xf5i9.execute-api.us-east-2.amazonaws.com/dev/users
curl --request GET https://sf7z9xf5i9.execute-api.us-east-2.amazonaws.com/dev/counselors
curl --request GET https://sf7z9xf5i9.execute-api.us-east-2.amazonaws.com/dev/appointments


