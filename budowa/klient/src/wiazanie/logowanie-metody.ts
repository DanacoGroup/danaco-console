// Wykaz metod wejścia w kaflach odsłony logowania: stan metody bierze się
// z odpowiedzi rdzenia, a kafel zakłada ją, zdejmuje albo nią wchodzi.
import { AuthMethodKind, Command, type AuthMethod } from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { tokenSesji } from '../protokol/token-sesji.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';
import { napis, urzadzenieTrwale } from './wejscie-katalog.ts';
import { zdanieOdmowy } from './wejscie-odmowa.ts';

/* Kolejność kafli jest umową z `okna/wejscie/skladniki/metody-logowania.js`;
   kafel trzeci — kod listem — nie ma odpowiednika w `AuthMethodKind`. */
const KAFLE: readonly (AuthMethodKind | '')[] = [AuthMethodKind.Pin, AuthMethodKind.Hello, ''];

const NAGLOWEK = 'Metody wejścia';

let metody: AuthMethod[] = [];

export function zapamietajMetody(wykaz: AuthMethod[] | undefined): void {
  if (wykaz === undefined) return;
  metody = wykaz;
  opiszKafle();
}

export function zwiazMetody(kanal: Kanal): void {
  oznaczKafle();
  opiszKafle();
  document.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const kafel = cel.closest<HTMLElement>('.au-metody [data-metoda]');
    if (kafel === null) return;
    zdarzenie.preventDefault();
    zdarzenie.stopPropagation();
    void wykonaj(kanal, kafel.dataset.metoda ?? '');
  }, true);
}

function oznaczKafle(): void {
  for (const grupa of document.querySelectorAll<HTMLElement>('.au-metody')) {
    grupa.querySelectorAll<HTMLElement>('.dn-kafel').forEach((kafel, numer) => {
      kafel.dataset.metoda = KAFLE[numer] ?? '';
    });
  }
}

function opiszKafle(): void {
  for (const kafel of document.querySelectorAll<HTMLElement>('.au-metody [data-metoda]')) {
    const rodzaj = kafel.dataset.metoda ?? '';
    const stan = kafel.querySelector('.dn-kafel-stan');
    if (stan === null) continue;
    const czynna = rodzaj !== '' && metodaRodzaju(rodzaj) !== undefined;
    stan.textContent = napis(`dostep.logowanie.metody.${czynna ? 'aktywna' : 'nieaktywna'}`);
    kafel.dataset.metodaCzynna = czynna ? 'tak' : 'nie';
  }
}

function metodaRodzaju(rodzaj: string): AuthMethod | undefined {
  const urzadzenie = urzadzenieTrwale();
  return metody.find(
    (metoda) => metoda.kind === rodzaj
      && (metoda.deviceId === undefined || metoda.deviceId === urzadzenie),
  );
}

/* Przed wydaniem sesji bramki kafel jest drogą wejścia, po jej wydaniu —
   czynnością ustawień: `auth.method.add` działa dopiero po zalogowaniu. */
async function wykonaj(kanal: Kanal, rodzaj: string): Promise<void> {
  if (rodzaj === '') {
    oglos(NAGLOWEK, 'Kod na adres e-mail nie jest metodą bramki w kontrakcie — '
      + 'adres służy potwierdzeniu konta i odzyskaniu dostępu.', 'ostrzezenie');
    return;
  }
  if (rodzaj === AuthMethodKind.Hello) {
    oglos(NAGLOWEK, 'Klucz systemowy staje się dostępny dopiero wtedy, gdy interfejs '
      + 'jest podawany z pochodzenia z domeną po https.', 'ostrzezenie');
    return;
  }
  if (tokenSesji() === '') {
    await wejdzPinem(kanal);
    return;
  }
  const stojaca = metodaRodzaju(rodzaj);
  if (stojaca === undefined) await zalozPin(kanal);
  else await zdejmijMetode(kanal, stojaca);
}

async function wejdzPinem(kanal: Kanal): Promise<void> {
  const sekret = sekretZOdslony();
  if (sekret === '') {
    oglos(NAGLOWEK, 'Wpisz PIN w polu hasła, zanim wejdziesz kodem PIN.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.AuthLogin, {
    method: AuthMethodKind.Pin,
    secret: sekret,
    deviceId: urzadzenieTrwale(),
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, zdanieOdmowy(wynik.blad), 'blad');
    return;
  }
  zapamietajMetody(wynik.wynik?.methods);
  document.dispatchEvent(new CustomEvent('wejscie-metoda-przyjeta', { bubbles: true }));
}

async function zalozPin(kanal: Kanal): Promise<void> {
  const sekret = sekretZOdslony();
  if (sekret === '') {
    oglos(NAGLOWEK, 'Wpisz PIN w polu hasła, zanim założysz go na tym urządzeniu.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.AuthMethodAdd, {
    kind: AuthMethodKind.Pin,
    deviceId: urzadzenieTrwale(),
    secret: sekret,
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, zdanieOdmowy(wynik.blad), 'blad');
    return;
  }
  zapamietajMetody(wynik.wynik?.methods);
  oglos(NAGLOWEK, 'Kod PIN obowiązuje na tym urządzeniu.');
}

async function zdejmijMetode(kanal: Kanal, metoda: AuthMethod): Promise<void> {
  if (metoda.anchor) {
    oglos(NAGLOWEK, 'Hasło jest kotwicą bramki — zdjąć go nie można.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.AuthMethodRemove, {
    methodId: metoda.id,
    deviceId: urzadzenieTrwale(),
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, zdanieOdmowy(wynik.blad), 'blad');
    return;
  }
  zapamietajMetody(wynik.wynik?.methods);
  oglos(NAGLOWEK, 'Metoda zdjęta z tego urządzenia.');
}

function sekretZOdslony(): string {
  const pole = document.querySelector<HTMLInputElement>(
    '.we-panel[data-widok-aktywny="tak"] .au-haslo input',
  );
  return pole?.value.trim() ?? '';
}
