/**
 * Droga wejścia — katalog treści, jedyne miejsce z tekstem widocznym dla
 * użytkownika. Klucze idą hierarchicznie: obszar do elementu, jak ścieżka.
 */

/** Węzeł katalogu: łańcuch, liczba, wykaz albo poddrzewo złożone z kolejnych węzłów tego samego rodzaju. */
export type WezelTresci = string | number | WezelTresci[] | { [klucz: string]: WezelTresci };

export const tresci = {
  jezyk: 'pl-PL',

  okno: {
    uruchamianie: 'Danaco Console — uruchamianie',
    dostep: 'Danaco Console — dostęp do konta',
    zwin: 'Zwiń okno',
    rozwin: 'Rozwiń okno',
    zamknij: 'Zamknij okno',
  },

  marka: {
    nazwa: 'Danaco Console',
    uruchamianie: {
      motto: 'AI Operating Environment',
      zalety: [
        {
          glowa: 'Cztery środowiska pracy',
          tresc: 'Rozmowa i wiedza, projekty, wytwarzanie oprogramowania oraz praca zespołu modeli.',
        },
        {
          glowa: 'Praca wraca w tym samym stanie',
          tresc: 'Zamknięcie aplikacji nie przerywa zadań — wracają wraz z kartami i kontekstem.',
        },
        {
          glowa: 'Jedno konto na wszystkie urządzenia',
          tresc: 'Ten sam dostęp na komputerze, tablecie i telefonie.',
        },
      ],
    },
    dostep: {
      motto: 'Platforma AI Workspace OS',
      zalety: [
        {
          glowa: 'Cztery środowiska pracy',
          tresc: 'Osobna przestrzeń robocza dla każdego rodzaju pracy: TalkIn, WorkSpace, CodeStudio, MultitaskingAI.',
        },
        {
          glowa: 'Wiele modeli nad jednym zleceniem',
          tresc: 'Koordynator rozdziela zlecenie między wykonawców; przebieg śledzisz w Execution Loop Window.',
        },
        {
          glowa: 'Izolacja kontekstu pod kontrolą',
          /* Zdanie mówi, co program robi, i nie porównuje się z niczym innym:
             okno programu nie jest miejscem na spór z cudzym rozwiązaniem. */
          tresc: 'Zakres kontekstu przekazywanego modelom ustala Operator, osobno dla każdej sesji.',
        },
      ],
    },
    przygotowanie: { motto: 'AI Operating Environment' },
    wydawca: 'Danaco Holding Group Sp. z o.o.',
    wydanie: 'wydanie {wersja} · {data}',
    wsparcie: 'support@danaco-core.pl',
  },

  uruchomienie: {
    nadtytul: 'Uruchomienie',
    etapy: [
      'Połączenie z serwerem',
      'Uzgodnienie wersji',
      'Rozpoznanie urządzenia',
      'Przygotowanie okna',
    ],
    stany: {
      oczekuje: '—',
      wToku: 'w toku',
      gotowe: 'gotowe',
      nawiazane: 'nawiązane',
      zgodna: 'zgodna',
      zaufane: 'zaufane',
      nieudane: 'nieudane',
    },
    laczenie: {
      tytul: 'Nawiązywanie połączenia',
      lid: 'Danaco Console łączy się z serwerem. Poczekaj.',
    },
    token: {
      tytul: 'Przywracanie sesji',
      lid: 'Danaco Console przywraca sesje otwarte na tym koncie.',
      /* Jedno zdanie o pominiętym logowaniu, nie dwa: baner nazywa skutek,
         a fraza pod nim mówi, co Operator ma teraz zrobić. */
      baner: {
        glowa: 'Rozpoznano zaufane urządzenie',
        tresc: 'Logowanie zostanie pominięte.',
      },
      fraza: 'Poczekaj — okno pracy otworzy się samo.',
    },
    blad: {
      tytul: 'Błąd połączenia',
      lid: 'Sprawdź połączenie sieciowe i spróbuj ponownie.',
      baner: {
        glowa: 'Serwer nie odpowiada',
        tresc: 'Ponowna próba za {odliczanie} s.',
      },
      fraza: 'Praca zapisana na serwerze jest bezpieczna.',
    },
    /* Rozjazd wersji protokołu jest stanem, którego okno musi umieć nazwać:
       powitanie jest jedynym miejscem, w którym klient odczytuje wersję rdzenia,
       a wersja niezgodna znaczy, że dalsza rozmowa pójdzie po omacku. */
    wersja: {
      glowa: 'Wersja protokołu serwera jest inna niż wersja tego programu.',
      tresc: 'Wersja serwera: {rdzen}. Wersja programu: {klient}. Zaktualizuj program.',
    },
  },

  dostep: {
    zakladki: ['Zaloguj się', 'Zarejestruj się'],

    logowanie: {
      tytul: 'Zaloguj się',
      lid: 'Podaj dane konta Operatora.',
      login: 'Login albo adres e-mail',
      haslo: 'Hasło',
      reset: { pytanie: 'Nie pamiętasz hasła?', czynnosc: 'Resetuj hasło' },
      metody: {
        naglowek: 'Inne metody logowania',
        pin: 'Kod PIN',
        klucz: 'Klucz systemowy',
        email: 'Kod na adres e-mail',
        aktywna: 'aktywna',
        nieaktywna: 'nieaktywna',
      },
      fraza: { pytanie: 'Nie masz konta?', czynnosc: 'Utwórz konto Operatora' },
    },

    logowanieBlad: {
      baner: {
        glowa: 'Nie rozpoznano danych logowania.',
        tresc: 'Sprawdź login i hasło.',
      },
      capsLock: 'Sprawdź, czy nie jest włączony Caps Lock.',
      fraza:
        'Dostęp można również odzyskać przy użyciu adresu e-mail konta.',
    },

    /* Odsłona zwłoki, nie zapory. Rdzeń nie odmawia kolejnej próby — nakłada na
       nią rosnącą zwłokę, którą okno odmierza i po której samo wraca do
       logowania. Zapisu o pięciu próbach i o godzinie tu nie ma, bo takiej
       reguły nie ma w rdzeniu. */
    logowanieWstrzymane: {
      /* Odliczanie stoi w jednym miejscu — w banerze. Tytuł nazywa sytuację,
         a nie powtarza czasu: ten sam licznik w dwóch miejscach ekranu czyta
         się jak dwa różne terminy. */
      tytul: 'Kolejna próba za chwilę',
      lid: 'Kolejna próba będzie możliwa po odczekaniu chwili.',
      baner: {
        glowa: 'Kolejna próba będzie możliwa za {czas}.',
        tresc: 'Dostęp można również odzyskać przy użyciu adresu e-mail konta.',
      },
      odzyskaj: 'Odzyskaj dostęp',
    },

    odzyskiwanieWstrzymane: {
      tytul: 'Kolejne wysłanie za chwilę',
      lid: 'Poprzednia wiadomość została wysłana przed chwilą.',
      baner: {
        glowa: 'Kolejne wysłanie będzie możliwe za {czas}.',
        tresc:
          'Jeżeli którykolwiek z wysłanych kodów dotarł, wprowadź go — każdy zachowuje ważność przez {minuty} minut od wysłania.',
      },
    },

    rejestracja: {
      tytul: 'Utwórz konto Operatora',
      lid: 'Adres e-mail posłuży do potwierdzenia konta i odzyskania dostępu.',
      login: 'Login',
      email: 'Adres e-mail',
      haslo: 'Hasło',
      hasloPowtorz: 'Powtórz hasło',
      fraza: { pytanie: 'Masz już konto?', czynnosc: 'Zaloguj się' },
    },

    /* Druga gałąź rejestracji — platforma bez
       konta nadawczego zakłada konto i wpuszcza hasłem, ale adresu nikt nie
       sprawdził. Ostrzeżenie jest WYMAGANE: adres jest jedyną drogą odzyskania
       konta, a Operator, który tego nie przeczyta przy rejestracji, dowie się
       w dniu, w którym będzie tego potrzebował. */
    kontoBezPotwierdzenia: {
      tytul: 'Konto Operatora zostało założone',
      lid: 'Do logowania służy hasło ustawione przed chwilą.',
      baner: {
        glowa: 'Adres {adres} pozostaje niepotwierdzony.',
        tresc:
          'Platforma nie ma konta nadawczego, więc wiadomość z kodem potwierdzającym nie została wysłana. Odzyskanie konta pocztą będzie możliwe dopiero po potwierdzeniu adresu.',
      },
      nota:
        'Konto nadawcze ustawia się w oknie Konfiguracji, w kategorii „Konto nadawcze platformy”.',
      wejdz: 'Wejdź do platformy',
    },

    sila: {
      puste: 'Siła hasła zostanie oceniona podczas wpisywania',
      /* Bez ocen w rodzaju „dostateczne” i „dobre”: stopnie mówią, ile warunków
         zostało do spełnienia, a nie jak program ocenia hasło. Wymagane są trzy;
         znak specjalny jest zalecany i podnosi tor, nie rozstrzygając przyjęcia. */
      stopnie: [
        'Hasło nie spełnia żadnego z warunków',
        'Hasło spełnia jeden warunek z trzech wymaganych',
        'Hasło spełnia dwa warunki z trzech wymaganych',
        'Hasło spełnia wymagania',
        'Hasło spełnia wymagania wraz z warunkiem zalecanym',
      ],
      zalecany: ' (zalecane)',
      warunki: {
        dlugosc: 'co najmniej {znaki} znaków',
        wielkosc: 'wielka i mała litera',
        cyfra: 'cyfra',
        znak: 'znak specjalny',
      },
    },

    haslo: {
      pokaz: 'Pokaż hasło',
      ukryj: 'Ukryj hasło',
    },

    sesja: {
      etykieta: 'Nie wylogowuj mnie na tym urządzeniu',
      opis: 'Sesja pozostanie aktywna do czasu wylogowania. Nie zaznaczaj na urządzeniu współdzielonym.',
    },

    kod: {
      tytul: 'Potwierdź adres e-mail',
      lid: 'Na adres {adres} wysłaliśmy kod potwierdzający. Zachowuje ważność przez {minuty} minut.',
      obszar: 'Kod potwierdzający z wiadomości',
      znak: 'Cyfra {numer} z {ile}',
      odliczanie: 'Kod traci ważność za',
      wklej: 'Wklej ze schowka',
      ponow: 'Wyślij ponownie',
      pomoc: {
        glowa: 'Nie ma wiadomości?',
        tresc:
          'Sprawdź folder wiadomości niechcianych. List przychodzi z adresu {nadawca} i dociera zwykle w ciągu minuty. Jeśli nie dotarł, wyślij go ponownie.',
      },
      pomocDane: { nadawca: 'noreply@danaco-core.pl' },
      ostrzezenie: {
        glowa: 'Kod potwierdzający wprowadza się wyłącznie w tym oknie.',
        tresc: 'Danaco Console nigdy nie prosi o niego przez telefon ani w wiadomości zwrotnej.',
      },
      zmienAdres: 'Zmień adres e-mail',
    },

    odzyskiwanie: {
      kroki: ['Adres e-mail', 'Potwierdzenie', 'Nowe hasło'],
      adres: {
        tytul: 'Odzyskaj dostęp do konta',
        lid: 'Podaj adres e-mail konta. Wyślemy na niego kod potwierdzający.',
        pole: 'Adres e-mail',
        ostrzezenie: {
          glowa: 'Zmiana hasła zakończy wszystkie sesje.',
          tresc: 'Na pozostałych urządzeniach będzie konieczne ponowne zalogowanie.',
        },
      },
      haslo: {
        tytul: 'Ustaw nowe hasło',
        lid: 'Nowe hasło zacznie obowiązywać natychmiast. Pozostałe urządzenia zostaną wylogowane.',
        nowe: 'Nowe hasło',
        powtorz: 'Powtórz nowe hasło',
      },
      powrot: 'Wróć do logowania',
    },
  },

  przygotowanie: {
    motto: 'AI Operating Environment',
    obszarPowlok:
      'Cztery powłoki platformy Danaco Console: serwer, środowiska pracy, moduły, interfejs',
    nadtytul: 'Uruchomienie',
    tytul: 'Przygotowanie środowiska pracy',
    lid: 'Danaco Console odtwarza stan pracy z ostatniego zamknięcia: otwarte karty sesji i kontekst projektów.',
    obszarEtapow: 'Postęp przygotowania środowiska pracy',
    etapy: ['Uwierzytelnienie', 'Przywracanie sesji z poprzedniej pracy'],
    miary: {
      rozpoznane: 'urządzenie rozpoznane',
      /* Sama liczba, bez „z ilu”: rdzeń oddaje wyłącznie karty odtworzone i nie
         podaje, ile ich było przed zamknięciem. Zapis „N z N” udawałby miarę
         postępu, która zawsze pokazuje komplet. */
      karty: 'odtworzone karty sesji: {odtworzone}',
      oczekuje: 'oczekuje',
    },
    postep: {
      etykieta: 'Przywracanie sesji',
      opisPaska: 'Postęp przywracania sesji',
    },
    nota: 'Przywracanie trwa po stronie serwera. Zamknięcie okna nie przerywa przywracania.',
  },

  dzialania: {
    zamknijAplikacje: 'Zamknij aplikację',
    ustawieniaPolaczenia: 'Ustawienia połączenia',
    sprobujPonownie: 'Spróbuj ponownie',
    zaloguj: 'Zaloguj się',
    zalogujPonownie: 'Zaloguj się ponownie',
    utworzKonto: 'Utwórz konto',
    potwierdzKonto: 'Potwierdź konto',
    wyslijKod: 'Wyślij kod potwierdzający',
    potwierdzHaslo: 'Potwierdź nowe hasło',
    potwierdzDroge: 'Potwierdź kod',
    wrocDoLogowania: 'Wróć do logowania',
    przerwijIWyloguj: 'Przerwij i wyloguj',
  },

  usterki: {
    wiele: 'Popraw zaznaczone dane.',
    naglowekKonto: 'Nie można utworzyć konta.',
    naglowekHaslo: 'Nie można ustawić hasła.',
    naglowekLogowanie: 'Nie można się zalogować.',
    naglowekKod: 'Nie można wysłać kodu potwierdzającego.',
    naglowekPotwierdzenie: 'Nie można potwierdzić adresu.',
    brakLoginu: {
      glowa: 'Podaj login albo adres e-mail.',
      tresc: 'Pole logowania jest puste.',
    },
    brakHasla: {
      glowa: 'Podaj hasło.',
      tresc: 'Pole hasła jest puste.',
    },
    brakDanych: {
      glowa: 'Podaj dane logowania.',
      tresc: 'Wpisz login albo adres e-mail oraz hasło do konta Operatora.',
    },
    brakAdresu: {
      glowa: 'Podaj adres e-mail konta.',
      tresc: 'Na ten adres wyślemy kod potwierdzający.',
    },
    brakDrogi: {
      glowa: 'Podaj kod potwierdzający z wiadomości.',
      tresc: 'Pole kodu potwierdzającego jest puste.',
    },
    loginZajety: {
      glowa: 'Ten login jest już zajęty.',
      tresc: 'Wybierz inny login. Adres e-mail może pozostać bez zmian.',
    },
    emailBledny: {
      glowa: 'Nieprawidłowy adres e-mail.',
      tresc: 'Sprawdź, czy adres zawiera znak @ oraz nazwę domeny, na przykład nazwa@firma.pl.',
    },
    /* Baner nie powtarza reguły — wykaz pod polem pokazuje ją w całości i
       zaznacza brakujące warunki na czerwono. Zdanie kieruje tam wzrok. */
    hasloSlabe: {
      glowa: 'Hasło nie zostało przyjęte.',
      tresc: 'Brakujące warunki są zaznaczone na czerwono pod polem hasła.',
    },
    haslaRozne: {
      glowa: 'Hasła nie są zgodne.',
      tresc: 'Wpisz to samo hasło w obu polach.',
    },
    /* Zdanie dla Operatora dobierane po kodzie odmowy. Serwer opisuje odmowę
       językiem swojego wnętrza i tę treść okno odkłada wyłącznie do dziennika;
       do widoku idzie zdanie stąd. Każdy kod z katalogu `ErrorCode` musi mieć
       tu pozycję — brak pozycji zostawiłby Operatora bez wyjaśnienia. */
    odmowaSerwera: {
      glowa: 'Nie udało się wykonać czynności.',
      tresc: 'Spróbuj ponownie. Jeżeli problem się powtarza, skontaktuj się z administratorem.',
      wedlugKodu: {
        validation_failed: {
          glowa: 'Podane dane są nieprawidłowe.',
          tresc: 'Popraw zaznaczone pola i spróbuj ponownie.',
        },
        not_found: {
          glowa: 'Nie znaleziono wskazanej pozycji.',
          tresc: 'Odśwież widok i spróbuj ponownie.',
        },
        not_authenticated: {
          glowa: 'Sesja nie jest zalogowana.',
          tresc: 'Zaloguj się ponownie.',
        },
        permission_denied: {
          glowa: 'Brak uprawnień do tej czynności.',
          tresc: 'Skontaktuj się z administratorem.',
        },
        conflict: {
          glowa: 'Nie można wykonać tej czynności w obecnym stanie.',
          tresc: 'Odśwież widok i sprawdź, czy czynność jest nadal potrzebna.',
        },
        channel_unavailable: {
          glowa: 'Model jest chwilowo niedostępny.',
          tresc: 'Spróbuj ponownie za chwilę.',
        },
        rate_limited: {
          glowa: 'Za dużo żądań w krótkim czasie.',
          tresc: 'Odczekaj chwilę i spróbuj ponownie.',
        },
        internal_error: {
          glowa: 'Wystąpił błąd po stronie serwera.',
          tresc: 'Spróbuj ponownie. Jeżeli problem się powtarza, skontaktuj się z administratorem.',
        },
      },
    },
    brakOdpowiedzi: {
      glowa: 'Serwer nie odpowiedział.',
      tresc: 'Połączenie zostało zerwane w trakcie wykonywania. Spróbuj ponownie.',
    },
  },

  komunikaty: {
    zamkniecie: {
      tytul: 'Zamknięcie aplikacji',
      tresc: 'Okno aplikacji zostanie zamknięte. Sesje pracujące w tle trwają dalej i wrócą przy kolejnym uruchomieniu.',
    },
    ustawienia: {
      tytul: 'Ustawienia połączenia',
      tresc: 'Ustawienia połączenia otwierają się w oknie konfiguracji, w sekcji „Sieć i serwer”.',
    },
    kodPonowiony: {
      tytul: 'Kod potwierdzający wysłany ponownie',
      tresc: 'Nowy kod został wysłany na adres konta. Poprzedni przestał obowiązywać.',
    },
    schowek: {
      tytul: 'Schowek niedostępny',
      tresc: 'Wpisz kod potwierdzający ręcznie.',
    },
  },
} as const;

export type KatalogTresci = typeof tresci;
