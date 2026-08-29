/* Katalog treści kreatora instalacji zawiera wyłącznie tekst widoczny dla użytkownika w kluczach hierarchicznych od obszaru do elementu, zapisany jako czysty JSON w skrypcie, bo okno bywa otwierane wprost z dysku. */
window.DanacoKreator = window.DanacoKreator || {};
window.DanacoKreator.tresci = {
  "_opis": "Katalog łańcuchów kreatora instalacji Danaco Console. Jedyne miejsce z tekstem widocznym dla użytkownika. Klucze idą hierarchicznie: obszar → element. Plik przekazuje się tłumaczowi bez dostępu do kodu — nie ma tu znaczników ani logiki. Wartości z {nawiasami} to miejsca na dane podstawiane w czasie działania; ich nazw nie tłumaczy się.",
  "_jezyk": "pl-PL",

  "okno": {
    "tytul": "Instalator programu Danaco Console",
    "zwin": "Zwiń okno",
    "rozwin": "Rozwiń okno",
    "zamknij": "Zamknij instalator"
  },

  "szyna": {
    "marka": "Danaco Console",
    "etykieta": "Kroki instalacji",
    "kroki": [
      "Wymagania",
      "Licencja",
      "Wersja programu",
      "Lokalizacja",
      "Instalacja",
      "Podsumowanie"
    ]
  },

  "bryla": {
    "etykieta": "Budowa modułu Danaco Console"
  },

  "krok1": {
    "nadtytul": "Krok 1 z 6",
    "tytul": "Kreator instalacji Danaco Console",
    "podtytul": "Przygotowanie do instalacji oprogramowania na tym urządzeniu",
    "akapit1": "Instalacja programu Danaco Console przebiega w sześciu krokach i trwa około trzech minut.",
    "akapit2": "Kolejne kroki obejmują akceptację warunków licencji, wybór wersji programu zgodnej z procesorem oraz wskazanie katalogu instalacji.",
    "wymagania": {
      "naglowek": "Wymagania i parametry instalacji",
      "wersja": { "etykieta": "Wersja programu", "wartosc": "2.0 (kompilacja 2026.08)" },
      "system": { "etykieta": "Wymagany system", "wartosc": "Windows 11 w wersji 22H2 lub nowszej" },
      "procesor": { "etykieta": "Obsługiwane procesory", "wartosc": "Intel, AMD, ARM" },
      "miejsce": { "etykieta": "Wymagane miejsce", "wartosc": "ok. 250 MB" },
      "uprawnienia": { "etykieta": "Uprawnienia", "wartosc": "konto użytkownika" }
    },
    "polaczenie": {
      "glowa": "Wymagane połączenie z internetem",
      "tresc": "Instalator pobiera składniki programu z serwera Danaco. Przerwanie połączenia zatrzyma instalację i cofnie zmiany."
    },
    "nawigacja": "Aby kontynuować, kliknij przycisk Dalej."
  },

  "krok2": {
    "nadtytul": "Krok 2 z 6",
    "tytul": "Umowa licencyjna",
    "podtytul": "Zapoznaj się z warunkami korzystania z programu Danaco Console.",
    "dokument": { "etykieta": "Warunki licencji — treść dokumentu" },
    "zgoda": {
      "etykieta": "Akceptuję warunki umowy licencyjnej",
      "opis": "Akceptacja warunków jest konieczna do przeprowadzenia instalacji."
    },
    "nawigacja": "Aby kontynuować, kliknij przycisk Dalej.",
    "nawigacjaBrakZgody": "Aby przejść dalej, zaakceptuj warunki umowy licencyjnej."
  },

  "krok3": {
    "nadtytul": "Krok 3 z 6",
    "tytul": "Wersja dla Twojego procesora",
    "podtytulX64": "Ten komputer ma procesor Intel lub AMD (x64). Zaznaczona jest wersja odpowiednia dla tego procesora.",
    "podtytulArm": "Ten komputer ma procesor ARM (ARM64). Zaznaczona jest wersja odpowiednia dla tego procesora.",
    "podtytulBrak": "Nie udało się rozpoznać procesora tego komputera. Wybierz wersję, którą chcesz zainstalować.",
    "wskazowkaWykryta": "Jeśli nie wiesz, którą wersję wybrać, zostaw zaznaczenie bez zmian. Wersja niezgodna z procesorem nie uruchomi się po instalacji.",
    "wskazowkaBrak": "Otwórz Ustawienia > System > Informacje i sprawdź pozycję Typ systemu. Wartość „x64” oznacza pierwszą opcję, wartość „ARM64” — drugą.",
    "wybor": { "etykieta": "Wersja programu" },
    "plakietka": "ZALECANE",
    "x64": {
      "nazwa": "Intel lub AMD (x64)",
      "opis": "Procesory Intel Core, Intel Celeron, Pentium oraz AMD Ryzen i Athlon. Ta wersja pasuje do niemal wszystkich komputerów z systemem Windows."
    },
    "arm": {
      "nazwa": "ARM (ARM64)",
      "opis": "Procesory Snapdragon X Elite, Snapdragon X Plus i Microsoft SQ. Występują w komputerach Copilot+\u00a0PC oraz w Surface Pro\u00a0X."
    },
    "niezgodnosc": "Wybrana wersja nie odpowiada procesorowi tego komputera. Program zainstaluje się, ale nie uruchomi.",
    "pomoc": "Jak sprawdzić procesor tego komputera?"
  },

  "krok4": {
    "nadtytul": "Krok 4 z 6",
    "tytul": "Lokalizacja i skróty",
    "podtytul": "Wskaż katalogi instalacji i określ sposób uruchamiania programu",
    "katalogProgramu": {
      "etykieta": "Katalog programu",
      "wartosc": "C:\\Users\\Operator\\AppData\\Local\\Programs\\Danaco Console",
      "opis": "Lokalizacja plików wykonywalnych. Instalacja w katalogu profilu użytkownika nie wymaga uprawnień administratora."
    },
    "katalogDanych": {
      "etykieta": "Katalog danych",
      "wartosc": "C:\\Users\\Operator\\AppData\\Roaming\\Danaco Console",
      "opis": "Lokalizacja projektów, ustawień i dzienników pracy. Dane pozostają na dysku po odinstalowaniu programu."
    },
    "zmien": "Zmień…",
    "miejsce": "Wymagane {wymagane} · dostępne {dostepne} na dysku {dysk}",
    "miejsceDane": { "wymagane": "250 MB", "dostepne": "84,2 GB", "dysk": "C:" },
    "skroty": {
      "naglowek": "Skróty",
      "pulpit": {
        "etykieta": "Utwórz skrót na pulpicie",
        "opis": "Umożliwia uruchamianie programu bez otwierania menu Start."
      },
      "start": {
        "etykieta": "Dodaj do menu Start",
        "opis": "Udostępnia program na liście aplikacji i w wyszukiwaniu systemowym."
      }
    }
  },

  "krok5": {
    "nadtytul": "Krok 5 z 6",
    "przebieg": {
      "tytul": "Instalowanie programu Danaco Console",
      "podtytul": "Zamknięcie okna przerwie instalację i cofnie wprowadzone zmiany"
    },
    "wycofywanie": {
      "tytul": "Cofanie zmian",
      "podtytul": "Instalator usuwa pliki i wpisy utworzone podczas instalacji."
    },
    "blad": {
      "tytul": "Nie udało się ukończyć instalacji",
      "podtytul": "Instalacja zatrzymała się na etapie rejestrowania składników. Zmiany zostały cofnięte — w komputerze nie pozostały pliki programu.",
      "szczegoly": "Kod błędu: 0x80070005\nOdmowa dostępu do katalogu C:\\Users\\Operator\\AppData\\Local\\Programs\\Danaco Console",
      "rada": "Zamknij inne programy i uruchom instalator ponownie jako administrator."
    },
    "postep": { "etykieta": "Postęp instalacji", "opisPaska": "Postęp instalacji" },
    "etapyNaglowek": "Przebieg instalacji",
    "stany": { "gotowe": "Gotowe", "wToku": "W toku", "oczekuje": "Oczekuje" },
    "etapy": [
      {
        "nazwa": "Sprawdzanie wymagań",
        "opis": "Weryfikacja wersji systemu, architektury procesora i miejsca na dysku."
      },
      {
        "nazwa": "Rozpakowywanie plików",
        "opis": "Zapis plików programu w katalogu docelowym.",
        "licznik": { "czasownik": "Rozpakowano", "cel": "250", "jednostka": "MB" },
        "pozostalo": "Pozostało około 2 minut"
      },
      {
        "nazwa": "Rejestrowanie składników",
        "opis": "Rejestracja bibliotek i skojarzeń plików w systemie.",
        "licznik": { "czasownik": "Zarejestrowano", "cel": "22", "rzecz": "składników" },
        "pozostalo": "Pozostało około minuty"
      },
      {
        "nazwa": "Weryfikacja podpisu",
        "opis": "Kontrola integralności zainstalowanych plików."
      }
    ],
    "numerEtapu": "Etap {numer} z {ile}"
  },

  "krok6": {
    "nadtytul": "Krok 6 z 6",
    "tytulGotowe": "Instalacja ukończona",
    "tytulOstrzezenia": "Instalacja ukończona z ostrzeżeniami",
    "podtytulGotowe": "Program Danaco Console jest gotowy do pracy",
    "podtytulOstrzezenia": "Program jest gotowy do pracy, ale nie wszystkie czynności powiodły się.",
    "ostrzezenie": {
      "glowa": "Nie utworzono skrótu na pulpicie",
      "tresc": "Program znajdziesz w menu Start."
    },
    "pierwszeUruchomienie": {
      "naglowek": "Przy pierwszym uruchomieniu",
      "tresc": "Pierwsze uruchomienie obejmuje utworzenie konta operatora, wybór modeli oraz wprowadzenie kluczy dostępu do kont usług."
    },
    "szczegoly": {
      "naglowek": "Szczegóły instalacji",
      "lokalizacja": { "etykieta": "Lokalizacja", "wartosc": "C:\\Users\\Operator\\AppData\\Local\\Programs\\Danaco Console" },
      "wersja": { "etykieta": "Wersja programu", "wartosc": "2.0 (kompilacja 2026.08)" },
      "procesor": { "etykieta": "Procesor" }
    },
    "przewodnik": {
      "etykieta": "Otwórz przewodnik konfiguracji",
      "opis": "Przewodnik prowadzi przez ustawienia wymagane do rozpoczęcia pracy."
    },
    "odinstalowanie": "Program można odinstalować w Ustawieniach systemu Windows, w sekcji Aplikacje."
  },

  "dzialania": {
    "anuluj": "Anuluj",
    "wstecz": "Wstecz",
    "dalej": "Dalej",
    "instaluj": "Instaluj",
    "zakoncz": "Zakończ",
    "uruchom": "Uruchom Danaco Console",
    "zamknij": "Zamknij",
    "sprobujPonownie": "Spróbuj ponownie",
    "kopiujSzczegoly": "Kopiuj szczegóły"
  },

  "dialogi": {
    "zamkniecie": {
      "tytul": "Zamknąć instalator?",
      "tresc": "Instalacja nie została rozpoczęta i żadne zmiany nie zostały wprowadzone na tym komputerze. Możesz uruchomić instalator ponownie w dowolnym momencie.",
      "zostan": "Wróć do instalacji",
      "wyjdz": "Zamknij instalator"
    },
    "anulowanie": {
      "tytul": "Anulować instalację?",
      "tresc": "Instalacja nie zostanie ukończona. Instalator cofnie zmiany wprowadzone do tej pory, więc w komputerze nie pozostaną pliki programu.",
      "zostan": "Kontynuuj instalację",
      "wyjdz": "Anuluj instalację"
    },
    "niezgodnosc": {
      "tytul": "Kontynuować z wybraną wersją?",
      "tresc": "Ten komputer ma procesor {wykryty}, a wybrana została wersja {wybrany}. Program zainstaluje się, ale nie uruchomi.",
      "popraw": "Zmień na zgodną wersję",
      "mimoTo": "Kontynuuj mimo to"
    },
    "pomocProcesor": {
      "tytul": "Jak sprawdzić procesor tego komputera",
      "kroki": [
        "Naciśnij klawisze Windows + I, aby otworzyć Ustawienia.",
        "Przejdź do sekcji System, a następnie wybierz Informacje.",
        "Odszukaj pozycję Typ systemu."
      ],
      "x64": "Wartość „System operacyjny 64-bitowy, procesor x64” oznacza wersję Intel lub AMD (x64).",
      "arm": "Wartość „System operacyjny 64-bitowy, procesor ARM” oznacza wersję ARM (ARM64).",
      "zamknij": "Zamknij"
    }
  }
};
