## A meal DB api with postgres

#### Endpoints
GET requests

* localhost:8000/meals
List all the meals
* localhost:8000/meal/{id}
List specific meal

POST request
* localhost:8000/meal/{id}
Add a new meal
```json
{
  "name": "some dish",
  "cuisine": "some cuisine",
  "category": "veg/nonveg"
}
```

PUT request

* localhost:8000/meal/{id}
Update a specific meal by ID

Delete request

* localhost:8000/meal/{id}
Delete a specific meal by ID
---