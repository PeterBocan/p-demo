# p-demo

To run the demo:

```sh
./docker.sh
docker run -p 6734:6734 demo-app
```

Create account: 
```sh
curl -X POST http://localhost:6734/accounts -d @account.json
```

Create a transaction:
```sh 
curl -X POST http://localhost:6734/transactions -d @transaction.json
```

Get account information:
```sh
curl http://localhost:6734/accounts/1 
```