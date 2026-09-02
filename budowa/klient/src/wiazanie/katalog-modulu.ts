// Generyczny katalog operacji modułu: buduje wykaz komend rodziny z kontraktu
// (jedyne źródło prawdy) i wywołuje wskazaną komendę z parametrami Operatora.
// Rodzina to przedrostek komendy przed pierwszą kropką, np. „research".

import { Command } from '../../../shared/contract.ts';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';

interface PozycjaKatalogu {
  komenda: Command;
  nazwa: string;
  grupa: string;
}

function operacjeRodziny(rodzina: string): PozycjaKatalogu[] {
  const przedrostek = `${rodzina}.`;
  const pozycje: PozycjaKatalogu[] = [];
  for (const [, wartosc] of Object.entries(Command)) {
    if (typeof wartosc !== 'string') continue;
    if (!wartosc.startsWith(przedrostek) || wartosc.endsWith('.unknown')) continue;
    const czlony = wartosc.slice(przedrostek.length).split('.');
    pozycje.push({
      komenda: wartosc as Command,
      nazwa: czlony.join(' · '),
      grupa: czlony[0] ?? 'ogólne',
    });
  }
  pozycje.sort((a, b) => a.komenda.localeCompare(b.komenda));
  return pozycje;
}

export function zwiazKatalogModulu(
  kanal: Kanal,
  idOkna: string,
  korzen: ParentNode,
  rodzina: string,
  nazwaModulu: string,
): Odsubskrybuj | null {
  const gniazdo = korzen.querySelector('.sta-obszar') ?? (korzen as Element).querySelector?.('*') ?? null;
  const cel = korzen.querySelector('.sta-obszar');
  if (cel === null) return null;

  const dokument = cel.ownerDocument ?? globalThis.document;
  const sekcja = dokument.createElement('div');
  sekcja.className = 'dn-katalog-modulu';

  const naglowek = dokument.createElement('div');
  naglowek.className = 'dn-etyk-mono';
  naglowek.textContent = `Operacje: ${nazwaModulu}`;
  sekcja.append(naglowek);

  const wybor = dokument.createElement('select');
  wybor.className = 'dn-pole';
  wybor.setAttribute('aria-label', `Operacja modułu ${nazwaModulu}`);
  const komendaPoWartosci = new Map<string, Command>();
  let grupaBiezaca: string | null = null;
  let optgroup: HTMLOptGroupElement | null = null;
  const pozycje = operacjeRodziny(rodzina);
  for (const pozycja of pozycje) {
    if (pozycja.grupa !== grupaBiezaca) {
      grupaBiezaca = pozycja.grupa;
      optgroup = dokument.createElement('optgroup');
      optgroup.label = pozycja.grupa;
      wybor.append(optgroup);
    }
    const wartosc = String(pozycja.komenda);
    komendaPoWartosci.set(wartosc, pozycja.komenda);
    const opcja = dokument.createElement('option');
    opcja.value = wartosc;
    opcja.textContent = pozycja.nazwa;
    (optgroup ?? wybor).append(opcja);
  }
  sekcja.append(wybor);

  if (pozycje.length === 0) {
    naglowek.textContent = `Moduł ${nazwaModulu} nie ma operacji w kontrakcie.`;
    cel.append(sekcja);
    return () => {
      sekcja.remove();
    };
  }

  const parametry = dokument.createElement('textarea');
  parametry.className = 'dn-pole';
  parametry.rows = 3;
  parametry.placeholder = 'Parametry JSON';
  parametry.setAttribute('aria-label', 'Parametry JSON');
  sekcja.append(parametry);

  const uruchom = dokument.createElement('button');
  uruchom.className = 'dn-btn dn-btn--sygnal dn-btn--sm';
  uruchom.type = 'button';
  uruchom.textContent = 'Uruchom operację';
  sekcja.append(uruchom);

  cel.append(sekcja);
  void gniazdo;

  const naKliknieciu = () => {
    void wykonaj();
  };
  uruchom.addEventListener('click', naKliknieciu);

  async function wykonaj(): Promise<void> {
    const komenda = komendaPoWartosci.get(wybor.value);
    if (komenda === undefined) return;
    let dodatkowe: Record<string, unknown> = {};
    const surowe = parametry.value.trim();
    if (surowe !== '') {
      try {
        const rozebrane = JSON.parse(surowe) as unknown;
        if (rozebrane !== null && typeof rozebrane === 'object') {
          dodatkowe = rozebrane as Record<string, unknown>;
        }
      } catch {
        oglos(nazwaModulu, 'Parametry operacji nie są poprawnym JSON.', 'ostrzezenie');
        return;
      }
    }
    const zadanie: Record<string, unknown> = { ...dodatkowe };
    if (idOkna !== '' && !('windowId' in zadanie)) zadanie['windowId'] = idOkna;
    const wynik = await wywolaj(kanal, komenda, zadanie as never);
    if (!wynik.udany) {
      oglos(nazwaModulu, wynik.blad?.message ?? 'Rdzeń odmówił wykonania operacji.', 'ostrzezenie');
    }
  }

  return () => {
    uruchom.removeEventListener('click', naKliknieciu);
    sekcja.remove();
  };
}
