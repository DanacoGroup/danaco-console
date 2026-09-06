// Wiązanie zmiany adresu konta: zamówienie, potwierdzenie kodem i wycofanie.
import { Command } from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';

const NAGLOWEK = 'Adres konta';

// Pozycja zapytania niesiona odsyłaczem z listu; nazwę składa rdzeń.
const PARAMETR_DROGI = 'droga';

// Trasa odsyłacza z listu; tę samą składa rdzeń w konfiguracji.
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

/*
postawUstawieniaZDrogi pokazuje okno Ustawień, gdy adres niesie drogę wycofania.

Odsyłacz z listu ostrzegawczego prowadzi wprost do trasy wycofania, a poza
Centrum nikt tej treści nie stawia — bez tego Operator klikał w list i trafiał
na okno wejścia, nie na sekcję Konto, o której mówi wiadomość.
*/
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
