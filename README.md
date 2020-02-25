# axios-test-data-server
This is a Go server that connects to ElephantSQL DB and returns the data you are testing for an app. 

## Resources
  
  Heroku build pack: https://github.com/kr/heroku-buildpack-go.git

  [Heroku-config-package](https://github.com/kardianos/govendor)

  [HTTP-ROUTER](https://github.com/gorilla/mux) - install: "go get -u github.com/gorilla/mux"

  [pq](https://github.com/lib/pq) - install: "go get github.com/lib/pq"

  [postman](https://www.getpostman.com/) - for testing endpoints

  [golang-install](https://golang.org/doc/install)

  [postgresql](https://www.postgresql.org/)

  [online-postgresql](https://www.elephantsql.com/)

### postgresql queries

  1. Creating test_items_weboost_cart table
  
      create TABLE test_items_weboost_cart (
      product_id int NOT NULL UNIQUE,  
      sub_total VARCHAR(5000),
      item_name VARCHAR(500) NOT NULL,
      item_handle VARCHAR(500) NOT NULL UNIQUE,
      item_sku VARCHAR(500) NOT NULL UNIQUE,
      item_price VARCHAR(500),
      item_image VARCHAR(500),
      PRIMARY KEY (product_id)
    );

  2. Creating a test_items_weboost_cart record
      INSERT INTO test_items_weboost_cart (
        product_id, 
        sub_total, 
        item_name,
        item_handle,
        item_sku,
        item_price,
        item_image
        ) VALUES (
          value1, 
          value2, 
          value3,
          value4,
          value5,
          value6,
          value7
        );
  
  3. Selecting errors

      select * from testdatacartweboost;
