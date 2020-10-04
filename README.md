### STORIES

[![Build Status](https://travis-ci.com/nsnikhil/stories.svg?token=r1U2n3nQxoEcNsRAxVeK&branch=master)](https://travis-ci.com/nsnikhil/stories)

#### TODO
- Add new relic transaction to database.
- Complete health server check and watch.
- Add more scenarios to integration test.

#### setup
```
make setup
```

#### start grpc server
```
make grpc-serve
```

#### start http server
```
make http-serve
```

#### test
```
make test
```

#### migration
```
make migrate
```

#### rollback
```
make rollback
```