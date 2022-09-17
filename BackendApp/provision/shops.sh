curl --location --request POST 'http://localhost:9192/shops' --header 'Authorization: Bearer token' --header 'Content-Type: application/json' --data-raw '{"name":"jikoni", "email":"jikoni@email.com", "number": "254710487831", "metadata": {"domain": "example.com"}}'
curl --location --request POST 'http://localhost:9192/shops' --header 'Authorization: Bearer token' --header 'Content-Type: application/json' --data-raw '{"name":"seasons", "email":"seasons@email.com", "number": "254710487832", "metadata": {"domain": "example.com"}}'
curl --location --request POST 'http://localhost:9192/shops' --header 'Authorization: Bearer token' --header 'Content-Type: application/json' --data-raw '{"name":"mess", "email":"mess@email.com", "number": "254710487833", "metadata": {"domain": "example.com"}}'
curl --location --request POST 'http://localhost:9192/shops' --header 'Authorization: Bearer token' --header 'Content-Type: application/json' --data-raw '{"name":"mathe", "email":"mathe@email.com", "number": "254710487834", "metadata": {"domain": "example.com"}}'
curl --location --request POST 'http://localhost:9192/shops' --header 'Authorization: Bearer token' --header 'Content-Type: application/json' --data-raw '{"name":"moriah", "email":"moriah@email.com", "number": "254710487835", "metadata": {"domain": "example.com"}}'


max=1000
for i in $(bash -c "echo {2..${max}}"); do curl --location --request POST 'http://localhost:9192/shops' --header 'Authorization: Bearer token' --header 'Content-Type: application/json' --data-raw '{"name":"shop'$i'", "email":"moriah'$i'@email.com", "number": "25471048783'$i'", "metadata": {"domain": "example.com"}}'&; done
