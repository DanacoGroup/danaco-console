# Ustrój budowy — role, koordynacja, nadzór

Dokument ustala, kto podejmuje jakie decyzje, jak sesje pracują równolegle bez
wchodzenia sobie w drogę i co musi się wydarzyć, zanim praca zostanie przyjęta.
Obowiązuje na wszystkich gałęziach i we wszystkich etapach.

## 1. Dlaczego ustrój, a nie same wytyczne

Poprzednie podejście do budowy upadło nie na pojedynczym błędzie, lecz na dryfie:
kolejne sesje dopowiadały nazwy i wartości, których nie było w źródłach, a nikt
niezależny tego nie sprawdzał. Ustrój odpowiada na tę przyczynę trzema
mechanizmami: **rozdziałem wykonawcy od kontrolera**, **pisemnym zakresem terenu
przed startem pracy** oraz **jednym rejestrem decyzji**, poza którym żadne
rozstrzygnięcie nie obowiązuje.

## 2. Role i odpowiedzialności

Rola nie jest osobą — jedna sesja pełni w danym momencie dokładnie jedną rolę
i nie łączy wykonania z kontrolą własnej pracy.

### 2.1 Właściciel produktu

Dariusz Naharnowicz. Rozstrzyga zakres produktu i miejsca oznaczone jako otwarte.
Przyjmuje etap i tym samym otwiera następny. Jest jedynym, kto zmienia plan
etapów, zakłada wpis w rejestrze decyzji i zgadza się na nowy dokument
w zamkniętym zbiorze. Żadna sesja nie rozstrzyga za Właściciela — brak nazwy,
wartości albo zachowania w źródłach jest usterką do zgłoszenia, nie swobodą
wykonawcy.

### 2.2 Prowadzący budowę

Sesja koordynująca, jedna naraz. Utrzymuje plan etapów, rejestr decyzji i rejestr
terenów. Dzieli etap na tereny, wyznacza zakres każdego z nich i pilnuje, żeby
tereny się nie nakładały. Otwiera i zamyka tereny, scala gałęzie do `main`.
Nie wytwarza treści merytorycznej — koordynacja i wykonanie w jednej sesji
znoszą kontrolę.

### 2.3 Wykonawca terenu

Sesja robocza, wiele naraz. Pracuje na jednym terenie we własnym drzewie
roboczym, wyłącznie na plikach swojego terenu. Wytwarza opracowania, prototypy
albo kod zgodnie z pisemnym zakresem terenu. Zgłasza Prowadzącemu każdą potrzebę
wyjścia poza teren zamiast sięgać po cudze pliki. Zgłasza braki w źródłach
zamiast je uzupełniać domysłem.

### 2.4 Kontroler

Sesja przeglądowa, uruchamiana po zgłoszeniu terenu do zamknięcia. Nigdy nie jest
tą samą sesją, która teren wykonała. Sprawdza wynik wobec kryteriów odbioru
terenu i zwraca uwagi — nie poprawia sam, bo poprawka kontrolera znosi kontrolę.
Wynik kontroli jest dwustanowy: teren przyjęty albo zwrócony z wykazem uchybień.

## 3. Podział na tereny

**Teren** to rozłączny wycinek drzewa przypisany jednej sesji na czas jednej
pracy. Podział na tereny jest jedynym mechanizmem pracy równoległej — dwie sesje
nigdy nie piszą do tego samego pliku.

Teren przed startem ma zapisane w [rejestrze terenów](rejestr-terenow.md):

1. **Nazwę** — rzeczownikową, opisującą przedmiot pracy, na przykład
   `zetony-i-typografia`, `okna-modulu-studio`.
2. **Wykaz plików** — pełne ścieżki albo katalog, którego teren dotyczy.
3. **Kryteria odbioru** — sprawdzalne zdania, po których wiadomo, że praca jest
   skończona.
4. **Gałąź i drzewo robocze** — `teren/<nazwa>` w katalogu `~/robocze/<nazwa>`.

Teren bez zapisanych kryteriów odbioru nie zostaje otwarty. Kryterium
niesprawdzalne („poprawić jakość", „ujednolicić wygląd") nie jest kryterium.

## 4. Praca równoległa

Sesja pracuje w osobnym drzewie roboczym Git, nie w drzewie głównym:

```bash
bash narzedzia/nowy-teren.sh <nazwa> <galaz-bazowa>
```

Skrypt zakłada gałąź `teren/<nazwa>` na wskazanej gałęzi bazowej, tworzy drzewo
robocze w `~/robocze/<nazwa>` i dowiązuje wspólne zasoby. Gałęzią bazową jest
`main` — niesie całość budowy. Wyjątkiem jest teren wyrastający z pracy innego
terenu; wtedy gałąź bazowa jest wskazana wprost w rejestrze terenów.

Reguły obowiązujące każdą sesję:

- Rewizja obejmuje wyłącznie pliki własnego terenu. Polecenie `git add -A`
  jest zabronione, bo w drzewie leży żywa praca innych sesji.
- Komunikat rewizji: tryb oznajmujący, jedno zdanie, do 70 znaków.
- Zmiana wykraczająca poza teren wraca do Prowadzącego jako zgłoszenie.
- Drzewo robocze po zamknięciu terenu jest usuwane — `git worktree remove`.

## 5. Bramki

Praca przechodzi dalej wyłącznie przez bramkę. Bramka jest zdarzeniem
sprawdzalnym, nie deklaracją.

### 5.1 Bramka wejścia terenu

Otwiera ją Prowadzący. Warunki: teren ma nazwę, wykaz plików i kryteria odbioru;
wykaz plików nie przecina się z żadnym terenem otwartym; gałąź bazowa jest
wskazana.

**Teren ruszający kontrakt otwiera się wyłącznie na podstawie pozycji rejestru
decyzji.** Kontrakt jest produktem — każda komenda to zdolność, którą platforma
odtąd obiecuje — więc dołożenie komendy jest rozstrzygnięciem o zakresie, a to
należy do Właściciela. Prawo zmiany kontraktu nie jest prawem Prowadzącego do
nadania. Warunek wprowadzony po tym, jak teren `pomiar-stron` dołożył trzy
zdolności bez takiej podstawy (pozycja 16 rejestru decyzji).

### 5.2 Bramka wyjścia terenu

Zamyka ją Kontroler. Warunki: wszystkie kryteria odbioru spełnione; rewizje
obejmują wyłącznie pliki terenu; weryfikacja maszynowa właściwa dla rodzaju
pracy zakończona wynikiem pozytywnym, z przytoczonym wynikiem uruchomienia.
Kontroler przytacza to, co sprawdził, i wprost nazywa to, czego sprawdzić się
nie dało.

### 5.3 Bramka etapu

Zamyka ją Właściciel. Warunki: wszystkie tereny etapu przyjęte; komplet
przewidziany planem etapów istnieje; rejestr decyzji nie zawiera pozycji
otwartych blokujących etap. Dopiero po zamknięciu bramki etapu Prowadzący scala
gałąź przebudowy do `main`.

## 6. Zapory przed dryfem

Poniższe reguły obowiązują bezwarunkowo i nie podlegają ocenie sytuacyjnej.

- **Bez dopowiadania.** Nazwa, wartość, okno, pole ani zachowanie bez pokrycia
  w źródłach nie wchodzi do produktu. Brak jest usterką źródła do zgłoszenia.
- **Jedno pojęcie — jedna nazwa.** Synonim w kodzie, dokumentacji albo etykiecie
  jest usterką. Nazwa opisuje funkcję, nie metaforę.
- **Bez wymyślonych oznaczeń.** Autorskie kody, sygnatury i numeracje są
  zabronione. Dopuszczalne są wyłącznie identyfikatory zakotwiczone w istniejących
  systemach: numer zgłoszenia, oznaczenie normy, rzeczywisty kod błędu platformy.
- **Bez narośli plikowych.** Nowa treść wchodzi edycją właściwego dokumentu
  kanonicznego. Notatki, podsumowania, plany robocze i kopie z przyrostkiem
  wersji nie należą do repozytorium.
- **Bez kroniki w treści.** Dokument i komentarz opisują stan obecny oraz jego
  uzasadnienie. Przebieg prac przechowuje historia rewizji.
- **Bez zapisu stanu przejściowego w dokumencie trwałym.** Ostrzeżenie w rodzaju
  „ten plik jest niesprawny", „nie czerp stąd ustaleń", „czeka na poprawę" nie
  należy do dokumentu prowadzenia. Stan przejściowy mija, a zapis zostaje —
  i od chwili poprawy zaczyna kłamać, kierując kolejne sesje przeciwko rzeczy,
  która jest już dobra. Tak zdarzyło się w poprzednim podejściu: adnotacje
  o niewiarygodności przeżyły powód swojego istnienia i wykonawcy omijali
  materiał od dawna poprawny. Stan przejściowy niesie teren, który się zamyka,
  albo zgłoszenie, które się rozstrzyga — nigdy dokument, który trwa.
- **Bez zmian zbiorczych.** Podmiana wzorcem po katalogach jest zabroniona —
  plik czyta się i zmienia celowaną zmianą.
- **Przed otwarciem terenu sprawdza się, czy przedmiot już nie istnieje.**
  Budowa toczy się na przejętym kodzie liczącym setki tysięcy wierszy, więc
  rzecz, która wygląda na brakującą, bywa już zrobiona. Sprawdzenie jest
  odczytem drzewa, nie przypomnieniem sobie: wykaz plików, przeszukanie nazw,
  uruchomienie tego, co znalezione. Teren otwarty na przedmiot już istniejący
  jest pracą do wyrzucenia — zdarzyło się to przy generatorze typów kontraktu,
  który stał gotowy w `budowa/shared/gen/` wraz z testem świeżości.
- **Pomiar odróżnia brak wyniku od wyniku pustego.** Skrypt, który może zwrócić
  zero albo pustkę, musi najpierw potwierdzić, że zmierzył to, co miał zmierzyć.
  Zero naruszeń na stronie, która się nie wczytała, nie jest wynikiem — jest
  brakiem pomiaru podanym jako wynik. Ta klasa błędu wystąpiła w jednym dniu
  trzykrotnie i za każdym razem instrument milczał, zamiast zgłosić przeszkodę:
  serwer podniesiony z niewłaściwego katalogu mierzył stronę nieistniejącą,
  przestarzały selektor wskazywał element sprzed zmiany, a zawieszony serwer
  zajął port i oddawał cudze drzewo. Wynik wyglądający dobrze jest groźniejszy
  od błędu, bo nie wzywa do sprawdzenia.

## 7. Zgłoszenia

**Zgłoszenie jest ostatecznością, nie odruchem.** Zanim sesja zgłosi cokolwiek,
wyczerpuje źródła: kontrakt, kod rdzenia, prototyp przyjęty, rejestr decyzji —
i mierzy uruchomieniem. Większość rzeczy, które wyglądają na braki, jest
odpowiedzią, której nikt nie doczytał. Zdarzyło się to w tej budowie wprost:
pozycja „pierwsze uruchomienie bez poczty" wisiała jako otwarta, a rdzeń miał
rozstrzygnięcie od początku wraz z uzasadnieniem — Prowadzący wziął nieaktualny
sprawdzian za usterkę. Pozycja otwarta, którą źródło zamyka, kosztuje Właściciela
uwagę i zatrzymuje front bez powodu.

Sesja zgłasza Prowadzącemu, gdy: brakuje wartości w źródle, dwa źródła są
sprzeczne, praca wymaga wyjścia poza teren albo kryterium odbioru okazuje się
niesprawdzalne. Zgłoszenie jest jednym akapitem: co stwierdzono, gdzie, jaki
skutek. Sesja nie wstrzymuje reszty pracy — wykonuje wszystko, co od zgłoszenia
nie zależy, i przekazuje teren z nazwanym brakiem.

Zgłoszenie nie wstrzymuje pracy. Sesja podaje rozstrzygnięcie, które przyjmuje,
wraz z powodem, i idzie dalej pod tym założeniem — wstrzymanie jest właściwe
wyłącznie wtedy, gdy praca pod dowolnym założeniem byłaby bezużyteczna, gdyby
założenie okazało się błędne.

Zgłoszenie rozstrzygnięte przez Właściciela trafia do rejestru decyzji. Poza
rejestrem żadne rozstrzygnięcie nie obowiązuje.
