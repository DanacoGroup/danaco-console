// Wiązanie zmiany adresu konta: zamówienie, potwierdzenie kodem i wycofanie.
import { AuthMethodKind, Command } from '../../../shared/contract.ts';
import type { AuthMethod } from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';
import { urzadzenieTrwale } from './wejscie-katalog.ts';

const NAGLOWEK = 'Adres konta';

const PARAMETR_DROGI = 'droga';

const SCIEZKA_WYCOFANIA = '/wycofaj-zmiane-adresu';

// Wycofanie idzie samą drogą z zapytania: Operator otwiera je z poczty.
export function zwiazZmianeAdresu(kanal: Kanal): void {
  postawUstawieniaZDrogi();
  void wycofajZDrogiWejscia(kanal);
  document.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    if (cel.closest('[data-adres-zamow]') !== null) void zamowZmiane(kanal);
    if (cel.closest('[data-adres-potwierdz]') !== null) void potwierdzZmiane(kanal);
    if (cel.closest('[data-pin-ustaw]') !== null) void ustawPin(kanal);
    const przelacznik = cel.closest<HTMLElement>('[data-metoda="pin"] [data-metoda-przel]');
    if (przelacznik !== null) void przelaczPin(kanal, przelacznik);
  }, true);
}

async function zamowZmiane(kanal: Kanal): Promise<void> {
  const adres = document.querySelector<HTMLInputElement>('[data-adres-nowy]');
  const haslo = document.querySelector<HTMLInputElement>('[data-adres-haslo]');
  if (adres === null || haslo === null) return;
  const wynik = await wywolaj(kanal, Command.AuthEmailChangeStart, {
    newEmail: adres.value,
    currentPassword: haslo.value,
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zamówienia zmiany adresu.', 'ostrzezenie');
    return;
  }
  haslo.value = '';
  // Odpowiedź rdzenia nie zdradza, czy adres jest wolny.
  oglos(NAGLOWEK, 'Kod poszedł na nowy adres, ostrzeżenie na dotychczasowy.');
}

// ustawPin zakłada metodę szybkiego wejścia na tym urządzeniu. PIN nie zastępuje
// hasła — otwiera bramkę obok niego, na urządzeniu, które Operator już zna.
async function ustawPin(kanal: Kanal): Promise<void> {
  const pole = document.querySelector<HTMLInputElement>('[data-pin-wartosc]');
  if (pole === null) return;
  const wynik = await wywolaj(kanal, Command.AuthMethodAdd, {
    kind: AuthMethodKind.Pin,
    deviceId: urzadzenieTrwale(),
    secret: pole.value,
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił założenia PIN-u.', 'ostrzezenie');
    return;
  }
  pole.value = '';
  nanieMetody(wynik.wynik?.methods ?? []);
  oglos(NAGLOWEK, 'PIN działa na tym urządzeniu.');
}

/*
przelaczPin zdejmuje PIN z urządzenia; założenie idzie polem, nie przełącznikiem.

Przełącznik przestawiony w prawo bez podanego PIN-u nie ma czym założyć metody,
więc wraca do pozycji wyjściowej i mówi, gdzie PIN się wpisuje.
*/
async function przelaczPin(kanal: Kanal, przelacznik: HTMLElement): Promise<void> {
  if (przelacznik.getAttribute('aria-checked') !== 'true') {
    oglos(NAGLOWEK, 'Podaj PIN w polu obok i naciśnij „Ustaw PIN”.');
    return;
  }
  const metoda = przelacznik.closest<HTMLElement>('[data-metoda]')?.dataset.metodaId;
  if (metoda === undefined || metoda === '') {
    oglos(NAGLOWEK, 'Rdzeń nie podał, którą metodę zdjąć.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.AuthMethodRemove, { methodId: metoda });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zdjęcia PIN-u.', 'ostrzezenie');
    return;
  }
  nanieMetody(wynik.wynik?.methods ?? []);
  oglos(NAGLOWEK, 'PIN zdjęty z tego urządzenia.');
}

/*
nanieMetody nanosi stan metod z rdzenia na przełączniki okna.

Prototyp niesie stan zastany — PIN wyłączony, e-mail włączony — a po zmianie
okno musi pokazywać to, co rdzeń faktycznie prowadzi, nie to, co narysował
prototyp. Identyfikator metody zostaje przy pozycji, bo zdjęcie go potrzebuje.
*/
function nanieMetody(metody: readonly AuthMethod[]): void {
  for (const pozycja of document.querySelectorAll<HTMLElement>('[data-metoda]')) {
    const rodzaj = pozycja.dataset.metoda;
    if (rodzaj === undefined) continue;
    const wpis = metody.find((m) => m.kind === rodzaj);
    const przelacznik = pozycja.querySelector<HTMLElement>('[data-metoda-przel]');
    if (przelacznik !== null) {
      przelacznik.setAttribute('aria-checked', wpis === undefined ? 'false' : 'true');
    }
    if (wpis === undefined) delete pozycja.dataset.metodaId;
    else pozycja.dataset.metodaId = wpis.id;
  }
}

// potwierdzZmiane domyka czynność kodem z listu.
async function potwierdzZmiane(kanal: Kanal): Promise<void> {
  const kod = document.querySelector<HTMLInputElement>('[data-adres-kod]');
  if (kod === null) return;
  const wynik = await wywolaj(kanal, Command.AuthEmailChangeConfirm, { code: kod.value });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił potwierdzenia zmiany adresu.', 'ostrzezenie');
    return;
  }
  kod.value = '';
  oglos(NAGLOWEK, `Konto stoi przy adresie ${wynik.wynik?.email ?? ''}.`);
}

// postawUstawieniaZDrogi pokazuje Ustawienia, gdy adres niesie drogę wycofania:
// poza Centrum nikt tej treści nie stawia.
function postawUstawieniaZDrogi(): void {
  if (window.location.pathname !== SCIEZKA_WYCOFANIA) return;
  const szablon = document.getElementById('dn-tresc-ustawienia');
  if (!(szablon instanceof HTMLTemplateElement)) return;
  const blok = szablon.content.firstElementChild;
  if (blok === null) return;
  const rama = document.querySelector<HTMLElement>('[data-rama-aplikacji]');
  if (rama === null) return;
  rama.hidden = false;
  rama.replaceChildren(blok.cloneNode(true));
}

// wycofajZDrogiWejscia działa bez kliknięcia: droga stoi w zapytaniu adresu.
async function wycofajZDrogiWejscia(kanal: Kanal): Promise<void> {
  const droga = new URLSearchParams(window.location.search).get(PARAMETR_DROGI);
  if (droga === null || droga === '') return;
  const wynik = await wywolaj(kanal, Command.AuthEmailChangeRevoke, { token: droga });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił wycofania zmiany adresu.', 'ostrzezenie');
    return;
  }
  if (wynik.wynik?.revoked === true) {
    oglos(NAGLOWEK, 'Zmiana adresu wycofana; konto zostaje przy adresie dotychczasowym.');
    return;
  }
  // Rdzeń nie rozróżnia drogi nieznanej, użytej i przeterminowanej.
  oglos(NAGLOWEK, 'Droga wycofania jest już nieważna.', 'ostrzezenie');
}
