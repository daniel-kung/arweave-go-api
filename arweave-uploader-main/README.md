API 


Deploy a file 
~~~~
curl -X POST localhost:8080/api/v1/images \
--form file=@/Users/Downloads/16.png
~~~~
Deploy a uri
~~~~
curl -X POST localhost:8080/api/v1/images/uri \
-H'Content-Type: application/json' \
-d '{
    "uri":"https://cdn.zmescience.com/wp-content/uploads/2020/11/sea-5213746_1920-1140x760.jpg"
}'