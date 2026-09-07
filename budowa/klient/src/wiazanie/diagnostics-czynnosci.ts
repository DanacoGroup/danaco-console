// Czynności okna Diagnostics: uruchomienie analizy błędów i odczyt dziennika
// rdzenia. Okno wypisywało dotąd same wykazy błędów i rekomendacji.
import { Command } from '../../../shared/contract.ts';
import type { LogEntry } from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';

const NAGLOWEK = 'Diagnostyka';

/*
zwiazCzynnosciDiagnostyki stawia pas nad panelem dziennika.

Pas stoi nad ciałem panelu, bo ciało jest wymieniane przy każdym odświeżeniu
wykazu; pas postawiony w nim znikałby razem z wpisami.
*/
export function zwiazCzynnosciDiagnostyki(
  kanal: Kanal,
  korzen: Element,
  idOkna: () => string,
  odswiez: () => void,
  przy: AddEventListenerOptions,
): void {
  postawPas(korzen);
  korzen.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const czynnosc = cel.closest<HTMLElement>('[data-diagnostyka]')?.dataset.diagnostyka;
    if (czynnosc === undefined) return;
    zdarzenie.stopPropagation();
    void wykonaj(kanal, korzen, czynnosc, idOkna(), odswiez);
  }, przy);
}

function postawPas(korzen: Element): void {
  const panel = korzen.querySelector('#panel-logi');
  const cialo = panel?.querySelector('.sta-okno-tresc');
  if (panel === null || panel === undefined || cialo === null || cialo === undefined) return;
  const pas = korzen.ownerDocument.createElement('div');
  pas.className = 'dn-pas-dzialan';
  const pole = korzen.ownerDocument.createElement('input');
  pole.type = 'text';
  pole.className = 'dn-pole dn-pole--sm';
  pole.placeholder = 'Wzorzec wyszukiwania w dzienniku';
  pole.setAttribute('aria-label', 'Wzorzec wyszukiwania w dzienniku');
  pole.dataset.dziennikWzorzec = '';
  pas.append(pole, przycisk(korzen, 'dziennik', 'Odczytaj dziennik'),
    przycisk(korzen, 'analiza', 'Uruchom analizę'));
  panel.insertBefore(pas, cialo);
}

function przycisk(korzen: Element, czynnosc: string, etykieta: string): HTMLButtonElement {
  const wezel = korzen.ownerDocument.createElement('button');
  wezel.type = 'button';
  wezel.className = 'dn-btn dn-btn--duch dn-btn--sm';
  wezel.dataset.diagnostyka = czynnosc;
  wezel.textContent = etykieta;
  return wezel;
}

async function wykonaj(
  kanal: Kanal,
  korzen: Element,
  czynnosc: string,
  idOkna: string,
  odswiez: () => void,
): Promise<void> {
  if (czynnosc === 'dziennik') return odczytajDziennik(kanal, korzen);
  if (czynnosc === 'analiza') return uruchomAnalize(kanal, idOkna, odswiez);
  /* Czynność spoza obsłużonych odmawia zamiast milczeć: przycisk bez gałęzi
     wyglądałby jak działający. */
  oglos(NAGLOWEK, `Czynność „${czynnosc}" nie jest prowadzona przez to okno.`, 'ostrzezenie');
}

/* Wynik przycięty granicą bufora jest nazwany: wykaz bez tej informacji
   wyglądałby na komplet dziennika. */
async function odczytajDziennik(kanal: Kanal, korzen: Element): Promise<void> {
  const cialo = korzen.querySelector('#panel-logi .sta-okno-tresc');
  if (cialo === null) return;
  const pole = korzen.querySelector<HTMLInputElement>('[data-dziennik-wzorzec]');
  const wzorzec = (pole?.value ?? '').trim();
  const wynik = await wywolaj(kanal, Command.DiagnosticsLogQuery, {
    ...(wzorzec === '' ? {} : { pattern: wzorzec }),
    limit: 100,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu dziennika.', 'ostrzezenie');
    return;
  }
  const wpisy = wynik.wynik.entries;
  cialo.replaceChildren(...wpisy.map((wpis) => wierszWpisu(cialo, wpis)));
  if (wpisy.length === 0) {
    oglos(NAGLOWEK, 'Żaden wpis dziennika nie spełnia tych warunków.');
    return;
  }
  if (wynik.wynik.truncated === true) {
    oglos(NAGLOWEK, 'Wynik przycięty granicą bufora — dziennik ma więcej wpisów.', 'ostrzezenie');
  }
}

function wierszWpisu(cialo: Element, wpis: LogEntry): HTMLElement {
  const wiersz = cialo.ownerDocument.createElement('div');
  wiersz.className = 'dn-wykaz-modulu-poz';
  const tresc = cialo.ownerDocument.createElement('span');
  tresc.textContent = wpis.message;
  const meta = cialo.ownerDocument.createElement('span');
  meta.className = 'dn-meta';
  meta.textContent = `${wpis.level} · ${wpis.source ?? 'rdzeń'}`;
  wiersz.append(tresc, meta);
  return wiersz;
}

/* Analiza obejmuje całe okno: zawężenie zakresem czasu ani wskazaniem błędów
   nie ma w oknie pola, a zmyślenie granic byłoby decyzją za Operatora. */
async function uruchomAnalize(
  kanal: Kanal,
  idOkna: string,
  odswiez: () => void,
): Promise<void> {
  if (idOkna === '') {
    oglos(NAGLOWEK, 'Rdzeń nie dał okna diagnostyki dla tej karty.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.DiagnosticsAnalyzeRun, { windowId: idOkna });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił uruchomienia analizy.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Analiza wykonana; rekomendacje przeliczone na nowo.');
  odswiez();
}
