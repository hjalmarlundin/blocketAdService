# blocketAdService 

Simple web service for adding and retreiving ads. Based on sqllite and gorilla mux.

## GET Examples
            curl --location --request GET 'http://localhost:10000/ads?sort_by=date.asc'
            curl --location --request GET 'http://localhost:10000/ads?sort_by=date.desc'
            curl --location --request GET 'http://localhost:10000/ads?sort_by=price.asc'
            curl --location --request GET 'http://localhost:10000/ads?sort_by=price.desc'


## POST Examples


            curl --location --request POST 'http://localhost:10000/ad' \
            --header 'Content-Type: application/json' \
            --data-raw '{
                "Subject" : "My add subject",
                "Body" : "The body",
                "Email" : "x@y.se"
            }


## DELETE Examples

            curl --location --request DELETE 'http://localhost:10000/ad/4'