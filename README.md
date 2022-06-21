## REST API project»
### Server API 
[https://cybertmtx.crabdance.com/items](https://cybertmtx.crabdance.com/items)
###
**// Get items**
```
curl https://cybertmtx.crabdance.com/items
```
**// Get string items**
```
curl https://cybertmtx.crabdance.com/stringitems
```
**// Add items**
```
curl -X "POST" -d '{"title":"Msc Apt","content":"Moscow","link":"https://ya.ru","lat":55.751244, "lon":37.618423}' "https://cybertmtx.crabdance.com/items"
curl -X "POST" -d '{"title":"NY Apt","content":"New York","link":"https://google.com","lat":40.650002, "lon":-73.949997}' "https://cybertmtx.crabdance.com/items"
curl -X "POST" -d '{"title":"Syd Apt","content":"Sydney","link":"https://yahoo.com","lat":-33.865143, "lon":151.209900 }' "https://cybertmtx.crabdance.com/items"
```
**// Delete items**
```
curl -X "DELETE" -d '{"id":2}' "https://cybertmtx.crabdance.com/items"
```
**// Delete all items**
```
curl -X "DELETE" "https://cybertmtx.crabdance.com/clear"
```
**// Sort items**
```
curl -X "GET" -d '{"title":"sy"}' "https://cybertmtx.crabdance.com/sortitems"
```

