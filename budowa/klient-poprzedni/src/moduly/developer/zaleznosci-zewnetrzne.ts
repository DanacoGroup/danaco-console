/**
 * Programy spoza instalki, na których stoją albo stanęłyby czynności modułu
 * Developer, wraz ze skutkiem ich braku dla operatora.
 */

/** Jedna zależność zewnętrzna: program, czynność na nim stojąca i skutek jego braku dla operatora okna. */
export interface ZaleznoscZewnetrzna {
  /** Klucz stały pozycji — po nim okno bierze swój podzbiór wykazu. */
  kod: string;
  /** Nazwa programu w brzmieniu wywołania. */
  program: string;
  /** Czynność modułu, która ten program uruchamia. */
  czynnosc: string;
  /** Co się dzieje, gdy programu nie ma na maszynie rdzenia. */
  skutekBraku: string;
  /** Czy czynność stoi już dziś, czy dopiero czekałaby na komendę. */
  stan: 'w-uzyciu' | 'po-dobudowie-komendy';
}

/**
 * Wykaz zależności modułu. Dwie pierwsze pozycje są czynne dzisiaj, pozostałe
 * czekają na komendy przyszłe.
 */
export const ZALEZNOSCI_ZEWNETRZNE: readonly ZaleznoscZewnetrzna[] = [
  {
    kod: 'git',
    program: 'git',
    czynnosc: 'Każda czynność Git Panelu — rdzeń składa z niej wiersz polecenia i uruchamia git.',
    skutekBraku:
      'Czynność wraca odmową wykonania z treścią systemu operacyjnego o nieodnalezionym programie. ' +
      'Repozytorium pozostaje nietknięte, a panel nie pozna gałęzi ani stanu zmian.',
    stan: 'w-uzyciu',
  },
  {
    kod: 'zadanie-budowania',
    program: 'program wpisany w pole „Zadanie” (na przykład go, npm, pytest, cargo, make)',
    czynnosc:
      'Uruchomienie budowania — pierwsze słowo zadania jest programem, reszta jego parametrami.',
    skutekBraku:
      'Przebieg nie rusza: rdzeń odpowiada odmową uruchomienia z treścią o nieodnalezionym ' +
      'programie, log pozostaje pusty i żaden przebieg nie trafia do dziennika.',
    stan: 'w-uzyciu',
  },
  {
    kod: 'serwer-jezyka',
    program:
      'serwer języka właściwy plikowi (gopls, typescript-language-server, pyright, rust-analyzer, clangd, jdtls)',
    czynnosc:
      'Podpowiedzi sygnatur, przejście do definicji, wykaz wystąpień, diagnostyka na żywo ' +
      'i refaktoryzacje semantyczne Code Editora.',
    skutekBraku:
      'Edytor zostaje polem tekstowym: treść pliku odczytuje i zapisuje bez zmian, ale żadna ' +
      'z powyższych czynności nie ma czego zapytać.',
    stan: 'po-dobudowie-komendy',
  },
  {
    kod: 'formater-linter',
    program:
      'formater i linter repozytorium (gofmt, prettier, black, rustfmt, golangci-lint, eslint, ruff)',
    czynnosc: 'Formatowanie dokumentu, formatowanie przy zapisie oraz podkreślenia analizy statycznej.',
    skutekBraku:
      'Treść zostaje w postaci wpisanej przez Operatora, a pasek statusu edytora nie ma skąd wziąć ' +
      'liczby zgłoszeń — pole liczby błędów mówi wtedy o braku pomiaru, nie o zerze błędów.',
    stan: 'po-dobudowie-komendy',
  },
  {
    kod: 'ripgrep',
    program: 'ripgrep',
    czynnosc: 'Wyszukiwanie wzorca w całym repozytorium (Grep) wraz z zamianą masową.',
    skutekBraku:
      'Wyszukiwanie globalne nie ma silnika; zostaje wyszukiwanie w treści pliku otwartego ' +
      'w edytorze, wykonywane po stronie przeglądarki.',
    stan: 'po-dobudowie-komendy',
  },
  {
    kod: 'adapter-dap',
    program: 'adapter debugowania właściwy językowi (delve, debugpy, js-debug, lldb, gdb)',
    czynnosc:
      'Sesja debugowania krokowego, punkty przerwania, podgląd zmiennych, stos wywołań i konsola wyrażeń.',
    skutekBraku:
      'Sesji debugowania nie da się rozpocząć. Punkty przerwania ustawione na marginesie pozostają ' +
      'oznaczeniem w edytorze, którego nikt nie wykona.',
    stan: 'po-dobudowie-komendy',
  },
  {
    kod: 'silnik-kontenerow',
    program: 'silnik kontenerów (Docker albo Podman) wraz z jego gniazdem',
    czynnosc:
      'Wykaz kontenerów i obrazów, uruchomienie i zatrzymanie kontenera, budowanie obrazu, ' +
      'stos usług oraz kontener deweloperski.',
    skutekBraku:
      'Zakładka Containers nie ma z czym rozmawiać: wykaz jest pusty nie dlatego, że kontenerów ' +
      'nie ma, lecz dlatego, że nie ma kogo o nie zapytać — a to dwa różne zdania.',
    stan: 'po-dobudowie-komendy',
  },
  {
    kod: 'serwer-bazy',
    program: 'serwer bazy danych wraz ze sterownikiem (PostgreSQL, MySQL, SQLite)',
    czynnosc: 'Przeglądarka schematu, konsola zapytań, edycja danych w siatce i migracje schematu.',
    skutekBraku:
      'Połączenie nie powstaje, a zakładka Data Console nie pokaże ani jednej tabeli. Poświadczenia ' +
      'połączenia podlegają zakresowi „konto i token per sesja” okna konfiguracji punktów izolacji.',
    stan: 'po-dobudowie-komendy',
  },
  {
    kod: 'skanery-bezpieczenstwa',
    program: 'skanery zależności i kodu (syft, grype, gitleaks, semgrep)',
    czynnosc:
      'Wykaz podatności zależności, skan sekretów przed zatwierdzeniem zmian, analiza wzorców ' +
      'podatności w kodzie i zestawienie licencji.',
    skutekBraku:
      'Zakładka Dependencies & Security nie zgłosi ani jednej pozycji. Brak zgłoszeń bez skanera ' +
      'nie znaczy repozytorium czystego — znaczy repozytorium niesprawdzone.',
    stan: 'po-dobudowie-komendy',
  },
];

/** Pozycje wykazu o wskazanych kodach, zwrócone w kolejności samego wykazu zależności modułu Developer. */
export function zaleznosci(kody: readonly string[]): readonly ZaleznoscZewnetrzna[] {
  const wybrane = new Set(kody);
  return ZALEZNOSCI_ZEWNETRZNE.filter((pozycja) => wybrane.has(pozycja.kod));
}

/**
 * Rysuje wykaz zależności zewnętrznych okna.
 *
 * Wykaz stoi w treści okna, nie w podpowiedzi przycisku: podpowiedź czyta ten,
 * kto już najechał na kontrolkę, a o wymaganym programie trzeba wiedzieć
 * wcześniej — przy planowaniu pracy, nie przy jej odmowie.
 */
export function rysujZaleznosci(pozycje: readonly ZaleznoscZewnetrzna[]): HTMLElement {
  const panel = document.createElement('div');
  panel.className = 'mdev-zaleznosci';

  const naglowek = document.createElement('p');
  naglowek.className = 'mdev-zaleznosci__naglowek';
  naglowek.textContent =
    'Programy spoza instalki Danaco Console, na których stoją czynności tego okna. ' +
    'Klient nie sprawdza, czy są na maszynie — mówi, czego wymagają i co się stanie bez nich.';
  panel.append(naglowek);

  const lista = document.createElement('ul');
  lista.className = 'mdev-wykaz';
  lista.setAttribute('aria-label', 'Zależności zewnętrzne okna');
  for (const pozycja of pozycje) {
    const wiersz = document.createElement('li');
    wiersz.className = 'mdev-pozycja';
    wiersz.dataset['zaleznosc'] = pozycja.stan;

    const nazwa = document.createElement('strong');
    nazwa.className = 'mdev-pozycja__tytul';
    nazwa.textContent = pozycja.program;

    const opis = document.createElement('span');
    opis.className = 'mdev-pozycja__opis';
    opis.textContent =
      `${pozycja.czynnosc} Bez tego programu: ${pozycja.skutekBraku}` +
      (pozycja.stan === 'w-uzyciu'
        ? ''
        : ' Czynność czeka ponadto na komendę, której kontrakt jeszcze nie niesie.');

    wiersz.append(nazwa, opis);
    lista.append(wiersz);
  }
  panel.append(lista);
  return panel;
}
