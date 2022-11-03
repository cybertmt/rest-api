### Server REST API  
[https://cybertmtx.crabdance.com/products](https://cybertmtx.crabdance.com/products)  
[https://cybertmtx.crabdance.com/stores](https://cybertmtx.crabdance.com/stores)  
[https://cybertmtx.crabdance.com/pricelist](https://cybertmtx.crabdance.com/pricelist)  
###
**--Products--**  
**Get all products**
```
curl https://cybertmtx.crabdance.com/products
```
**Create product**
```
curl -X "POST" -d '{"prod_name":"Асперин","prod_tr_name":"Asperin","prod_desc1":"Асперин: параметры"}' "https://cybertmtx.crabdance.com/products"
curl -X "POST" -d '{"prod_name":"Панадол","prod_tr_name":"Panadol","prod_desc1":"Панадол: параметры"}' "https://cybertmtx.crabdance.com/products"
curl -X "POST" -d '{"prod_name":"Парацетамол","prod_tr_name":"Paracetamol","prod_desc1":"Парацетамол: параметры"}' "https://cybertmtx.crabdance.com/products"
```
**Delete product by id**
```
curl -X "DELETE" -d '{"prod_id":2}' "https://cybertmtx.crabdance.com/products"
```
**Delete all products**
```
curl -X "DELETE" "https://cybertmtx.crabdance.com/clearproducts"
```
**Search sorted products by substring**
```
curl -X "POST" -d '{"prod_name":"П"}' "https://cybertmtx.crabdance.com/sortproducts"
```
**--Stores--**  
**Get all stores**
```
curl https://cybertmtx.crabdance.com/stores
```
**Create store**
```
curl -X "POST" -d '{"store_name":"Ригла","store_address":"Гончарный пр., 6, стр. 1, Москва","store_email":"info@rigla.ru","store_phone":"8 (800) 777-03-03","store_lat":55.739399,"store_lon":37.649848}' "https://cybertmtx.crabdance.com/stores"
curl -X "POST" -d '{"store_name":"Здоров.ру","store_address":"ул. Шаболовка, 34, стр. 3, Москва","store_email":"info@zdorov.ru","store_phone":"+7 (495) 363-35-00","store_lat":55.718311,"store_lon":37.607876}' "https://cybertmtx.crabdance.com/stores"
curl -X "POST" -d '{"store_name":"Горздрав","store_address":"Большая Переяславская ул., 11, Москва","store_email":"info@gorzdrav.ru","store_phone":"+7 (499) 653-62-77","store_lat":55.784470,"store_lon":37.641093}' "https://cybertmtx.crabdance.com/stores"
```
**Delete store by id**
```
curl -X "DELETE" -d '{"store_id":2}' "https://cybertmtx.crabdance.com/stores"
```
**Delete all stores**
```
curl -X "DELETE" "https://cybertmtx.crabdance.com/clearstores"
```
**--Prices--**  
**Create or Update price**
```
curl -X "POST" -d '{"prod_id":1,"store_id":1,"price":10.50}' "https://cybertmtx.crabdance.com/prices"
curl -X "POST" -d '{"prod_id":1,"store_id":2,"price":11.30}' "https://cybertmtx.crabdance.com/prices"
curl -X "POST" -d '{"prod_id":1,"store_id":3,"price":10.00}' "https://cybertmtx.crabdance.com/prices"
curl -X "POST" -d '{"prod_id":2,"store_id":1,"price":120.50}' "https://cybertmtx.crabdance.com/prices"
curl -X "POST" -d '{"prod_id":2,"store_id":2,"price":140.30}' "https://cybertmtx.crabdance.com/prices"
curl -X "POST" -d '{"prod_id":2,"store_id":3,"price":135.00}' "https://cybertmtx.crabdance.com/prices"
curl -X "POST" -d '{"prod_id":3,"store_id":1,"price":40.50}' "https://cybertmtx.crabdance.com/prices"
curl -X "POST" -d '{"prod_id":3,"store_id":2,"price":45.30}' "https://cybertmtx.crabdance.com/prices"
curl -X "POST" -d '{"prod_id":3,"store_id":3,"price":39.00}' "https://cybertmtx.crabdance.com/prices"
```
**Delete price by prod_id and store_id**
```
curl -X "DELETE" -d '{"prod_id":2, "store_id":2}' "https://cybertmtx.crabdance.com/prices"
```
**Delete all Prices**
```
curl -X "DELETE" "https://cybertmtx.crabdance.com/clearprices"
```
**Get Price List (sort: prod_name -> price -> store_name)**
```
curl https://cybertmtx.crabdance.com/pricelist
```
**Get Price List by product name**
```
curl -X "POST" -d '{"prod_name":"Панадол"}' "https://cybertmtx.crabdance.com/productprice"
```
**SignUp/Add user**
```
curl -X "POST" -d '{"useremail":"user1@yax.ru","password":"user1password"}' "https://cybertmtx.crabdance.com/signup"
```
**SignIn/Login user**
```
curl -X "POST" -d '{"useremail":"user1@yax.ru","password":"user1password"}' "https://cybertmtx.crabdance.com/signin"
curl -X "POST" -d '{"useremail":"user1@yax.ru","password":"user1passwo"}' "https://cybertmtx.crabdance.com/signin"
```
