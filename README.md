# go-config-api

### Endpoints

Following are the endpoints are implemented:

| Name   | Method      | URL                          |
| ------ | ----------- | ---------------------------- |
| List   | `GET`       | `/configs`                   |
| Create | `POST`      | `/configs`                   |
| Get    | `GET`       | `/configs/{name}`            |
| Update | `PUT/PATCH` | `/configs/{name}`            |
| Delete | `DELETE`    | `/configs/{name}`            |
| Query  | `GET`       | `/search?metadata.key=value` |

#### Query

The query endpoint **MUST** return all configs that satisfy the query argument.

Query example-1:

```sh
curl http://config-service/search?metadata.monitoring.enabled=true
```

Response example:

```json
[
  {
    "name": "datacenter-1",
    "metadata": {
      "monitoring": {
        "enabled": "true"
      },
      "limits": {
        "cpu": {
          "enabled": "false",
          "value": "300m"
        }
      }
    }
  },
  {
    "name": "datacenter-2",
    "metadata": {
      "monitoring": {
        "enabled": "true"
      },
      "limits": {
        "cpu": {
          "enabled": "true",
          "value": "250m"
        }
      }
    }
  },
]
```


Query example-2:

```sh
curl http://config-service/search?metadata.allergens.eggs=true
```

Response example-2:

```json
[
  {
    "name": "burger-nutrition",
    "metadata": {
      "calories": 230,
      "fats": {
        "saturated-fat": "0g",
        "trans-fat": "1g"
      },
      "carbohydrates": {
          "dietary-fiber": "4g",
          "sugars": "1g"
      },
      "allergens": {
        "nuts": "false",
        "seafood": "false",
        "eggs": "true"
      }
    }
  }
]
```

#### Schema

- **Config**
  - Name (string)
  - Metadata (nested key:value pairs where both key and value are strings of arbitrary length)

### Configuration

Application serves the API on the port defined by the environment variable `SERVE_PORT`.
The application fails if the environment variable is not defined.
