// Operacje platformowe: rodziny komend bez własnego okna modułu (nastawy,
// poczta, izolacja, pamięć, kolejka i pozostałe globalne) dostają katalog
// otwierany z Centrum. Każda pozycja wywołuje własną komendę rodziny.

import type { Kanal } from '../protokol/kanal.ts';
import { zwiazKatalogModulu } from './katalog-modulu.ts';

const RODZINY_PLATFORMY: readonly [string, string][] = [
  ['config', 'Nastawy'],
  ['settings', 'Ustawienia'],
  ['mail', 'Poczta'],
  ['isolation', 'Izolacja'],
  ['memory', 'Pamięć'],
  ['queue', 'Kolejka zadań'],
  ['orchestration', 'Orkiestracja'],
  ['schedule', 'Harmonogram'],
  ['automation', 'Automatyzacje'],
  ['extension', 'Rozszerzenia'],
  ['image', 'Obraz'],
  ['speech', 'Mowa'],
  ['knowledge', 'Wiedza'],
  ['access', 'Dostępy'],
  ['identity', 'Tożsamość'],
  ['role', 'Role'],
  ['account', 'Konta'],
  ['session', 'Sesje'],
  ['channel', 'Kanały modeli'],
  ['health', 'Kondycja'],
  ['alert', 'Alerty'],
  ['aod', 'Nakładka ekranowa'],
  ['history', 'Historia sesji'],
  ['diagnostics', 'Diagnostyka'],
  ['provenance', 'Prowenancja'],
  ['notification', 'Powiadomienia'],
  ['clipboard', 'Schowek'],
  ['snippet', 'Wstawki tekstowe'],
  ['tools', 'Narzędzia modelu'],
  ['media', 'Multimedia'],
  ['usage', 'Zużycie'],
  ['archive', 'Archiwum'],
  ['context', 'Kontekst'],
  ['launcher', 'Wyzwalacz'],
  ['monitor', 'Monitor'],
  ['document', 'Dokumenty'],
  ['subagent', 'Podagenci'],
  ['team', 'Zespoły'],
  ['mobile', 'Mobile'],
  ['component', 'Komponenty'],
  ['window', 'Okna'],
  ['model', 'Model'],
  ['retention', 'Retencja'],
  ['action', 'Akcje'],
  ['advisor', 'Doradca'],
];

let otwarte: HTMLElement | null = null;

export function otworzOperacjePlatformy(kanal: Kanal): void {
  const dokument = globalThis.document;
  if (otwarte !== null) {
    otwarte.remove();
    otwarte = null;
    return;
  }
  const nakladka = dokument.createElement('div');
  nakladka.className = 'dn-modal-nakladka';
  nakladka.setAttribute('role', 'dialog');
  nakladka.setAttribute('aria-label', 'Operacje platformy');

  const okno = dokument.createElement('div');
  okno.className = 'dn-modal dn-modal--szeroki';

  const belka = dokument.createElement('div');
  belka.className = 'dn-modal-belka';
  const tytul = dokument.createElement('b');
  tytul.textContent = 'Operacje platformy';
  const zamknij = dokument.createElement('button');
  zamknij.className = 'dn-btn-ikona';
  zamknij.type = 'button';
  zamknij.setAttribute('aria-label', 'Zamknij');
  zamknij.textContent = '✕';
  belka.append(tytul, zamknij);
  okno.append(belka);

  const obszar = dokument.createElement('div');
  obszar.className = 'sta-obszar';
  okno.append(obszar);
  nakladka.append(okno);
  dokument.body.append(nakladka);
  otwarte = nakladka;

  for (const [rodzina, nazwa] of RODZINY_PLATFORMY) {
    zwiazKatalogModulu(kanal, '', obszar, rodzina, nazwa);
  }

  const zamknijPanel = (): void => {
    nakladka.remove();
    otwarte = null;
  };
  zamknij.addEventListener('click', zamknijPanel);
  nakladka.addEventListener('click', (zdarzenie) => {
    if (zdarzenie.target === nakladka) zamknijPanel();
  });
}
