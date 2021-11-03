# Météo des aéroports

To build all files in bin/ :

```bash
make build
```

To build and launch the api framework :

```bash
make http
```

To build and launch both subscribers (to fulfill redis database and csv file) :

```bash
make sub
```

To build and launch a probe measurement :

```bash
make probe
```

To launch the database :

1. Firstly, launch the server to allow database requests and centralize the several communication modes :

```bash
redis-server
```

2. Secondly, launch the redis client to interact with the redis database

```bash
redis-cli
```

## Composants

Protocoles utilisés : MQTT & HTTP

Database NoSQL Key-Value : REDIS (ou autre base NoSQL)

Broker MQTT : Moquitto (ou autre broker)

## Priorités

- Capteur
- Brokker MQTT
- DB Redis (+ API Rest)

## Compétences à acquérir

- DataBase REDIS
- Protocole MQTT
- Langage GO
- Inteface Paho GO
