curl --location --request POST 'http://localhost:9191/orders' --header 'Authorization: Bearer token' --header 'Content-Type: application/json' --data-raw '{"name":"order1", "vendor":"01GD6E46KKHZYQDYKYSCHHN2CY", "price": 100,"status": "ordered", "place": "inhouse", "metadata": {"domain": "example.com"}}'
curl --location --request POST 'http://localhost:9191/orders' --header 'Authorization: Bearer token' --header 'Content-Type: application/json' --data-raw '{"name":"order2", "vendor":"01GD6E46KKHZYQDYKYSCHHN2CY", "price": 150,"status": "ordered", "place": "delivery", "metadata": {"domain": "example.com"}}'
curl --location --request POST 'http://localhost:9191/orders' --header 'Authorization: Bearer token' --header 'Content-Type: application/json' --data-raw '{"name":"order3", "vendor":"01GD6E46N4FS8SB8PSRRHPZ9TJ", "price": 150,"status": "paid", "place": "delivery", "metadata": {"domain": "example.com"}}'
curl --location --request POST 'http://localhost:9191/orders' --header 'Authorization: Bearer token' --header 'Content-Type: application/json' --data-raw '{"name":"order4", "vendor":"01GD6E46N4FS8SB8PSRRHPZ9TJ", "price": 300,"status": "paid", "place": "delivery", "metadata": {"domain": "example.com"}}'
curl --location --request POST 'http://localhost:9191/orders' --header 'Authorization: Bearer token' --header 'Content-Type: application/json' --data-raw '{"name":"order5", "vendor":"01GD6E46N4FS8SB8PSRRHPZ9TJ", "price": 300,"status": "ordered", "place": "delivery", "metadata": {"domain": "example.com"}}'
curl --location --request POST 'http://localhost:9191/orders' --header 'Authorization: Bearer token' --header 'Content-Type: application/json' --data-raw '{"name":"order6", "vendor":"01GD6E46KKHZYQDYKYSCHHN2CY", "price": 200,"status": "ordered", "place": "delivery", "metadata": {"domain": "example.com"}}'
curl --location --request POST 'http://localhost:9191/orders' --header 'Authorization: Bearer token' --header 'Content-Type: application/json' --data-raw '{"name":"order7", "vendor":"01GD6E46KKHZYQDYKYSCHHN2CY", "price": 200,"status": "paid", "place": "inhouse", "metadata": {"domain": "example.com"}}'
curl --location --request POST 'http://localhost:9191/orders' --header 'Authorization: Bearer token' --header 'Content-Type: application/json' --data-raw '{"name":"order8", "vendor":"01GD6E46P13RH2TKVKN2N2Y6GN", "price": 100,"status": "paid", "place": "inhouse", "metadata": {"domain": "example.com"}}'
curl --location --request POST 'http://localhost:9191/orders' --header 'Authorization: Bearer token' --header 'Content-Type: application/json' --data-raw '{"name":"order9", "vendor":"01GD6E46P13RH2TKVKN2N2Y6GN", "price": 70,"status": "paid", "place": "inhouse", "metadata": {"domain": "example.com"}}'
curl --location --request POST 'http://localhost:9191/orders' --header 'Authorization: Bearer token' --header 'Content-Type: application/json' --data-raw '{"name":"order10", "vendor":"01GD6E46P13RH2TKVKN2N2Y6GN", "price": 20,"status": "paid", "place": "inhouse", "metadata": {"domain": "example.com"}}'



max=1000
for i in $(bash -c "echo {2..${max}}"); do curl --location --request POST 'http://localhost:9191/orders' --header 'Authorization: Bearer token' --header 'Content-Type: application/json' --data-raw '{"name":"order10", "vendor":"01GD6E46P13RH2TKVKN2N2Y6GB'$i'", "price": 20,"status": "paid", "place": "inhouse", "metadata": {"domain": "example.com"}}'&; done
