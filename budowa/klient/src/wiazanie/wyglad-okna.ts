// Motyw i gęstość interfejsu: wybór idzie do rdzenia i wraca przy otwarciu.
import { ConfigScope, Command } from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';

const NAGLOWEK = 'Wygląd';

const KLUCZ_MOTYWU = 'wyglad.motyw';
const KLUCZ_GESTOSCI = 'wyglad.gestosc';

const NAZWY_MOTYWU: Readonly<Record<string, string>> = {
  light: 'jasny',
  dark: 'ciemny',
  system: 'systemowy',
};

const NAZWY_GESTOSCI: Readonly<Record<string, string>> = {
  zwarta: 'zwarta',
  luzna: 'luźna',
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
  }, true);
}

// Brak nastawy zostawia stan z prototypu — to on niesie wartość domyślną.
async function wczytajWyglad(kanal: Kanal): Promise<void> {
  for (const klucz of [KLUCZ_MOTYWU, KLUCZ_GESTOSCI]) {
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

function nanies(klucz: string, wartosc: string): void {
  const motyw = klucz === KLUCZ_MOTYWU;
  const pozycja = document.getElementById(motyw ? 'us-motyw' : 'us-gestosc');
  const wybor = motyw ? 'data-motyw' : 'data-gestosc';
  if (pozycja !== null) {
    for (const zakladka of pozycja.querySelectorAll<HTMLElement>(`[${wybor}]`)) {
      const wlasna = zakladka.getAttribute(wybor) === wartosc;
      zakladka.setAttribute('aria-checked', wlasna ? 'true' : 'false');
    }
  }
  const stan = document.getElementById(motyw ? 'us-motyw-stan' : 'us-gestosc-stan');
  const nazwy = motyw ? NAZWY_MOTYWU : NAZWY_GESTOSCI;
  if (stan !== null) stan.textContent = `wartość: ${nazwy[wartosc] ?? wartosc}`;
  if (!motyw) {
    document.documentElement.dataset.gestosc = wartosc;
    return;
  }
  // Motyw systemowy zdejmuje znacznik — arkusz idzie wtedy za systemem.
  if (wartosc === 'system') delete document.documentElement.dataset.theme;
  else document.documentElement.dataset.theme = wartosc;
}
