# Météo des aéroports

To build all files in bin/ :

```shell
make build
```

To build and launch a probe measurement :

```shell
make probe
```

To build and launch a subscriber :

```shell
make sub
```

To build and launch the server :

```shell
make http
```

To launch the database :

1. Firstly, launch the server to allow database requests and centralize the several communication modes :

```
redis serveur
```

2. Secondly, launch the redis client to allow redis to communicate

```
redis cli
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
