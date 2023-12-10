# Go-Server-Sort-Multiple-Arrays

#Pull Docker image:
docker pull nishantsrivastav260/go-sorting-server:1.01

# Run Docker image
docker run -d -p 8080:8080 nishantsrivastav260/go-sorting-server:1.01

# To run the code through curl command example:
curl -X POST -H "Content-Type: application/json" -d '{"to_sort": [[3, 54, 56, 24, 36], [7, 5, 5, 6], [7, 23, 9]]}' http://localhost:8080/process-single
curl -X POST -H "Content-Type: application/json" -d '{"to_sort": [[3, 54, 56, 24, 36], [7, 5, 5, 6], [7, 23, 9]]}' http://localhost:8080/process-concurrent
