# go-practise-project
Scaling example project in Golang 

### Steps to Setup

```
docker-compose up
```

Used Socket.IO (v2 protocol)

```
http://localhost:8000?user_id=4
```

need to add this url to socket.io and it will connect for that user, Events associated with the server are :
<ul>
  <li> startPlayingMode</li>
  <li> fetchPopularMode</li>
</ul>

for StartPlayingMode the Message is

```
{
    "area_code": 212,
    "current_login_time": {{$timestamp}},
    "game_mode": "multiplayer"
}
```

for FetchPopularMode we have to pass the,

```
{
    "area_code": 212,
    "current_login_time": {{$timestamp}},
}
```
 
Once the user gets connected the session will create and once he disconnectes he will get logout.



