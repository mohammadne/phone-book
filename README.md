# PhoneBook

## Deploy Postgresql locally∏

```bash
cd deployments

docker compose up -d
docker compose ps

docker exec -it deployments-db-1 psql -U PHONEBOOK_USER -W PHONEBOOK_DB
\dt
exit

docker compose down
```
