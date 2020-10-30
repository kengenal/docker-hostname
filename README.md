# Whoami

Display container hostname, system environments and headers

## Get started 
```console
term@term:~$ docker run -p 80:8080 kengenal/whoami
```

Now you should see container hostname on [localhost:8080](http://localhost:8080)

## Add system env
You can add new environments in query string like this:

``` console
http://localhost:8080?<key>=<value>
```