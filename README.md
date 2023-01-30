# Préveneur de Nuit

Préveneur de Nuit est une app qui permet de prévenir un Channel Slack que ça va être tout noir 1 minute avant la tombée de la nuit.

## Dépendances

L'app utilise l'api meteo-concept

Pour l'utiliser, il faut se créer un compte sur https://api.meteo-concept.com/ pour obtenir un token

L'app utilise aussi l'api Slack. Elle a besoin d'un token pour fonctionner, ainsi qu'une URL et un ID de channel sur lequel poster.

## SETUP

Certaines variables d'env sont nécessaires

```
API_METEO_CONCEPT_TOKEN 
API_METEO_CONCEPT_MAX_TRIES
API_METEO_CONCEPT_INSEE_LOCATION
SLACK_CHANNEL
SLACK_API_URL
SLACK_BOT_TOKEN
```
