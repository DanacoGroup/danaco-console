/* ============================================================================
   PRZEPŁYW WEJŚCIA — katalog treści

   Jedyne miejsce z tekstem widocznym dla użytkownika. Klucze idą hierarchicznie:
   obszar → element. Plik przekazuje się tłumaczowi bez dostępu do kodu — nie ma
   tu znaczników ani logiki.

   Wartości z {nawiasami} to miejsca na dane podstawiane w czasie działania;
   ich nazw nie tłumaczy się.

   Dlaczego skrypt, a nie plik `.json`: okno bywa otwierane wprost z dysku,
   a przeglądarka blokuje wtedy pobieranie plików towarzyszących. Zawartość
   pozostaje czystym JSON-em — klucze w cudzysłowach, bez przecinka po ostatniej
   pozycji — więc narzędzia tłumaczy czytają ją tak samo.
   ============================================================================ */
window.DanacoWejscie = window.DanacoWejscie || {};
window.DanacoWejscie.tresci = {
  "_jezyk": "pl-PL",

  "okno": {
    "uruchamianie": "Danaco Console — uruchamianie",
    "dostep": "Danaco Console — dostęp do konta",
    "zwin": "Zwiń okno",
    "rozwin": "Rozwiń okno",
    "zamknij": "Zamknij okno"
  },

  "marka": {
    "nazwa": "Danaco Console",
    "uruchamianie": {
      "motto": "AI Operating Environment",
      "zalety": [
        { "glowa": "Cztery środowiska pracy",
          "tresc": "Rozmowa i wiedza, projekty, wytwarzanie oprogramowania oraz praca zespołu modeli." },
        { "glowa": "Praca wraca w tym samym stanie",
          "tresc": "Zamknięcie aplikacji nie przerywa zadań — wracają wraz z kartami i kontekstem." },
        { "glowa": "Jedno konto na wszystkie urządzenia",
          "tresc": "Ten sam dostęp na komputerze, tablecie i telefonie." }
      ]
    },
    "dostep": {
      "motto": "Platforma AI Workspace OS",
      "zalety": [
        { "glowa": "Cztery środowiska pracy",
          "tresc": "Osobna przestrzeń robocza dla każdego rodzaju pracy: TalkIn, WorkSpace, CodeStudio, MultitaskingAI." },
        { "glowa": "Wiele modeli nad jednym zleceniem",
          "tresc": "Koordynator rozdziela zlecenie między wykonawców; przebieg śledzisz w Execution Loop Window." },
        { "glowa": "Izolacja kontekstu pod kontrolą",
          "tresc": "Punkty izolacji kontekstu ustala Operator, nie sztywna reguła systemu." }
      ]
    },
    "wydawca": "Danaco Holding Group Sp. z o.o.",
    "wydanie": "wydanie 2.0",
    "wsparcie": "support@danaco-core.pl"
  },

  "uruchomienie": {
    "nadtytul": "Uruchomienie",
    "etapy": ["Połączenie z serwerem", "Uzgodnienie wersji", "Rozpoznanie urządzenia", "Przygotowanie okna"],
    "stany": {
      "oczekuje": "—",
      "wToku": "w toku",
      "gotowe": "gotowe",
      "nawiazane": "nawiązane",
      "zgodna": "zgodna",
      "zaufane": "zaufane",
      "nieudane": "nieudane"
    },
    "laczenie": {
      "tytul": "Nawiązywanie połączenia",
      "lid": "Danaco Console łączy się z serwerem. Poczekaj."
    },
    "token": {
      "tytul": "Przywracanie sesji",
      "lid": "Danaco Console przywraca sesje otwarte na tym koncie.",
      "baner": {
        "glowa": "Rozpoznano zaufane urządzenie",
        "tresc": "Uruchomienie bez logowania."
      },
      "fraza": "Rejestracja i logowanie zostają pominięte, bo urządzenie jest zaufane."
    },
    "blad": {
      "tytul": "Błąd połączenia",
      "lid": "Sprawdź połączenie sieciowe i spróbuj ponownie.",
      "baner": {
        "glowa": "Serwer nie odpowiada",
        "tresc": "Ponowna próba za {sekundy} sekund."
      },
      "banerDane": { "sekundy": 15 },
      "fraza": "Praca zapisana na serwerze jest bezpieczna."
    }
  },

  "dostep": {
    "zakladki": ["Zaloguj się", "Zarejestruj się"],

    "logowanie": {
      "tytul": "Zaloguj się",
      "lid": "Podaj dane konta Operatora.",
      "login": "Login albo adres e-mail",
      "haslo": "Hasło",
      "reset": { "pytanie": "Nie pamiętasz hasła?", "czynnosc": "Resetuj hasło" },
      "metody": {
        "naglowek": "Inne metody logowania",
        "pin": "Kod PIN",
        "klucz": "Klucz systemowy",
        "email": "Kod na adres e-mail",
        "aktywna": "aktywna",
        "nieaktywna": "nieaktywna"
      },
      "fraza": { "pytanie": "Nie masz konta?", "czynnosc": "Utwórz konto Operatora" }
    },

    "logowanieBlad": {
      "baner": {
        "glowa": "Nie rozpoznano danych logowania.",
        "tresc": "Pozostały dwie próby, po nich logowanie zostanie wstrzymane na godzinę."
      },
      "capsLock": "Sprawdź, czy nie jest włączony Caps Lock.",
      "fraza": "Po pięciu nieudanych próbach logowanie zostaje wstrzymane na godzinę. Dostęp można wtedy odzyskać przez adres e-mail."
    },

    "rejestracja": {
      "tytul": "Utwórz konto Operatora",
      "lid": "Adres e-mail posłuży do potwierdzenia konta i odzyskania dostępu.",
      "login": "Login",
      "email": "Adres e-mail",
      "haslo": "Hasło",
      "hasloPowtorz": "Powtórz hasło",
      "fraza": { "pytanie": "Masz już konto?", "czynnosc": "Zaloguj się" }
    },

    "sila": {
      "puste": "Siła hasła zostanie oceniona podczas wpisywania",
      "stopnie": [
        "Hasło nie spełnia żadnego z warunków",
        "Hasło słabe — spełnia jeden warunek",
        "Hasło dostateczne — spełnia dwa warunki",
        "Hasło dobre — spełnia trzy warunki",
        "Hasło mocne — spełnia wszystkie cztery warunki"
      ],
      "warunki": {
        "dlugosc": "co najmniej 12 znaków",
        "wielkosc": "wielka i mała litera",
        "cyfra": "cyfra",
        "znak": "znak specjalny"
      }
    },

    "sesja": {
      "etykieta": "Nie wylogowuj mnie na tym urządzeniu",
      "opis": "Sesja pozostanie aktywna do czasu wylogowania. Nie zaznaczaj na urządzeniu współdzielonym."
    },

    "kod": {
      "tytul": "Potwierdź adres e-mail",
      "lid": "Na adres {adres} został wysłany sześciocyfrowy kod. Zachowuje ważność przez {minuty} minut.",
      "lidDane": { "adres": "operator@danaco-group.pl", "minuty": 10 },
      "obszar": "Kod potwierdzający — sześć znaków",
      "znak": "Znak {numer} z sześciu",
      "odliczanie": "Kod traci ważność za",
      "wklej": "Wklej kod ze schowka",
      "ponow": "Wyślij kod ponownie",
      "pomoc": {
        "glowa": "Nie ma wiadomości?",
        "tresc": "Sprawdź folder wiadomości niechcianych. Kod przychodzi z adresu {nadawca} i dociera zwykle w ciągu minuty. Jeśli nie dotarł, wyślij go ponownie."
      },
      "pomocDane": { "nadawca": "noreply@danaco-core.pl" },
      "ostrzezenie": {
        "glowa": "Kod wprowadza się wyłącznie w tym oknie.",
        "tresc": "Danaco Console nigdy nie prosi o kod przez telefon ani w wiadomości zwrotnej."
      },
      "zmienAdres": "Zmień adres e-mail"
    },

    "odzyskiwanie": {
      "kroki": ["Adres e-mail", "Potwierdzenie", "Nowe hasło"],
      "adres": {
        "tytul": "Odzyskaj dostęp do konta",
        "lid": "Podaj adres e-mail konta. Zostanie na niego wysłany sześciocyfrowy kod potwierdzający.",
        "pole": "Adres e-mail",
        "ostrzezenie": {
          "glowa": "Zmiana hasła kończy wszystkie sesje.",
          "tresc": "Na pozostałych urządzeniach trzeba zalogować się ponownie."
        }
      },
      "haslo": {
        "tytul": "Ustaw nowe hasło",
        "lid": "Nowe hasło zacznie obowiązywać od razu. Pozostałe urządzenia zostaną wylogowane.",
        "nowe": "Nowe hasło",
        "powtorz": "Powtórz nowe hasło"
      },
      "powrot": "Wróć do logowania"
    }
  },

  "dzialania": {
    "zamknijAplikacje": "Zamknij aplikację",
    "ustawieniaPolaczenia": "Ustawienia połączenia",
    "sprobujPonownie": "Spróbuj ponownie",
    "zaloguj": "Zaloguj się",
    "zalogujPonownie": "Zaloguj się ponownie",
    "utworzKonto": "Utwórz konto",
    "potwierdzKonto": "Potwierdź konto",
    "potwierdzKod": "Potwierdź kod",
    "wyslijKod": "Wyślij kod potwierdzający",
    "potwierdzHaslo": "Potwierdź nowe hasło"
  },

  "usterki": {
    "wiele": "Popraw zaznaczone dane.",
    "naglowekKonto": "Nie można utworzyć konta.",
    "naglowekHaslo": "Nie można ustawić hasła.",
    "naglowekLogowanie": "Nie można się zalogować.",
    "naglowekKod": "Nie można wysłać kodu.",
    "brakLoginu": {
      "glowa": "Podaj login albo adres e-mail.",
      "tresc": "Pole logowania jest puste."
    },
    "brakHasla": {
      "glowa": "Podaj hasło.",
      "tresc": "Pole hasła jest puste."
    },
    "brakDanych": {
      "glowa": "Podaj dane logowania.",
      "tresc": "Wpisz login albo adres e-mail oraz hasło do konta Operatora."
    },
    "brakAdresu": {
      "glowa": "Podaj adres e-mail konta.",
      "tresc": "Na ten adres zostanie wysłany kod potwierdzający."
    },
    "loginZajety": {
      "glowa": "Ten login jest już zajęty.",
      "tresc": "Wybierz inny login. Adres e-mail może pozostać bez zmian."
    },
    "emailBledny": {
      "glowa": "Nieprawidłowy adres e-mail.",
      "tresc": "Sprawdź, czy adres zawiera znak @ oraz nazwę domeny, na przykład nazwa@firma.pl."
    },
    "hasloSlabe": {
      "glowa": "Hasło nie spełnia wymagań.",
      "tresc": "Spełnij wszystkie cztery warunki podane pod polem hasła."
    },
    "haslaRozne": {
      "glowa": "Hasła nie są zgodne.",
      "tresc": "Wpisz to samo hasło w obu polach."
    }
  },

  "komunikaty": {
    "zamkniecie": {
      "tytul": "Zamknięcie aplikacji",
      "tresc": "Okno aplikacji zostanie zamknięte. Sesje pracujące w tle trwają dalej i wrócą przy kolejnym uruchomieniu."
    },
    "ustawienia": {
      "tytul": "Ustawienia połączenia",
      "tresc": "Ustawienia połączenia otwierają się w oknie konfiguracji, w zakresie „Sieć i serwer”."
    },
    "kodPonowiony": {
      "tytul": "Kod wysłany ponownie",
      "tresc": "Nowy kod wysłany na adres konta. Poprzedni przestał obowiązywać."
    },
    "schowek": {
      "tytul": "Schowek niedostępny",
      "tresc": "Wpisz kod ręcznie w sześciu polach."
    }
  }
};
