URL = http://url-target.com
CONCURRENCY = 2
REQUESTS = 10
HTTP_METHOD = GET

build:
	docker build -t go-http-stress .

run:
	docker run --rm go-http-stress:latest -u $(URL) -r $(REQUESTS) -c $(CONCURRENCY) -X $(HTTP_METHOD)
