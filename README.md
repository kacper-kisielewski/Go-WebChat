# Go-WebChat

A real-time webchat built with Go

## Running in container

First of all create ``.env`` file in project's directory with the following content

```txt
JWT_SECRET=<random jwt key>
```

Generate a random key with the following command

``openssl rand -hex 128``

You can also customize application's settings in [settings.go](https://github.com/Kacperek1337/Go-WebChat/blob/master/settings/settings.go)

Finally you may run application in the container

``docker-compose up --build``

Wait a bit then open a web browser and navigate to [http://**localhost**:8000](http://localhost:8000)

> If cookies are not being set during the login process make sure you went to ``localhost`` not ``127.0.0.1``
