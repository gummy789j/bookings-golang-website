# Bookings Website Project

* Built in Go version 

### Implemented Functions
* Make Reservation
* Search Availability
* Sending mail for notification
* Authentication
* Secure back end Administration

<br>

### Dependencies
* Uses the [chi router](https://github.com/go-chi/chi)
* Uses [alex edwards SCS](https://github.com/alexedwards/scs/v2)
* Uses [notie](https://github.com/jaredreich/notie)
* Uses [nosurf](https://github.com/justinas/nosurf)
* Uses [pgx Postgresql driver](https://github.com/jackc/pgx)
* Uses [Go validator](https://github.com/asaskevich/govalidator)

<br>

### How to run this website on your own server

1. clone project
```
$ git clone https://github.com/gummy789j/bookings_golang_website

```

2. build
```
$ go build -o bookings cmd/web/*.go

```

3. migrate Postgresql (you have to add soda as your enviromental parameters)
```
$ soda migrate

```

4. run
```
$ ./run.sh

```

P.S. Still writing test..


