// Motyw i gęstość interfejsu: wybór idzie do rdzenia i wraca przy otwarciu.
import { ConfigScope, Command } from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';

const NAGLOWEK = 'Wygląd';

const KLUCZ_MOTYWU = 'wyglad.motyw';
const KLUCZ_GESTOSCI = 'wyglad.gestosc';
// Klucz czytany przez rdzeń przy wysyłce listu o zakończonym przebiegu.
const KLUCZ_LISTOW = 'poczta.przebiegi.kiedy';

const NAZWY_MOTYWU: Readonly<Record<string, string>> = {
  light: 'jasny',
  dark: 'ciemny',
  system: 'systemowy',
};

const NAZWY_GESTOSCI: Readonly<Record<string, string>> = {
  zwarta: 'zwarta',
  luzna: 'luźna',
};

const NAZWY_LISTOW: Readonly<Record<string, string>> = {
  zawsze: 'zawsze',
  niepowodzenia: 'tylko niepowodzenia',
  nigdy: 'nigdy',
};

export function zwiazWyglad(kanal: Kanal): void {
  void wczytajWyglad(kanal);
  document.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const motyw = cel.closest<HTMLElement>('[data-motyw]');
    if (motyw !== null) void ustawWyglad(kanal, KLUCZ_MOTYWU, motyw.dataset.motyw ?? '');
    const gestosc = cel.closest<HTMLElement>('[data-gestosc]');
    if (gestosc !== null) void ustawWyglad(kanal, KLUCZ_GESTOSCI, gestosc.dataset.gestosc ?? '');
    const listy = cel.closest<HTMLElement>('[data-listy]');
    if (listy !== null) void ustawWyglad(kanal, KLUCZ_LISTOW, listy.dataset.listy ?? '');
  }, true);
}

// Brak nastawy zostawia stan z prototypu — to on niesie wartość domyślną.
async function wczytajWyglad(kanal: Kanal): Promise<void> {
  for (const klucz of [KLUCZ_MOTYWU, KLUCZ_GESTOSCI, KLUCZ_LISTOW]) {
    const wynik = await wywolaj(kanal, Command.ConfigGet, {
      key: klucz,
      scope: ConfigScope.Application,
    });
    // Odpowiedź niesie wykaz wpisów, nie jedną wartość.
    const wpis = wynik.udany ? wynik.wynik?.entries.find((e) => e.key === klucz) : undefined;
    const wartosc = typeof wpis?.value === 'string' ? wpis.value : '';
    if (wartosc !== '') nanies(klucz, wartosc);
  }
}

// Wybór nanosi się po potwierdzeniu: nastawa, której rdzeń nie przyjął, kłamie.
async function ustawWyglad(kanal: Kanal, klucz: string, wartosc: string): Promise<void> {
  if (wartosc === '') return;
  const wynik = await wywolaj(kanal, Command.ConfigSet, {
    key: klucz,
    value: wartosc,
    scope: ConfigScope.Application,
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zapisu nastawy.', 'ostrzezenie');
    return;
  }
  nanies(klucz, wartosc);
}

interface Pozycja {
  grupa: string;
  wybor: string;
  stan: string;
  nazwy: Readonly<Record<string, string>>;
}

const POZYCJE: Readonly<Record<string, Pozycja>> = {
  [KLUCZ_MOTYWU]: {
    grupa: 'us-motyw', wybor: 'data-motyw', stan: 'us-motyw-stan', nazwy: NAZWY_MOTYWU,
  },
  [KLUCZ_GESTOSCI]: {
    grupa: 'us-gestosc', wybor: 'data-gestosc', stan: 'us-gestosc-stan', nazwy: NAZWY_GESTOSCI,
  },
  [KLUCZ_LISTOW]: {
    grupa: 'us-listy-przebiegow', wybor: 'data-listy', stan: 'us-listy-stan', nazwy: NAZWY_LISTOW,
  },
};

function nanies(klucz: string, wartosc: string): void {
  const opis = POZYCJE[klucz];
  if (opis === undefined) return;
  const pozycja = document.getElementById(opis.grupa);
  if (pozycja !== null) {
    for (const zakladka of pozycja.querySelectorAll<HTMLElement>(`[${opis.wybor}]`)) {
      const wlasna = zakladka.getAttribute(opis.wybor) === wartosc;
      zakladka.setAttribute('aria-checked', wlasna ? 'true' : 'false');
    }
  }
  const stan = document.getElementById(opis.stan);
  if (stan !== null) stan.textContent = `wartość: ${opis.nazwy[wartosc] ?? wartosc}`;
  if (klucz === KLUCZ_GESTOSCI) {
    document.documentElement.dataset.gestosc = wartosc;
    return;
  }
  // Listy o przebiegach zmieniają zachowanie rdzenia, nie wygląd okna.
  if (klucz !== KLUCZ_MOTYWU) return;
  // Motyw systemowy zdejmuje znacznik — arkusz idzie wtedy za systemem.
  if (wartosc === 'system') delete document.documentElement.dataset.theme;
  else document.documentElement.dataset.theme = wartosc;
}
