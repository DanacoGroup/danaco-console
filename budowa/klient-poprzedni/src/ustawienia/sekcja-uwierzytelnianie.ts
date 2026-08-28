import type { AuthChangeReason, AuthMethod } from '../../../shared/contract';
import { opisOdmowyBledu } from '../komponenty/odmowa';
import { przycisk } from '../modele/kontrolki-formularza';
import type { Kanal } from '../protokol/kanal';
import { odczytajSesje } from '../uwierzytelnienie/sesja-bramki';
import type { SekcjaUstawien } from './sekcje';
import { tozsamoscUrzadzenia } from './tozsamosc-urzadzenia';
import { utworzWierszWymoguLogowania, type WierszWymogu } from './wiersz-wymogu-logowania';
import { utworzZrodloBramki, type ZrodloBramki } from './zrodlo-bramki';

/** Sekcja uwierzytelniania w oknie ustawień wykonuje cztery komendy bramki po zalogowaniu i pokazuje wymóg logowania, którego zmiana należy do okna konfiguracji. */
export function utworzSekcjeUwierzytelnianie(kanal: Kanal): SekcjaUstawien {
  const zrodlo: ZrodloBramki = utworzZrodloBramki(kanal);

  const element = document.createElement('div');
  element.className = 'du-sekcja';

  const odpowiedz = document.createElement('p');
  odpowiedz.className = 'du-odpowiedz';
  odpowiedz.hidden = true;

  const wykaz = document.createElement('ul');
  wykaz.className = 'du-metody';

  const oWykazie = document.createElement('p');
  oWykazie.className = 'dn-pole-opis du-granica';
  oWykazie.textContent =
    'Wykaz metod przychodzi z odpowiedzi na czynność oraz ze zdarzenia ' +
    'auth.changed — kontrakt nie ma komendy odczytu metod (auth.method.list). ' +
    'Do pierwszej zmiany wykaz jest pusty i nie znaczy to, że metod nie ma. ' +
    'Czynność wykonana w drugim oknie dolatuje tu sama, bez odświeżania.';

  const wymog: WierszWymogu = utworzWierszWymoguLogowania(kanal);

  function powiedz(tresc: string, powodzenie: boolean): void {
    odpowiedz.textContent = tresc;
    odpowiedz.hidden = tresc === '';
    odpowiedz.dataset['powodzenie'] = String(powodzenie);
  }

  function pokazMetody(metody: readonly AuthMethod[]): void {
    wykaz.replaceChildren(
      ...metody.map((metoda) => {
        const wpis = document.createElement('li');
        wpis.className = 'du-metoda';
        wpis.dataset['metoda'] = metoda.kind;

        const nazwa = document.createElement('span');
        nazwa.className = 'du-metoda__nazwa';
        nazwa.textContent = metoda.label ?? metoda.kind;
        wpis.append(nazwa);

        if (metoda.anchor) {
          // Kotwicy zdjąć się nie da; zamiast przycisku pewnej odmowy stoi zdanie, dlaczego go nie ma.
          const kotwica = document.createElement('span');
          kotwica.className = 'dn-plakietka du-metoda__kotwica';
          kotwica.textContent = 'kotwica bramki — nie do zdjęcia';
          wpis.append(kotwica);
          return wpis;
        }

        const zdejmij = przycisk('Zdejmij', 'dn-btn dn-btn--sm dn-btn--niebezpieczny');
        zdejmij.addEventListener('click', () => void zdejmijMetode(metoda));
        wpis.append(zdejmij);
        return wpis;
      }),
    );
  }

  const pin = poleTajne('PIN urządzenia', 'du-pole-pin');
  const nazwaUrzadzenia = poleJawne('Nazwa urządzenia', 'du-pole-urzadzenie');
  const zaloz = przycisk('Załóż PIN na tym urządzeniu', 'dn-btn dn-btn--sm dn-btn--atrament');

  const haslo = poleTajne('Hasło bieżące', 'du-pole-haslo');
  const noweHaslo = poleTajne('Hasło nowe', 'du-pole-haslo');
  const zmien = przycisk('Zmień hasło bramki', 'dn-btn dn-btn--sm');

  const przedluz = przycisk('Przedłuż sesję bramki', 'dn-btn dn-btn--sm');

  // Powód niedostępności metody systemowej — słowami rdzenia, cytatem zamiast własnego zdania.
  const oHello = document.createElement('p');
  oHello.className = 'dn-pole-opis du-granica';
  oHello.textContent =
    'Windows Hello jest dziś niedostępne, a rdzeń odmawia tak: „metoda hello ' +
    '(Windows Hello przez WebAuthn) nie jest zbudowana: rozstrzygnięcie ' +
    'Właściciela z 13.08.2026 odkłada ją do chwili wystawienia platformy na ' +
    'serwerze, a WebAuthn i tak wywodzi rp_id z pochodzenia dokumentu — pod ' +
    'http://127.0.0.1 poprawnego rp_id nie ma. Dziś działają hasło i PIN". ' +
    'Samo wystawienie pod domeną zatem nie wystarczy — metody nie ma po stronie ' +
    'rdzenia.';

  element.append(
    wykaz,
    oWykazie,
    pin,
    nazwaUrzadzenia,
    zaloz,
    haslo,
    noweHaslo,
    zmien,
    przedluz,
    odpowiedz,
    wymog.element,
    oHello,
  );

  async function zalozPin(): Promise<void> {
    const sekret = pin.value.trim();
    const urzadzenie = tozsamoscUrzadzenia();
    powiedz('Zakładanie PIN-u…', true);
    const wynik = await zrodlo.zalozMetode({
      kind: 'pin',
      deviceId: urzadzenie,
      ...(sekret === '' ? {} : { secret: sekret }),
      ...(nazwaUrzadzenia.value.trim() === ''
        ? {}
        : { deviceName: nazwaUrzadzenia.value.trim() }),
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      powiedz(opisOdmowyBledu('Założenie PIN-u', wynik.blad), false);
      return;
    }
    pin.value = '';
    pokazMetody(wynik.wynik.methods);
    powiedz('PIN założony na tym urządzeniu.', true);
  }

  async function zdejmijMetode(metoda: AuthMethod): Promise<void> {
    powiedz('Zdejmowanie metody…', true);
    const wynik = await zrodlo.zdejmijMetode(metoda.id, metoda.deviceId);
    if (!wynik.udany || wynik.wynik === undefined) {
      powiedz(opisOdmowyBledu('Zdjęcie metody', wynik.blad), false);
      return;
    }
    pokazMetody(wynik.wynik.methods);
    // Rozstrzyga odpowiedź rdzenia, nie samo powodzenie wywołania, bo rdzeń może metody nie zdjąć.
    powiedz(
      wynik.wynik.removed
        ? 'Metoda zdjęta z urządzenia.'
        : 'Rdzeń przyjął wywołanie, ale metody nie zdjął — zdjęcie się nie odbyło.',
      wynik.wynik.removed,
    );
  }

  async function zmienHaslo(): Promise<void> {
    powiedz('Zmiana hasła bramki…', true);
    const wynik = await zrodlo.zmienHaslo(haslo.value, noweHaslo.value);
    if (!wynik.udany || wynik.wynik === undefined) {
      powiedz(opisOdmowyBledu('Zmiana hasła', wynik.blad), false);
      return;
    }
    haslo.value = '';
    noweHaslo.value = '';
    const uniewaznione = wynik.wynik.revokedSessions ?? 0;
    powiedz(
      wynik.wynik.changed
        ? `Hasło bramki zmienione. Sesji unieważnionych: ${uniewaznione}.`
        : 'Rdzeń przyjął wywołanie, ale hasła nie zmienił.',
      wynik.wynik.changed,
    );
  }

  async function przedluzSesje(): Promise<void> {
    powiedz('Przedłużanie sesji bramki…', true);
    const wynik = await zrodlo.przedluzSesje(odczytajSesje()?.token ?? '');
    if (!wynik.udany || wynik.wynik === undefined) {
      powiedz(opisOdmowyBledu('Przedłużenie sesji', wynik.blad), false);
      return;
    }
    powiedz('Sesja bramki przedłużona.', true);
  }

  zaloz.addEventListener('click', () => void zalozPin());
  zmien.addEventListener('click', () => void zmienHaslo());
  przedluz.addEventListener('click', () => void przedluzSesje());

  // Zmiana metod wykonana gdzie indziej dolatuje tutaj zdarzeniem niosącym komplet metod.
  const odsubskrybuj = zrodlo.naZmianeBramki((zmiana) => {
    if (zmiana.methods !== undefined) pokazMetody(zmiana.methods);
    powiedz(zdanieZmiany(zmiana.reason), true);
  });

  return {
    element,

    // Odświeżenie nie ma czego odczytać w części metod i mówi to wprost, zamiast milczeć.
    odswiez() {
      wymog.odswiez();
      powiedz(
        'Wykazu metod nie ma czym odczytać: kontrakt nie niesie komendy odczytu ' +
          'metod bramki. Wykaz odświeży się sam przy najbliższej zmianie — także ' +
          'wykonanej w innym oknie. Wymóg logowania odczytano ponownie.',
        true,
      );
    },

    rozlacz: () => {
      odsubskrybuj();
      wymog.rozlacz();
    },
  };
}

/** Zdanie o zmianie przyniesionej zdarzeniem, bez którego zmiana wykonana w innym oknie byłaby niema dla operatora. */
function zdanieZmiany(powod: AuthChangeReason): string {
  switch (powod) {
    case 'methodAdded':
      return 'Metoda wejścia założona — wykaz odświeżony zdarzeniem rdzenia.';
    case 'methodRemoved':
      return 'Metoda wejścia zdjęta — wykaz odświeżony zdarzeniem rdzenia.';
    case 'passwordChanged':
      return 'Hasło bramki zmienione — zgłosił to rdzeń zdarzeniem auth.changed.';
    case 'sessionRevoked':
      return 'Sesja bramki unieważniona — zgłosił to rdzeń zdarzeniem auth.changed.';
    default:
      // Powód spoza kontraktu nie wywraca sekcji ani nie jest przemilczany, nazywamy go dosłownie.
      return `Rdzeń zgłosił zmianę stanu uwierzytelnienia: ${String(powod)}.`;
  }
}

/** Pole sekretu jest osobną wytwórnią, bo typ pola rozstrzyga o widoczności wpisywanej treści na ekranie. */
function poleTajne(etykieta: string, klasa: string): HTMLInputElement {
  const pole = document.createElement('input');
  pole.type = 'password';
  pole.className = `dn-pole ${klasa}`;
  pole.placeholder = etykieta;
  pole.setAttribute('aria-label', etykieta);
  return pole;
}

function poleJawne(etykieta: string, klasa: string): HTMLInputElement {
  const pole = document.createElement('input');
  pole.type = 'text';
  pole.className = `dn-pole ${klasa}`;
  pole.placeholder = etykieta;
  pole.setAttribute('aria-label', etykieta);
  return pole;
}

// Tożsamość urządzenia mieszka w osobnym pliku; sekcja urządzeń czyta tę samą wartość.
