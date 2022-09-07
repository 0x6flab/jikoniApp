curl --location --request POST 'http://localhost:9191/orders' --header 'Authorization: Bearer token' --header 'Content-Type: application/json' --data-raw '{"name":"order1", "vendor":"jikoni", "price": 100,"status": "ordered", "place": "inhouse", "metadata": {"domain": "example.com"}}'
curl --location --request POST 'http://localhost:9191/orders' --header 'Authorization: Bearer token' --header 'Content-Type: application/json' --data-raw '{"name":"order2", "vendor":"jikoni", "price": 150,"status": "ordered", "place": "delivery", "metadata": {"domain": "example.com"}}'
curl --location --request POST 'http://localhost:9191/orders' --header 'Authorization: Bearer token' --header 'Content-Type: application/json' --data-raw '{"name":"order3", "vendor":"seasons", "price": 150,"status": "paid", "place": "delivery", "metadata": {"domain": "example.com"}}'
curl --location --request POST 'http://localhost:9191/orders' --header 'Authorization: Bearer token' --header 'Content-Type: application/json' --data-raw '{"name":"order4", "vendor":"seasons", "price": 300,"status": "paid", "place": "delivery", "metadata": {"domain": "example.com"}}'
curl --location --request POST 'http://localhost:9191/orders' --header 'Authorization: Bearer token' --header 'Content-Type: application/json' --data-raw '{"name":"order5", "vendor":"seasons", "price": 300,"status": "ordered", "place": "delivery", "metadata": {"domain": "example.com"}}'
curl --location --request POST 'http://localhost:9191/orders' --header 'Authorization: Bearer token' --header 'Content-Type: application/json' --data-raw '{"name":"order6", "vendor":"jikoni", "price": 200,"status": "ordered", "place": "delivery", "metadata": {"domain": "example.com"}}'
curl --location --request POST 'http://localhost:9191/orders' --header 'Authorization: Bearer token' --header 'Content-Type: application/json' --data-raw '{"name":"order7", "vendor":"jikoni", "price": 200,"status": "paid", "place": "inhouse", "metadata": {"domain": "example.com"}}'
curl --location --request POST 'http://localhost:9191/orders' --header 'Authorization: Bearer token' --header 'Content-Type: application/json' --data-raw '{"name":"order8", "vendor":"mess", "price": 100,"status": "paid", "place": "inhouse", "metadata": {"domain": "example.com"}}'
curl --location --request POST 'http://localhost:9191/orders' --header 'Authorization: Bearer token' --header 'Content-Type: application/json' --data-raw '{"name":"order9", "vendor":"mess", "price": 70,"status": "paid", "place": "inhouse", "metadata": {"domain": "example.com"}}'
curl --location --request POST 'http://localhost:9191/orders' --header 'Authorization: Bearer token' --header 'Content-Type: application/json' --data-raw '{"name":"order10", "vendor":"mess", "price": 20,"status": "paid", "place": "inhouse", "metadata": {"domain": "example.com"}}'
