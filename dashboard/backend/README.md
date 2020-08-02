## Конфигурация

|Параметр           |Пример        |
|:------------------|:-------------|
|SERVICE_HTTP_PORT  |8080          |
|LOGGER_LEVEL       |info          |
|JAEGER_URI         |127.0.0.1:6831|
|CLICKHOUSE_URI     |127.0.0.1:9500|
|CLICKHOUSE_USER    |user          |
|CLICKHOUSE_DATABASE|password      |


Проксируем запросы к базе <br />
Проверяем запросы к базе <br />
Если что то не так - шлем нахрен

---
Request: 
```json
{
    "params": [
        {
            "type": "key/column",
            "key": "p1",
            "value": {
                "item": "v1",                    // ?
                "list": ["v1_1", "v1_2", "v1_3"] // IN (?, ?, ?)
            },
            "operator": ">="
        }
    ],
    "limit": 1000,
    "offset": 5000
}
```

ЕСЛИ value == list && (operator == LIKE || operator == NOT LIKE)

SELECT * FROM internal_logs WHERE (p1 LIKE v1_1 OR p1 LIKE v1_2 OR p1 LIKE v1_3)

---
Пример получения данных из базы:
```sql
SELECT * FROM internal_logs 
      ARRAY JOIN `params_string.keys` AS key_string,     -- только если есть json поля
      arrayEnumerate(`params_string.keys`) AS idx_string -- только если есть json поля
      ARRAY JOIN `params_float.keys` AS key_float,       -- только если есть json поля
      arrayEnumerate(`params_float.keys`) AS idx_float   -- только если есть json поля

    WHERE (key_string='n2' AND arrayElement(`params_string.values`, idx_string) LIKE '%v3%') AND
          key_float='n3' AND arrayElement(`params_float.values`, idx_float) > 10
    ORDER BY time DESC
    LIMIT 5000, 1000;
```

---
Response:
```json
{
    "errors": [
        {
            "code": 123,
            "detail": "details"
        }
    ],
    "data": [
        {
            "time": "",
            "nsec": "",
            "namespace": "",
            "source": "",
            "host": "",
            "level": "",
            "trace_id": "",
            "message": "",
            "params": {},
            "build_commit": "",
            "config_hash": ""
        }
    ]
}
```

## Методы

api/v1/entry/list

api/v1/suggest/namespace 
api/v1/suggest/source
api/v1/suggest/host      --> "value": "v1" (host LIKE %v1%) --> "data": ["a", "b", "c"]