// Katalog operacji Studia wstawiony do Tools Panel: każdą pozycję Operator
// uruchamia własną komendą studio.*. Żądanie składa się z identyfikatora
// dokumentu bieżącego i parametrów wpisanych w polu panelu.

import { Command } from '../../../shared/contract.ts';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';
import { NAZWY_GRUP, OPERACJE_STUDIA, type GrupaOperacji } from './studio-operacje.ts';

export function zwiazKatalogOperacji(
  kanal: Kanal,
  idOkna: string,
  korzen: ParentNode,
): Odsubskrybuj | null {
  const panel = korzen.querySelector('#panel-tools .sta-okno-tresc');
  if (panel === null) return null;

  const dokument = korzen.ownerDocument ?? globalThis.document;
  const sekcja = dokument.createElement('div');
  sekcja.className = 'st-katalog-operacji';

  const naglowek = dokument.createElement('div');
  naglowek.className = 'dn-etyk-mono';
  naglowek.textContent = 'Katalog operacji dokumentu';
  sekcja.append(naglowek);

  const wybor = dokument.createElement('select');
  wybor.className = 'dn-pole';
  wybor.setAttribute('aria-label', 'Operacja Studia');
  const komendaPoWartosci = new Map<string, Command>();
  let grupaBiezaca: GrupaOperacji | null = null;
  let optgroup: HTMLOptGroupElement | null = null;
  for (const operacja of OPERACJE_STUDIA) {
    if (operacja.grupa !== grupaBiezaca) {
      grupaBiezaca = operacja.grupa;
      optgroup = dokument.createElement('optgroup');
      optgroup.label = NAZWY_GRUP[operacja.grupa];
      wybor.append(optgroup);
    }
    const wartosc = String(operacja.komenda);
    komendaPoWartosci.set(wartosc, operacja.komenda);
    const opcja = dokument.createElement('option');
    opcja.value = wartosc;
    opcja.textContent = operacja.nazwa;
    (optgroup ?? wybor).append(opcja);
  }
  sekcja.append(wybor);

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

  panel.append(sekcja);

  let idDokumentu = '';
  async function dokumentBiezacy(): Promise<string> {
    if (idDokumentu !== '') return idDokumentu;
    const wynik = await wywolaj(kanal, Command.StudioDocumentOpen, { windowId: idOkna });
    if (!wynik.udany || wynik.wynik === undefined) return '';
    idDokumentu = wynik.wynik.document.id;
    return idDokumentu;
  }

  const uruchomOperacje = () => {
    void wykonaj();
  };
  uruchom.addEventListener('click', uruchomOperacje);

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
        oglos('Studio', 'Parametry operacji nie są poprawnym JSON.', 'ostrzezenie');
        return;
      }
    }
    const dok = await dokumentBiezacy();
    const zadanie: Record<string, unknown> = { windowId: idOkna, ...dodatkowe };
    if (dok !== '') zadanie['documentId'] = dok;
    const wynik = await wywolaj(kanal, komenda, zadanie as never);
    if (!wynik.udany) {
      oglos('Studio', wynik.blad?.message ?? 'Rdzeń odmówił wykonania operacji.', 'ostrzezenie');
    }
  }

  return () => {
    uruchom.removeEventListener('click', uruchomOperacje);
    sekcja.remove();
  };
}
