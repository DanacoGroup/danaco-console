import { opisOdmowyBledu, type PowodOdmowy } from '../komponenty/odmowa';

/**
 * Trzecia część odmowy nakładki: czym Operator to zmieni.
 *
 * Dwie pierwsze części — co się nie udało i dlaczego — składa
 * `komponenty/odmowa.ts`. Tu dochodzi rada zależna od czynności nakładki,
 * której wspólny komponent znać nie może.
 *
 * Część odmów rdzenia w rodzinie `aod.*` jest zachowaniem poprawnym, nie
 * usterką: odpięcie procesu spoza wykazu przypiętych zwraca `not_found`,
 * a puste `processId` przy przypięciu albo odpięciu — `validation_failed`.
 * Bez zdania trzeciego Operator czyta je jak awarię.
 *
 * Nakładka nie stawia bramki przed wywołaniem: pustego pola nie blokuje,
 * przycisku nie wyszarza i nie pyta o potwierdzenie. Rozstrzyga rdzeń,
 * a nakładka nazywa jego odpowiedź po ludzku.
 */

/** Rada domyślna, gdy kod odmowy nie ma osobnego zdania. */
const RADA_OGOLNA = 'Odczytaj stan nakładki ponownie i powtórz czynność.';

const RADY: Record<string, Record<string, string>> = {
  przypiecie: {
    validation_failed:
      'Wpisz identyfikator procesu w pole obok wykazu — rdzeń nie zgaduje, który proces obserwować.',
    not_found:
      'Rdzeń nie zna procesu o tym identyfikatorze — sprawdź go w wykazie procesów i wpisz ponownie.',
  },
  odpiecie: {
    validation_failed:
      'Wskaż proces do odpięcia — puste pole nie mówi rdzeniowi, które przypięcie zdjąć.',
    not_found:
      'Tego procesu nie ma już wśród przypiętych — nic nie trzeba odpinać. ' +
      'Przypięcie mogło zostać zdjęte z innego urządzenia; odczytaj stan ponownie, aby zobaczyć wykaz zgodny z rdzeniem.',
  },
  rozmowa: {
    validation_failed:
      'Wpisz treść wiadomości; okno zostaw puste, jeśli ma ją przyjąć okno ogniskowane.',
    not_found:
      'Rdzeń nie znalazł okna rozmowy — zostaw pole okna puste, żeby wiadomość poszła do okna ogniskowanego, albo otwórz okno rozmowy.',
  },
  glos: {
    validation_failed:
      'Wpisz treść polecenia — rdzeń nie przyjmie wywołania bez transkrypcji, bo nagrania nakładka nie robi.',
    not_found:
      'Rdzeń nie znalazł okna asystenta — zostaw pole okna puste, żeby polecenie poszło do okna ogniskowanego.',
  },
  podpowiedzi: {
    not_found: 'Rdzeń nie ma dziś podpowiedzi dla tego okna — wykaz zostaje pusty, to stan poprawny.',
  },
  stan: {
    not_found:
      'Rdzeń nie widzi nakładki na tym urządzeniu — otwórz sesję i odczytaj stan ponownie.',
  },
  kontekst: {
    not_found:
      'Rdzeń nie zna okna, o którego kontekst pytamy — ognisko mogło się zmienić; odczytaj stan ponownie.',
  },

  // ── kolejka decyzji: odmowy komend spoza rodziny `aod.*` ──────────────────
  // Rady mówią o czynności, której odmówiono, a nie o błędzie — część tych
  // odmów jest zachowaniem poprawnym.

  procesy: {
    not_found:
      'Rdzeń nie zna żadnego procesu — telemetria jest pusta. To stan poprawny, gdy nic nie biegnie; odczytaj ponownie po uruchomieniu pracy.',
    validation_failed:
      'Odczyt procesów idzie bez argumentów i rdzeń go tak przyjmuje — jeśli odmawia, zgłoś to jako usterkę rdzenia, bo nakładka nie ma czego poprawić.',
  },
  wstrzymanie: {
    not_found:
      'Rdzeń nie zna tej kolejki, tego okna albo tej sesji — byt mógł się zamknąć, odkąd nakładka go zobaczyła. Odczytaj kolejkę decyzji ponownie.',
    validation_failed:
      'Rdzeń nie przyjął wskazania bytu do wstrzymania — odczytaj kolejkę decyzji ponownie i powtórz czynność na świeżym wpisie.',
  },
  konfiguracja: {
    not_found:
      'Rdzeń nie zna okna, dla którego konfiguracja miała się otworzyć — okno mogło zostać zamknięte. Odczytaj kolejkę decyzji ponownie.',
    validation_failed:
      'Rdzeń nie przyjął zasięgu okna dla tej konfiguracji — otwórz konfigurację z listwy Ustawienia i wybierz zasięg tam.',
  },
  przejecie: {
    not_found:
      'Rdzeń nie zna okna ani sesji, które miały zostać przejęte — byt mógł się zamknąć. Odczytaj kolejkę decyzji ponownie.',
    validation_failed:
      'Rdzeń nie przyjął wskazania okna do przejęcia. Ognisko wymaga identyfikatora klienta z powitania połączenia, którego nakładka dziś nie dostaje — wyjęcie okna z pętli działa niezależnie od niego.',
  },
};

/**
 * Składa pełne zdanie odmowy: co się nie udało, dlaczego (wprost z rdzenia)
 * i czym Operator to zmieni.
 *
 * @param czynnosc nazwa czynności widziana przez Operatora, np. „Odpięcie procesu”.
 * @param obszar klucz rad — `przypiecie`, `odpiecie`, `rozmowa`, `glos`,
 *   `podpowiedzi`, `stan`, `kontekst`, `procesy`, `wstrzymanie`, `konfiguracja`,
 *   `przejecie`.
 */
export function opisOdmowyAod(
  czynnosc: string,
  obszar: string,
  blad?: PowodOdmowy | null,
): string {
  const kod = (blad?.code ?? '').trim();
  const rada = RADY[obszar]?.[kod] ?? RADA_OGOLNA;
  return `${opisOdmowyBledu(czynnosc, blad)}. ${rada}`;
}
