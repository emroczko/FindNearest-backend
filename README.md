# FindNearMe

## TODO 

### logika databaseClienta (który jest de facto też input serwisem):

1. Logika inputu, czyli co user może wrzucić (np. typ sklepu, trzeba to sparsować i odpowiednio wrzucić odpowiednie zapytanie do bazy danych) 

2. Oddzielenie modelu bazy danych od zwracanego modelu; docelowo w odpowiedzi z tego serwisu ma być tablica obiektów z polami:

objectType
objectName
objectInfo (optional)
distance 
distanceToOtherObject
time
timeToOtherObject
ewentualnie route to other object, a jeśli serwis w Javie będzie szybki to i route to object
