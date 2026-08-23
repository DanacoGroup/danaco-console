import type { Kanal } from '../protokol/kanal';
import { podepnijMostMotywu, type MostMotywu } from './most-motywu';
import type { SekcjaUstawien } from './sekcje';
import { utworzWierszNastawy, type WierszNastawy } from './wiersz-nastawy';

/**
 * Sekcja „Wygląd i język" — dziś jedna nastawa ze sterem i jedno nazwane
 * miejsce puste.
 *
 * Motyw to ta sama nastawa, co na pasku górnym. Wiersz nie ma własnego stanu
 * ani własnego zapisu: bierze most (`most-motywu.ts`), jedynego właściciela tej
 * nastawy po stronie klienta, i jest jego drugim sterem — pierwszym jest
 * przełącznik na pasku. Zmiana tutaj przestawia pasek, zmiana na pasku
 * przestawia ten wiersz, a zmiana w drugim oknie dolatuje do obu zdarzeniem
 * `config.changed`.
 *
 * Opcje wyboru (`''` — preferencja systemu, `light`, `dark`) i ich etykiety
 * przychodzą z katalogu rdzenia, nie z tego pliku.
 *
 * Język interfejsu jest tu miejscem nazwanym, nie kontrolką. Katalog ustawień
 * rdzenia (`settings.definition.list`) nie niesie klucza języka interfejsu —
 * jedyny klucz językowy, `mowa_jezyk`, dotyczy rozpoznawania mowy. Kontrakt nie
 * ma komendy zmiany języka, a klient nie ma warstwy tłumaczeń, więc przełącznik
 * nie miałby dokąd pójść. Zamiast niego stoi zdanie mówiące, czego brakuje.
 */
export function utworzSekcjeWyglad(kanal: Kanal): SekcjaUstawien {
  const most: MostMotywu = podepnijMostMotywu(kanal);

  const element = document.createElement('div');
  element.className = 'du-sekcja';

  const miejsceMotywu = document.createElement('div');
  miejsceMotywu.className = 'du-miejsce-wiersza';

  const stanOdczytu = document.createElement('p');
  stanOdczytu.className = 'du-odpowiedz';
  stanOdczytu.textContent = 'Odczyt katalogu ustawień rdzenia…';

  const oJezyku = document.createElement('p');
  oJezyku.className = 'dn-pole-opis du-granica';
  oJezyku.textContent =
    'Języka interfejsu nie ma czym zmienić: katalog ustawień rdzenia nie niesie ' +
    'ani jednego klucza języka interfejsu (jedyny klucz językowy, mowa_jezyk, ' +
    'dotyczy rozpoznawania mowy), kontrakt nie ma komendy zmiany języka, a klient ' +
    'nie ma warstwy tłumaczeń. Kontrolki tu nie ma, bo nie miałaby dokąd pójść.';

  element.append(stanOdczytu, miejsceMotywu, oJezyku);

  let wiersz: WierszNastawy | null = null;

  /**
   * Buduje wiersz motywu z definicji katalogu — raz, przy pierwszym udanym
   * odczycie. Kolejne zmiany wartości wchodzą przez `ustaw`, żeby przerysowanie
   * nie zwijało wykazu pod ręką Operatora ani nie gubiło ogniska.
   */
  function pokaz(): void {
    const definicja = most.definicja();
    const odmowa = most.odmowa();

    if (definicja === undefined) {
      stanOdczytu.hidden = false;
      stanOdczytu.dataset['powodzenie'] = String(odmowa === '');
      stanOdczytu.textContent =
        odmowa === '' ? 'Odczyt katalogu ustawień rdzenia…' : odmowa;
      return;
    }

    stanOdczytu.hidden = true;

    if (wiersz === null) {
      wiersz = utworzWierszNastawy({
        definicja,
        wartosc: most.wybor(),
        wykonaj: (wartosc) => most.ustaw(wartosc === 'light' || wartosc === 'dark' ? wartosc : ''),
      });
      miejsceMotywu.replaceChildren(wiersz.element);
    } else {
      wiersz.ustaw(most.wybor());
      wiersz.zdanie(odmowa);
    }
  }

  // Most rozgłasza każdą zmianę wartości: własną, z paska i z drugiego okna.
  // Sekcja nie pyta o nią ponownie — nadąża nasłuchem.
  const odsubskrybuj = most.naZmiane(() => pokaz());
  pokaz();

  return {
    element,

    odswiez() {
      stanOdczytu.hidden = false;
      stanOdczytu.dataset['powodzenie'] = 'true';
      stanOdczytu.textContent = 'Odczyt katalogu ustawień rdzenia…';
      void most.odczytaj().then(() => pokaz());
    },

    // Most przeżywa zamknięcie okna, bo motyw obowiązuje także wtedy, gdy okno
    // stoi zamknięte — sekcja zdejmuje wyłącznie własny nasłuch, nie most.
    rozlacz: () => odsubskrybuj(),
  };
}
