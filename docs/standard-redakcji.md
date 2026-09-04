# Standard redakcji komentarza i opisu

Dokument określa wymagania stawiane komentarzom w plikach kodu oraz opisom
w dokumentacji tej budowy. Wymagania obowiązują wszystkie pliki niosące treść
pisaną, niezależnie od języka i warstwy.

## Zasięg wymagań

Wymagania budowy komentarza dotyczą plików kodu. Opracowanie tekstowe niesie
treść ciągłą, w której znak wyliczenia wypada za znacznik komentarza, więc miara
nagłówka i punktu zwraca dla niego wynik poza zasięgiem, a jakość opracowania
rozstrzyga odczyt.

Wymagania językowe obowiązują również opracowania, z jednym wyłączeniem:
wskazanie dokumentu, z którego opracowanie czerpie, jest w nim treścią, a nie
odesłaniem zamiast treści. Zakaz odsyłania dotyczy komentarza w pliku kodu,
ponieważ komentarz ma być samodzielnym nośnikiem informacji.

## Objętość komentarza w pliku

Udział komentarzy nie przekracza pięciu procent objętości pliku, a pojedynczy
komentarz trzystu pięćdziesięciu znaków. Miarę egzekwuje walidator dyscypliny
inżynierskiej i to on rozstrzyga; zapisy podające inną wartość są nieaktualne.

Pliki migracji podlegają redakcji komentarza. Rdzeń liczy sumę kontrolną
z treści znormalizowanej — bez komentarzy, bez wierszy pustych i bez odstępu
końcowego (`store.normalizujTrescMigracji`), więc zmiana samego komentarza nie
rusza sumy i nie zatrzymuje bazy istniejącej. Sam schemat kroku już
zastosowanego zostaje nietknięty: zmiana DDL zmienia sumę i jest odmową startu.

Pliki wytworzone generatorem są poza zasięgiem redakcji. Ich treść powstaje
z generatora, więc ręczna zmiana ginie przy następnym wytworzeniu, a do czasu
tego wytworzenia zapis w repozytorium przeczy swojemu źródłu.

## Budowa komentarza

Każdy plik zawiera nagłówek. Nagłówek niesie co najmniej jedno pełne zdanie
i mieści się w przedziale od stu do dwustu pięćdziesięciu znaków. Wymaganie
dotyczy każdego nagłówka z osobna i nie zależy od liczby wierszy kodu.

Komentarz umieszczony w treści kodu mieści się w stu znakach. Liczba takich
komentarzy nie jest ograniczona, o ile każdy niesie objaśnienie kluczowe albo
krytyczne dla zrozumienia działania.

## Język

Obowiązuje polszczyzna zawodowa i formalna, w trybie oznajmującym. Komentarz
podaje informację o stanie obowiązującym, a nie przebieg prac ani ich historię.

Wypowiedź składa się z pełnych zdań i pełnych słów. Skróty i parafrazy są
niedopuszczalne, podobnie jak terminy niedookreślone.

Terminologia pochodzi wyłącznie ze słownika zawodowego dziedziny. Wprowadzanie
własnych określeń, oznaczeń i postanowień jest niedopuszczalne.

## Wyrażenia niedopuszczalne

Komentarz nie odsyła czytelnika do opracowań tekstowych, dokumentów planowania
budowy ani wykazów ustaleń. Komentarz nie przywołuje numeracji wprowadzonej na
potrzeby prac ani znaczników roboczych.

Oznaczenia złożone z liter i cyfr nie należą do języka zawodowego i nie są
dopuszczalne. Niedopuszczalne są również odwołania osobowe oraz przywoływanie
ustaleń podjętych w toku prac.

Nazwanie narzędzia, biblioteki, arkusza stylu, programu zewnętrznego, migracji
albo pliku źródłowego, z którym kod współpracuje, jest treścią techniczną
i wymagań nie narusza. Tak samo oznaczenie normy, algorytmu albo protokołu
zewnętrznego, ponieważ jest zakotwiczone poza budową.

Zakaz odwołań osobowych obejmuje przywołanie roli prowadzącej budowę oraz roli
rozstrzygającej zakres produktu, a wraz z nimi ustaleń zapadłych w toku prac.
Wyraz nazywający rolę wewnątrz produktu, na przykład właściciela konta albo
urządzenia, należy do słownika dziedziny i pozostaje dopuszczalny.

## Treść wykraczająca poza komentarz

Uzasadnienie dłuższe niż nagłówek przenosi się do dokumentacji technicznej
warstwy, której dotyczy. Komentarz w pliku kodu pozostaje samodzielnym nośnikiem
informacji: podaje rzecz, zamiast kierować po nią gdzie indziej.

Powiązanie kodu z dokumentacją niesie porządek nazewniczy: rozdział dokumentacji
nosi nazwę pliku, którego dotyczy.

## Niezmienność kodu

Redakcja komentarza nie zmienia kodu. Dowodem jest zgodność strumienia jednostek
składniowych pliku, pomijającego komentarze, przed redakcją i po niej. Narzędzie
`narzedzia/redakcja.go` podaje ten strumień oraz mierzy pozostałe wymagania.
