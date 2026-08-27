# Uzasadnienia warstwy projektowej

Uzasadnienia i zastrzeżenia projektowe przeniesione z komentarzy plików warstwy
projektowej, których treść przekracza dopuszczalną długość nagłówka. Rozdziały
noszą nazwę pliku, którego dotyczą.

### design/zasoby/powloka.js

Plik jest jedynym źródłem prawdy dla powłoki prototypów i niesie markup stałych
elementów okna: szynę nawigacji, belkę tytułową, pasek edycji (wstążkę), pasek
stanu oraz okno „Dostosuj paski”. Każdy prototyp podaje wyłącznie treść okna
centralnego w `<template id="dn-tresc-okna">` oraz punkt montażu
`<div id="dn-powloka-montaz">` — skrypt wstrzykuje powłokę wokół tej treści,
odtwarzając strukturę okna jeden do jednego. Działa wprost z dysku przez
`file://`, bez zapytań sieciowych.

### design/zasoby/tresci/licencja.js

Treść jest wynikiem działania generatora `zbuduj-licencje.py`, uruchamianego
na źródłowym tekście umowy licencyjnej. Plik nie jest przeznaczony do
bezpośredniej edycji — zmiana wchodzi w źródle i przechodzi przez generator
ponownie. Treść stoi w skrypcie, a nie w osobnym pliku do pobrania, ponieważ
okno bywa otwierane wprost z dysku, a przeglądarka blokuje wówczas pobieranie
plików towarzyszących.

## design/animacje/danaco-splash-3d.html

Plik zawiera ekran startowy trójwymiarowy, wyświetlany jako nakładka pomiędzy
uruchomieniem aplikacji a pojawieniem się jej właściwego widoku. Udostępnia
trzy warianty animacji w jednym pliku, wybierane parametrem adresu `v`:

- wariant `a`, „Złożenie rdzenia" — cząstki składają stos warstw platformy,
  stos zapada się w jeden moduł tworzący znak marki, moduł rozszerza się
  w kadr aplikacji;
- wariant `b`, „Rozkładanie interfejsu" — panele konsoli rozkładają się
  przestrzennie na wzór origami i układają w szkielet docelowego widoku;
- wariant `c`, „Przelot do rdzenia" — przelot kamery przez bramy kolejnych
  etapów startu (moduły, środowisko, klucze, sesja) aż do rdzenia.

Cykl odtwarzania przebiega przez fazy: wprowadzenie (budowa sceny), oczekiwanie
w pętli na gotowość aplikacji, zakończenie (odsłonięcie widoku docelowego)
i wywołanie zakończenia.

Osadzenie w powłoce aplikacji (Electron albo aplikacja jednostronicowa) odbywa
się przez umieszczenie elementu przykrywającego cały widok, wewnątrz którego
osadza się ramkę wskazującą ten plik z parametrem wariantu. Obiekt sterujący
udostępniany przez ramkę pozwala: rozpocząć odtwarzanie ze wskazanym wariantem
i minimalnym czasem trwania, ustawić opisowy podpis bieżącego etapu, wskazać
numer podświetlonej bramy dla wariantu przelotu, zgłosić gotowość aplikacji
rozpoczynającą odsłonięcie widoku docelowego oraz wymusić natychmiastowe
przejście do odsłonięcia. Zawartość pliku można też osadzić bezpośrednio
w oknie aplikacji, bez pośredniczącej ramki.

Parametry adresu obsługiwane przez plik: wybór wariantu, włączenie panelu
demonstracyjnego, włączenie pętli podglądu, czas w milisekundach do
automatycznego zgłoszenia gotowości oraz wyłączenie tła.

Ekran nie korzysta z zależności zewnętrznych: rysowanie odbywa się na płótnie
dwuwymiarowym z własną projekcją przestrzenną.
