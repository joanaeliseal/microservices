module github.com/joanaeliseal/microservices/order

go 1.26.1

require github.com/joanaeliseal/microservices-proto/golang/payment v0.0.0-00010101000000-000000000000
replace github.com/joanaeliseal/microservices-proto/golang/payment = > ../../
    microservices-proto/golang/payment