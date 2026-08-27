import type { Kanal } from '../protokol/kanal';
import { podepnijMostMotywu, type MostMotywu } from './most-motywu';
import type { SekcjaUstawien } from './sekcje';
import { utworzWierszNastawy, type WierszNastawy } from './wiersz-nastawy';

/** Sekcja wygląd i język niesie dziś jedną nastawę motywu ze sterem, dzieloną z paskiem górnym, i jedno nazwane miejsce puste dla języka interfejsu, którego kontrakt nie obsługuje. */
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

  // Buduje wiersz motywu z definicji katalogu raz, przy pierwszym udanym odczycie.
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

  // Most rozgłasza każdą zmianę wartości: własną, z paska i z drugiego okna; sekcja nadąża nasłuchem.
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

    // Most przeżywa zamknięcie okna, bo motyw obowiązuje także wtedy, gdy okno stoi zamknięte.
    rozlacz: () => odsubskrybuj(),
  };
}
