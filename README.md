# Météo des aéroports

To build all files in bin/ :

```shell
make build
```

To run probe file :

```shell
make probe
```

To run subscriber file :

```shell
make sub
```

## Composants

Protocoles utilisés : MQTT & HTTP

Database NoSQL Key-Value : REDIS (ou autre base NoSQL)

Broker MQTT : Moquitto (ou autre broker)

## Priorités

-   Capteur
-   Brokker MQTT
-   DB Redis (+ API Rest)

## Compétences à acquérir

-   DataBase REDIS
-   Protocole MQTT
-   Langage GO
-   Inteface Paho GO
