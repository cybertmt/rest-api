## REST API project»

// Get items
curl http://cybertmtx.crabdance.com:7531/items

// Add item
curl -X "POST" -d '{"title":"Msc Apt","content":"Moscow","link":"https://ya.ru","lat":55.751244, "lon":37.618423}' "http://cybertmtx.crabdance.com:7531/items"
curl -X "POST" -d '{"title":"NY Apt","content":"New York","link":"https://google.com","lat":40.650002, "lon":-73.949997}' "http://cybertmtx.crabdance.com:7531/items"
curl -X "POST" -d '{"title":"Syd Apt","content":"Sydney","link":"https://yahoo.com","lat":-33.865143, "lon":151.209900 }' "http://cybertmtx.crabdance.com:7531/items"

// Delete item
curl -X "DELETE" -d '{"id":2}' "http://cybertmtx.crabdance.com:7531/items"

